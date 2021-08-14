package controllers

import (
	// "github.com/dockboxhq/cli/cmd"
	"context"

	"github.com/dockboxhq/server/services"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/gin-gonic/gin"
)

type DockboxController struct{}

type CreatePayload struct {
	Url string `json:"url"`
}

func (ws DockboxController) Create(c *gin.Context) {
	var payload CreatePayload
	err := c.BindJSON(&payload)
	if err != nil || payload.Url == "" {
		c.JSON(400, gin.H{"error": "Invalid payload"})
		return
	}
	cli := services.DockerCli

	ctx := context.Background()

	createResponse, errCreate := cli.ContainerCreate(ctx, &container.Config{
		Image:        "ubuntu:latest",
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		OpenStdin:    true,
		WorkingDir:   "/app",
		Entrypoint:   []string{"/bin/bash"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/Users/sriharivishnu/Desktop/dev/dockbox/cli/temp/problem-solving-javascript",
				Target: "/app",
			},
		},
	}, nil, nil, "")

	if errCreate != nil {
		c.JSON(500, gin.H{"error": errCreate})
		return
	}

	containerID := createResponse.ID

	errStart := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if errCreate != nil {
		c.JSON(500, gin.H{"error": errStart})
		return
	}

	// cmd.RunCreateCommand()
	c.JSON(200, gin.H{"id": containerID, "message": "Successfully created dockbox"})
}
