package models

import (
	"github.com/tkmcclellan/kocha/internal/util"
	"gorm.io/gorm"
)

type Chapter struct {
	gorm.Model
	Title   string
	Uri     string
	Read    bool
	MangaID uint
	// Manga   Manga
}

func (Chapter) TableName() string {
	return "chapters"
}

func init() {
	db, cancel := util.Database()
	defer cancel()
	db.AutoMigrate(&Chapter{})
}

func (c *Chapter) Create() *Chapter {
	db, cancel := util.Database()
	defer cancel()
	db.Create(&c)
	return c
}

func (c *Chapter) Save() *Chapter {
	db, cancel := util.Database()
	defer cancel()
	db.Save(&c)
	return c
}

func (c *Chapter) Manga() *Manga {
	db, cancel := util.Database()
	defer cancel()

	var manga Manga
	db.Model(&manga).Where("ID = ?", c.MangaID).First(&manga)
	return &manga
}

func (c *Chapter) Dirname() string {
	return util.CleanString(c.Title)
}
