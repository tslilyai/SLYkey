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

func NewBlockQueue(initSize uint64) BlockQueue {
	ret := BlockQueue{}
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
		realloc := make([]Block, len(bq.queue) + bq.size)
		copy(realloc, bq.queue[bq.head:])
		copy(realloc[len(bq.queue) - bq.head:], bq.queue[:bq.head])
		bq.head = 0
		bq.tail = len(bq.queue)
		bq.queue = realloc
	}
	bq.queue[bq.tail] = n
	bq.tail = (bq.tail + 1) % len(bq.queue)
	bq.count++
}

func (bq *BlockQueue) Pop() Block {
	if bq.count == 0 {
		return Block{}
	}
	b := q.queue[q.head]
	q.head = (q.head + 1) % len(q.queue)
	q.count--
	return b
}
