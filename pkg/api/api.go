// Package api provides HTTP handlers for the scheduler REST API.
package api

import "net/http"

// Init registers all API route handlers.
func Init() {
	http.HandleFunc("/api/signin", signInHandler)
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
}
