package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/factorysh/postier/internal/pkg/handlers"
	"github.com/factorysh/postier/pkg/history"
)

// DefaultHistoryEndpoint is used as a fallback value for history endpoint
const DefaultHistoryEndpoint = "/postier-history"

// DefaultTimeout is used as a fallback for shutdown timeout value
const DefaultTimeout = 2 * time.Second

// NewServer inits and return a configured postier server
func NewServer(listenURL, historyEndpoint string) *http.Server {
	m := history.NewMemory()

	mux := http.NewServeMux()
	mux.Handle("/", handlers.Handle(&m))
	mux.Handle(historyEndpoint, handlers.HandleHistory(&m))

	srv := &http.Server{
		Addr:    listenURL,
		Handler: mux,
	}

	return srv
}

// Start inits and starts a new server with signal handling for gracefull stop
func Start(listenURL, historyEndpoint string) {
	remote, waiter := StartWithControls(listenURL, historyEndpoint)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		remote <- true
	}()

	<-waiter
}

// StartWithControls inits and starts the server, returning two controls channels (remote for shutdow, waiter for shutdown feedback)
func StartWithControls(listenURL, historyEndpoint string) (chan bool, chan bool) {
	srv := NewServer(listenURL, historyEndpoint)
	remote := make(chan bool)
	waiter := make(chan bool)

	go func() {
		log.Printf("Server listening on %s with history endpoint set to %s", srv.Addr, historyEndpoint)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
		log.Println("Server gracefully stopped")
	}()

	go func() {
		<-remote

		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}

		waiter <- true
	}()

	return remote, waiter
}
