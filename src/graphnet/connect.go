package graphnet

import (
	"bufio"
	"fmt"
	"net"
)

// ConnPool keeps track of all connections in an array; they are initially
// dumped in the unregistered list and moved to the correct index in the
// registered array when their index is received
type ConnPool struct {
	unregistered []*Client
	registered   []*Client
}

// Client is a struct to keep track of a single connection from this node
// to another node
type Client struct {
	outgoing   chan 	string
	reader     *bufio.Reader
	writer     *bufio.Writer
	conn       net.Conn
	Connection *Client
}

func (client *Client) Read() {
	for {
		line, err := client.reader.ReadString('\n')

		if err == nil {
			if client.Connection != nil {
				client.Connection.outgoing <- line
			}
			fmt.Println(err)
		} else {
			break 
		}
	}

	client.conn.Close()
	//delete(AllClients, client)
	
	if client.Connection != nil {
		client.Connection.Connection = nil
	}
	client = nil
}

func (client *Client) Write() {
	for data := range client.outgoing {
		client.writer.WriteString(data)
		client.writer.Flush()
	}
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

func NewClient(connection net.Conn) *Client {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &Client{
		outgoing: 	make(chan string), 
		conn:		connection,
		reader:		reader, 
		writer:		writer,
	}
	client.Listen()

	return client
}
