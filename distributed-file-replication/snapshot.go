package main

import "fmt"

func (n *Node) StartSnapshot() {

	fmt.Println("Snapshot started by Node", n.ID)

	for id := range n.Peers {

		if id == n.ID {
			continue
		}

		fmt.Println("Sending MARKER to Node", id)

		var ack bool

		n.CallPeer(id, "RPCHandler.Marker", n.ID, &ack)
	}
}

func (h *RPCHandler) Marker(sender int, ack *bool) error {

	fmt.Println("Received snapshot marker")

	PrintStatus(h.node)

	*ack = true

	return nil
}