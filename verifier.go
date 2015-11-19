package main

import (
	"crypto"
	"crypto/rsa"
	"encoding/json"
	"fmt"
)

// TODO(Serena): add to ValidateTxn a local database copy

// TODO(Serena): put these functions in block.go
// Validates transactions and simultaneously updates the database
// for all transactions in a block.
func ValidateTxnsAndUpdateDatabase(b Block) error {
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
			// update database with the register request
			Database[txn.Email] = txn.PublicKey
			return nil
		}
		// else this is an update
		if err := rsa.VerifyPKCS1v15(&lastPubKey, crypto.SHA256, jsonBytes, txn.Signature); err != nil {
			if txn.Type != Update {
				return fmt.Errorf("Cannot register if you already are in the database")
			}
			return fmt.Errorf("Signature on new transaction does not match")
		}
		// update database with the update request
		Database[txn.Email] = txn.PublicKey
	}
	return nil
}

// Validates both the block hash and the transactions
func ValidateBlockAndUpdateDatabase(b Block) error {
	if err := b.ValidateHash(); err != nil {
		return err
	}
	if err := ValidateTxnsAndUpdateDatabase(b); err != nil {
		return err
	}
	return nil
}

func VerifyBlockChainAndUpdateDatabase() error {
	// for each block in block chain
	seqNum := uint64(0)
	// Sync: need to lock whole blockchain here by creating
	// getBlock function instead.
	for i := 0; i < len(BlockChain); i++ {
		if block, ok := BlockChain[seqNum]; ok {
			// If that returned an error, raise error here.
			// Proper go error handling?
			if err := ValidateBlockAndUpdateDatabase(block); err != nil {
				return fmt.Errorf("Verifier could not complete: invalid block")
			}
			seqNum++
		}
	}
	return nil
}

func RunVerifier() {
	VerifyBlockChainAndUpdateDatabase()
}
