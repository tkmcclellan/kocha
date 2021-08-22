package providers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tkmcclellan/kocha/internal/models"
)

type MangaKakalot struct{}

func (m MangaKakalot) Search(name string) (SearchResult, error) {
	var result SearchResult

	res, err := http.Get(fmt.Sprintf("https://ww.mangakakalot.tv/search/%s", name))
	if err != nil {
		return result, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return result, errors.New("mangakakalot search failed")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return result, errors.New("mangakakalot search failed")
	}

	doc.Find("div.story_item").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		var manga models.Manga

		manga.Title = s.Find("h3.story_name > a").Text()
		href, exists := s.Find("h3.story_name > a").Attr("href")
		if !exists {
			log.Fatal(errors.New("manga missing link"))
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
		manga.Authors = authors

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
			log.Fatal(err)
			return
		}
		manga.Updated = update_time

		result.Manga = append(result.Manga, manga)
	})

	r1 := regexp.MustCompile(`[0-9]+`)
	total_pages, err := strconv.ParseUint(r1.FindString(doc.Find("a.page_last").Text()), 10, 64)
	if err != nil {
		log.Fatal(err)
		return result, err
	}
	result.total_pages = total_pages

	current_page, err := strconv.ParseUint(doc.Find("a.page_select").Text(), 10, 64)
	result.current_page = current_page

	return result, nil
}
