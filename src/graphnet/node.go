// Package graphnet includes the network datastructures and utility functions
// for the multi-node algorithm

// Server

package graphnet

import "net"

type Node struct {
	ip     net.IP
	port   int
}

func (node *Node) InitializeNode(portP int) {
	// Sets the default port value to 8000
	if portP == 0 && node.port == 0 {
		node.port = 8000
	}

	// Loops through all the available interfaces on the machine 
	// And assigns the IP address and default port to the Node Struct
	ifaces, _ := net.Interfaces()

	for _, i := range ifaces {
		addrs, _ := i.Addrs()

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				node.ip = v.IP
			case *net.IPAddr:
				node.ip = v.IP
			}
		}
	}
}