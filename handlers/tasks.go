package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/IvanMishnev/go_final_project/database"
	"github.com/IvanMishnev/go_final_project/models"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	const limit = 15
	var rows *sql.Rows
	var err error
	search := r.URL.Query().Get("search")
	if search != "" {
		date, err := time.Parse("02.01.2006", search)
		if err == nil {
			rows, err = database.DB.Query("SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?", date.Format("20060102"), limit)
			if err != nil {
				JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			search = "%" + search + "%"
			rows, err = database.DB.Query("SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?", search, search, limit)
			if err != nil {
				JSONError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		rows, err = database.DB.Query("SELECT * FROM scheduler ORDER BY date LIMIT ?", limit)
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tasks := []models.Task{}
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, t)
	}

	mTasks := map[string][]models.Task{
		"tasks": tasks,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mTasks)
}
