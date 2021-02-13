package main

import (
	"encoding/binary"
	"flag"
	"graphalgo/color/distributed"
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

	// this stores algorithm state
	state := distributed.WorkerState{}

	// create node connection pool
	ncp := graphnet.NewNodeConnPool()

	// worker message handlers
	dispatchTab := make(map[byte]graphnet.Dispatch)
	dispatchTab[graphnet.MSG_VERTEX_INFO] = func(vertexInfo []byte) {
		logger.Printf("Indexes %d have been updated to %d.",
			binary.LittleEndian.Uint32(vertexInfo[4:]),
			binary.LittleEndian.Uint32(vertexInfo[:4]))
	}
	dispatchTab[graphnet.MSG_NODE_FINISHED] = func(nodeIndex []byte) {
		logger.Printf("Node %d has finished processing.\n",
			nodeIndex[0])
	}
	dispatchTab[graphnet.MSG_NODE_ROUND_FINISHED] = func(nodeIndex []byte) {
		logger.Printf("Node %d has finished a round.\n",
			nodeIndex[0])
	}
	dispatchTab[graphnet.MSG_NODE_INDEX_COUNT] = func(indexCount []byte) {
		logger.Printf("Node %d, %d total nodes.",
			indexCount[0], indexCount[1])
		state.NodeIndex = int(indexCount[0])
		state.NodeCount = int(indexCount[1])
	}
	dispatchTab[graphnet.MSG_NODE_ADDRESS] = func(ip []byte) {
		logger.Printf("Node %d has IP of %d.%d.%d.%d and port of %d",
			ip[0], ip[1], ip[2], ip[3], ip[4], ip[5:])
		ipv4 := net.IP(ip[1:5]).String()
		port := strconv.Itoa(int(binary.LittleEndian.Uint16(ip[5:])))
		conn, err := net.Dial("tcp",ipv4+":"+port)
		if err != nil {
			logger.Fatal(err)
		}
		node := graphnet.NewNodeConn(conn, logger, dispatchTab)
		ncp.AddUnregistered(node)
		node.Index = int(ip[0])
	}

	// receive incoming connections from lower-indexed nodes
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Fatal(err)
		}

		nodeConn := graphnet.NewNodeConn(conn, logger, dispatchTab)
		ncp.AddUnregistered(nodeConn)
	}

	// reorder nodes so that they're in the correct order
	ncp.Register()

	// start coloring
	//distributed.ColorDistributed()
}
