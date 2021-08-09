package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sriharivishnu/dockbox/server/controllers"
)

func SetUpV1(router *gin.Engine) {
	websocket := new(controllers.WebsocketController)

	v1 := router.Group("v1")
	v1.GET("/ws/:id", websocket.ContainerConnect)
}
