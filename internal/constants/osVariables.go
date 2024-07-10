package constants

import (
	"log"
	"os"
	"path/filepath"
)

var Password string
var TockenSecret string
var DBfile string
var Port string

func СonstInit() {
	Password = os.Getenv("TODO_PASSWORD")
	if Password == "" {
		Password = "103"
	}

	TockenSecret = os.Getenv("TODO_TOCKEN_SECRET")
	if TockenSecret == "" {
		TockenSecret = "secret"
	}

	DBfile = os.Getenv("TODO_DBFILE")
	if DBfile == "" {
		appPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		DBfile = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}

	Port = os.Getenv("TODO_PORT")
	if Port == "" {
		Port = "7540"
	}
}
