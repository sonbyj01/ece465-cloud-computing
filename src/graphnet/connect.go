package graphnet

import (
	"bufio"
	"fmt"
	"net"
)

//var AllClients map[*Client] int
var AllClients = make(map[*Client]int)

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
	delete(AllClients, client)
	
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
