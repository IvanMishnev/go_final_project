package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/IvanMishnev/go_final_project/database"
	"github.com/IvanMishnev/go_final_project/models"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := database.TaskDB.GetTasks(search)
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
	}

	mTasks := map[string][]models.Task{
		"tasks": tasks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mTasks)
}
