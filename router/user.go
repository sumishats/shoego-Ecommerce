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

	r.GET("/order-success/:order_id", handlers.GetOrderSuccess)

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
		//referral 
		userProtected.GET("/referral", handlers.GetReferralDetails)
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
		//wishlist management

		userProtected.POST("/wishlist", handlers.AddToWishlist)
		userProtected.GET("/wishlist", handlers.GetWishlist)
		userProtected.DELETE("/wishlist/:product_id", handlers.RemoveFromWishlist)

	}
	{
		//coupon management
		userProtected.POST("/coupons/apply", handlers.ApplyCoupon)
		userProtected.DELETE("/coupons/remove", handlers.RemoveCoupon)

	}
	{
		//checkout management
		userProtected.GET("/checkout", handlers.GetCheckoutPage)
		userProtected.POST("/checkout/place-order", handlers.PlaceCODOrder)
		userProtected.POST("/checkout/wallet",handlers.PlaceWalletOrder)

		//razorpay
		userProtected.POST("/checkout/razorpay", handlers.CreateRazorpayOrder)
		userProtected.POST("/checkout/razorpay/verify", handlers.VerifyPayment)

	}
	{
		//order management
		userProtected.GET("/orders", handlers.GetUserOrders)
		userProtected.GET("/orders/:order_id", handlers.GetOrderDetail)
		userProtected.PATCH("/orders/:order_id/cancel", handlers.CancelOrder)
		userProtected.PATCH("/orders/:order_id/items/:item_id/cancel", handlers.CancelOrderItem)
		userProtected.PATCH("/orders/:order_id/return", handlers.ReturnOrder)
		userProtected.PATCH("/orders/:order_id/items/:item_id/return", handlers.ReturnOrderItem)
		userProtected.GET("/orders/:order_id/invoice", handlers.DownloadInvoice)
	}
	{
		//wallet
		userProtected.GET("/wallet", handlers.GetWallet)
		userProtected.GET("/wallet/history", handlers.GetWalletHistory)
	}
	{
		// user profile management
		userProtected.GET("/profile", handlers.GetProfile)
		userProtected.PUT("/profile/edit", handlers.EditProfile)
		userProtected.PUT("/profile/change-password", handlers.ChangePassword)
		userProtected.POST("/profile/request-email-change", handlers.RequestEmailChange)
		userProtected.POST("/profile/verify-email-change", handlers.VerifyEmailChange)

		//address

		userProtected.GET("/address", handlers.GetAddresses)
		userProtected.POST("/address", handlers.AddAddress)
		userProtected.PUT("/address/:id", handlers.EditAddress)
		userProtected.DELETE("/address/:id", handlers.DeleteAddress)
	}

	return r
}
