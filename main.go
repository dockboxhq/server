package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dockboxhq/server/server"
	db "github.com/dockboxhq/server/services"
	config "github.com/dockboxhq/server/utils"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	config.PopulateConfig()

	// environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	db.Init()
	server.Init()
}
