package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDB() {
	var dbFile string
	envDbPath := os.Getenv("TODO_DBFILE")
	if envDbPath != "" {
		dbFile = envDbPath
	} else {
		appPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		dbFile = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}

	_, err := os.Stat(dbFile)
	if err != nil {
		os.Create(dbFile)
	}

	DB, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal("Failed to connect to database.\n", err)
		os.Exit(1)
	}

	log.Println("connected to DB")

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS scheduler(
		id INTEGER PRIMARY KEY,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(256) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(256) NOT NULL DEFAULT ""
		);`)
	if err != nil {
		log.Fatal(err)
	}

}
