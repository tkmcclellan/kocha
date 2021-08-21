package providers

import "errors"

type MangaKakolot struct{}

func (m MangaKakolot) Search(name string) (SearchResult, error) {
	return SearchResult{name: "test"}, errors.New("unimplemented")
}
