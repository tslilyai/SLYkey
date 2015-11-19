package main

import "fmt"

func VerifyBlockChainAndUpdateDatabase() error {
	// for each block in block chain
	seqNum := uint64(0)
	// Sync: need to lock whole blockchain here by creating
	// getBlock function instead.
	for i := 0; i < len(BlockChain); i++ {
		if block, ok := BlockChain[seqNum]; ok {
			// If that returned an error, raise error here.
			// Proper go error handling?
			if err := block.ValidateHash(); err != nil {
				return fmt.Errorf("Verifier could not complete: invalid block hash")
			}
			if err := block.ValidateTxn(); err != nil {
				return fmt.Errorf("Verifier could not complete: invalid block transactions")
			}
			updateDatabase(&block)
			seqNum++
		}
	}
	return nil
}
