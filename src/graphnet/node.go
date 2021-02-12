// Package graphnet includes the network datastructures and utility functions
// for the multi-node algorithm
package graphnet

import (
	"bufio"
	"fmt"
	"net"
)

type Node struct {
	Outgoing 	chan string
	Reader		*bufio.Reader
	Writer		*bufio.Writer
	Conn 		net.Conn
	Connection 	*Node
}

// Read will continuously be looking for an input from the socket connection and print it out
func (node *Node) Read() {
	for {
		fmt.Println("Reading...")
		line, err := node.Reader.ReadString('\n')
		fmt.Println(line)
		if err == nil {
			if node.Connection != nil {
				node.Connection.Outgoing <- line
			}
		} else {
			break
		}
	}
	node.Conn.Close()
	if node.Connection != nil {
		node.Connection.Connection = nil
	}
	node = nil
}

// Write will continuously be looking for an input to print out to the socket connection
func (node *Node) Write() {
	for data := range node.Outgoing {
		node.Writer.WriteString(data)
		node.Writer.Flush()
	}
}

func (node *Node) Listen() {
	go node.Read()
	go node.Write()
}

func NewNode(connection net.Conn) *Node {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	node := &Node{
		Outgoing: 	make(chan string),
		Conn: 		connection,
		Reader: 	reader,
		Writer:		writer,
	}
	node.Listen()
	return node
}

func (node *Node) sendVertexMessage(msg VertexMessage) {
	// TODO: implement this
}

func (node *Node) getVertexChannel() chan VertexMessage {
	// TODO: implement this
	bufSize := 64
	return make(chan VertexMessage, bufSize)
}
