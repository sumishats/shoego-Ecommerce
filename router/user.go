package router

import (
	"shoego/handlers"
	"shoego/middileware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(r *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {
	r.POST("/signup", handlers.Signup)
	r.POST("/verify-otp", handlers.VerifyOTP)
	r.POST("/resend-otp", handlers.ResendOTP)
	r.POST("/login", handlers.UserLoginWithPassword)
	r.POST("/forgot-password", handlers.ForgotPassword)
	r.POST("/reset-password", handlers.ResetPassword)

	r.GET("/auth/google/login", handlers.GoogleLogin)
	r.GET("/auth/google/callback", handlers.GoogleCallback) //response from google


	r.Group("/users")
	r.Use(middileware.AuthMiddleware())

	return r
}
