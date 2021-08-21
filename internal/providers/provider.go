package providers

import "errors"

type Provider interface {
	Search(name string) (SearchResult, error)
}

type SearchResult struct {
	name string
}

func FindProvider(provider string) (Provider, error) {
	switch provider {
	case "mangakakolot":
		m := MangaKakolot{}
		return m, nil
	default:
		return nil, errors.New("invalid provider")
	}
}
