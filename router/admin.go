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

	return r
}
