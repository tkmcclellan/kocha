package models

import (
	util "github.com/tkmcclellan/kocha/internal"

	"log"
	"time"
)

type Manga struct {
	Id       string
	Title    string
	Uri      string
	Authors  []string
	Dlmode   string
	Updated  time.Time
	Provider string
}

func (m *Manga) Chapters() []Chapter {
	db, err := util.GetDatabase()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query(`SELECT * FROM chapters WHERE manga_id = ?`, m.Id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	chapters := []Chapter{}

	for rows.Next() {
		var (
			id    int
			title string
			uri   string
			path  string
			read  int
		)
		if err := rows.Scan(&id, &title, &uri, &path, &read); err != nil {
			log.Fatal(err)
		} else {
			chapters = append(chapters, Chapter{id, title, uri, path, read == 1})
		}
	}

	return chapters
}
