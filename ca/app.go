package ca

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/stretchr/graceful"
	"github.com/tslilyai/SLYcoin"
)

var (
	port        = flag.Int("port", 8080, "HTTP port")
	privKeyFile = flag.String("key", "rsa_pub", "Private Key File")
	ErrPost     = "must use POST"
	ErrDecode   = "bad transaction json data"
	ErrBadEmail = "bad email address"
	ErrBadType  = "must be a request transaction"
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
	priv := getPrivateKey()
	app := &app{server, priv}

	handler.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		app.registerReq(w, r)
	})

	return app, nil
}

type app struct {
	server     *graceful.Server
	privateKey *rsa.PrivateKey
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

	// sign the transaction, with the hash of the transaction json (without signature)
	s, err := rsa.SignPKCS1v15(rand.Reader, a.privateKey, crypto.SHA256, []byte(r.Body))
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

func getPrivateKey() *rsa.PrivateKey {
	// Extract the PEM-encoded data block
	pemData := ioutil.ReadFile(*privKeyFile)
	block, _ := pem.Decode(pemData)
	if block == nil {
		log.Fatalf("bad key data: %s", "not PEM-encoded")
	}
	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		log.Fatalf("unknown key type %q, want %q", got, want)
	}

	// Decode the RSA private key
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("bad private key: %s", err)
	}
	return priv
}
