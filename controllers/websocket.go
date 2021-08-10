package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sriharivishnu/dockbox/server/socket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketController struct{}

func (ws WebsocketController) ContainerConnect(c *gin.Context) {
	id, _ := c.Params.Get("id")
	DOCKER_HOST := os.Getenv("DOCKER_SERVER_HOST")
	log.Println(DOCKER_HOST)
	// time.Sleep(5 * time.Second)

	backendURL := fmt.Sprintf("ws://%s/containers/%s/attach/ws?logs=0&stream=1&stdin=1&stdout=1&stderr=1", DOCKER_HOST, id)
	log.Println(backendURL)
	dockerURL, err := url.Parse(backendURL)
	if err != nil {
		c.JSON(500, gin.H{"message": "Could not reach backend server"})
		return
	}

	proxy := &socket.WebsocketProxy{
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

}
