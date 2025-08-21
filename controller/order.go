package controller

import (
	"appGO/config"
	"appGO/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOrdersHandler(c *gin.Context) {
	userID := c.GetUint("userID") // from JWT middleware

	var orders []model.Order
	db := config.DB

	if err := db.Preload("OrderItems.Item.Category").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert DB models to response structs
	var ordersResp []model.OrderResponse
	for _, o := range orders {
		userSummary := model.UserSummary{
			ID:    o.User.ID,
			Name:  o.User.Name,
			Email: o.User.Email,
		}

		var itemsResp []model.OrderItemResponse
		for _, oi := range o.OrderItems {
			itemSimple := model.ItemSimple{
				ID:          oi.Item.ID,
				Name:        oi.Item.Name,
				Description: oi.Item.Description,
				Price:       oi.Item.Price,
				ImageURL:    oi.Item.ImageURL,
				Category: model.CategorySimple{
					ID:       oi.Item.Category.ID,
					Name:     oi.Item.Category.Name,
					ImageURL: oi.Item.Category.ImageURL,
				},
			}
			itemsResp = append(itemsResp, model.OrderItemResponse{
				ID:       oi.ID,
				OrderID:  oi.OrderID,
				ItemID:   oi.ItemID,
				Item:     itemSimple,
				Quantity: oi.Quantity,
				Price:    oi.Price,
			})
		}

		ordersResp = append(ordersResp, model.OrderResponse{
			ID:              o.ID,
			UserID:          o.UserID,
			User:            userSummary,
			TotalPrice:      o.TotalPrice,
			PaymentMethod:   o.PaymentMethod,
			Status:          o.Status,
			ShippingAddress: o.ShippingAddress,
			Phone:           o.Phone,
			CreatedAt:       o.CreatedAt,
			OrderItems:      itemsResp,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "orders fetched successfully",
		"orders":  ordersResp,
	})
}
