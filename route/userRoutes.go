package route

import (
	"appGO/controller"
	"appGO/middleware"
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

	protected := router.Group("/user")
protected.Use(middleware.AuthMiddleware())
 {
    protected.GET("/home", controller.Home)
    protected.GET("/profile", controller.UserProfile)
}

}
