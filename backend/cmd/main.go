package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/automation/backend/db"
	cors "github.com/automation/backend/internal/middlewares"
	"github.com/pressly/goose/v3"
)

func main() {
	dsn := "./app.db" // default SQLite file
	if envDSN := os.Getenv("SQLITE_DSN"); envDSN != "" {
		dsn = envDSN
	}

	db.InitDB(dsn)
	defer db.GetDB().Close()

	// Run Goose migrations
	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db.GetDB(), "./migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Simple example handler
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	// Apply CORS middleware
	handler := cors.EnableCORS(mux)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
