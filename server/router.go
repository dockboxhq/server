package server

import (
	"time"

	"github.com/dockboxhq/server/controllers"
	"github.com/dockboxhq/server/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000/*", "http://dockbox.ca/*", "https://dockbox.ca/*", "https://*.dockbox.ca/*", "http://*.dockbox.ca/*"},
		AllowMethods:     []string{"PUT", "POST", "DELETE", "GET"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowWildcard:    true,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	health := new(controllers.HealthController)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from dockbox API"})
	})
	router.GET("/health", health.Status)
	router.Use(middlewares.AuthMiddleware())

	SetUpV1(router)
	return router

}
