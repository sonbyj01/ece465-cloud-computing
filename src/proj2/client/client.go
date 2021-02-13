package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"graphalgo/color/distributed"
	"graphnet"
	"net"
	"proj2/common"
	"strconv"
	"sync"
)

// main is the driver to be built into the executable for the client
func main() {
	logger, logFile := common.CreateLogger("worker")
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

	// waitgroup to wait on getting node index
	var nodeIndexWg sync.WaitGroup

	// waitgroup for handshake completion: only frees once everything is set up
	var setupWg sync.WaitGroup

	// worker message handlers
	dispatchTab := make(map[byte]graphnet.Dispatch)
	dispatchTab[graphnet.MSG_VERTEX_INFO] = func(vertexInfo []byte,
		_ *graphnet.NodeConn) {

		logger.Printf("Indexes %d have been updated to %d.",
			binary.LittleEndian.Uint32(vertexInfo[4:]),
			binary.LittleEndian.Uint32(vertexInfo[:4]))
	}
	dispatchTab[graphnet.MSG_NODE_FINISHED] = func(nodeIndex []byte,
		_ *graphnet.NodeConn) {

		logger.Printf("Node %d has finished processing.\n",
			nodeIndex[0])
	}
	dispatchTab[graphnet.MSG_NODE_ROUND_FINISHED] = func(nodeIndex []byte,
		_ *graphnet.NodeConn) {

		logger.Printf("Node %d has finished a round.\n",
			nodeIndex[0])
	}

	// receive total number of nodes, begin listening for nodes to dial
	setupWg.Add(2)
	nodeIndexWg.Add(1)
	dispatchTab[graphnet.MSG_NODE_INDEX_COUNT] = func(indexCount []byte,
		nodeConn *graphnet.NodeConn) {

		defer setupWg.Done()
		defer nodeIndexWg.Done()

		// update logger prefix
		logger.SetPrefix(fmt.Sprintf("worker %d: ", indexCount[0]))

		logger.Printf("Got node index %d, %d total nodes.",
			indexCount[0], indexCount[1])
		state.NodeIndex = int(indexCount[0])
		state.NodeCount = int(indexCount[1])

		// receive incoming connections from lower-indexed nodes
		// add two items to the waitgroup per node: one for the initial
		// connection, one for the extra message indicating which node it is
		setupWg.Add(2 * (state.NodeIndex-1))

		// set index of server to 0
		nodeConn.Index = 0
	}

	// receive address of higher-indexed node, dial it
	dispatchTab[graphnet.MSG_NODE_ADDRESS] = func(ip []byte,
		_ *graphnet.NodeConn) {

		logger.Printf("Node %d has IP of %d.%d.%d.%d and port of %d",
			ip[0], ip[1], ip[2], ip[3], ip[4], ip[5:])
		ipv4 := net.IP(ip[1:5]).String()
		port := strconv.Itoa(int(binary.LittleEndian.Uint16(ip[5:])))
		conn, err := net.Dial("tcp", ipv4+":"+port)
		if err != nil {
			logger.Fatal(err)
		}
		nodeConn := graphnet.NewNodeConn(conn, logger, dispatchTab)
		ncp.AddUnregistered(nodeConn)
		nodeConn.Index = int(ip[0])

		// send current node index to dialee, but must make sure current
		// node has been notified of index first
		nodeIndexWg.Wait()
		buf := make([]byte, 1)
		buf[0] = byte(state.NodeIndex)
		nodeConn.WriteBytes(graphnet.MSG_DIALER_INDEX, buf)
	}

	// receive node index of dialee
	dispatchTab[graphnet.MSG_DIALER_INDEX] = func(nodeIndex []byte,
		nodeConn *graphnet.NodeConn) {

		defer setupWg.Done()
		logger.Printf("Received dial from %d\n", nodeIndex[0])
		nodeConn.Index = int(nodeIndex[0])
	}

	// begin listening
	for i := 0; state.NodeIndex == 0 || i < state.NodeIndex; i++ {
		conn, err := listener.Accept()
		if err != nil {
			logger.Fatal(err)
		}

		nodeConn := graphnet.NewNodeConn(conn, logger, dispatchTab)
		ncp.AddUnregistered(nodeConn)
		setupWg.Done()
	}

	// wait for all handshake actions to complete
	setupWg.Wait()
	logger.Printf("Handshake complete\n")

	// reorder nodes so that they're in the correct order
	ncp.Register()
	if ncp.Index != state.NodeIndex {
		// just an extra check: these two values should be redundant
		logger.Fatal("ncpIndex and state.NodeIndex should match")
	}

	// start coloring
	//distributed.ColorDistributed()
}
