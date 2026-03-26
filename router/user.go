package router

import (
	"shoego/handlers"
	"shoego/middleware"

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

	// product and category
	r.GET("/products", handlers.GetUserProducts)
	r.GET("/products/:id", handlers.GetUserProductDetails)
	r.GET("/categories", handlers.GetUserCategories)

	userProtected := r.Group("/")
	userProtected.Use(middleware.AuthMiddleware())
	{
		userProtected.POST("/logout", handlers.Logout)
	}

	{
		//check product availability
		userProtected.GET("/products/:id/validate", handlers.ValidateUserProductAvailability)
	}

	{
		//cart management

		userProtected.POST("/cart", handlers.AddToCart)
		userProtected.GET("/cart", handlers.GetCart)
		userProtected.PATCH("/cart/:product_id", handlers.UpdateCartQuantity)
		userProtected.DELETE("/cart/:product_id", handlers.RemoveCartItem)
		userProtected.GET("/cart/validate", handlers.ValidateCartCheckout)
	}

	{
		// user profile management
		userProtected.GET("/profile", handlers.GetProfile)
		userProtected.PUT("/profile/edit", handlers.EditProfile)
		userProtected.PUT("/profile/change-password", handlers.ChangePassword)
		userProtected.POST("/profile/request-email-change", handlers.RequestEmailChange)
		userProtected.POST("/profile/verify-email-change", handlers.VerifyEmailChange)

		//address
		userProtected.POST("/address", handlers.AddAddress)
		userProtected.PUT("/address/:id", handlers.EditAddress)
		userProtected.DELETE("/address/:id", handlers.DeleteAddress)
	}

	return r
}
