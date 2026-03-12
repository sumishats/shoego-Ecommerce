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

//	@title			shoego
//	@version		1.0
//	@description	shoego API Documentation
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

	db, err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Database connection Errors %v", err)
	}

	router := gin.Default()
	routes.UserRoutes(router.Group("/"), db)
	routes.AdminRoutes(router.Group("/admin"), db)

	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = router.Run("localhost:8080")
	if err != nil {

		log.Fatalf("Local host Errors %v", err)

	}

}
