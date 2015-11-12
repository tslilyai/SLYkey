package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

const numTrans = 5

type Block struct {
	Transactions []Transaction
	SeqNum       uint64
	ProofOfWork  uint64
	Hash         uint64
}

var (
	CurrentBlock = Block{
		Transactions: nil,
		SeqNum:       1,
		ProofOfWork:  0,
		Hash:         0,
	}
	// initialize all blockchains with dummy block of seqnum 0
	BlockChain map[uint64]Block = map[uint64]Block{
		0: Block{
			Transactions: nil,
			SeqNum:       0,
			ProofOfWork:  0,
			Hash:         0,
		},
	}
)

// add a transaction to a block, beginning work on a block if the node is not currently working on a block
func addToBlock(t Transaction) {
	CurrentBlock.Transactions = append(CurrentBlock.Transactions, t)
}

// clear the current block to fill up another block with transactions
func clearCurrentBlock() {
	CurrentBlock = Block{
		Transactions: nil,
		SeqNum:       CurrentBlock.SeqNum + 1,
		ProofOfWork:  0,
	}
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
