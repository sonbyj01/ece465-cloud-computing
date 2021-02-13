package main

import (
	"bufio"
	"flag"
	"graphnet"
	"net"
	"os"
	"proj2/common"
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

	// establish a connection with each node from configuration file
	for i, address := range addresses {
		logger.Printf("Establishing connection with %s...\n", address)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Printf("Connection established with %s.\n", address)
		nodeConn := graphnet.NewNodeConn(conn, logger)
		ncp.AddUnregistered(nodeConn)

		// send information about other nodes to this node
		for j, address2 := range addresses {
			if i == j {
				continue
			}

			*nodeConn.Channel() <- address2
		}
		*nodeConn.Channel() <- "Done\n"
	}

	// this shouldn't have any effect for the server, since all nodes were
	// added in order
	ncp.Register()

	// start coloring
	//distributed.ColorDistributedServer()
}