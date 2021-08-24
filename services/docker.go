package services

import (
	"context"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

var DockerCli *client.Client

func CreateContainerForDockbox(mountPath string) (string, error) {
	_, err := os.Stat(mountPath)
	if err != nil {
		return "", err
	}

	cli := DockerCli
	ctx := context.Background()

	createResponse, errCreate := cli.ContainerCreate(ctx, &container.Config{
		Image:        "ubuntu:latest",
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		OpenStdin:    true,
		WorkingDir:   "/app",
		Cmd:          []string{"/bin/bash", "-c", "while true; do bash; done"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: mountPath,
				Target: "/app",
			},
		},
	}, nil, nil, "")

	if errCreate != nil {
		return "", errCreate
	}

	containerID := createResponse.ID

	return containerID, nil
}

func StartContainer(containerID string) error {
	errStart := DockerCli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if errStart != nil {
		return errStart
	}
	return nil
}

func StopContainer(containerID string) error {
	errStop := DockerCli.ContainerKill(context.Background(), containerID, "SIGKILL")
	if errStop != nil {
		log.Fatalf("Error stopping container: %s", containerID)
	} else {
		log.Printf("Successfully stopping container: %s", containerID)
	}
	return errStop
}

func GetContainerStatus(containerID string) (*types.ContainerState, error) {
	ctx := context.Background()
	res, err := DockerCli.ContainerInspect(ctx, containerID)
	if err != nil {
		if client.IsErrNotFound(err) {
			return nil, nil
		}
		log.Fatalf("Unexpected error when retrieving status of container: %v", err)
		return nil, err
	}
	return res.State, nil
}

func init() {
	var err error
	DockerCli, err = client.NewClientWithOpts()

	if err != nil {
		log.Fatalln("Failed to connect to docker")
		return
	}
	log.Println("Connected to docker")
}
