package graphnet

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"sync"
)

// Node Struct
type Node struct {
	Reader		*bufio.Reader
	Writer		*bufio.Writer
	Conn		net.Conn
}

// checkError - minimal error checking, essentially just panic
func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

// Read - processes messages from sockets and decodes gob into VertexMessage struct
func (node *Node) Read() {
	temp := make([]byte, 500)

	for {
		_, err := node.Conn.Read(temp)
		checkError(err)

		tempBuff := bytes.NewBuffer(temp)
		tempStruct := new(VertexMessage)

		dec := gob.NewDecoder(tempBuff)
		dec.Decode(tempStruct)
		fmt.Println(tempStruct)
	}
}

// SendVertexMessage - sends type VertexMessage through socket using gob
func (node *Node) SendVertexMessage(msg VertexMessage) {
	tempBuff := new(bytes.Buffer)
	enc := gob.NewEncoder(tempBuff)
	enc.Encode(msg)
	node.Conn.Write(tempBuff.Bytes())
}

// Listen - go routine to continuously listen for messages from socket
func (node *Node) Listen() {
	go node.Read()
}

// ListenConnections - listens for dial connections via port and then appends to an array of nodes
func ListenConnections(port *int, wg *sync.WaitGroup) {
	fmt.Printf("Listening on port %d...", *port)
	allNodes := make(map[*Node]int)
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(*port))
	checkError(err)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		checkError(err)
		node := NewNode(conn)
		allNodes[node] = 1
	}
}

// EstablishConnection - will create a node and attach a socket to 'address'
// 'address' - [IPv4]:[port]
func EstablishConnection(address string, wg *sync.WaitGroup) *Node {
	fmt.Printf("Establishing connection with %s...", address)
	conn, err := net.Dial("tcp", address)
	checkError(err)

	node := NewNode(conn)
	return node
}

// NewNode - creates a new node with bufio.reader/writer and socket connection attached
func NewNode(connection net.Conn) *Node {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	node := &Node {
		Conn:		connection,
		Reader: 	reader,
		Writer: 	writer,
	}

	node.Listen()
	return node
}