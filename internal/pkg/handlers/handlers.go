package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/factorysh/postier/pkg/history"
)

// Handle is a catch all
func Handle(m *history.Memory) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Mutex.Lock()
		defer m.Mutex.Unlock()

		if r.Method != http.MethodPost {
			m.NonPostCounter++
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			m.ReadBodyErrorCounter++
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error reading body content"))
			return
		}

		m.Requests = append(m.Requests, history.Request{
			TS:      time.Now(),
			Host:    r.Host,
			URL:     r.URL.Path,
			Method:  r.Method,
			Headers: r.Header,
			Body:    string(content),
		})

		w.WriteHeader(http.StatusOK)
	})
}

// HandleHistory respond to a get request on /history
func HandleHistory(m *history.Memory) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		m.Mutex.Lock()
		defer m.Mutex.Unlock()

		data, err := json.Marshal(m)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
}
