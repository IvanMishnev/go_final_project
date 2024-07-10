package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

	task, err = ValidateTask(task)
	if err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := database.TaskDB.AddTask(task)
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

	task, err := database.TaskDB.GetTask(id)
	if err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
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
		JSONError(w, errors.New("wrong id").Error(), http.StatusBadRequest)
		return
	}

	task, err = ValidateTask(task)
	if err != nil {
		JSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.TaskDB.EditTask(task)
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
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

	err = database.TaskDB.DoneTask(id)
	if err != nil {
		JSONError(w, err.Error(), http.StatusInternalServerError)
		return
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

	err = database.TaskDB.DeleteTask(id)
	if err != nil {
		JSONError(w, "delete from DB failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})
}
