package kocha

import (
	"github.com/tkmcclellan/kocha/internal/models"
	"github.com/tkmcclellan/kocha/internal/util"
)

func List() []models.Manga {
	db, cancel := util.Database()
	defer cancel()
	var manga []models.Manga
	_ = db.Find(&manga)
	return manga
}
