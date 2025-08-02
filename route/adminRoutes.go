package route

import (
	"appGO/controller"
	"appGO/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(r *gin.Engine) {
	//  Public admin routes 
	public := r.Group("/admin")
	{
		public.POST("/signup", controller.AdminSignupController)
		
	}

	//  Protected admin routes (Auth + AdminOnly middleware)
	protected := r.Group("/admin")
	protected.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		protected.GET("/dashboard", controller.AdminDashboard)
		// Item Routes
		protected.POST("/item", controller.AddItem)
		protected.POST("/multipleItems", controller.AddMultipleItems)
		protected.PUT("/item/:id",controller.UpdateItem)
		protected.DELETE("/item/:id",controller.DeleteItem)
		// Category Routes
		protected.POST("/addcategory", controller.AddCategory)
		protected.GET("/categories", controller.GetCategoriesWithItems)
		protected.PUT("/category/:id", controller.UpdateCategory)
		protected.DELETE("/category/:id", controller.DeleteCategory)

	}
}
