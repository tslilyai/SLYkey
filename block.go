package main

import "container/list"

type TransType int

const (
	NumTransactions           = 5
	Register        TransType = 1 + iota
	Update
)

type Transaction struct {
	Type       TransType
	Email      string
	Public_key string
	Signature  string
}

type Block struct {
	Type         TransType
	Transactions [NumTransactions]Transaction
	Email        string
	SeqNum       int
	ProofOfWork  int
}

var BlockChain *list.List

// maps from email to public key
var Database map[string]string
