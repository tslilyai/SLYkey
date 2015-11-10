package main

type TransType int

const (
	CASig                     = "signature of CA"
	NumTransactions           = 5
	Register        TransType = 1 + iota
	Update
)

type Transaction struct {
	Type      TransType
	Email     string
	PublicKey string
	Signature string
}

func GetPublicKey(string email) string {
	return Database[email]
}

func RegisterPublicKey(string key, string email) *Transaction {
	// XXX we want to broadcast this somehow
	// we also want to add this to our "block" that we're working on?
	return &Transaction{
		Type:      Register,
		Email:     email,
		PublicKey: key,
		Signature: CASig,
	}
}
