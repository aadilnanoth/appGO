package route

import (
	"appGO/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the API"})
	})

	// Grouped user routes
	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/signup", controller.SignupController)
		userRoutes.POST("/login", controller.LoginController)
	}
}
