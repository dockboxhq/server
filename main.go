package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sriharivishnu/dockbox/server/db"
	"github.com/sriharivishnu/dockbox/server/server"
)

func main() {
	// environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	db.Init()
	server.Init()
}
