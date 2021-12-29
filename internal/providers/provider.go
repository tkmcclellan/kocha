package providers

import (
	"errors"

	"github.com/tkmcclellan/kocha/internal/models"
)

var ProviderList = make(map[string]Provider)

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
	foundProvider := ProviderList[provider]
	if foundProvider != nil {
		return foundProvider, nil
	} else {
		return nil, errors.New("Invalid provider")
	}
}

func ListProviders() []string {
	i := 0
	keys := make([]string, len(ProviderList))

	for k := range ProviderList {
		keys[i] = k
		i++
	}

	return keys
}
