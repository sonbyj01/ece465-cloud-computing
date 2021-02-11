package main

import (
	"graphnet"
)

func main() {
	allClients = make(map[*Client]int)
	listener, _ := net.Listen("tcp", ":8080")

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err.Error())
		}

		client := NewClient(conn)

		for clientList, _ := range allClients {
			if clientList.connection == nil {
				client.connection = clientList
				clientList.connection = client
				fmt.Println("Connected")
			}
		}
		allClients[client] = 1
		fmt.Println(len(allClients))
	}
}