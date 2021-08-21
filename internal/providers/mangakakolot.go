package providers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

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

		fmt.Printf("%#v\n", manga)

		result.Manga = append(result.Manga, manga)
	})

	return result, errors.New("unimplemented")
}
