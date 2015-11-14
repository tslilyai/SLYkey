package main

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
)

const (
	MAX_PEERS = 64
)

type NodeServer struct {
	qMu         sync.Mutex // mutex for BlockQueue
	mMu         sync.Mutex // mutex for block chain (map)
	dead        int32
	rpcListener net.Listener
	peers       []string
	blkQueue    *BlockQueue
}

// addr: local address?
// peers: bunch of other servers?
// XXX should call StartRPCServer, NewBlockQueue, etc.
func NewNodeServer(addr string, peers []string) *NodeServer {
	ns := &NodeServer{
		// initialize fields here
		blkQueue: NewBlockQueue(),
	}
	ns.StartRPCServer(addr)
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

// compute the proof of work and then broadcast the block
// we also update the database and block chain with the block transactions
// should we just have a loop that calls this?
// XXX drop block if we receive another block?
func (ns *NodeServer) WorkOnBlock(pBlock Block, c chan Block) error {
	for {
		if CurrentBlock.Transactions != nil {
			b := CurrentBlock
			clearCurrentBlock()
			b.SetProofOfWork(pBlock.Hash, c)
			select {
			// we found a block in the channel, so continue/start over
			case pBlock := <-c:
				continue
			default:
				ns.blkQueue.Push(b)
			}
		}
	}
	return nil
}
