package kocha

import (
	"github.com/tkmcclellan/kocha/internal/providers"

	"fmt"
	"os"
)

func Search(name string, providerType string, page uint64) providers.SearchResult {
	provider, err := providers.FindProvider(providerType)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	searchResults, err := provider.Search(name, page)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	return searchResults
}
