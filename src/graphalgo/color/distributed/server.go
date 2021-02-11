package distributed

import (
	"graphnet"
	"io"
)

// InitConnections will reach out to all of the clients specified in configFile
func InitConnections(configFile io.Reader) []graphnet.Node {
	// TODO
	return make([]graphnet.Node, 0)
}

// SendConnections will send connection info about nodes to all of the nodes
// so they can establish connections with themselves
func SendConnections() {
	// TODO
}

// SendGraph will partition graph into Subgraph instances and send off to the
// nodes
func SendGraph() {
	// TODO
}

// ColorDistributedServer tells each client to start running ColorDistributed
// and listens for completion
func ColorDistributedServer() {
	// TODO
}