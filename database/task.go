package database

import (
	"errors"
	"time"

	"github.com/IvanMishnev/go_final_project/internal/nextdate"
	"github.com/IvanMishnev/go_final_project/models"
)

func (store TaskStore) AddTask(task models.Task) (int64, error) {
	res, err := store.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (store TaskStore) GetTask(id int) (models.Task, error) {
	row := store.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)

	var t models.Task
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil || t.ID == "" {
		return t, err
	}

	return t, nil
}

func (store TaskStore) EditTask(task models.Task) error {
	res, err := store.db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("task has not been updated")
	}

	return nil
}

func (store TaskStore) DoneTask(id int) error {
	row := store.db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	var t models.Task
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return err
	}

	if t.Repeat == "" {
		_, err := store.db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			return err
		}
	} else {
		t.Date, err = nextdate.NextDate(time.Now(), t.Date, t.Repeat)
		if err != nil {
			return err
		}
		_, err = store.db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
			t.Date, t.Title, t.Comment, t.Repeat, t.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (store TaskStore) DeleteTask(id int) error {
	_, err := store.db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
