package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/IvanMishnev/go_final_project/database"
	"github.com/IvanMishnev/go_final_project/models"
)

func AddTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		JSONError(w, "JSON deserealization error: "+err.Error(), http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		JSONError(w, "missing title", http.StatusBadRequest)
		return
	}
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}
	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		JSONError(w, "wrong date format: "+err.Error(), http.StatusBadRequest)
		return
	}
	if date.Before(time.Now().Truncate(24 * time.Hour)) {
		if task.Repeat == "" {
			task.Date = time.Now().Format("20060102")
		} else {
			nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				JSONError(w, err.Error(), http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}
	} else {
		if task.Repeat != "" {
			_, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				JSONError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	res, err := database.DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id": id,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		JSONError(w, errors.New("wrong id").Error(), http.StatusBadRequest)
		return
	}

	row := database.DB.QueryRow("SELECT * FROM scheduler WHERE id = ?", id)
	var t models.Task
	err = row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil || t.ID == "" {
		JSONError(w, errors.New("task not found").Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t)
}

func EditTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		JSONError(w, "JSON deserealization error: "+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = strconv.Atoi(task.ID)
	if err != nil {
		JSONError(w, "wrong id", http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		JSONError(w, "missing title", http.StatusBadRequest)
		return
	}
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}
	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		JSONError(w, "wrong date format: "+err.Error(), http.StatusBadRequest)
		return
	}
	if date.Before(time.Now().Truncate(24 * time.Hour)) {
		if task.Repeat == "" {
			task.Date = time.Now().Format("20060102")
		} else {
			nextDate, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				JSONError(w, err.Error(), http.StatusBadRequest)
				return
			}
			task.Date = nextDate
		}
	} else {
		if task.Repeat != "" {
			_, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				JSONError(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	res, err := database.DB.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		JSONError(w, "task has not been updated", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}

func DoneTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		JSONError(w, errors.New("wrong id").Error(), http.StatusBadRequest)
		return
	}

	row := database.DB.QueryRow("SELECT * FROM scheduler WHERE id = ?", id)
	var t models.Task
	err = row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		JSONError(w, "task not found in DB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if t.Repeat == "" {
		_, err := database.DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			JSONError(w, "delete from DB failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		t.Date, err = NextDate(time.Now(), t.Date, t.Repeat)
		if err != nil {
			JSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = database.DB.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
			t.Date, t.Title, t.Comment, t.Repeat, t.ID)
		if err != nil {
			JSONError(w, "update task failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		JSONError(w, errors.New("wrong id").Error(), http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		JSONError(w, "delete from DB failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}
