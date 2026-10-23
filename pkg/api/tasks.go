package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/March-Man/go_32_SuperFinal/pkg/db"
)

// TasksResp represents the JSON response containing a list of tasks.
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler handles the GET /api/tasks request with optional search.
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	var tasks []*db.Task
	var err error

	if search != "" {
		tasks, err = db.TasksSearch(search, 50)
	} else {
		tasks, err = db.Tasks(50)
	}

	if err != nil {
		log.Printf("ошибка получения задач: %v", err)
		writeJSON(w, map[string]string{"error": "Внутренняя ошибка"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks}, http.StatusOK)
}

// taskHandler dispatches /api/task requests based on the HTTP method.
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// getTaskHandler handles the GET /api/task?id=<id> request.
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, task, http.StatusOK)
}

// updateTaskHandler handles the PUT /api/task request to update an existing task.
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "ошибка десериализации: " + err.Error()}, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"}, http.StatusBadRequest)
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

	if err := db.UpdateTask(&task); err != nil {
		log.Printf("ошибка обновления задачи: %v", err)
		writeJSON(w, map[string]string{"error": "Внутренняя ошибка"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{}, http.StatusOK)
}

// deleteTaskHandler handles the DELETE /api/task?id=<id> request.
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"}, http.StatusBadRequest)
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}
	writeJSON(w, map[string]string{}, http.StatusOK)
}
