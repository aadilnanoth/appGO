package controller

import (
	"appGO/config"
	"appGO/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Add item to wishlist
func AddToWishlist(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		return
	}

	var input struct {
		ItemID uint `json:"item_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find or create wishlist
	var wishlist model.Wishlist
	if err := config.DB.Where("user_id = ?", userID).First(&wishlist).Error; err != nil {
		wishlist = model.Wishlist{UserID: userID}
		if err := config.DB.Create(&wishlist).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create wishlist"})
			return
		}
	}

	// Check if already in wishlist
	var existing model.WishlistItem
	if err := config.DB.Where("wishlist_id = ? AND item_id = ?", wishlist.ID, input.ItemID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Item already in wishlist"})
		return
	}

	// Add to wishlist
	item := model.WishlistItem{WishlistID: wishlist.ID, ItemID: input.ItemID}
	if err := config.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add item to wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added to wishlist", "wishlist_item": item})
}

// Get all wishlist items
func GetWishlist(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		return
	}

	var wishlist model.Wishlist
	if err := config.DB.Preload("Items.Item").Where("user_id = ?", userID).First(&wishlist).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"wishlist": []model.WishlistItem{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wishlist": wishlist.Items})
}

// Remove item from wishlist
func RemoveFromWishlist(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		return
	}

	id := c.Param("id")
	var item model.WishlistItem
	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wishlist item not found"})
		return
	}

	// Optional: Check if it belongs to the user
	var wishlist model.Wishlist
	if err := config.DB.First(&wishlist, item.WishlistID).Error; err == nil && wishlist.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if err := config.DB.Delete(&model.WishlistItem{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not remove item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from wishlist"})
}
