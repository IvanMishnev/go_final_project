package database

import (
	"database/sql"
	"time"

	"github.com/IvanMishnev/go_final_project/internal/constants"
	"github.com/IvanMishnev/go_final_project/models"
)

func (store TaskStore) GetTasks(search string) ([]models.Task, error) {
	limit := constants.IssueLimit

	var rows *sql.Rows
	var err error

	if search != "" {
		date, err := time.Parse("02.01.2006", search)
		if err == nil {
			rows, err = store.db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?", date.Format(constants.DateFormat), limit)
			if err != nil {
				return nil, err
			}
		} else {
			search = "%" + search + "%"
			rows, err = store.db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?", search, search, limit)
			if err != nil {
				return nil, err
			}
		}
	} else {
		rows, err = store.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", limit)
		if err != nil {
			return nil, err
		}
	}

	tasks := []models.Task{}
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
