package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"generated-application/internal/models"
)

// CreatePost creates a new Post
func (h *Handler) CreatePost(c *gin.Context) {
	var post models.Post
	
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := models.CreatePost(h.DB, &post); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Post created successfully",
		Data:    post,
	})
}

// GetPost retrieves a Post by ID
func (h *Handler) GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	post, err := models.GetPostByID(h.DB, id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Post not found"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: post})
}

// GetAllPosts retrieves all Posts
func (h *Handler) GetAllPosts(c *gin.Context) {
	posts, err := models.GetAllPosts(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: posts})
}

// UpdatePost updates a Post
func (h *Handler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	post.ID = id
	if err := models.UpdatePost(h.DB, &post); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Post updated successfully",
		Data:    post,
	})
}

// DeletePost deletes a Post
func (h *Handler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	if err := models.DeletePost(h.DB, id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Post deleted successfully"})
}
