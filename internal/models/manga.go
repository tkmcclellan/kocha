package models

import "time"

type Manga struct {
	Title   string
	Authors []string
	Updated time.Time
}
