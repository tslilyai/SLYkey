package main

type Block struct {
	Transactions []Transaction
	SeqNum       uint64
	ProofOfWork  uint64
}

var BlockChain map[uint64]Block

// maps from email to public key
var Database map[string]string
