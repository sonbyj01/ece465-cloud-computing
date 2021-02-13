package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"graphnet"
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
		"File containing the node configurations")
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

	// create node connection pool
	ncp := graphnet.NewNodeConnPool()

	// create server dispatch table
	// TODO: fill with actual handlers
	dispatchTab := make(map[byte]graphnet.Dispatch)

	var wg sync.WaitGroup
	nWorkers := len(addresses)
	wg.Add(nWorkers)
	dispatchTab[graphnet.MSG_NODE_FINISHED] = func(nodeIndex []byte) {
		wg.Done()
		logger.Printf("Node %d has finished processing.\n",
			nodeIndex[0])
	}

	// establish a connection with each node from configuration file
	for i, address := range addresses {
		logger.Printf("Establishing connection with %s...\n", address)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Printf("Connection established with %s.\n", address)
		nodeConn := graphnet.NewNodeConn(conn, logger, dispatchTab)
		nodeConn.Index = i + 1
		ncp.AddUnregistered(nodeConn)
	}

	// send information about all nodes to each node
	buf := make([]byte, 7)
	for i, nodeConn := range ncp {
		// send node index and count to worker
		buf[0] = byte(i+1)
		buf[1] = byte(nWorkers+1)
		nodeConn.WriteBytes(graphnet.MSG_NODE_INDEX_COUNT, buf[:2])

		// send addresses of higher indexed nodes to node
		for j := i + 1; j < len(addresses); j++ {
			ipComponents := strings.Split(addresses[j], ":")
			buf[0] = byte(j+1)
			copy(buf[1:5], net.ParseIP(ipComponents[0]))

			port, err := strconv.Atoi(ipComponents[1])
			if err != nil {
				logger.Fatal(err)
			}
			binary.LittleEndian.PutUint16(buf[5:7], uint16(port))
			nodeConn.WriteBytes(graphnet.MSG_NODE_ADDRESS, buf[:7])
		}

		nodeConn.WriteBytes(graphnet.MSG_HANDSHAKE_DONE, buf[:0])
	}

	// this shouldn't have any effect for the server, since all nodes were
	// added in order
	ncp.Register()

	// TODO: send subgraphs

	// TODO: start coloring
	// distributed.ColorDistributedServer()
}
