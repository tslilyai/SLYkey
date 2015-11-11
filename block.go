package main

import (
	"crypto/sha256"
	"encoding/binary"
)

const numTrans = 5

type Block struct {
	Transactions []Transaction
	SeqNum       uint64
	ProofOfWork  uint64
	Hash         uint64
}

var (
	BlockChain   map[uint64]Block
	CurrentBlock = Block{
		Transactions: nil,
		SeqNum:       0,
		ProofOfWork:  0,
		Hash:         0,
	}
)

// add a transaction to a block, beginning work on a block if the block becomes full
func addToBlock(b *Block, t Transaction) {
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
	setProofOfWork(b)
	broadcast(b)
	updateDatabase(b)
	return nil
}

// compute and set the proof of work and hash of the block
func setProofOfWork(b *Block) {
	chainHash := uint64(0)
	if b.SeqNum != 0 {
		chainHash = BlockChain[b.SeqNum-1].Hash
	}
	// resulting hash must begin with 28 zeros
	target := uint64(^uint64(0) >> 28)
	buf := make([]byte, 16)
	nonce := uint64(0)
	hash := uint64(^uint(0))
	for hash > target {
		binary.PutUvarint(buf, chainHash+nonce)
		checksum := sha256.Sum256(buf)
		hash = binary.BigEndian.Uint64(checksum[0:32])
		nonce++
	}
	b.ProofOfWork = nonce
	b.Hash = hash
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
		Transactions: nil,
		SeqNum:       CurrentBlock.SeqNum + 1,
		ProofOfWork:  0,
	}
}
