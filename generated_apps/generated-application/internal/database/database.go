package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Initialize initializes the database connection and runs migrations
func Initialize(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		databaseURL = "./app.db"
	}

	db, err := sql.Open("sqlite3", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Println("Database initialized successfully")
	return db, nil
}

// runMigrations runs database migrations
func runMigrations(db *sql.DB) error {
	migrations := []string{
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to execute migration: %v", err)
		}
	}

	return nil
}
