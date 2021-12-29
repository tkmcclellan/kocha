package providers

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cheggaaa/pb"
	"github.com/tkmcclellan/kocha/internal/models"
	"github.com/tkmcclellan/kocha/internal/util"
)

type MangaKakalot struct{}

func (m MangaKakalot) Search(name string, page uint64) (SearchResult, error) {
	var result SearchResult

	res, err := http.Get(fmt.Sprintf("https://ww.mangakakalot.tv/search/%s?page=%d", name, page))
	if err != nil {
		return result, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return result, errors.New("mangakakalot search failed")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
		return result, errors.New("mangakakalot search failed")
	}

	doc.Find("div.story_item").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		var manga models.Manga

		manga.Provider = "mangakakalot"
		manga.Title = s.Find("h3.story_name > a").Text()
		href, exists := s.Find("h3.story_name > a").Attr("href")
		if !exists {
			fmt.Println(errors.New("manga missing link"))
			return
		}

		manga.Uri = "https://ww.mangakakalot.tv" + href

		info := s.Find("span")
		r1 := regexp.MustCompile(`\n(\s+)`)
		authors := strings.Split(
			r1.ReplaceAllString(
				strings.Trim(strings.Split(info.Nodes[0].FirstChild.Data, ":")[1], " \n"),
				"\n",
			),
			"\n",
		)
		manga.Authors = strings.Join(authors, ",")

		// TODO: Someone please save me from this hell
		raw_time := strings.Trim(strings.SplitN(info.Nodes[1].FirstChild.Data, ":", 2)[1], " ")
		r2 := regexp.MustCompile(`(.*),`)
		date := strings.Trim(r2.FindString(raw_time), ",")
		date_split := strings.Split(date, " ")
		date = fmt.Sprintf("%s %s", date_split[1], date_split[0])
		r3 := regexp.MustCompile(`,([0-9]{4})`)
		year := strings.Trim(r3.FindString(raw_time), ",")
		r4 := regexp.MustCompile(`-\s(.*)`)
		time_string := strings.TrimPrefix(r4.FindString(raw_time), "- ")
		time_suffix := strings.Split(time_string, " ")[1]

		update_time, err := time.Parse(fmt.Sprintf("02 Jan 2006 15:04 %s", time_suffix), fmt.Sprintf("%s %s %s", date, year, time_string))
		if err != nil {
			fmt.Println(err)
			return
		}
		manga.Updated = update_time

		result.Manga = append(result.Manga, manga)
	})

	r1 := regexp.MustCompile(`[0-9]+`)
	total_pages, err := strconv.ParseUint(r1.FindString(doc.Find("a.page_last").Text()), 10, 64)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	result.TotalPages = total_pages

	current_page, err := strconv.ParseUint(doc.Find("a.page_select").Text(), 10, 64)
	result.CurrentPage = current_page

	return result, nil
}

func (mangakakalot MangaKakalot) DownloadChapter(chapter models.Chapter, completed chan bool) {
	res, err := http.Get(chapter.Uri)
	if err != nil {
		fmt.Println(err)
		completed <- false
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("downloading chapter failed")
		completed <- false
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
		completed <- false
		return
	}

	selection := doc.Find("img.img-loading")
	imageCompleted := make(chan bool)
	numImages := selection.Length()
	dirname := filepath.Join(chapter.Manga().Dirname(), chapter.Dirname())
	selection.Each(func(i int, s *goquery.Selection) {
		src := s.AttrOr("data-src", "")
		formattedTitle := fmt.Sprintf("%06d", i)
		go util.DownloadImage(dirname, formattedTitle, src, imageCompleted)
	})
	for i := 0; i < numImages; i++ {
		<-imageCompleted
	}
	completed <- true

	err = util.ChapterToPdf(dirname)
	if err != nil {
		panic(err)
	}
}

func (mangakakalot MangaKakalot) DownloadManga(manga *models.Manga) error {
	res, err := http.Get(manga.Uri)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("mangakakalot search failed")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
		return errors.New("mangakakalot search failed")
	}

	chapterList := []models.Chapter{}
	doc.Find("#chapter > div > div.chapter-list > div > span > a").Each(func(i int, s *goquery.Selection) {
		title, exists := s.Attr("title")
		if !exists {
			title = s.Text()
		}
		uri := "https://ww.mangakakalot.tv" + s.AttrOr("href", "")
		chapter := models.Chapter{
			Title:   strings.Trim(title, " "),
			Uri:     uri,
			Read:    false,
			MangaID: manga.ID,
		}

		// Believe me that the other ways of getting the list in reverse order
		// were uglier than this. Hopefully this isn't super slow
		chapterList = append([]models.Chapter{chapter}, chapterList...)
	})

	// Probably very inefficient, but it ensures chapters are created in order
	for _, chapter := range chapterList {
		chapter.Create()
	}

	var downloadList []models.Chapter
	if manga.Dlmode == "dynamic" {
		downloadList = chapterList[0:2]
	} else if manga.Dlmode == "all" {
		downloadList = chapterList
	} else if manga.Dlmode == "none" {
		downloadList = []models.Chapter{}
	} else {
		panic("invalid download mode")
	}

	downloadCount := len(downloadList)
	if downloadCount > 0 {
		bar := pb.StartNew(downloadCount)
		completed := make(chan bool)
		for _, chapter := range downloadList {
			go mangakakalot.DownloadChapter(chapter, completed)
		}

		for range downloadList {
			<-completed
			bar.Increment()
		}
		bar.Finish()

	}
	return nil
}
