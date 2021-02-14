package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"graphnet"
	"math"
	"net"
	"os"
	"proj2/common"
	"strconv"
	"strings"
	"sync"
)

// main is the driver to be built into the executable for the server
func main() {
	// create logger
	logger, logFile := common.CreateLogger("server")
	defer func() {
		err := logFile.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// Takes in command line flag(s)
	configFile := flag.String("config", "",
		"Node addresses file")
	graphFile := flag.String("graph", "",
		"Graph description file")
	flag.Parse()
	if *configFile == "" {
		logger.Fatal("No configuration file.")
	}

	// read node configuration file
	// https://stackoverflow.com/questions/8757389
	file, err := os.Open(*configFile)
	if err != nil {
		logger.Fatal(err)
	}
	fileScanner := bufio.NewScanner(file)
	addresses := make([]string, 0)
	for fileScanner.Scan() {
		logger.Printf("Reading config: %s\n", fileScanner.Text())
		addresses = append(addresses, fileScanner.Text())
	}
	nWorkers := len(addresses)

	// create node connection pool
	ncp := graphnet.NewNodeConnPool()

	// create server dispatch table
	dispatchTab := make(map[byte]graphnet.Dispatch)

	// buffer for sending
	buf := make([]byte, 8)

	// wait for all handshakes to finish; workers will send ack upon finishing
	// handshake
	var handshakeWg sync.WaitGroup
	handshakeWg.Add(nWorkers)
	dispatchTab[graphnet.MSG_HANDSHAKE_DONE] = func(buf []byte,
		_ *graphnet.NodeConn) {

		handshakeWg.Done()
		logger.Printf("Node %d has completed handshake.\n", buf[0])
	}

	// handler for MSG_NODE_FINISHED: when all nodes finished, finish
	var wg sync.WaitGroup
	wg.Add(nWorkers)
	dispatchTab[graphnet.MSG_NODE_FINISHED] = func(nodeIndex []byte,
		_ *graphnet.NodeConn) {

		wg.Done()
		logger.Printf("Node %d has finished processing.\n",
			nodeIndex[0])
	}

	// establish a connection with each node from configuration file
	for i, address := range addresses {
		logger.Printf("Establishing connection with %s (node %d)...\n",
			address, i+1)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Printf("Connection established with %s (node %d).\n",
			address, i+1)
		nodeConn := graphnet.NewNodeConn(conn, logger, dispatchTab)
		nodeConn.Index = i + 1
		ncp.AddUnregistered(nodeConn)
	}

	// send information about all nodes to each node
	for i, nodeConn := range ncp.Conns {
		logger.Printf("Sending handshake to node %d\n", i+1)

		// send node index and count to worker
		buf[0] = byte(i + 1)
		buf[1] = byte(nWorkers + 1)
		nodeConn.WriteBytes(graphnet.MSG_NODE_INDEX_COUNT, buf[:2],
			false)

		// send addresses of higher indexed nodes to node
		for j := i + 1; j < len(addresses); j++ {
			ipComponents := strings.Split(addresses[j], ":")
			buf[0] = byte(j + 1)

			copy(buf[1:5], net.ParseIP(ipComponents[0]))

			port, err := strconv.Atoi(ipComponents[1])
			if err != nil {
				logger.Fatal(err)
			}
			binary.LittleEndian.PutUint16(buf[5:7], uint16(port))
			nodeConn.WriteBytes(graphnet.MSG_NODE_ADDRESS, buf[:7],
				false)
		}
	}

	// make sure that nodes are in order; shouldn't really have an effect
	// for the server
	ncp.Register()

	// divide graph into equally-sized subgraphs (padding with empty nodes
	// as necessary) and streaming
	// TODO: should probably move this to its own function
	file, err = os.Open(*graphFile)
	if err != nil {
		logger.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	nVertices, err := strconv.Atoi(scanner.Text())
	if err != nil {
		logger.Fatal(err)
	}
	verticesPerNode := int(math.Ceil(float64(nVertices) / float64(nWorkers)))
	for i, nodeConn := range ncp.Conns {
		if i == 0 {
			continue
		}

		// begin writing file
		nodeConn.WriteBytes(graphnet.MSG_SUBGRAPH, buf[:0], true)

		nodeConn.WriteBytes(graphnet.MSG_CONT,
			[]byte(fmt.Sprintf("%d", verticesPerNode)), true)

		// send vertices, pad with disconnected vertices if necessary
		for j := 0; j < verticesPerNode; j++ {
			nodeConn.WriteBytes(graphnet.MSG_CONT, []byte("\n"), true)
			if scanner.Scan() {
				nodeConn.WriteBytes(graphnet.MSG_CONT, scanner.Bytes(),
					true)
			} else {
				nodeConn.WriteBytes(graphnet.MSG_CONT, []byte("0;"),
					true)
			}
		}

		// end write file
		nodeConn.WriteBytes(graphnet.DELIM_EOF, buf[:0], false)
	}

	// wait for all nodes to finish handshake
	handshakeWg.Wait()
	logger.Println("All nodes have completed handshake.")

	// send signal to start coloring
	ncp.Broadcast(graphnet.MSG_BEGIN_COLORING, buf[:0])

	// wait until all nodes finished coloring; this will activate when nWorkers
	// MSG_NODE_FINISHED are received
	wg.Wait()

	// TODO: collect subgraphs and verify coloring

	logger.Printf("Done.")
}
