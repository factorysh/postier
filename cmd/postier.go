package main

import (
	"log"
	"os"

	"github.com/factorysh/postier/pkg/server"
)

const defaultListenURL = "0.0.0.0:8042"

func main() {
	listenURL := os.Getenv("LISTEN_URL")
	if listenURL == "" {
		log.Printf("Warning, no LISTEN_URL provided, using default one (%s)\n", defaultListenURL)
		listenURL = defaultListenURL
	}

	historyEndpoint := os.Getenv("HISTORY_ENDPOINT")
	if historyEndpoint == "" {
		log.Printf("Warning, no history endpoint provided, using default one (%s)\n", server.DefaultHistoryEndpoint)
		historyEndpoint = server.DefaultHistoryEndpoint
	}

	server.Start(listenURL, historyEndpoint)
}
