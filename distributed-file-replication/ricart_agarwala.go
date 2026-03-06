package main

import "fmt"

type RequestMessage struct {
	NodeID    int
	Timestamp int
}

// RPCHandler.RequestCS is called by a remote node that wants the critical section.
// We reply immediately unless we have higher priority, in which case we defer.
func (h *RPCHandler) RequestCS(req RequestMessage, reply *bool) error {

	n := h.node

	n.mu.Lock()

	sendNow := false

	if n.State == "RELEASED" {
		sendNow = true
	} else if n.State == "WANTED" {
		// Lower timestamp wins; tie-break by lower node ID
		if req.Timestamp < n.RequestTimestamp ||
			(req.Timestamp == n.RequestTimestamp && req.NodeID < n.ID) {
			sendNow = true
		} else {
			n.DeferredRequests = append(n.DeferredRequests, req.NodeID)
			fmt.Println("Node", n.ID, "deferred reply to Node", req.NodeID)
		}
	} else {
		// HELD — always defer
		n.DeferredRequests = append(n.DeferredRequests, req.NodeID)
		fmt.Println("Node", n.ID, "deferred reply to Node", req.NodeID)
	}

	n.mu.Unlock()

	if sendNow {
		fmt.Println("Node", n.ID, "granted reply to Node", req.NodeID)
	}

	*reply = sendNow
	return nil
}

// RPCHandler.GrantCS is called by a peer that previously deferred our CS request
// and is now releasing the critical section.
func (h *RPCHandler) GrantCS(senderID int, reply *bool) error {

	n := h.node

	n.mu.Lock()
	n.ReplyCount++
	needed := len(n.Peers) - 1
	ready := n.ReplyCount >= needed && n.State == "WANTED"
	ch := n.csReadyCh
	n.mu.Unlock()

	fmt.Println("Node", n.ID, "received deferred grant from Node", senderID)

	if ready && ch != nil {
		select {
		case ch <- struct{}{}:
		default:
		}
	}

	*reply = true
	return nil
}

func (n *Node) RequestCS() {

	n.mu.Lock()
	n.LogicalClock++
	n.State = "WANTED"
	n.RequestTimestamp = n.LogicalClock
	n.ReplyCount = 0
	n.csReadyCh = make(chan struct{}, 1)
	n.mu.Unlock()

	fmt.Println("Node", n.ID, "requesting critical section")

	needed := len(n.Peers) - 1

	for id := range n.Peers {

		if id == n.ID {
			continue
		}

		fmt.Println("Sending REQUEST to Node", id)

		var reply bool
		req := RequestMessage{NodeID: n.ID, Timestamp: n.RequestTimestamp}

		err := n.CallPeer(id, "RPCHandler.RequestCS", req, &reply)

		if err != nil {
			// Peer is down — treat as implicit grant
			fmt.Println("Node", id, "unreachable, counting as grant")
			n.mu.Lock()
			n.ReplyCount++
			n.mu.Unlock()
		} else if reply {
			fmt.Println("Reply granted by Node", id)
			n.mu.Lock()
			n.ReplyCount++
			n.mu.Unlock()
		} else {
			fmt.Println("Reply deferred by Node", id)
		}
	}

	// Block until all deferred replies arrive via RPCHandler.GrantCS
	n.mu.Lock()
	alreadyReady := n.ReplyCount >= needed
	n.mu.Unlock()

	if !alreadyReady {
		fmt.Println("Node", n.ID, "waiting for deferred replies...")
		<-n.csReadyCh
	}

	n.mu.Lock()
	n.State = "HELD"
	n.mu.Unlock()

	fmt.Println("Node", n.ID, "entered critical section")
}

func (n *Node) EnterCS() {

	n.mu.RLock()
	state := n.State
	n.mu.RUnlock()

	if state == "HELD" {
		fmt.Println("Node", n.ID, "is already in the critical section")
	} else {
		fmt.Println("Node", n.ID, "is not in the critical section (state:", state+")")
	}
}

func (n *Node) ExitCS() {

	n.mu.Lock()

	if n.State != "HELD" {
		n.mu.Unlock()
		fmt.Println("Cannot exit critical section: node is not inside it")
		return
	}

	n.State = "RELEASED"
	n.RequestTimestamp = 0
	n.ReplyCount = 0
	deferred := n.DeferredRequests
	n.DeferredRequests = nil

	n.mu.Unlock()

	fmt.Println("Node", n.ID, "exited critical section")

	// Send queued-up grants now that we are released
	for _, peerID := range deferred {
		var reply bool
		err := n.CallPeer(peerID, "RPCHandler.GrantCS", n.ID, &reply)
		if err != nil {
			fmt.Println("Failed to send grant to Node", peerID, ":", err)
		} else {
			fmt.Println("Node", n.ID, "sent grant to Node", peerID)
		}
	}
}