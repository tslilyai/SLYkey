package main

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

// add a transaction to a block, beginning work on a block if the block becomes full
// XXX or if the node is not currently working on a block
func addToBlock(b *Block, t Transaction) {
	if len(b.Transactions) == numTrans {
		block := CurrentBlock
		clearCurrentBlock()
		process(&block)
	} else {
		b.Transactions = append(b.Transactions, t)
	}
}

// clear the current block to fill up another block with transactions
func clearCurrentBlock() {
	CurrentBlock = Block{
		Transactions: nil,
		SeqNum:       CurrentBlock.SeqNum + 1,
		ProofOfWork:  0,
	}
}
