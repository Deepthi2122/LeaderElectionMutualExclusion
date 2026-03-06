package main

import "fmt"

func (n *Node) ReplicateFile(name string, content string) {

	fmt.Println("Starting replication for file", name)

	acks := 1

	for id := range n.Peers {

		if id == n.ID {
			continue
		}

		var reply bool

		err := n.CallPeer(id, "RPCHandler.Replicate", name, &reply)

		if err == nil && reply {
			acks++
		}
	}

	if acks >= (len(n.Peers)/2)+1 {

		n.WriteFile(name, content)

		fmt.Println("Replication committed")

	} else {

		fmt.Println("Replication failed")
	}
}

func (h *RPCHandler) Replicate(file string, reply *bool) error {

	h.node.WriteFile(file, "replicated")

	*reply = true

	return nil
}