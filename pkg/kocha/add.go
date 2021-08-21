package kocha

import (
	"strings"

	"github.com/tkmcclellan/kocha/internal/providers"

	"errors"
)

func Add(p string, d string, n string) error {
	if n == "" {
		return errors.New("missing manga title")
	}

	provider, err := providers.FindProvider(p)
	if err != nil {
		return err
	}

	search_string := strings.ReplaceAll(n, " ", "+")

	_, err = provider.Search(search_string)
	if err != nil {
		return err
	}

	return nil
}
