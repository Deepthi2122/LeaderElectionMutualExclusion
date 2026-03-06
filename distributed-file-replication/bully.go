package main

import "fmt"

func (n *Node) StartElection() {

	fmt.Println("Node", n.ID, "starting election")

	higherResponded := false

	for id := range n.Peers {

		if id <= n.ID {
			continue
		}

		fmt.Println("Sending election message to Node", id)

		var reply bool

		err := n.CallPeer(id, "RPCHandler.Election", n.ID, &reply)

		if err == nil && reply {

			higherResponded = true
			fmt.Println("Higher node", id, "is alive. Waiting for new leader.")
			return
		}
	}

	if !higherResponded {
		n.BecomeLeader()
	}
}

func (n *Node) BecomeLeader() {

	n.mu.Lock()
	n.LeaderID = n.ID
	n.mu.Unlock()

	fmt.Println("Node", n.ID, "became LEADER")

	for id := range n.Peers {

		if id == n.ID {
			continue
		}

		var ack bool

		n.CallPeer(id, "RPCHandler.Leader", n.ID, &ack)
	}
}

func (h *RPCHandler) Election(senderID int, reply *bool) error {

	fmt.Println("Received election message from Node", senderID)

	*reply = true

	go h.node.StartElection()

	return nil
}

func (h *RPCHandler) Leader(leaderID int, ack *bool) error {

	h.node.mu.Lock()
	h.node.LeaderID = leaderID
	h.node.mu.Unlock()

	fmt.Println("Node", leaderID, "is the new leader")

	*ack = true

	return nil
}