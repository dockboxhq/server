package server

import (
	"github.com/dockboxhq/server/utils"
	"github.com/gin-gonic/gin"
)

func Init() {
	config := utils.Config
	if config.ENVIRONMENT == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := NewRouter()
	r.Run(":" + config.PORT)
}
