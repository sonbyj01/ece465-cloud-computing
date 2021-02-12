package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// main is the driver to be built into the executable for the server
func main() {
	// Takes in command line flag(s)
	// https://stackoverflow.com/questions/45117892/passing-cli-arguments-to-excutables-with-go-run
	// https://gobyexample.com/command-line-flags
	configFile := flag.String("config", "", "File name that contains the node configurations")
	flag.Parse()
	if *configFile == "" {
		panic("No configuration file.")
	}

	// read node configuration file
	// https://stackoverflow.com/questions/8757389/reading-a-file-line-by-line-in-go/16615559#16615559
	file, err := os.Open(*configFile)
	check(err)
	fileScanner := bufio.NewScanner(file)
	addresses := make([]string, 0)
	for fileScanner.Scan() {
		fmt.Println("Config: ", fileScanner.Text())
		addresses = append(addresses, fileScanner.Text())
	}

	// establish a connection with each node from configuration file
	// https://dev.to/alicewilliamstech/getting-started-with-sockets-in-golang-2j66
	for _, address := range addresses {
		fmt.Println("Establishing connection with ", address, " ...")
		conn, err := net.Dial("tcp", address)
		if err != nil {
			panic(err)
		}

		reader := bufio.NewReader(os.Stdin)

		for {
			fmt.Print("Text to send: ")
			input, _ := reader.ReadString('\n')
			fmt.Println(input)
			message, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println("Server relay: ", message)
		}
	}

	// start coloring
	//distributed.ColorDistributedServer()
}