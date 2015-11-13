package ca

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/stretchr/graceful"
	"github.com/tslilyai/SLYcoin"
)

var (
	port        = flag.Int("port", 8080, "HTTP port")
	ErrPost     = "must use POST"
	ErrDecode   = "bad transaction json data"
	ErrBadEmail = "bad email address"
	ErrBadType  = "must be a request transaction"
	// XXX get from database or something
	privateKey = rsa.PrivateKey
)

// App defines an application that can be run
type App interface {
	Run() error
}

// NewApp returns a new application to run
func NewApp() (App, error) {
	handler := http.NewServeMux()
	server := &graceful.Server{
		Server: &http.Server{
			Addr:    ":" + strconv.Itoa(*port),
			Handler: handler,
		},
		Timeout: 2 * time.Second,
	}

	app := &app{server}

	handler.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		app.registerReq(w, r)
	})

	return app, nil
}

type app struct {
	server *graceful.Server
}

// Run starts the app
func (a *app) Run() error {
	return a.server.ListenAndServe()
}

func (a *app) registerReq(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, ErrPost, http.StatusMethodNotAllowed)
		return
	}
	// attempt to decode data
	var data main.Transaction
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, ErrDecode, http.StatusBadRequest)
	}
	// validate the registration type
	if data.Type != "register" {
		http.Error(w, ErrBadType, http.StatusBadRequest)
	}
	// validate the email
	if !validateEmail(data.Email) {
		http.Error(w, ErrBadEmail, http.StatusBadRequest)
	}

	// sign the transaction, with the hash of the public key
	s, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, []byte(data.Email))
	if err != nil {
		http.Error(w, "Could not sign registration request", http.StatusInternalServerError)
	}

	// return the signature of the hash of the email
	fmt.Fprintf(w, s)
}

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}
