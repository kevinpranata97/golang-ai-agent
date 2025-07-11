package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"database/sql"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) AddTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Add task endpoint"})
}

func (h *Handler) ListTasks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List tasks endpoint"})
}

func (h *Handler) UpdateTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update task endpoint"})
}

func (h *Handler) DeleteTask(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Delete task endpoint"})
}


