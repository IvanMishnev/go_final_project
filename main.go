package main

import (
	"github.com/IvanMishnev/go_final_project/api"
	"github.com/IvanMishnev/go_final_project/database"
	"github.com/IvanMishnev/go_final_project/internal/constants"
)

func main() {
	constants.СonstInit()

	database.TaskDB.Connect()
	defer database.TaskDB.Close()

	api.StartServer()
}
