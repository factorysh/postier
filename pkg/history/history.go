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
}

// NewMemory inits a new memory tape
func NewMemory() Memory {
	return Memory{
		Requests: make([]Request, 0),
	}
}
