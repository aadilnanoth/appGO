package middleware

import (
	"appGO/config"
	"appGO/model"
	"appGO/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		var user struct {
			ID         uint
			IsVerified bool
			Role       string
		}

		err = config.DB.Model(&model.User{}).
			Select("id", "is_verified", "role").
			Where("email = ?", claims.Email).
			Scan(&user).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			c.Abort()
			return
		}

		if !user.IsVerified {
			c.JSON(http.StatusForbidden, gin.H{"error": "Please verify your account via OTP"})
			c.Abort()
			return
		}

		// Set user data in context for handlers to access
		c.Set("userID", user.ID)
		c.Set("userEmail", claims.Email)
		c.Set("Role", user.Role)

		fmt.Printf("Token Received: %s\n", tokenStr)
		fmt.Printf("Claims: %+v\n", claims)

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("Role")
		role, ok := roleVal.(string)
		if !exists || !ok || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access only"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func UserOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("Role")
		role, ok := roleVal.(string)
		if !exists || !ok || role != "user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "User access only"})
			c.Abort()
			return
		}
		c.Next()
	}
}
