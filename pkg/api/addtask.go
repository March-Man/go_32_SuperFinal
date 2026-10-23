package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/March-Man/go_32_SuperFinal/pkg/db"
)

// addTaskHandler handles the POST /api/task request to create a new task.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "ошибка десериализации: " + err.Error()}, http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		log.Printf("ошибка добавления задачи: %v", err)
		writeJSON(w, map[string]string{"error": "Внутренняя ошибка"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)}, http.StatusOK)
}

// checkDate validates and corrects the task's date based on the repeat rule.
func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("неверный формат даты: %s", task.Date)
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("неверное правило повторения: %w", err)
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(dateFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}

// writeJSON writes the given data as a JSON response with the appropriate headers and status code.
func writeJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
