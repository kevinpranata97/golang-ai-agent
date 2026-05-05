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
		// Post routes
		api.GET("/posts", h.GetAllPosts)
		api.GET("/posts/:id", h.GetPost)
		api.POST("/posts", h.CreatePost)
		api.PUT("/posts/:id", h.UpdatePost)
		api.DELETE("/posts/:id", h.DeletePost)

	}
}
