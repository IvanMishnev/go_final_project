package main

import (
	"os"

	"github.com/IvanMishnev/go_final_project/api"
	"github.com/IvanMishnev/go_final_project/database"
)

func main() {
	if os.Getenv("TODO_PASSWORD") == "" {
		os.Setenv("TODO_PASSWORD", "103")
	}

	database.ConnectDB()
	api.StartServer()
}
