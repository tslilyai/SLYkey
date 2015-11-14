package main

const (
	MAX_QUEUE = 64
)

// Go doesn't give us a queue
// so this is our queue...
// Many thanks to https://gist.github.com/moraes/2141121
type BlockQueue struct {
	queue []Block
	head  uint64
	tail  uint64
	size  uint64
	count uint64
}

func NewBlockQueue(initSize uint64) *BlockQueue {
	ret := &BlockQueue{}
	ret.queue = make([]Block, initSize)
	ret.head = 0
	ret.tail = 0
	ret.size = initSize
	ret.count = 0
	return ret
}

func (bq *BlockQueue) Count() uint64 {
	return bq.count
}

func (bq *BlockQueue) Push(b Block) {
	if bq.head == bq.tail && bq.count > 0 {
		realloc := make([]Block, uint64(len(bq.queue))+bq.size)
		copy(realloc, bq.queue[bq.head:])
		copy(realloc[uint64(len(bq.queue))-bq.head:], bq.queue[:bq.head])
		bq.head = 0
		bq.tail = uint64(len(bq.queue))
		bq.queue = realloc
	}
	bq.queue[bq.tail] = b
	bq.tail = (bq.tail + 1) % uint64(len(bq.queue))
	bq.count++
}

func (bq *BlockQueue) Pop() Block {
	if bq.count == 0 {
		return Block{}
	}
	b := bq.queue[bq.head]
	bq.head = (bq.head + 1) % uint64(len(bq.queue))
	bq.count--
	return b
}
