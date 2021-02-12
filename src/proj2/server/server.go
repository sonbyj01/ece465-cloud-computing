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
	// https://stackoverflow.com/questions/45117892
	// https://gobyexample.com/command-line-flags
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

	// establish a connection with each node from configuration file
	// https://dev.to/alicewilliamstech/getting-started-with-sockets-in-golang-2j66
	for i, address := range addresses {
		logger.Printf("Establishing connection with %s...\n", address)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Printf("Connection established with %s.\n", address)

		client := graphnet.NewClient(conn)

		// send information about other nodes to this node
		for j := range addresses {
			if i == j {
				continue
			}


		}

		//reader := bufio.NewReader(os.Stdin)
		//
		//for {
		//	fmt.Print("Text to send: ")
		//	input, _ := reader.ReadString('\n')	// request
		//	fmt.Fprintf(conn, input)
		//	fmt.Println(input)
		//	//message, _ := bufio.NewReader(conn).ReadString('\n') // response
		//	//fmt.Println("Server relay: ", message)
		//}
	}

	// start coloring
	//distributed.ColorDistributedServer()
}