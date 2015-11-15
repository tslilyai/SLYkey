package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

const (
	MAX_PEERS = 64
)

type NodeServer struct {
	qMu           sync.Mutex // mutex for BlockQueue
	mMu           sync.Mutex // mutex for block chain (map)
	dead          int32
	rpcListener   net.Listener
	peers         []string
	workerChannel chan Block
	blkQueue      *BlockQueue
}

func assert(condition bool, func_name string) {
	if !condition {
		log.Fatal("Assertion failure in function %v", func_name)
	}
}

// addr: local address?
// peers: bunch of other servers?
// XXX should call StartRPCServer, NewBlockQueue, etc.
func NewNodeServer(addr string, peers []string) *NodeServer {
	ns := &NodeServer{
		// initialize fields here
		dead:          0,
		peers:         peers,
		workerChannel: make(chan Block, 16),
		blkQueue:      NewBlockQueue(1),
	}
	ns.StartRPCServer(addr)
	go ns.WorkOnBlock(BlockChain[0])
	return ns
}

func (ns *NodeServer) isdead() bool {
	return atomic.LoadInt32(&ns.dead) != 0
}

func (ns *NodeServer) Shutdown() {
	atomic.StoreInt32(&ns.dead, 1)
}

// RPC methods here!!
func (ns *NodeServer) SendBlock(remote string, block Block) bool {
	args := SendBlockArgs{block}
	reply := SendBlockReply{}
	args.Block = block
	ok := RPCCall(remote, "ns.RecvIncomingBlock", args, &reply)
	return ok
}

func (ns *NodeServer) RecvIncomingBlock(args *SendBlockArgs, reply *SendBlockReply) error {
	ns.qMu.Lock()
	defer ns.qMu.Unlock()

	// discard incoming blocks if our queue is full
	if ns.blkQueue.Count() >= MAX_QUEUE {
		return nil
	}

	ns.blkQueue.Push(args.Block)
	reply.Status = ErrOK
	return nil
}

func (ns *NodeServer) RequestBlock(remote string, seqNum uint64) (bool, Block) {
	args := RequestBlockArgs{}
	reply := RequestBlockReply{}
	args.SeqNum = seqNum
	ok := RPCCall(remote, "ns.RemoteBlockLookup", args, &reply)
	if ok && reply.Status == ErrFound {
		return true, reply.Block
	} else {
		return false, Block{}
	}
}

func (ns *NodeServer) RemoteBlockLookup(args *RequestBlockArgs, reply *RequestBlockReply) error {
	// XXX we need a mutex (called mMu right now) for the block chain
	ns.mMu.Lock()
	defer ns.mMu.Unlock()

	if block, ok := BlockChain[args.SeqNum]; ok {
		reply.Block = block
		reply.Status = ErrFound
		return nil
	}
	return fmt.Errorf(ErrNotFound)
}

// End of RPC methods

// Yihe's processing thread:
// checks the queue for incomings and notify worker thread
// if necessary
func (ns *NodeServer) ProcessBlock() error {
	for !ns.isdead() {
		ns.mMu.Lock()
		ns.qMu.Lock()
		if ns.blkQueue.Count() > 0 {
			b := ns.blkQueue.Pop()
			ns.qMu.Unlock()
			if ns.blockSanityCheck(b) == false {
				// skip garbage blocks
				continue
			}
			our_block, ok := BlockChain[b.SeqNum]
			if !ok {
				// fill up block chain here
				if ns.processUnseenBlock(b) {
					maxb := BlockChain[b.SeqNum]
					ns.mMu.Unlock()
					// signal worker of block b.SeqNum; could block so unlock
					ns.workerChannel <- maxb
				} else {
					ns.mMu.Unlock()
				}
			} else {
				if !ns.BlockCompare(b, our_block) {
					// peerCheckAndFixBlock will return the new max
					// sequence after it fixes the block chian
					// returns 0 if our block is valid
					max := ns.peerCheckAndFixBlock(b)
					if max > 0 {
						maxb := BlockChain[max]
						ns.mMu.Unlock()
						// signal worker of new job here
						ns.workerChannel <- maxb
					} else {
						ns.mMu.Unlock()
					}
				} else {
					ns.mMu.Unlock()
				}
			}
		} else {
			ns.qMu.Unlock()
			ns.mMu.Unlock()
		}
	}
	return nil
}

func (ns *NodeServer) HashEq(h1 [32]byte, h2 [32]byte) bool {
	for idx, b := range h1 {
		if b != h2[idx] {
			return false
		}
	}
	return true
}

// making sure the block isn't random garbage...
func (ns *NodeServer) blockSanityCheck(b Block) bool {
	hashNum := binary.BigEndian.Uint64(b.Hash[0:32])
	// XXX I wish I could just do b.GetHash() without passing in parent's hash
	// parent's hash should be part of the block, I think
	return ns.HashEq(b.GetHash(), b.Hash) && hashNum <= uint64(^uint64(0)>>NumZeros)
}

func (ns *NodeServer) BlockCompare(b1 Block, b2 Block) bool {
	assert(ns.blockSanityCheck(b1) && ns.blockSanityCheck(b2), "ns.BlockCompare")
	return ns.HashEq(b1.Hash, b2.Hash)
}

// Precondition: mMu acquired
func (ns *NodeServer) processUnseenBlock(b Block) bool {
	// block b is not seen before, make sure the largest block
	// we have is correct
	seq := b.SeqNum - 1
	for {
		_, ok := BlockChain[seq]
		if ok {
			break
		} else {
			seq = seq - 1
		}
	}

	// now we have seq as the highest block we have in the BlockChain
	maxb := BlockChain[seq]
	// make sure our highest block is still valid
	// if not, throw away that block and back track
	for {
		matches, peer_block := ns.peerCheckBlock(maxb)
		if !matches {
			seq--
			maxb = BlockChain[seq]
		} else {
			break
		}
	}

	// now we know everything so far seems valid; fill up our BlockChain upto b.SeqNum - 1
	seq++
	for {
		exists, peer_block := ns.peerRequestBlock(seq)
		if !exists || peer_block.ValidateHash() != nil || peer_block.ValidateTxn() != nil {
			// can't form a valid block chain, give up
			return false
		}
		ns.recordBlock(b)
		BlockChain[seq] = peer_block
		if seq >= b.SeqNum-1 {
			break
		}
		seq++
	}

	// check if b can be based on top of us
	if err := b.ValidateHash(); err != nil {
		if err := b.ValidateTxn(); err != nil {
			ns.recordBlock(b)
			BlockChain[b.SeqNum] = b
			return true
		}
	}
	return false
}

// Precondition: mMu acquired
func (ns *NodeServer) peerCheckAndFixBlock(b Block) uint64 {
	our_block := BlockChain[b.SeqNum]
	conflict := false
	seq := b.SeqNum
	for {
		matches, peer_block := ns.peerCheckBlock(our_block)
		if !matches {
			seq--
			our_block = BlockChain[seq]
			conflict = true
		} else {
			break
		}
	}
	if !conflict {
		return 0
	}
	// if conflicts, seq is the forking position
	seq++
	for {
		exists, peer_block := ns.peerRequestBlock(seq)
		if !exists {
			break
		}
		if peer_block.ValidateHash() == nil && peer_block.ValidateTxn() == nil {
			// adds block to database + blockchain
			ns.recordBlock(peer_block)
			BlockChain[seq] = peer_block
			seq++
		} else {
			// we have a problem here, just return the current seq
			// number
			return seq
		}
	}
	return (seq - 1)
}

func (ns *NodeServer) peerCheckBlock(ob Block) (bool, Block) {
	exists, peer_block := ns.peerRequestBlock(ob.SeqNum)
	if !exists || ns.BlockCompare(ob, peer_block) == true {
		return true, Block{}
	}
	return false, peer_block
}

func (ns *NodeServer) peerRequestBlock(seq uint64) (bool, Block) {
	// perform sanity checks on all received blocks
	// only return blocks that represent received majority
	var wg sync.WaitGroup
	var mu sync.Mutex
	blockMap := make(map[[32]byte]Block)
	blockCount := make(map[[32]byte]uint64)
	wg.Add(len(ns.peers))
	for _, peer := range ns.peers {
		go func(peer string) {
			found, blk := ns.RequestBlock(peer, seq)
			if found && ns.blockSanityCheck(blk) {
				mu.Lock()
				blockMap[blk.Hash] = blk
				v, ok := blockCount[blk.Hash]
				if !ok {
					v = 1
				} else {
					v++
				}
				blockCount[blk.Hash] = v
				mu.Unlock()
			}
			wg.Done()
		}(peer)
	}
	wg.Wait()

	// collect responses and pick the "majority"
	var max_count uint64 = 0
	var max_count_hash [32]byte = [32]byte{}
	for hash, count := range blockCount {
		if count > max_count {
			max_count_hash = hash
			max_count = count
		}
	}
	if max_count_hash == [32]byte{} {
		return false, Block{}
	}
	return true, blockMap[max_count_hash]
}

// compute the proof of work and then add the block to our queue
func (ns *NodeServer) WorkOnBlock(pBlock Block) error {
	for {
		if CurrentBlock.Transactions != nil {
			b := CurrentBlock
			clearCurrentBlock()
			b.SeqNum = pBlock.SeqNum + 1
			b.SetProofOfWork(pBlock.Hash, ns.workerChannel)
			select {
			// we found a block in the channel, so continue/start over
			case pBlock := <-ns.workerChannel:
				continue
			default:
				ns.blkQueue.Push(b)
			}
		}
	}
	return nil
}
