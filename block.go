package main

import "container/list"

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
