package main

import (
	"bufio"
	"flag"
	"fmt"
	"graphnet"
	"net"
	"proj2/common"
	"strconv"
)

// https://dev.to/alicewilliamstech/getting-started-with-sockets-in-golang-2j66
func handleConnection(conn net.Conn) {
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		fmt.Println("Client left.")
		conn.Close()
		return
	}

	fmt.Println("Client Message: ", string(buffer[:len(buffer)-1]))
	conn.Write(buffer)
	handleConnection(conn)
}

// main is the driver to be built into the executable for the client
func main() {
	logger, logFile := common.CreateLogger("client")
	defer func() {
		err := logFile.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// get port number to listen on
	port := flag.Int("port", 0, "Port to listen on")
	flag.Parse()
	if *port == 0 {
		panic("No port specified")
	}

	// begin listening for incoming connections
	logger.Printf("Listening on port %d...", *port)
	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(*port))
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		err := listener.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// initialize array of connections
	allNodes := make(map[*graphnet.Node]int)

	// listen for connections from others
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Fatal(err)
		}

		client := graphnet.NewClient(conn)
		//node := graphnet.NewNode(conn)
		//for nodeList, _ := range allNodes {
		//	if nodeList.Connection == nil {
		//		node.Connection = nodeList
		//		nodeList.Connection = node
		//		fmt.Println("Connected")
		//	}
		//}
		//allNodes[node] = 1
	}

	// start coloring
	//distributed.ColorDistributed()
}
