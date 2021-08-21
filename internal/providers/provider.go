package providers

import (
	"errors"

	"github.com/tkmcclellan/kocha/internal/models"
)

type Provider interface {
	Search(name string) (SearchResult, error)
}

type SearchResult struct {
	Manga        []models.Manga
	total_pages  uint64
	current_page uint64
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
