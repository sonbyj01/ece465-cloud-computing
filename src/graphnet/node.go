// Package graphnet includes the network datastructures and utility functions
// for the multi-node algorithm
package graphnet

import "net"

type Node struct {
	ip     net.IPAddr
	port   int

	// TODO: socket
}
