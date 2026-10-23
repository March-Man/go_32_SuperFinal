package api

import (
	"log"
	"net/http"
	"time"

	"github.com/March-Man/go_32_SuperFinal/pkg/db"
)

// doneTaskHandler handles the POST /api/task/done?id=<id> request.
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			log.Printf("ошибка удаления задачи: %v", err)
			writeJSON(w, map[string]string{"error": "Внутренняя ошибка"}, http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{}, http.StatusOK)
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		log.Printf("ошибка обновления даты: %v", err)
		writeJSON(w, map[string]string{"error": "Внутренняя ошибка"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{}, http.StatusOK)
}
