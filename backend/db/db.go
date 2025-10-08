package db

import (
	"database/sql"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var db *sql.DB

func InitDB(dsn string) {
	var err error
	db, err = sql.Open("sqlite3", dsn)
	if err != nil {
		panic(err)
	}

	// Enable WAL mode for better concurrency
	if _, err = db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		panic(err)
	}

	// Set max open connections if provided
	if maxOpenConns, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS")); err == nil {
		db.SetMaxOpenConns(maxOpenConns)
	}

	// Set max idle connections if provided
	if maxIdleConns, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS")); err == nil {
		db.SetMaxIdleConns(maxIdleConns)
	}

	// Enforce foreign key constraints
	if _, err = db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		panic(err)
	}
}

func GetDB() *sql.DB {
	return db
}
