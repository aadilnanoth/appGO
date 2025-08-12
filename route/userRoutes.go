package route

import (
	"appGO/controller"
	"appGO/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Public route: Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the API"})
	})

	// Public user routes
	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/signup", controller.SignupController)
		userRoutes.POST("/signup/otp", controller.VerifyOTP)
		userRoutes.POST("/login", controller.LoginController)
	}

	// Protected user routes (JWT required)
	protected := r.Group("/user")
	protected.Use(middleware.AuthMiddleware(), middleware.UserOnly())
	{
		protected.GET("/home", controller.Home)
	    protected.GET("/categories",controller.GetCategoriesWithItems)
		protected.GET("/profile", controller.UserProfile)
		protected.POST("/cart",controller.AddOrUpdateCartItem)
		protected.POST("/wishlist",controller.AddToWishlist)
		protected.GET("/cart", controller.GetCartItems)
		protected.GET("/whishlist", controller.GetWishlist)
		protected.DELETE("/cart/:id", controller.RemoveCartItem)
		protected.DELETE("/whishlist/:id", controller.RemoveFromWishlist)
	}

	
}
