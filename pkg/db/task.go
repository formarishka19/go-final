package db

import (
	"database/sql"
	"errors"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

func Tasks(limit int) ([]*Task, error) {

	tasks := make([]*Task, 0)
	query := `SELECT id, date, title, comment, repeat from scheduler ORDER BY date DESC LIMIT :limit`
	rows, err := DB.Query(query, sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}

func UpdateTask(task *Task) error {

	query := `UPDATE scheduler SET date = :date,  title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	res, err := DB.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`task not found`)
	}
	return nil
}

func DeleteTask(task *Task) error {
	query := `DELETE FROM scheduler WHERE id = :id`
	res, err := DB.Exec(query, sql.Named("id", task.ID))
	if err != nil {
		return fmt.Errorf(`error in task deleting from db`)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("task not found for delete from  db")
	}
	return nil
}

func GetTask(t *Task) error {
	query := `SELECT id, date, title, comment, repeat from scheduler WHERE id = :id`
	err := DB.QueryRow(query, sql.Named("id", t.ID)).Scan(
		&t.ID,
		&t.Date,
		&t.Title,
		&t.Comment,
		&t.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf(`no task with requested id %w`, err)
		}
		return fmt.Errorf(`incorrect query for get task from db %w`, err)
	}
	return nil
}
