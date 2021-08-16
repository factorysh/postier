package history

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

// Requests represent a list of requests structs
type Requests []Request

// FilterURL filters requests that matches a pattern found in the url
func (r Requests) FilterURL(pattern string) Requests {
	if pattern == "" {
		return r
	}

	filtered := make(Requests, 0)
	for _, req := range r {
		if strings.Contains(req.URL, pattern) {
			filtered = append(filtered, req)
		}
	}

	return filtered
}

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
	Requests Requests     `json:"requests"`
}

// NewMemory inits a new memory tape
func NewMemory() Memory {
	return Memory{
		Requests: make(Requests, 0),
	}
}
