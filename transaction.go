package main

import (
	"crypto/rsa"
	"fmt"
)

type TransType int

const (
	CASig              = "signature of CA"
	Register TransType = 1 + iota
	Update
)

type Transaction struct {
	Type      TransType
	Email     string
	PublicKey rsa.PublicKey
	Signature string
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
	trans := &Transaction{
		Type:      Register,
		Email:     email,
		PublicKey: key,
		Signature: CASig,
	}
	// how to implement?
	// XXX we want to broadcast this somehow and have other nodes (including us) add this transaction to their blocks
	// we also want to add this to our "block" that we're working on?
	trans.SendTransaction()

	return nil
}
