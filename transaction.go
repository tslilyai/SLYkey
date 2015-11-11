package main

type TransType int

const (
	CASig              = "signature of CA"
	Register TransType = 1 + iota
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

// Returns error on failure, nil on success
// Registers a public key transaction, signed by the CA
func RegisterPublicKey(string key, string email) error {
	trans := &Transaction{
		Type:      Register,
		Email:     email,
		PublicKey: key,
		Signature: CASig,
	}
	// XXX we want to broadcast this somehow
	// we also want to add this to our "block" that we're working on?
	return 0
}
