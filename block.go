package main

import "crypto/rsa"

const numTrans = 5

type Block struct {
	Transactions []Transaction
	SeqNum       uint64
	ProofOfWork  uint64
}

var BlockChain map[uint64]Block
var CurrentBlock Block
var Database map[string]rsa.PublicKey

// add a transaction to a block, beginning work on a block if the block becomes full
func addToBlock(t Transaction) {
	if len(b.Transactions) == numTrans {
		block := CurrentBlock
		clearCurrentBlock()
		process(&block)
	} else {
		b.Transactions = append(b.Transactions, t)
	}
}

// process computes the proof of work and then broadcasts the block
// XXX drop block if we receive another block?
// we also then update the database and block chain with the block transactions
func process(b *Block) error {
	b.ProofOfWork = computeProofOfWork(b)
	broadcast(b)
	updateDatabase(b)
	return nil
}

// compute the proof of work
func computeProofOfWork(b *Block) uint64 {
	return 0
}

// broadcast the block to all nodes
func broadcast(b *Block) error {
	return nil
}

// receive blocks from other nodes
func receive() error {
	return nil
}

// request block of number seqNum from other nodes
func request(seqNum uint64) *Block {
	return nil
}

// validate block with respect to the block chain
// also validates by checking that the block's parent's hash matches
// the hash of the parent (seqNum - 1)
func validate(b *Block) error {
	return nil
}

// clear the current block to fill up another block with transactions
func clearCurrentBlock() {
	CurrentBlock = Block{
		Transactions: Nil,
		SeqNum:       CurrentBlock.SeqNum + 1,
		ProofOfWork:  0,
	}
	return nil
}
