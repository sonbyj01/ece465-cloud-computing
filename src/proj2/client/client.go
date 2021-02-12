package main

import (
	"bufio"
	"flag"
	"fmt"
	"graphnet"
	"os"
	"sync"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// main is the driver to be built into the executable for the client
func main() {
	configFile := flag.String("config", "", "File name that contains the node configurations")
	port := flag.Int("port", 0, "Listening port number")
	flag.Parse()
	if *port == 0 {
		panic("No port specified")
	}
	if *configFile == "" {
		panic("No configuration file.")
	}

	file, err := os.Open(*configFile)
	check(err)
	fileScanner := bufio.NewScanner(file)
	addresses := make([]string, 0)
	for fileScanner.Scan() {
		fmt.Println("Config: ", fileScanner.Text())
		addresses = append(addresses, fileScanner.Text())
	}

	var wg sync.WaitGroup
	fmt.Println("Listening")
	wg.Add(1)
	go graphnet.ListenConnections(port, &wg)

	fmt.Println("Establishing")
	wg.Add(1)
	go graphnet.EstablishConnections(addresses, *port, &wg)

	wg.Wait()

	// start coloring
	//distributed.ColorDistributed()
}