package controller

import (
	"appGO/config"
	"appGO/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Extract user ID from Gin context, assuming your auth middleware sets "userID"
func getUserIDFromContext(c *gin.Context) uint {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		// User not authenticated, you can return 0 or abort here
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return 0
	}
	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		c.Abort()
		return 0
	}
	return userID
}

// Add or update cart item
func AddOrUpdateCartItem(c *gin.Context) {
	var input struct {
		ItemID   uint `json:"item_id" binding:"required"`
		Quantity int  `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := getUserIDFromContext(c)
	if userID == 0 {
		return // Already handled unauthorized in getUserIDFromContext
	}

	// Find or create cart for user
	var cart model.Cart
	err := config.DB.Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		// Cart not found, create new
		cart = model.Cart{UserID: userID}
		if err := config.DB.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create cart"})
			return
		}
	}

	// Find cart item by cartID and itemID
	var cartItem model.CartItem
	err = config.DB.Where("cart_id = ? AND item_id = ?", cart.ID, input.ItemID).First(&cartItem).Error
	if err != nil {
		// Not found, create new cart item
		cartItem = model.CartItem{
			CartID:   cart.ID,
			ItemID:   input.ItemID,
			Quantity: input.Quantity,
		}
		if err := config.DB.Create(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add item to cart"})
			return
		}
	} else {
		// Exists, update quantity
		cartItem.Quantity = input.Quantity
		if err := config.DB.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update cart item"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart updated", "cart_item": cartItem})
}

// Get all cart items for authenticated user (with item details)
func GetCartItems(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		return // unauthorized handled
	}

	var cart model.Cart
	err := config.DB.Preload("Items.Item").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		// No cart found, return empty list
		c.JSON(http.StatusOK, gin.H{"cart": []model.CartItem{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cart.Items})
}

// Remove a cart item by ID
func RemoveCartItem(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		return // unauthorized handled
	}

	id := c.Param("id")

	// Optional: Verify the cart item belongs to user's cart before deleting
	var cartItem model.CartItem
	if err := config.DB.Preload("Cart").First(&cartItem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	if cartItem.Cart.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this item"})
		return
	}

	if err := config.DB.Delete(&model.CartItem{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cart item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item removed"})
}
