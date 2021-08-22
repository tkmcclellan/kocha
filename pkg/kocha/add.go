package kocha

import (
	"github.com/tkmcclellan/kocha/internal/providers"

	"errors"
)

func Add(p providers.Provider, id string) error {
	if id == "" {
		return errors.New("missing manga id")
	}

	return nil
}
