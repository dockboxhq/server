package server

import (
	"github.com/dockboxhq/server/controllers"
	"github.com/gin-gonic/gin"
)

func SetUpV1(router *gin.Engine) {
	websocket := new(controllers.WebsocketController)
	dockbox := new(controllers.DockboxController)

	v1 := router.Group("v1")
	v1.GET("/ws/:id", websocket.ContainerConnect)

	dockboxGroup := v1.Group("dockbox")
	dockboxGroup.POST("", dockbox.Create)
}
