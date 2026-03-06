package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const storageDir = "node_storage"

func (n *Node) InitStorage() {

	nodeDir := filepath.Join(storageDir, fmt.Sprintf("node%d", n.ID))

	os.MkdirAll(nodeDir, os.ModePerm)
}

func (n *Node) WriteFile(name string, content string) {

	nodeDir := filepath.Join(storageDir, fmt.Sprintf("node%d", n.ID))

	filePath := filepath.Join(nodeDir, name)

	version := n.Files[name] + 1

	os.WriteFile(filePath, []byte(content), 0644)

	n.Files[name] = version

	fmt.Println("Node", n.ID, "stored", name, "version", version)
}

func (n *Node) ListFiles() {

	for file, version := range n.Files {

		fmt.Println(file, "version", version)
	}
}