package providers

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cheggaaa/pb"
	"github.com/tkmcclellan/kocha/internal/models"
	"github.com/tkmcclellan/kocha/internal/util"
)

type MangaKakalot struct{}

func (m MangaKakalot) Search(name string, page uint64) (SearchResult, error) {
	var result SearchResult

	doc, err := FetchDocument(fmt.Sprintf("https://ww.mangakakalot.tv/search/%s?page=%d", name, page))
	if err != nil {
		fmt.Println(err)
		return result, errors.New("mangakakalot search failed")
	}

	var parseError error
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
			parseError = err
			return
		}
		manga.Updated = update_time

		result.Manga = append(result.Manga, manga)
	})

	if parseError != nil {
		fmt.Println(parseError)
		return result, errors.New("Failed to parse manga")
	}

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

func (mangakakalot MangaKakalot) DownloadChapter(chapter models.Chapter) {
	doc, err := FetchDocument(chapter.Uri)
	if err != nil {
		fmt.Println(err)
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

	err = util.ChapterToPdf(dirname)
	if err != nil {
		fmt.Println("here")
		panic(err)
	}
}

// Download chapters of a manga according to the download mode
func (mangakakalot MangaKakalot) DownloadManga(manga *models.Manga) error {
	doc, err := FetchDocument(manga.Uri)
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
		var wg sync.WaitGroup
		bar := pb.StartNew(downloadCount)
		for _, chapter := range downloadList {
			wg.Add(1)
			go func(chapter models.Chapter) {
				mangakakalot.DownloadChapter(chapter)
				wg.Done()
				bar.Increment()
			}(chapter)
		}

		wg.Wait()
		bar.Finish()
	}
	return nil
}

// Get a manga from a provided URL
func (m MangaKakalot) GetManga(url string, dlmode string) (models.Manga, error) {
	var manga models.Manga

	doc, err := FetchDocument(url)
	if err != nil {
		fmt.Println(err)
		return manga, errors.New("mangakakalot search failed")
	}

	var parseError error
	doc.Find("ul.manga-info-text").Each(func(i int, s *goquery.Selection) {
		manga.Provider = "mangakakalot"
		manga.Title = s.Find("li > h1").Text()
		manga.Uri = url
		manga.Dlmode = dlmode

		authors := []string{}
		first_child := s.Children().Nodes[1].FirstChild
		last_child := s.Children().Nodes[1].LastChild
		position := first_child.NextSibling
		for position != last_child {
			if position.Type == 3 {
				authors = append(authors, position.FirstChild.Data)
			}
			position = position.NextSibling
		}
		manga.Authors = strings.Join(authors, ",")

		raw_time := strings.Trim(strings.SplitN(s.Children().Nodes[3].FirstChild.Data, ":", 2)[1], " ")
		date_string := raw_time[:strings.Index(raw_time, "-")-1]
		date_string = strings.ReplaceAll(date_string, ",", " ")
		time_string := raw_time[strings.Index(raw_time, "-")+1:]
		parsed_time := strings.SplitN(time_string, " ", -1)
		time_suffix := parsed_time[len(parsed_time)-1]
		update_time, err := time.Parse(fmt.Sprintf("Jan 02 2006 15:04 %s", time_suffix), fmt.Sprintf("%s %s", date_string, time_string))
		if err != nil {
			parseError = err
			return
		}
		manga.Updated = update_time
	})

	if parseError != nil {
		fmt.Println(parseError)
		return manga, errors.New("Failed to parse manga")
	}

	return manga, nil
}
