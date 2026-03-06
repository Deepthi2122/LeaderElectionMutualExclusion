package main

import (
	"fmt"
	"strconv"
)

func IntToString(i int) string {

	return strconv.Itoa(i)
}

func PrintStatus(n *Node) {

	n.mu.RLock()
	id := n.ID
	leader := n.LeaderID
	state := n.State
	clock := n.LogicalClock
	files := n.Files
	n.mu.RUnlock()

	fmt.Println("---- NODE STATUS ----")
	fmt.Println("Node ID:", id)
	fmt.Println("Leader:", leader)
	fmt.Println("State:", state)
	fmt.Println("Clock:", clock)
	fmt.Println("Files:", files)
	fmt.Println("---------------------")
}

func (n *Node) StartHeartbeat() {

	fmt.Println("Heartbeat running for Node", n.ID)
}