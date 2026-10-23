// Package main is the entry point of the TODO scheduler application.
package main

import (
	"log"
	"os"

	"github.com/March-Man/go_32_SuperFinal/pkg/db"
	"github.com/March-Man/go_32_SuperFinal/pkg/server"
)

// main initializes the database and starts the HTTP server.
func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}

	defer db.GetDB().Close()

	log.Printf("База данных готова: %s", dbFile)

	webDir := "./web"
	if err := server.Start(webDir); err != nil {
		log.Fatal(err)
	}
}
