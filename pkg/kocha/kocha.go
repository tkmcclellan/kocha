package kocha

import (
	"log"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
	util "github.com/tkmcclellan/kocha/internal"
)

func Init() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
		return err
	}
	dirpath := path.Join(homedir, ".kocha")

	_, err = os.Stat(dirpath)
	if os.IsNotExist(err) {
		os.Mkdir(dirpath, 0775)
	}

	db, err := util.GetDatabase()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	createTables := `
	create table if not exists manga (
		id int primary key,
		title varchar(255) not null,
		uri varchar(255) not null,
		authors varchar(255),
		download_mode varchar(255) not null,
		provider varchar(255) not null,
		current_chapter int default 0
	);

	create table if not exists chapters (
		id int primary key,
		manga_id int,
		title varchar(255) not null,
		uri varchar(255) not null,
		path varchar(255) not null,
		read int not null check (read in (0, 1)),
		foreign key(manga_id) references manga(id)
	);
	`

	_, err = db.Exec(createTables)

	if err != nil {
		log.Printf("%q: %s\n", err, createTables)
	}

	return nil
}
