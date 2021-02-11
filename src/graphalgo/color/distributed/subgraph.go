package distributed

import (
	"graph"
	"graphnet"
)

type Subgraph struct {
	graph.Graph
	pos, iBegin, iEnd int 	// graph index/position, vertex start/end indices
	stored map[int]int		// stored neighbor vertex values
}

// sendToNodeCP sends a control message to node n
func (sg *Subgraph) sendToNodeCP(node *graphnet.Node,
	msg graphnet.VertexMessageType) {

	//graphnet.VertexMessage{msg}
}
