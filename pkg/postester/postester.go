package postester

import (
	"fmt"
	"net"

	"github.com/factorysh/postier/pkg/history"
	"github.com/factorysh/postier/pkg/server"
)

const localhost = "127.0.0.1"

// Env wraps public variables usable when postier is used in golang testing environment
type Env struct {
	Port            int
	URL             string
	HistoryEndpoint string
	remote          chan bool
	waiter          chan bool
	memory          *history.Memory
}

// Cleanup is used to gracefully shutdown the postier env
func (e Env) Cleanup() {
	e.remote <- true
	<-e.waiter
}

// History is used to fetch all requests send to postier endpoint
func (e Env) History() history.Requests {
	e.memory.Mutex.Lock()
	defer e.memory.Mutex.Unlock()
	return e.memory.Requests
}

// StartTesting is use to start postier in a golang testing context
func StartTesting() (*Env, error) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:0", localhost))
	if err != nil {
		return nil, err
	}

	addr := l.Addr().(*net.TCPAddr)

	r, w, m := server.StartWithControls(l, server.DefaultHistoryEndpoint)

	return &Env{
		Port:            addr.Port,
		URL:             fmt.Sprintf("http://%s:%d", localhost, addr.Port),
		HistoryEndpoint: server.DefaultHistoryEndpoint,
		remote:          r,
		waiter:          w,
		memory:          m,
	}, nil
}
