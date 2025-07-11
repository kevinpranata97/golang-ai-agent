package routes

import (
	"github.com/gin-gonic/gin"
	"generated-application/internal/handlers"
)

// Setup configures all routes
func Setup(r *gin.Engine, h *handlers.Handler) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")
	{
		api.POST("/tasks", h.AddTask)
		api.GET("/tasks", h.ListTasks)
		api.PUT("/tasks/:id", h.UpdateTask)
		api.DELETE("/tasks/:id", h.DeleteTask)
	}
}
