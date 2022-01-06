package providers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tkmcclellan/kocha/internal/models"
)

type Provider interface {
	Search(name string, page uint64) (SearchResult, error)
	DownloadManga(manga *models.Manga) error
	DownloadChapter(chapter models.Chapter, completed chan bool)
	GetManga(url string, dlmode string) (models.Manga, error)
}

type SearchResult struct {
	Manga       []models.Manga
	TotalPages  uint64
	CurrentPage uint64
}

func FindProvider(provider string) (Provider, error) {
	switch provider {
	case "mangakakalot":
		m := MangaKakalot{}
		return m, nil
	default:
		return nil, errors.New("invalid provider")
	}
}

// Return a supported provider from a given URL. This could
// possibly be done better but this was the first solution
// I came up with
func FindProviderFromUrl(url string) (Provider, error) {
	switch {
	case strings.HasPrefix(url, "https://ww1.mangakakalot.tv"), strings.HasPrefix(url, "http://ww1.mangakakalot.tv"):
		fallthrough
	case strings.HasPrefix(url, "https://ww.mangakakalot.tv"), strings.HasPrefix(url, "http://ww.mangakakalot.tv"):
		m := MangaKakalot{}
		return m, nil
	default:
		return nil, errors.New("invalid provider")
	}
}

// General function to query a document from a given URL
func FetchDocument(url string) (document *goquery.Document, err error) {
	res, err := http.Get(url)
	if err != nil {
		return document, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return document, errors.New("Failed to fetch document")
	}

	document, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
		return document, errors.New("Failed to parse document")
	}
	return document, err
}
