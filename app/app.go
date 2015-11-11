/*
Package app defines a new application server and registers endpoints.
*/
package app

import (
	"net/http"
	"time"

	"github.com/stretchr/graceful"
)

// App defines the interface for an application
type App interface {
	Run() error
}

type app struct {
	server  *graceful.Server // server to listen to external/internal requests
	handler *http.ServeMux   // handle requests to external/internal endpoints
}

// NewApp returns a new application
func NewApp() (App, error) {
	a := &app{}
	a.server = &graceful.Server{
		Server: &http.Server{
			Addr:    ":80",
			Handler: a.handler,
		},
		Timeout: time.Duration(500 * time.Millisecond),
	}
	//	a.handler.Handle("/endpoint" /* handler here */)
	return a, nil
}

// Run runs the application
func (a *app) Run() error {
	return a.server.ListenAndServe()
}
