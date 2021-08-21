package util

import (
	"database/sql"
	"log"
	"path"
	"os"
)

func GetDatabase() (*sql.DB, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dbpath := path.Join(homedir, ".kocha", "manga.db")
	return sql.Open("sqlite3", dbpath)
}

