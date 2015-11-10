/*
Package app defines a new application server and registers endpoints.
*/
package app

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
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
	a.setupServer()
	a.server = &graceful.Server{
		Server: &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: a.handler,
		},
		Timeout: time.Duration(serverTimeout) * time.Millisecond,
	}
	a.handler.Handle("/endpoint", handlers.HealthHandler(/*args here*/))
}	return a, nil
}

// Run runs the application
func (a *app) Run() error {
	for _, statsd := range a.statsdClients {
		defer statsd.Close()
	}
	return a.server.ListenAndServe()
}
