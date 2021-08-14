package services

import (
	"log"

	"github.com/docker/docker/client"
)

var DockerCli *client.Client

func init() {
	var err error
	DockerCli, err = client.NewClientWithOpts()

	if err != nil {
		log.Fatalln("Failed to connect to docker")
		return
	}
	log.Println("Connected to docker")
}
