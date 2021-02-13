package main

import (
	"flag"
	"graphnet"
	"net"
	"proj2/common"
	"strconv"
)

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

	// create node connection pool
	ncp := graphnet.NewNodeConnPool()

	// TODO: create dispatch table
	var dispatchTab map[byte]graphnet.Dispatch

	// receive incoming connections from lower-indexed nodes
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Fatal(err)
		}

		nodeConn := graphnet.NewNodeConn(conn, logger, dispatchTab)
		ncp.AddUnregistered(nodeConn)

		for test := range *nodeConn.Channel() {
			logger.Printf("Received %s\n", test)
		}
	}

	// TODO: dial connections to higher-indexed nodes

	// reorder nodes so that they're in the correct order
	ncp.Register()

	// start coloring
	//distributed.ColorDistributed()
}
