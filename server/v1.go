package server

import (
	"github.com/dockboxhq/server/controllers"
	"github.com/gin-gonic/gin"
)

func SetUpV1(router *gin.Engine) {
	dockbox := new(controllers.DockboxController)

	v1 := router.Group("v1")

	dockboxGroup := v1.Group("dockbox")
	dockboxGroup.GET("/ws/:id", dockbox.Connect)
	dockboxGroup.POST("", dockbox.Create)
}
