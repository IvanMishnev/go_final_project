package handlers

import (
	"errors"
	"time"

	"github.com/IvanMishnev/go_final_project/internal/constants"
	"github.com/IvanMishnev/go_final_project/internal/nextdate"
	"github.com/IvanMishnev/go_final_project/models"
)

func ValidateTask(task models.Task) (models.Task, error) {
	if task.Title == "" {
		return models.Task{}, errors.New("missing title")
	}
	if task.Date == "" {
		task.Date = time.Now().Format(constants.DateFormat)
	}
	date, err := time.Parse(constants.DateFormat, task.Date)
	if err != nil {
		return models.Task{}, errors.New("wrong date format")
	}

	if date.Before(time.Now().Truncate(24 * time.Hour)) {
		if task.Repeat == "" {
			task.Date = time.Now().Format(constants.DateFormat)
		} else {
			nextDate, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				return models.Task{}, err
			}
			task.Date = nextDate
		}
	} else {
		if task.Repeat != "" {
			_, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				return models.Task{}, err
			}
		}
	}

	return task, nil
}
