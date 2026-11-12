package main

import (
	"log"
	"shoego/config"
	"shoego/database"
	"shoego/docs"
	routes "shoego/router"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Shoego API
//	@version		1.0
//	@description	API for ecommerce website
//	@securityDefinitions.apiKey	Bearer
//	@in							header
//	@name						Authorization
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//	@host			localhost:8080
//	@BasePath		/
//
// @schemes	http
func main() {

	//swgger info to test api easy
	docs.SwaggerInfo.Title = "Shoego API"
	docs.SwaggerInfo.Description = "Shoego API Documentation"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http"}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading the config file")
	}

	config.InitGoogleOAuth(
		cfg.GOOGLE_CLIENT_ID,
		cfg.GOOGLE_CLIENT_SECRET,
		cfg.GOOGLE_REDIRECT_URL,
	)
	config.InitRazorpay(cfg.RAZORPAY_KEY_ID, cfg.RAZORPAY_KEY_SECRET)

	db, err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Database connection Errors %v", err)
	}

	router := gin.Default()

	// Configure trusted proxy
	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatal(err)

		router.Static("/images", "./images")

		router.Static("/static", "./template")

		routes.UserRoutes(router.Group("/"), db)
		routes.AdminRoutes(router.Group("/admin"), db)

		router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		if err := router.Run("0.0.0.0:8080"); err != nil {
			log.Fatalf("Local host Errors %v", err)
		}
	}
}
