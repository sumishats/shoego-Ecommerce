package router

import (
	"shoego/handlers"
	"shoego/middileware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRoutes(r *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	r.POST("/admin", handlers.AdminLogin)

	//middileware
	r.Use(middileware.AuthorizationMiddleware())

	return r 
}
