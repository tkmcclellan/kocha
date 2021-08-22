package util

import (
	"database/sql"
	"log"
	"os"
	"os/exec"
	"path"
)

func GetDatabase() (*sql.DB, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dbpath := path.Join(homedir, ".kocha", "manga.db")
	return sql.Open("sqlite3", dbpath)
}

func MonitorStdin(ch chan string) {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		ch <- string(b)
	}
}

func Cleanup() {
	exec.Command("stty", "-F", "/dev/tty", "sane").Run()
}

func MonitorSigint(ch chan os.Signal) {
	<-ch
	Cleanup()
	os.Exit(1)
}
