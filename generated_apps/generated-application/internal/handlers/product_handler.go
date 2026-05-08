package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"generated-application/internal/models"
)

// CreateProduct creates a new Product
func (h *Handler) CreateProduct(c *gin.Context) {
	var product models.Product
	
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := models.CreateProduct(h.DB, &product); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Product created successfully",
		Data:    product,
	})
}

// GetProduct retrieves a Product by ID
func (h *Handler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	product, err := models.GetProductByID(h.DB, id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Product not found"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: product})
}

// GetAllProducts retrieves all Products
func (h *Handler) GetAllProducts(c *gin.Context) {
	products, err := models.GetAllProducts(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: products})
}

// UpdateProduct updates a Product
func (h *Handler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	product.Id = id
	if err := models.UpdateProduct(h.DB, &product); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product updated successfully",
		Data:    product,
	})
}

// DeleteProduct deletes a Product
func (h *Handler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	if err := models.DeleteProduct(h.DB, id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Product deleted successfully"})
}
