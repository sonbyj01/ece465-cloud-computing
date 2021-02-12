package main

import (
	"bufio"
	"flag"
	"fmt"
	"graphnet"
	"net"
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
	networkInterface := flag.String("intf", "", "File name that contains the node configurations")
	port := flag.Int("port", 0, "Listening port number")
	flag.Parse()
	if *networkInterface == "" {
		panic("No interface specified")
	}
	if *port == 0 {
		panic("No port specified")
	}

	fmt.Println("Listening on port ", *port, " ...")
	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(*port))
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	allNodes := make(map[*graphnet.Node]int)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		node := graphnet.NewNode(conn)
		for nodeList, _ := range allNodes {
			if nodeList.Connection == nil {
				node.Connection = nodeList
				nodeList.Connection = node
				fmt.Println("Connected")
			}
		}
		allNodes[node] = 1
	}

	// start coloring
	//distributed.ColorDistributed()
}