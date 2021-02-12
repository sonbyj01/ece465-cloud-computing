// Package graphnet includes the network datastructures and utility functions
// for the multi-node algorithm
package graphnet

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

var allNodes = make(map[*Node]int)

func checkError(e error) {
	if e != nil {
		//panic(e)
		panic(e)
	}
}

type Node struct {
	Outgoing		chan 	string			// test/debug purpose
	OutgoingData	chan 	VertexData
	OutgoingMessage	chan 	VertexMessage
	Reader			*bufio.Reader
	Writer			*bufio.Writer
	Conn 			net.Conn
	Connection 		*Node
}

// Read will continuously be looking for an input from the socket connection and print it out
// https://dchua.com/2017/06/23/sending-your-structs-across-the-wire-(tcp-connection)/
func (node *Node) Read() {
	defer node.Conn.Close()
	tmp := make([]byte, 500)
	for {
		fmt.Println("Reading...")
		tmpbuff := bytes.NewBuffer(tmp)
		tmpstruct := new(VertexMessage)

		dec := gob.NewDecoder(tmpbuff)
		dec.Decode(tmpstruct)
		fmt.Println(tmpstruct)

		//var msg VertexMessage
		//dec := gob.NewDecoder(node.Conn)
		//dec.Decode(&msg)
		//fmt.Println(msg)

		//fmt.Println("Reading...")
		//line, err := node.Reader.ReadString('\n')
		//fmt.Println(line)
		//if err == nil {
		//	if node.Connection != nil {
		//		node.Connection.Outgoing <- line
		//	}
		//} else {
		//	break
		//}
	}
	node.Conn.Close()
	if node.Connection != nil {
		node.Connection.Connection = nil
	}
	node = nil
}

// Write will continuously be looking for an input to print out to the socket connection
func (node *Node) Write() {
	// --- Test ---
	vertexData := make([]VertexData, 10)

	msg := VertexMessage{
		Data:	vertexData,
	}

	node.sendVertexMessage(msg)
	// --- Test ---
	for {
		fmt.Println("Writing...")
		//reader := bufio.NewReader(os.Stdin)
		//fmt.Print("Text to send: ")
		//input, _ := reader.ReadString('\n')
		//node.Writer.WriteString(input)
		//node.Writer.Flush()

	}
	//for data := range node.Outgoing {
	//	node.Writer.WriteString(data)
	//	node.Writer.Flush()
	//}
}

func (node *Node) Listen() {
	go node.Read()
	go node.Write()
}

// Workers - listens for dial connection from server/other workers
func ListenConnections(port *int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Listening on port", *port, " ...")
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(*port))
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		node := NewNode(conn)
		
		for nodeList, _ := range allNodes {
			if nodeList.Connection == nil {
				node.Connection = nodeList
				nodeList.Connection = node
				fmt.Println("Connected")
			}
		}
		allNodes[node] = 1
	}
}

func EstablishConnections(addresses []string, port int, wg *sync.WaitGroup) {
	defer wg.Done()

	// last node in array, there are no more connections that need to be established
	if addresses != nil {fmt.Println("Last Node")}

	addrs, _ := net.InterfaceAddrs()
	//fmt.Println("Address:", strings.Split(addrs[1].String(), "/"))

	// establishes connection with node at 'address'
	for _, address := range addresses {
		if address == strings.Split(addrs[1].String(), "/")[0]+":"+strconv.Itoa(port) {
			continue
		}

		fmt.Println("Establishing connection with", address, "...")
		conn, err := net.Dial("tcp", address)
		checkError(err)
		//if err != nil {
		//	continue
		//}

		node := NewNode(conn)

		for {
			for nodeList, _ := range allNodes {
				if nodeList.Connection == nil {
					node.Connection = nodeList
					nodeList.Connection = node
					fmt.Println("Connected")
				}
			}
			allNodes[node] = 1
		}
	}
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
	defer node.Conn.Close()
	bin_buf := new(bytes.Buffer)
	enc := gob.NewEncoder(bin_buf)
	enc.Encode(msg)
	node.Conn.Write(bin_buf.Bytes())
}

func (node *Node) getVertexChannel() chan VertexMessage {
	// TODO: implement this
	bufSize := 64
	return make(chan VertexMessage, bufSize)
}
