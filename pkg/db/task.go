package db

import (
	"database/sql"
	"fmt"
	"time"
)

// Task represents a task in the scheduler.
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask inserts a new task into the scheduler table and returns its ID.
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления задачи: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID: %w", err)
	}

	return id, nil
}

// Tasks returns a list of upcoming tasks sorted by date.
func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса задач: %w", err)
	}
	defer rows.Close()

	tasks := make([]*Task, 0)

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("ошибка сканирования: %w", err)
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации: %w", err)
	}

	return tasks, nil
}

// TasksSearch searches for tasks by a substring in the title or comment, or by a specific date.
func TasksSearch(search string, limit int) ([]*Task, error) {
	if t, err := time.Parse("02.01.2006", search); err == nil {
		dateStr := t.Format("20060102")
		return tasksByDate(dateStr, limit)
	}

	like := "%" + search + "%"
	query := `SELECT id, date, title, comment, repeat FROM scheduler
	          WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`

	rows, err := db.Query(query, like, like, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTasks(rows)
}

// tasksByDate returns tasks scheduled for a specific date.
func tasksByDate(date string, limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler
	          WHERE date = ? ORDER BY date LIMIT ?`

	rows, err := db.Query(query, date, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanTasks(rows)
}

// scanTasks scans sql.Rows into a slice of tasks.
func scanTasks(rows *sql.Rows) ([]*Task, error) {
	tasks := make([]*Task, 0)
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, rows.Err()
}

// GetTask returns a task by its ID.
func GetTask(id string) (*Task, error) {
	var task Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, fmt.Errorf("ошибка получения задачи: %w", err)
	}

	return &task, nil
}

// UpdateTask updates an existing task in the scheduler table.
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества строк: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// DeleteTask deletes a task by its ID.
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества строк: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// UpdateDate updates only the date of a task.
func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	res, err := db.Exec(query, next, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления даты: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества строк: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
