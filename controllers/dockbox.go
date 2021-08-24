package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/dockboxhq/server/models"
	"github.com/dockboxhq/server/services"
	"github.com/dockboxhq/server/socket"
	"github.com/dockboxhq/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-getter"
	"github.com/karrick/godirwalk"
	"github.com/lithammer/shortuuid/v3"
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

	var dockboxID = uuid.New().String()
	var customID = shortuuid.New()

	errGetData := getRepositoryData(payload.Url, filepath.Join(utils.Config.MOUNT_POINT, dockboxID))
	if errGetData != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Error fetching data from repository: %v. Please check your URL and try again.", errGetData)})
		return
	}

	containerID, err := services.CreateContainerForDockbox(filepath.Join(utils.Config.MOUNT_POINT, dockboxID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	dockbox, err := models.CreateDockbox(
		dockboxID,
		payload.Url,
		sql.NullString{String: containerID, Valid: true},
		sql.NullString{String: customID, Valid: true},
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		log.Println(err.Error())
		return
	}

	log.Printf("%v\n", dockbox)

	c.JSON(200, gin.H{"id": dockbox.ID, "message": "Successfully created dockbox"})
}

func getRepositoryData(url string, dest string) error {
	if strings.Contains(url, "github") || strings.Contains(url, "gitlab") && !strings.HasPrefix(url, "git::") {
		url = "git::" + url
	}
	client := &getter.Client{
		Ctx:  context.Background(),
		Dst:  dest,
		Src:  url,
		Mode: getter.ClientModeAny,
		Detectors: []getter.Detector{
			&getter.GitHubDetector{},
			&getter.GitDetector{},
			&getter.S3Detector{},
		},
		//provide the getter needed to download the files
		Getters: map[string]getter.Getter{
			"git":   &getter.GitGetter{},
			"http":  &getter.HttpGetter{},
			"https": &getter.HttpGetter{},
		},
	}
	if err := client.Get(); err != nil {
		fmt.Fprintf(os.Stderr, "Error while fetching code from %s: %v", client.Src, err)
		return err
	}
	return nil
}

func (d DockboxController) Connect(c *gin.Context) {
	id, _ := c.Params.Get("id")

	dockbox, err := models.GetDockboxById(id)
	if err != nil {
		log.Fatalf("Dockbox with id %s not found", id)
		c.JSON(404, gin.H{"message": "Dockbox not found", "error": err.Error()})
		return
	}

	status, statusErr := services.GetContainerStatus(dockbox.ContainerId.String)
	// The container ID is null; create a new container
	if !dockbox.ContainerId.Valid || status == nil {
		containerID, err := services.CreateContainerForDockbox(filepath.Join(utils.Config.MOUNT_POINT, dockbox.ID))
		if err != nil {
			c.JSON(404, gin.H{"message": "Could not start dockbox", "error": err.Error()})
			return
		}
		dockbox, err = models.UpdateDockbox(dockbox.ID, sql.NullString{String: containerID, Valid: true})
		if err != nil {
			c.JSON(404, gin.H{"message": "Could not start dockbox", "error": err.Error()})
			return
		}
		// Get new status
		status, statusErr = services.GetContainerStatus(dockbox.ContainerId.String)
	}

	if statusErr != nil {
		c.JSON(404, gin.H{"message": "Error retrieving dockbox status", "error": err.Error()})
		return
	}

	if !status.Running || !status.Restarting {
		errStart := services.StartContainer(dockbox.ContainerId.String)
		if errStart != nil {
			c.JSON(404, gin.H{"message": "Error starting connection", "error": errStart.Error()})
			return
		}
		errStartManager := services.ContainerManager.StartedContainer(dockbox.ContainerId.String)
		if errStartManager != nil {
			c.JSON(404, gin.H{"message": "Error creating connection", "error": errStartManager.Error()})
			return
		}
	}

	backendURL := fmt.Sprintf("ws://%s/containers/%s/attach/ws?logs=0&stream=1&stdin=1&stdout=1&stderr=1", utils.Config.DOCKER_SERVER_HOST, dockbox.ContainerId.String)
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
		OnConnect: func() {
			errConnect := services.ContainerManager.NewConnection(dockbox.ContainerId.String)
			if errConnect != nil {
				log.Fatalf("Error occurred with new connection: %v\n", errConnect)
			}
		},
		OnClose: func() {
			log.Printf("Closed connection")
			errClose := services.ContainerManager.RemoveConnection(dockbox.ContainerId.String)
			if err != nil {
				log.Fatalf("Error occurred with close connection: %v\n", errClose)
			}
		},
	}

	proxy.ServeHTTP(c.Writer, c.Request)

}

type Resource struct {
	Name     string
	Path     string
	Hash     string
	Type     string
	Children []*Resource
}

func (d DockboxController) GetFilesystem(c *gin.Context) {
	id, _ := c.Params.Get("id")

	if id == "" {
		c.AbortWithStatusJSON(500, map[string]string{"error": "no ID supplied"})
	}
	dirPath := filepath.Join(utils.Config.MOUNT_POINT, id)
	if _, err := os.Stat(dirPath); err != nil {
		c.AbortWithStatusJSON(500, map[string]string{"error": "ID not found"})
	}

	rootNode := &Resource{
		Name:     filepath.Base(dirPath),
		Path:     "/app",
		Children: make([]*Resource, 0),
	}
	curRoot := rootNode

	godirwalk.Walk(dirPath,
		&godirwalk.Options{
			Callback: func(osPathname string, d *godirwalk.Dirent) error {
				//Compute hash
				//Add Resource to children
				//update root hash
				name := filepath.Base(osPathname)
				resource := &Resource{
					Name: name,
					Path: filepath.Join(curRoot.Path, name),
				}
				curRoot.Children = append(curRoot.Children, resource)
				return nil
			},
			ErrorCallback: func(path string, err error) godirwalk.ErrorAction {
				log.Printf("Error accessing file: %s", path)
				return godirwalk.SkipNode
			},
			PostChildrenCallback: func(path string, entry *godirwalk.Dirent) error {

				return nil
			},
		},
	)

}
