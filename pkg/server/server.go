// Package server provides the HTTP server setup for the scheduler application.
package server

import (
	"log"
	"net/http"
	"os"

	"github.com/March-Man/go_32_SuperFinal/pkg/api"
)

// Start launches the HTTP server, registering API handlers and serving static files.
func Start(webDir string) error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	addr := ":" + port

	api.Init()

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Сервер запущен на http://localhost%s", addr)
	return http.ListenAndServe(addr, nil)
}
