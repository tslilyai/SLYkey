package main

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
)

type TransType int

type Transaction struct {
	Type      TransType     `json:"type"`
	Email     string        `json:"email"`
	PublicKey rsa.PublicKey `json:"public_key"`
	Signature []byte
}

const (
	Register TransType = 1 + iota
	Update
)

var (
	Database map[string]rsa.PublicKey
	// XXX need to have actual key
	CAKey = rsa.PublicKey{
		N: big.NewInt(3),
		E: 3,
	}
	CAurl = "https://ca.com/register"
)

func updateDatabase(b *Block) {
	// we should have already checked if txn and signatures are valid
	for _, txn := range b.Transactions {
		Database[txn.Email] = txn.PublicKey
	}
}

func GetPublicKey(email string) rsa.PublicKey {
	return Database[email]
}

// Returns error on failure, nil on success
// Registers a public key transaction, signed by the CA
func RegisterPublicKey(key rsa.PublicKey, email string) error {
	// value already in map, don't reregister
	if _, ok := Database[email]; ok {
		return fmt.Errorf("You have already registered for a public key")
	}
	jsonBytes, err := json.Marshal(&Transaction{Type: Register, Email: email, PublicKey: key})
	if err != nil {
		log.Print(err)
	}

	// make request
	res, err := http.Post(CAurl, "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		log.Print(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	// check request response and log errors
	if res.StatusCode != 200 {
		log.Print(err)
	}
	trans := Transaction{
		Type:      Register,
		Email:     email,
		PublicKey: key,
		Signature: body,
	}
	// add this to our "block" that we're working on
	addToBlock(trans)
	return nil
}

// Returns error on failure, nil on success
// Updates a public key, signed by the user
// signature should be signed on an empty byte hash
func UpdatePublicKey(key rsa.PublicKey, sig []byte, email string) error {
	// value already in map, don't reregister
	oldKey, ok := Database[email]
	if !ok {
		return fmt.Errorf("You have never registered for a public key")
	}
	if err := rsa.VerifyPKCS1v15(&oldKey, 0, []byte{}, sig); err != nil {
		return fmt.Errorf("bad signature")
	}
	trans := Transaction{
		Type:      Update,
		Email:     email,
		PublicKey: key,
		Signature: sig,
	}
	// add this to our "block" that we're working on
	addToBlock(trans)
	return nil
}
