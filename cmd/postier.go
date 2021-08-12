package main

import (
	"log"
	"net/http"
	"os"

	"github.com/factorysh/postier/internal/pkg/handlers"
	"github.com/factorysh/postier/pkg/history"
)

var defaultListenURL = "0.0.0.0:8042"
var defaultHistoryEndpoint = "/postier-history"

func main() {

	listenURL := os.Getenv("LISTEN_URL")
	if listenURL == "" {
		log.Printf("Warning, no LISTEN_URL provided, using default one (%s)\n", defaultListenURL)
		listenURL = defaultListenURL
	}

	// hist endpoint is used to select the GET url used for retreiving history data
	historyEndpoint := os.Getenv("HISTORY_ENDPOINT")
	if historyEndpoint == "" {
		log.Printf("Warning, no HISTORY_ENDPOINT provided, using default one (%s)\n", defaultHistoryEndpoint)
		historyEndpoint = defaultHistoryEndpoint
	}
	log.Printf("Info, history url %s\n", historyEndpoint)

	m := history.NewMemory()

	http.Handle("/", handlers.Handle(&m))
	http.Handle(historyEndpoint, handlers.HandleHistory(&m))

	http.ListenAndServe("0.0.0.0:8042", nil)
}
