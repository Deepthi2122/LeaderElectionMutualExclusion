package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func StartCommandListener(n *Node) {

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print(">> ")

		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		switch cmd {

		case "election":
			n.StartElection()

		case "cs":
			go n.RequestCS()

		case "enter":
			n.EnterCS()

		case "exitcs":
			n.ExitCS()

		case "replicate":
			n.ReplicateFile("report.txt", "distributed file")

		case "snapshot":
			n.StartSnapshot()

		case "files":
			n.ListFiles()

		case "status":
			PrintStatus(n)

		default:
			fmt.Println("Unknown command")
		}
	}
}