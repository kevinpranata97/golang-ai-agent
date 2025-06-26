package handlers

import (
	"database/sql"
)

// Handler contains the database connection and other dependencies
type Handler struct {
	DB *sql.DB
}

// New creates a new handler instance
func New(db *sql.DB) *Handler {
	return &Handler{
		DB: db,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
