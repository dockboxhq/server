package server

import (
	"os"

	"github.com/gin-gonic/gin"
)

type config struct {
	PORT        string
	ENVIRONMENT string
}

func getConfig() *config {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8000"
	}
	ENVIRONMENT := os.Getenv("ENVIRONMENT")
	if ENVIRONMENT == "" {
		ENVIRONMENT = "development"
	}
	if ENVIRONMENT != "development" && ENVIRONMENT != "production" {
		panic("Unknown environment variable set. Must be one of development or production")
	}
	return &config{
		ENVIRONMENT: ENVIRONMENT,
		PORT:        PORT,
	}
}

func Init() {
	config := getConfig()
	if config.ENVIRONMENT == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := NewRouter()
	r.Run(":" + config.PORT)
}
