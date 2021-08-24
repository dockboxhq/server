package services

import (
	"context"
	"io/ioutil"
	"log"
	"path/filepath"

	cli "github.com/dockboxhq/cli/cmd"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/karrick/godirwalk"
)

var DockerCli *client.Client

func analyze(mountPath string) (imagename string) {
	stats := make(map[string]int)
	godirwalk.Walk(mountPath,
		&godirwalk.Options{
			Callback: func(osPathname string, d *godirwalk.Dirent) error {
				log.Println(osPathname, d.Name())
				if !d.IsDir() {
					file_extension := filepath.Ext(d.Name())
					language_name := cli.ExtensionToLanguage[file_extension]
					if len(language_name) > 0 {
						stats[language_name] += 1
					}
				}
				return nil
			},
			ErrorCallback: func(path string, err error) godirwalk.ErrorAction {
				log.Printf("Error accessing file: %s", path)
				return godirwalk.SkipNode
			},
			Unsorted: true,
		},
	)

	sorted := cli.SortMap(stats)
	for i := len(sorted) - 1; i >= 0; i-- {
		if val, ok := cli.LanguageToImageMapper[sorted[i].Key]; ok {
			return val.Image
		}
	}

	return cli.LanguageToImageMapper["unknown"].Image

}

func pullImage(imageName string) error {
	ctx := context.Background()
	if _, _, err := DockerCli.ImageInspectWithRaw(ctx, imageName); err == nil {
		return nil
	}
	reader, err := DockerCli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return nil
	}
	defer reader.Close()
	if _, err := ioutil.ReadAll(reader); err != nil {
		return err
	}

	return nil
}

func CreateContainerForDockbox(mountPath string) (string, error) {
	cli := DockerCli
	ctx := context.Background()

	imageName := analyze(mountPath)
	err := pullImage(imageName)
	if err != nil {
		return "", err
	}

	createResponse, errCreate := cli.ContainerCreate(ctx, &container.Config{
		Image:        imageName,
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
