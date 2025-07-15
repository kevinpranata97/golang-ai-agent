package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import for SQLite3 driver
)

// Initialize initializes the database connection
func Initialize(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")

	// Create tables

	createUserTableSQL := "CREATE TABLE IF NOT EXISTS 