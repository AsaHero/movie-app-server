package api

import (
	"net/http"

	"github.com/AsaHero/movie-app-server/delivery/api/docs"
	"github.com/AsaHero/movie-app-server/delivery/api/handlers"
	"github.com/AsaHero/movie-app-server/delivery/api/handlers/auth"
	"github.com/AsaHero/movie-app-server/delivery/api/handlers/movies"
	"github.com/AsaHero/movie-app-server/delivery/api/middlewares"
	"github.com/AsaHero/movie-app-server/pkg/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//  @title      		MovieAppServer
//  @version    		0.0.1
//  @description  		Documentation for "MovieAppServer" API
//  @termsOfService  	http://swagger.io/terms/

// @securityDefinitions.basic 	BasicAuth
// @securityDefinitions.apikey 	ApiKeyAuth
// @in              			header
// @name           				Authorization
// @description     			Basic Auth "Authorization: Basic <base64 encoded username:password>"

func NewRouter(cfg *config.Config, opt *handlers.HandlerOptions) *gin.Engine {
	r := gin.Default()

	// CORS configuration
	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // нужно изменить в продакшене
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.ClientIP()

		c.Next()
	})

	// Set base path /api/v1
	router := r.Group(middlewares.APIPrefix)

	auth.New(router.Group("/auth"), opt)
	movies.New(router.Group("/movies"), opt)

	// Swagger Route
	docs.SwaggerInfo.BasePath = middlewares.APIPrefix
	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
