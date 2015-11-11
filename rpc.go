package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"syscall"
)

// RPCCall helper function:
// Does what the name says :)
// Parameters:
//   remote: address of the remote server
//   name:   full name of the method, like "ns.ReceiveBlock"
//   args:   struct containing arguments; must match on the caller-callee sides
//   reply:  *reference* to the struct to store response
// Return value:
//   true if successful, false if errors occurred
// Example:
//   // Arg and Reply are some custom struct types
//   args Arg{}
//   reply Reply{}
//   err := RPCCall("/tmp/slycoin-server1.sock", "ns.ReceiveBlock", args, &reply)
func RPCCall(remote string, name string, args interface{}, reply interface{}) bool {
	c, err := rpc.Dial("unix", remote)
	if err != nil {
		err1 := err.(*net.OpError)
		if err1.Err != syscall.ENOENT && err1.Err != syscall.ECONNREFUSED {
			fmt.Printf("RPC Dial() failed: %v\n", err1)
		}
		return false
	}
	defer c.Close()

	err = c.Call(name, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

// XXX assuming we will have a isdead() method to tell rpc when to stop
// using unix domain socket for RPC right now
// ** Call this function upon NodeServer initialization **
// XXX nodeserver isn't defined?
func (ns *NodeServer) StartRPCServer(addr string) bool {
	rpcs := rpc.NewServer()
	rpcs.Register(ns)

	os.Remove(addr)
	l, e := net.Listen("unix", addr)
	if e != nil {
		log.Fatal("listen error: ", e)
		return false
	}
	// ns.rpcListener : net.Listener
	ns.rpcListener = l
	// Start RPC listener thread
	go func(addr string) {
		for ns.isdead() == false {
			conn, err := ns.l.Accept()
			// XXX px isn't defined?
			if err == nil && px.isdead() == false {
				go rpcs.ServeConn(conn)
			} else if err == nil {
				conn.Close()
			}
			if err != nil && ns.isdead() == false {
				log.Printf("%v accept: %v\n", addr, err.Error())
			}
		}
	}(addr)

	return true
}
