package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/automation/backend/db"
	cors "github.com/automation/backend/internal/middlewares"
	"github.com/automation/backend/internal/routes"
	"github.com/gin-gonic/gin"
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

	// if err := goose.Up(GetDB(), "./migrations"); err != nil {
	// 	log.Fatalf("failed to run migrations: %v", err)
	// }

	r := gin.Default()
	r.Use(cors.CORSMiddleware())
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Grouped routes
	api := r.Group("/api")
	routes.RegisterResumeRoutes(api)
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
