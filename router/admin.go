package router

import (
	"shoego/handlers"
	"shoego/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRoutes(r *gin.RouterGroup, db *gorm.DB) *gin.RouterGroup {

	r.POST("/login", handlers.AdminLogin)

	adminProtected := r.Group("/")
	adminProtected.Use(middleware.AuthorizationMiddleware())

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
		//product variants
		adminProtected.POST("/products/:id/variants", handlers.AddVariants)
		adminProtected.GET("/products/:id/variants", handlers.GetVariants)
		adminProtected.PUT("/variants/:id", handlers.EditVariant)
		adminProtected.DELETE("/variants/:id", handlers.DeleteVariant)

	}
	{
		//cateory management
		adminProtected.POST("/categories", handlers.AddCategory)
		adminProtected.PUT("/categories/:id", handlers.EditCategory)
		adminProtected.DELETE("/categories/:id", handlers.DeleteCategory)
		adminProtected.GET("/categories", handlers.GetCategories)
	}
	{
		//offer
		adminProtected.POST("/offers/product", handlers.CreateProductOffer)
		adminProtected.GET("/offers/product", handlers.GetAllProductOffers)
		adminProtected.PUT("/offers/:id", handlers.UpdateProductOffer)
		adminProtected.DELETE("/offers/product/:id", handlers.DeleteProductOffer)

		adminProtected.POST("/category-offers", handlers.CreateCategoryOffer)
		adminProtected.GET("/category-offers", handlers.GetAllCategoryOffers)
		adminProtected.PUT("/category-offers/:id", handlers.UpdateCategoryOffer)
		adminProtected.DELETE("/category-offers/:id", handlers.DeleteCategoryOffer)
	}
	{
		//order
		adminProtected.GET("/orders", handlers.GetAdminOrders)
		adminProtected.GET("/orders/:id", handlers.GetAdminOrderDetail)
		adminProtected.PATCH("/orders/:id/status", handlers.UpdateAdminOrderStatus)
		//inventory
		adminProtected.GET("/inventory", handlers.GetAdminInventory)
		adminProtected.PATCH("/inventory/:id/stock", handlers.UpdateAdminProductStock)
	}
	{
		//dashboard
		adminProtected.GET("/dashboard/best-selling-products", handlers.GetAdminBestSellingProducts)
		adminProtected.GET("/dashboard/best-selling-categories", handlers.GetAdminBestSellingCategories)
		adminProtected.GET("/dashboard/sales-chart", handlers.GetAdminSalesChart)
	}
	{
		//coupon managemnt
		adminProtected.POST("/coupons", handlers.CreateCoupon)
		adminProtected.GET("/coupons", handlers.GetAllCoupons)
		adminProtected.PUT("/coupons/:id", handlers.UpdateCoupon)
		adminProtected.DELETE("/coupons/:id", handlers.DeleteCoupon)
	}
	{
		//sales report
		adminProtected.GET("/sales-report", handlers.GetSalesReport)
		adminProtected.GET("/sales-report/pdf", handlers.DownloadSalesReportPDF)

	}
	{
		//admin logout
		adminProtected.POST("/logout", handlers.AdminLogout)

	}
	return r
}
