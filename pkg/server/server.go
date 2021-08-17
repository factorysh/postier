package server

import (
	"context"
	"log"
	"net"
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
func NewServer(historyEndpoint string) (*http.Server, *history.Memory) {
	m := history.NewMemory()

	mux := http.NewServeMux()
	mux.Handle("/", handlers.Handle(&m))
	mux.Handle(historyEndpoint, handlers.HandleHistory(&m))

	srv := &http.Server{
		Handler: mux,
	}

	return srv, &m
}

// Start inits and starts a new server with signal handling for gracefull stop
func Start(listenURL, historyEndpoint string) error {
	l, err := net.Listen("tcp", listenURL)
	if err != nil {
		return err
	}
	defer l.Close()

	remote, waiter, _ := StartWithControls(l, historyEndpoint)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		remote <- true
	}()

	<-waiter

	return nil
}

// StartWithControls inits and starts the server, returning two controls channels (remote for shutdow, waiter for shutdown feedback) and a pointer to history data
func StartWithControls(listener net.Listener, historyEndpoint string) (chan bool, chan bool, *history.Memory) {
	srv, memory := NewServer(historyEndpoint)
	remote := make(chan bool)
	waiter := make(chan bool)

	go func() {
		log.Printf("Server listening on %s with history endpoint set to %s", listener.Addr(), historyEndpoint)
		if err := srv.Serve(listener); err != http.ErrServerClosed {
			log.Fatal(err)
		}
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

	return remote, waiter, memory
}
