package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/IvanMishnev/go_final_project/internal/constants"
	_ "github.com/mattn/go-sqlite3"
)

type TaskStore struct {
	db *sql.DB
}

var TaskDB TaskStore

func (store *TaskStore) Connect() {
	dbFile := constants.DBfile

	_, err := os.Stat(dbFile)
	if err != nil {
		os.Create(dbFile)
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal("Failed to connect to database.\n", err)
	}

	log.Println("connected to DB")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler(
		id INTEGER PRIMARY KEY,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(256) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(256) NOT NULL DEFAULT ""
		);`)
	if err != nil {
		log.Fatal(err)
	}

	store.db = db
}

func (store *TaskStore) Close() {
	err := store.db.Close()
	if err != nil {
		log.Fatal(err)
	}
}
