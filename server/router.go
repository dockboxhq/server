package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sriharivishnu/dockbox/server/controllers"
	"github.com/sriharivishnu/dockbox/server/middlewares"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from dockbox API"})
	})
	router.Use(middlewares.AuthMiddleware())

	v1 := router.Group("v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello from dockbox API"})
		})
	}
	return router

}
