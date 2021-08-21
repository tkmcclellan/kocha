package kocha

import (
	"database/sql"
	"log"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

func Init() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dirpath := path.Join(homedir, ".kocha")

	_, err = os.Stat(dirpath)
	if os.IsNotExist(err) {
		os.Mkdir(dirpath, 0775)
	}
	dbpath := path.Join(dirpath, "manga.db")

	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTables := ``

	_, err = db.Exec(createTables)

	if err != nil {
		log.Printf("%q: %s\n", err, createTables)
	}

	return err
}
