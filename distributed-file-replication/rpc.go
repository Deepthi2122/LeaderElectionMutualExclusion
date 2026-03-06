package main

import (
	"log"
	"net"
	"net/rpc"
)

type RPCHandler struct {
	node *Node
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

	client, err := rpc.Dial("tcp", addr)

	if err != nil {
		return err
	}

	defer client.Close()

	return client.Call(method, args, reply)
}