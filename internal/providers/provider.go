package providers

import (
	"errors"

	"github.com/tkmcclellan/kocha/internal/models"
)

type Provider interface {
	Search(name string, page uint64) (SearchResult, error)
	DownloadManga(manga *models.Manga) error
	DownloadChapter(chapter models.Chapter, completed chan bool)
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
