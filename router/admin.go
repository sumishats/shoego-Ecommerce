package router

import (
	"shoego/handlers"
	"shoego/middileware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRoutes(r *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	// public admin route
	r.POST("/login", handlers.AdminLogin)

	// protected admin routes
	adminProtected := r.Group("/")
	adminProtected.Use(middileware.AuthorizationMiddleware())

	{
		//user management 
		adminProtected.GET("/users", handlers.GetUsers)
		adminProtected.PATCH("/block-user/:id", handlers.BlockUser)
		adminProtected.PATCH("/unblock-user/:id", handlers.UnblockUser)
	}
	{
		// product management

		adminProtected.POST("/products", handlers.AddProduct)
		adminProtected.PUT("/products/:id", handlers.EditProduct)
		adminProtected.DELETE("/products/:id", handlers.DeleteProduct)
		adminProtected.GET("/products", handlers.GetProducts)
	}
	{
		//cateory management
		adminProtected.POST("/categories", handlers.AddCategory)
		adminProtected.PUT("/categories/:id", handlers.EditCategory)
		adminProtected.DELETE("/categories/:id", handlers.DeleteCategory)
		adminProtected.GET("/categories", handlers.GetCategories)
	}

	return r
}
