package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"graph"
	"graphalgo/color/distributed"
	"graphnet"
	"net"
	"proj2/common"
	"runtime"
	"strconv"
	"sync"
)

// main is the driver to be built into the executable for the client
func main() {
	logger, logFile := common.CreateLogger("worker ", "color")
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
	ws := distributed.NewWorkerState()

	// waitgroup to wait on getting node index
	var nodeIndexWg sync.WaitGroup

	// waitgroup for handshake completion: only frees once all handshake
	// actions are set up (handshake doesn't include sending subgraph)
	var setupWg sync.WaitGroup

	// waitgroup to begin coloring; needs to be unlocked by start signal
	// being sent by server and subgraph being totally received and processed
	var startColoringWg sync.WaitGroup
	startColoringWg.Add(2)

	// buf for messages
	buf := make([]byte, 8)

	// worker message handlers
	dispatchTab := make(map[byte]graphnet.Dispatch)
	dispatchTab[graphnet.MSG_VERTEX_INFO] = func(vertexInfo []byte,
		_ *graphnet.NodeConn) {

		color := int(binary.LittleEndian.Uint32(vertexInfo[:4]))
		index := int(binary.LittleEndian.Uint32(vertexInfo[4:]))

		// TODO: probably want to remove this
		logger.Printf("Indexes %d have been updated to %d.",
			index, color)

		ws.Stored[index] = color
	}
	dispatchTab[graphnet.MSG_NODE_FINISHED] = func(nodeIndex []byte,
		_ *graphnet.NodeConn) {

		logger.Printf("Node %d has finished processing.\n",
			nodeIndex[0])

		// decrease total number of nodes remaining in the pool;
		// the timing of this (waiting for the DetectWg lock) means that it
		// has already incremented the ColorWg semaphore
		ws.DetectWg.Wait()
		ws.ColorWg.Done()
		ws.NodeCount--
	}
	dispatchTab[graphnet.MSG_NODE_ROUND_FINISHED] = func(nodeIndex []byte,
		_ *graphnet.NodeConn) {

		logger.Printf("Node %d has finished a round.\n",
			nodeIndex[0])

		// decrease the number of nodes we are waiting for
		ws.DetectWg.Wait()
		ws.ColorWg.Done()
	}

	// receive total number of nodes, begin listening for nodes to dial
	setupWg.Add(2)
	nodeIndexWg.Add(1)
	dispatchTab[graphnet.MSG_NODE_INDEX_COUNT] = func(indexCount []byte,
		nodeConn *graphnet.NodeConn) {

		defer setupWg.Done()
		defer nodeIndexWg.Done()

		// update logger prefix
		logger.SetPrefix(fmt.Sprintf("worker%d:\t", indexCount[0]))

		logger.Printf("Got node index %d, %d total nodes.",
			indexCount[0], indexCount[1])
		ws.NodeIndex = int(indexCount[0])
		ws.NodeCount = int(indexCount[1])

		// add number of tasks: expect two connections from each lower node
		// (nodeIndex-1 times) as well as one connection to each higher node
		// (nodeCount-nodeIndex-1 times)
		setupWg.Add(2*(ws.NodeIndex-1) +
			ws.NodeCount - ws.NodeIndex - 1)

		// set index of server to 0
		nodeConn.Index = 0
	}

	// receive address of higher-indexed node, dial it
	dispatchTab[graphnet.MSG_NODE_ADDRESS] = func(ip []byte,
		_ *graphnet.NodeConn) {

		defer setupWg.Done()

		ipv4 := net.IP(ip[1:5]).String()
		port := strconv.Itoa(int(binary.LittleEndian.Uint16(ip[5:])))

		logger.Printf("Node %d has address of %s:%s\n",
			ip[0], ipv4, port)

		conn, err := net.Dial("tcp", ipv4+":"+port)
		if err != nil {
			logger.Fatal(err)
		}
		nodeConn := graphnet.NewNodeConn(conn, logger, dispatchTab)
		ws.ConnPool.AddUnregistered(nodeConn)
		nodeConn.Index = int(ip[0])

		// send current node index to dialee, but must make sure current
		// node has been notified of index first
		nodeIndexWg.Wait()
		buf := make([]byte, 1)
		buf[0] = byte(ws.NodeIndex)
		nodeConn.WriteBytes(graphnet.MSG_DIALER_INDEX, buf, false)
	}

	// receive node index of dialee
	dispatchTab[graphnet.MSG_DIALER_INDEX] = func(nodeIndex []byte,
		nodeConn *graphnet.NodeConn) {

		defer setupWg.Done()
		logger.Printf("Received dial from %d\n", nodeIndex[0])
		nodeConn.Index = int(nodeIndex[0])
	}

	// receive subgraph
	dispatchTab[graphnet.MSG_SUBGRAPH] = func(buf []byte,
		_ *graphnet.NodeConn) {

		defer startColoringWg.Done()
		logger.Printf("Receiving subgraph...\n")
		ws.Subgraph, err = graph.Load(bytes.NewReader(buf))
		if err != nil {
			logger.Fatal(err)
		}

		// calculate start, end vertices
		nodeIndexWg.Wait()
		ws.VertexBegin = (ws.NodeIndex - 1) * len(ws.Subgraph.Vertices)
		ws.VertexEnd = ws.VertexBegin + len(ws.Subgraph.Vertices)
		logger.Printf("Finished receiving subgraph (vertices %d-%d).\n",
			ws.VertexBegin, ws.VertexEnd-1)
	}

	// start coloring
	dispatchTab[graphnet.MSG_BEGIN_COLORING] = func(_ []byte,
		_ *graphnet.NodeConn) {

		logger.Println("Received signal to begin coloring.")
		startColoringWg.Done()
	}

	// begin listening; expecting NodeIndex incoming dials (one from server,
	// NodeIndex-1 from lower-indexed workers)
	for i := 0; ws.NodeIndex == 0 || i < ws.NodeIndex; i++ {
		conn, err := listener.Accept()
		if err != nil {
			logger.Fatal(err)
		}
		logger.Printf("Accepted incoming connection %s<-%s\n",
			conn.LocalAddr().String(), conn.RemoteAddr().String())

		nodeConn := graphnet.NewNodeConn(conn, logger, dispatchTab)
		ws.ConnPool.AddUnregistered(nodeConn)
		setupWg.Done()

		// first connection should be server; wait for node to receive its index
		nodeIndexWg.Wait()
	}

	// wait for all handshake tasks to complete, send ack to server
	logger.Println("Waiting on setupWg...")
	setupWg.Wait()
	logger.Printf("Handshake complete\n")
	buf[0] = byte(ws.NodeIndex)
	ws.ConnPool.Conns[0].WriteBytes(graphnet.MSG_HANDSHAKE_DONE, buf[:1],
		false)

	// reorder nodes so that they're in the correct order
	ws.ConnPool.Register()
	if ws.ConnPool.Index != ws.NodeIndex {
		// just an extra check: these two values should be redundant
		logger.Fatalf("Connection pool index and state NodeIndex "+
			"should match. %d != %d\n", ws.ConnPool.Index, ws.NodeIndex)
	}

	// begin coloring
	logger.Println("Waiting on startColoringWg...")
	startColoringWg.Wait()
	logger.Printf("Beginning coloring...\n")
	distributed.ColorDistributed(ws, 10, runtime.NumCPU()*2, logger)

	logger.Printf("Done.")
}
