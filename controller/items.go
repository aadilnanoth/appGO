package controller

import (
	"appGO/config"
	"appGO/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Add  single Item to Category
func AddItem(c *gin.Context) {
	var item model.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := config.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create item"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item added", "item": item})
}

// add multipleItems to category
func AddMultipleItems(c *gin.Context) {
	var items []model.Item
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}
	if err := config.DB.Create(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create items"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Items added", "items": items})
}



func UpdateItem(c *gin.Context) {
	id := c.Param("id")
	var item model.Item

	// Find existing item
	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Bind JSON input to itemUpdate struct
	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		ImageURL    string  `json:"image_url"`
		CategoryID  uint    `json:"category_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Update fields
	item.Name = input.Name
	item.Description = input.Description
	item.Price = input.Price
	item.ImageURL = input.ImageURL
	item.CategoryID = input.CategoryID

	if err := config.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item updated", "item": item})
}


func DeleteItem(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&model.Item{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}
