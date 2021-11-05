package kocha

import (
	"github.com/tkmcclellan/kocha/internal/models"
	"github.com/tkmcclellan/kocha/internal/util"
)

func Remove(manga *models.Manga) {
	err := util.DeleteDir(manga.Dirname())
	if err != nil {
		panic(err)
	}
	manga.Delete()
}
