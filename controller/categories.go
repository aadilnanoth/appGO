package controller

import (
	"appGO/config"
	"appGO/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Add Category
func AddCategory(c *gin.Context) {
	var cat model.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := config.DB.Create(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category added", "category": cat})
}


func GetCategoriesWithItems(c *gin.Context) {
	var categories []model.Category

	// Preload Items for each Category
	if err := config.DB.Preload("Items").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories and items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category model.Category

	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var input struct {
		Name     string `json:"name"`
		ImageURL string `json:"image_url"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	category.Name = input.Name
	category.ImageURL = input.ImageURL

	if err := config.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated", "category": category})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	// Optional: check if category has dependent items and prevent deletion

	if err := config.DB.Delete(&model.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
