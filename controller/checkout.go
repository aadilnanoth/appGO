package controller

import (
	"appGO/config"
	"appGO/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CheckoutRequest struct {
	ShippingAddress string `json:"shipping_address" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
}

func CheckoutHandler(c *gin.Context) {
	userID := c.GetUint("userID") // from JWT middleware

	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB

	// 1. Fetch cart items for user
	var cartItems []model.CartItem
	if err := db.Preload("Item.Category").
		Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("carts.user_id = ?", userID).
		Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch cart"})
		return
	}

	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart is empty"})
		return
	}

	// 2. Calculate total price
	var total float64
	for _, citem := range cartItems {
		total += float64(citem.Quantity) * citem.Item.Price
	}

	// 3. Create order with transaction
	order := model.Order{
		UserID:          userID,
		TotalPrice:      total,
		PaymentMethod:   "COD",
		Status:          "Pending",
		ShippingAddress: req.ShippingAddress,
		Phone:           req.Phone,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for _, citem := range cartItems {
			orderItem := model.OrderItem{
				OrderID:  order.ID,
				ItemID:   citem.ItemID,
				Quantity: citem.Quantity,
				Price:    citem.Item.Price,
			}

			if err := tx.Create(&orderItem).Error; err != nil {
				return err
			}

			// Deduct stock safely
			if err := tx.Model(&model.Item{}).
				Where("id = ? AND stock >= ?", citem.ItemID, citem.Quantity).
				Update("stock", gorm.Expr("stock - ?", citem.Quantity)).Error; err != nil {
				return err
			}
		}

		// Clear cart
		if err := tx.Where("cart_id IN (SELECT id FROM carts WHERE user_id = ?)", userID).
			Delete(&model.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "checkout failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order placed successfully", "order_id": order.ID})
}
