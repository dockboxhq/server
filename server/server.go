package server

import (
	"context"
	"log"

	"github.com/dockboxhq/server/services"
	"github.com/dockboxhq/server/utils"
	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
)

func populateContainers() {
	log.Println("Populating containers")
	containers, err := services.DockerCli.ContainerList(context.Background(), types.ContainerListOptions{})

	if err != nil {
		log.Fatalf("Error retrieving containers: %v\n", err)
	}
	for _, container := range containers {
		services.ContainerManager.StartedContainer(container.ID)
	}
}

func Init() {
	config := utils.Config
	if config.ENVIRONMENT == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	populateContainers()

	r := NewRouter()
	r.Run(":" + config.PORT)
}
