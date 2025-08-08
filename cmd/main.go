package main

import (
	"downloader/internal/config"
	"downloader/internal/server"
	"log"
)

func main() {
	newConfig, err := config.LoadConfig("config.yaml")

	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	server.NewServer(newConfig)
}
