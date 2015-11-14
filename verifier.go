package main

// TODO(Serena): put these functions in block.go
func validateTransactionsAndUpdateDatabase (b *Block) error {
	// TODO(Serena):
	// validate block transactions and update database simultaneously.
	// Iterate from the beginning of the block. 
	// Assumes that the database is correct and up to date for the // all blocks prior to the current block. 
	// Thus, uses the database values as the most recent user key.
	// Also validates CA signature for creation of a new key.
}

func validateBlockAndUpdateDatabase(b *Block) error {
	// validate block hash
	validateTransactionsAndUpdateDatabase(Block)
	// Handle errors correctly that bubble up from helper function?
}

func verifyBlockChainAndUpdateDatabase() {
	// for each block in block chain
	seqNum = uint64(0)
	// Sync: need to lock whole blockchain here by creating
	// getBlock function instead.
	while block, ok := BlockChain[args.SeqNum]; ok {
		// If that returned an error, raise error here. 
		// Proper go error handling?
		validateBlockAndUpdateDatabase(block)
		seqNum++
	}
}

func runVerifier () {
	verifyBlockChainAndUpdateDatabase()
}

