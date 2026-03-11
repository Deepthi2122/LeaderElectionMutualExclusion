package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"
)

type RPCHandler struct {
	node *Node
}

func (h *RPCHandler) Ping(senderID int, reply *bool) error {
	*reply = true
	return nil
}

func (n *Node) StartRPC() {

	handler := &RPCHandler{node: n}

	rpc.Register(handler)

	addr := ":" + IntToString(n.Port)

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("RPC server running on", addr)

	rpc.Accept(listener)
}

func (n *Node) CallPeer(peerID int, method string, args interface{}, reply interface{}) error {

	addr := n.Peers[peerID]

	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		return err
	}

	client := rpc.NewClient(conn)
	defer client.Close()

	call := client.Go(method, args, reply, nil)
	select {
	case <-call.Done:
		return call.Error
	case <-time.After(1 * time.Second):
		return fmt.Errorf("rpc call timeout to %s", addr)
	}
}
