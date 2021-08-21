package kocha

import (
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

	provider.Search(n)

	return nil
}
