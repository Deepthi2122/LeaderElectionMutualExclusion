package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

func parsePeers(s string) map[int]string {
	peers := make(map[int]string)
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		id, err := strconv.Atoi(strings.TrimSpace(kv[0]))
		if err != nil {
			continue
		}
		peers[id] = strings.TrimSpace(kv[1])
	}
	return peers
}

func main() {

	id := flag.Int("id", 1, "Node ID")
	port := flag.Int("port", 8001, "Port")
	peersFlag := flag.String("peers", "", "Peer addresses as id=host:port,id=host:port,...")

	flag.Parse()

	peers := parsePeers(*peersFlag)

	node := NewNode(*id, *port, peers)

	node.InitStorage()

	fmt.Println("Starting Node", node.ID)

	go node.StartRPC()

	go node.StartHeartbeat()

	go node.StartElection()

	go StartCommandListener(node)

	select {}
}
