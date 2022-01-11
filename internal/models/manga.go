package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/tkmcclellan/kocha/internal/util"
	"gorm.io/gorm"
)

type Manga struct {
	gorm.Model
	Title          string
	Uri            string
	Authors        string
	Dlmode         string
	CurrentChapter int
	Updated        time.Time
	Provider       string
	//Chapters       []Chapter
}

func (Manga) TableName() string {
	return "manga"
}

func init() {
	db, cancel := util.Database()
	defer cancel()
	db.AutoMigrate(&Manga{})
}

func (m *Manga) Save() *Manga {
	db, cancel := util.Database()
	defer cancel()
	db.Save(&m)
	return m
}

func (m *Manga) Exists() *Manga {
	db, cancel := util.Database()
	defer cancel()
	args := map[string]interface{}{"title": m.Title, "uri": m.Uri, "authors": m.Authors}
	var manga Manga
	result := db.Where(args).First(&manga)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return &manga
	}
}

func (m *Manga) Chapters(params map[string]interface{}) []Chapter {
	var chapters []Chapter
	db, cancel := util.Database()
	defer cancel()
	query := db.Model(&Chapter{}).Where("manga_id = ?", m.ID)
	if params != nil {
		query = query.Where(params)
	}
	query.Find(&chapters)
	return chapters
}

func (m *Manga) Create() *Manga {
	db, cancel := util.Database()
	defer cancel()
	db.Create(&m)
	return m
}

func (m *Manga) Delete() {
	db, cancel := util.Database()
	defer cancel()
	var chapter Chapter
	db.Where("manga_id = ?", m.ID).Delete(&chapter)
	var manga Manga
	db.Where("id = ?", m.ID).Delete(&manga)
}

func (m *Manga) ToReadable() string {
	return fmt.Sprintf("%s - %s - %s", m.Title, m.Authors, m.Updated)
}

func (m *Manga) Dirname() string {
	return util.CleanString(fmt.Sprintf("%d_%s", m.ID, m.Title))
}

func (m *Manga) ToRow() []string {
	id := fmt.Sprintf("%d", m.ID)
	return []string{id, m.Title, m.Authors, m.Uri, m.Updated.String(), m.Dlmode}
}
