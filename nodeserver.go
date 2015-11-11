package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// XXX define nodeserver here?

// process computes the proof of work and then broadcasts the block
// XXX drop block if we receive another block?
// we also then update the database and block chain with the block transactions
func process(b *Block) error {
	setProofOfWork(b)
	broadcast(b)
	updateDatabase(b)
	return nil
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

// compute and set the proof of work and hash of the block
func setProofOfWork(b *Block) {
	chainHash := BlockChain[b.SeqNum-1].Hash
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

// validate block with respect to the block chain
// also validates by checking that the block's parent's hash matches
// the hash of the parent (seqNum - 1)
func validate(b *Block) error {
	// VALIDATE BLOCK'S HASH
	buf := make([]byte, 16)
	binary.PutUvarint(buf, BlockChain[b.SeqNum-1].Hash+b.ProofOfWork)
	target := uint64(^uint64(0) >> 28)
	checksum := sha256.Sum256(buf)
	hash := binary.BigEndian.Uint64(checksum[0:32])
	// ensure checksum begins with 28 0s
	if hash > target {
		return fmt.Errorf("invalid proof of work, hash does not begin with 28 0s")
	}
	// ensure that block hash matches hash of parent + proof of work
	if hash != b.Hash {
		return fmt.Errorf("hash does not match that of parent")
	}

	// VALIDATE WITH RESPECT TO BLOCK CHAIN
	// TODO
	return nil
}
