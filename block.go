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

const (
	NumZeros = 28
	NumTries = 5000
)

type Block struct {
	Transactions []Transaction
	SeqNum       uint64
	ProofOfWork  []byte
	Hash         [sha256.Size]byte
	ParentHash   [sha256.Size]byte
}

var (
	CurrentBlock = Block{
		Transactions: nil,
		SeqNum:       1,
		ProofOfWork:  []byte{},
		Hash:         [sha256.Size]byte{},
	}
	// initialize all blockchains with dummy block of seqnum 0
	BlockChain map[uint64]Block = map[uint64]Block{
		0: Block{
			Transactions: nil,
			SeqNum:       0,
			ProofOfWork:  []byte{},
			Hash:         [sha256.Size]byte{},
		},
	}
)

// gets the hash of the block when the proof of work is already computed
func (b *Block) GetHash() [sha256.Size]byte {
	toHash := append(b.strToHash(b.ParentHash), b.ProofOfWork...)
	checksum := sha256.Sum256(toHash)
	return checksum
}

// computes the string of parenthash + transaction json strings
func (b *Block) strToHash(parentHash [sha256.Size]byte) []byte {
	var (
		toHash    = parentHash[:]
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
func (b *Block) SetProofOfWork(parentHash [sha256.Size]byte, c chan Block) {
	toHash := b.strToHash(parentHash)

	// resulting hash must begin with numZero zeros
	target := uint64(^uint64(0) >> NumZeros)
	nonceBuf := make([]byte, 8)
	nonce := uint64(0)
	checksum := [sha256.Size]byte{}
	hashNum := uint64(^uint(0))

	// check every NumTries times to see if we received a block in the queue
	// drop block and try again with updated parent block
	ctr := 0
	dropBlock := false
	pBlock := Block{}
	for hashNum > target && !dropBlock {
		if ctr%NumTries == 0 {
			dropBlock = false
			for {
				select {
				case pBlock := <-c:
					dropBlock := true
					continue
				default:
					break
				}
			}
		}
		binary.PutUvarint(nonceBuf, nonce)
		tryHash := append(toHash, nonceBuf...)
		checksum = sha256.Sum256(toHash)
		hashNum = binary.BigEndian.Uint64(checksum[0:sha256.Size])
		nonce++
	}
	if dropBlock {
		b.SeqNum = pBlock.SeqNum + 1
		go b.SetProofOfWork(pBlock.Hash, c)
		return
	}
	b.ProofOfWork = nonceBuf
	b.Hash = checksum
	b.ParentHash = parentHash
}

// verify proof of work -- invariant: the parent exists in the map
// 		- check that the block's parent's hash matches the hash of the parent block (seqNum - 1)
func (b *Block) ValidateHash() error {
	// VALIDATE BLOCK'S HASH (Proof of Work)
	if b.ParentHash != BlockChain[b.SeqNum-1].Hash {
		return fmt.Errorf("invalid parent block hash")
	}
	toHash := b.GetHash()
	target := uint64(^uint64(0) >> NumZeros)
	checksum := sha256.Sum256(toHash[:])

	// ensure checksum begins with 28 0s
	if binary.BigEndian.Uint64(checksum[0:sha256.Size]) > target {
		return fmt.Errorf("invalid proof of work, hash does not begin with 28 0s")
	}
	// ensure that block hash matches hash of parent + proof of work
	if bytes.Equal(checksum[:], b.Hash[:]) {
		return fmt.Errorf("hash does not match that of parent")
	}
	return nil
}

func (b *Block) ValidateTxn() error {
	for _, txn := range b.Transactions {
		// get the bytes to hash
		jsonBytes, err := json.Marshal(&txn)

		// did not find previous transaction of this user
		// must be a registration and signed by the CA
		lastPubKey, ok := Database[txn.Email]
		if !ok {
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
		if err := rsa.VerifyPKCS1v15(&lastPubKey, crypto.SHA256, jsonBytes, txn.Signature); err != nil {
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
