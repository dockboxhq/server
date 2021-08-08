package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sriharivishnu/dockbox/server/controllers"
	"github.com/sriharivishnu/dockbox/server/middlewares"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	health := new(controllers.HealthController)
	// websocket := new(controllers.WebsocketController)

	router.GET("/health", health.Status)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from dockbox API"})
	})
	router.Use(middlewares.AuthMiddleware())

	v1 := router.Group("v1")
	v1.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from dockbox API"})
	})
	v1.GET("/ws/:id", func(c *gin.Context) {
		id, _ := c.Params.Get("id")
		backendURL := fmt.Sprintf("ws://localhost:2375/containers/%s/attach/ws?logs=0&stream=1&stdin=1&stdout=1&stderr=1", id)
		dockerURL, err := url.Parse(backendURL)
		if err != nil {
			c.JSON(500, gin.H{"message": "Could not reach backend server"})
			return
		}

		proxy := &WebsocketProxy{
			Backend: func(req *http.Request) *url.URL {
				return dockerURL
			},
			Upgrader: &websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		}
		proxy.ServeHTTP(c.Writer, c.Request)
		// websocket.Start(c.Writer, c.Request)
		// c.JSON(200, gin.H{"message": "Hello from dockbox API"})

	})
	return router

}
