package main

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

const NumZeros = 28

type Block struct {
	Transactions []Transaction
	SeqNum       uint64
	ProofOfWork  []byte
	Hash         []byte
}

var (
	CurrentBlock = Block{
		Transactions: nil,
		SeqNum:       1,
		ProofOfWork:  []byte{},
		Hash:         []byte{},
	}
	// initialize all blockchains with dummy block of seqnum 0
	BlockChain map[uint64]Block = map[uint64]Block{
		0: Block{
			Transactions: nil,
			SeqNum:       0,
			ProofOfWork:  []byte{},
			Hash:         []byte{},
		},
	}
)

// gets the hash of the block when the proof of work is already computed
func (b *Block) GetHash(parentHash []byte) []byte {
	toHash := append(b.strToHash(parentHash), b.ProofOfWork...)
	checksum := sha256.Sum256(toHash)
	return checksum[:]
}

// computes the string of parenthash + transaction json strings
func (b *Block) strToHash(parentHash []byte) []byte {
	var (
		toHash    = parentHash
		jsonBytes []byte
		err       error
	)

	for _, txn := range b.Transactions {
		jsonBytes, err = json.Marshal(&txn)
		toHash = append(toHash, jsonBytes...)
	}
	return toHash
}

// compute and set the proof of work and hash of the block
// we will want to hash the (block transactions + parent hash + pow/nonce)
func (b *Block) setProofOfWork(parentHash []byte) {
	toHash := b.strToHash(parentHash)

	// resulting hash must begin with numZero zeros
	target := uint64(^uint64(0) >> NumZeros)
	nonceBuf := make([]byte, 8)
	nonce := uint64(0)

	checksum := [32]byte{}
	hashNum := uint64(^uint(0))
	for hashNum > target {
		binary.PutUvarint(nonceBuf, nonce)
		tryHash := append(toHash, nonceBuf...)
		checksum = sha256.Sum256(toHash)
		hashNum = binary.BigEndian.Uint64(checksum[0:32])
		nonce++
	}
	b.ProofOfWork = nonceBuf
	b.Hash = checksum[:]
}

// verify proof of work -- invariant: the parent exists in the map
// 		- check that the block's parent's hash matches the hash of the parent block (seqNum - 1)
func (b *Block) validateHash() error {
	// VALIDATE BLOCK'S HASH (Proof of Work)
	toHash := b.GetHash(BlockChain[b.SeqNum-1].Hash)
	target := uint64(^uint64(0) >> NumZeros)
	checksum := sha256.Sum256(toHash)

	// ensure checksum begins with 28 0s
	if binary.BigEndian.Uint64(checksum[0:32]) > target {
		return fmt.Errorf("invalid proof of work, hash does not begin with 28 0s")
	}
	// ensure that block hash matches hash of parent + proof of work
	if bytes.Equal(checksum[:], b.Hash) {
		return fmt.Errorf("hash does not match that of parent")
	}
	return nil
}

func (b *Block) validateTxn() error {
	// VALIDATE WITH RESPECT TO BLOCK CHAIN by finding the most recent transaction pertaining to the user,
	// checking the signature with the transaction to be added or verifying the CASig if there has been no
	// user-made transaction thus far.
	var lastTxn Transaction
	for _, txn := range b.Transactions {
		// get most recent block
		found := false
		for i := len(BlockChain); i > 0; i-- {
			for _, t := range b.Transactions {
				if t.Email == txn.Email {
					lastTxn = t
					found = true
				}
			}
			if found {
				break
			}
		}

		// get the bytes to hash
		jsonBytes, err := json.Marshal(&txn)

		// did not find previous transaction of this user
		// must be a registration and signed by the CA
		if !found {
			if txn.Type != Register {
				return fmt.Errorf("Cannot update a nonexistent public key")
			}
			// verify the CA signed this request

			if err := rsa.VerifyPKCS1v15(&CAKey, crypto.SHA256, jsonBytes, txn.Signature); err != nil {
				return fmt.Errorf("Not signed by the CA")
			}
			return nil
		}
		// else this is an update
		if err := rsa.VerifyPKCS1v15(&lastTxn.PublicKey, crypto.SHA256, jsonBytes, lastTxn.Signature); err != nil {
			if txn.Type != Update {
				return fmt.Errorf("Cannot register if you already are in the database")
			}
			return fmt.Errorf("Signature on new transaction does not match")
		}
	}
	return nil
}

// add a transaction to a block, beginning work on a block if the node is not currently working on a block
func addToBlock(t Transaction) {
	CurrentBlock.Transactions = append(CurrentBlock.Transactions, t)
}

// clear the current block to fill up another block with transactions
func clearCurrentBlock() {
	CurrentBlock = Block{
		Transactions: nil,
		SeqNum:       CurrentBlock.SeqNum + 1,
		ProofOfWork:  []byte{},
	}
}
