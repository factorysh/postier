package history

import (
	"net/http"
	"sync"
	"time"
)

// Request represent data about a received request
type Request struct {
	TS      time.Time   `json:"timestamp"`
	Host    string      `json:"host"`
	URL     string      `json:"url"`
	Headers http.Header `json:"headers"`
	Method  string      `json:"method"`
	Body    string      `json:"body"`
}

// Memory is an array of requests
type Memory struct {
	Mutex    sync.RWMutex `json:"-"`
	Requests []Request    `json:"requests"`
	// Count all requests that are not post requests
	NonPostCounter int `json:"non_post_requests"`
	// Count all body read errors
	ReadBodyErrorCounter int `json:"read_body_errors"`
}

// NewMemory inits a new memory tape
func NewMemory() Memory {
	return Memory{
		Requests: make([]Request, 0),
	}
}
