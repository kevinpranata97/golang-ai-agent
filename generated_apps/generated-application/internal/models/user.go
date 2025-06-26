package models

import (
	"time"
	"database/sql"
)

// User represents the User entity
type User struct {
	Id int `json:"id" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Created_at time.Time `json:"created_at" validate:"required"`
}

// CreateUser creates a new User in the database
func CreateUser(db *sql.DB, user *User) error {
	query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	
	result, err := db.Exec(query, 