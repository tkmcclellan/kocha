package kocha

import (
	"fmt"

	"github.com/tkmcclellan/kocha/internal/models"
	"github.com/tkmcclellan/kocha/internal/providers"
)

func Add(manga *models.Manga, dlmode string) error {
	manga.Dlmode = dlmode
	if manga.Exists() != nil {
		fmt.Println("Manga already added!")
		return nil
	}
	manga.Create()
	provider, err := providers.FindProvider(manga.Provider)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = provider.DownloadManga(manga)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func AddFromUrl(url string, dlmode string) error {
	provider, err := providers.FindProviderFromUrl(url)
	if err != nil {
		fmt.Println(err)
		return err
	}

	manga, err := provider.GetManga(url, dlmode)
	if err != nil {
		fmt.Println(err)
		return err
	}

	manga.Create()

	err = provider.DownloadManga(&manga)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
