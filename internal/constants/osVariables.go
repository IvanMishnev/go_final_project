package constants

import (
	"os"
)

var Password string
var TockenSecret string
var DBfile string
var Port string

func СonstInit() {
	Password = os.Getenv("TODO_PASSWORD")

	TockenSecret = os.Getenv("TODO_TOCKEN_SECRET")
	if TockenSecret == "" {
		TockenSecret = "secret"
	}

	DBfile = os.Getenv("TODO_DBFILE")
	if DBfile == "" {
		DBfile = "./scheduler.db"
	}

	Port = os.Getenv("TODO_PORT")
	if Port == "" {
		Port = "7540"
	}
}
