package distributed

import (
	"graph"
	"graphnet"
)

// WorkerState holds the algorithm state for a worker node
type WorkerState struct {
	Subgraph    graph.Graph
	NodeIndex   int         // node index in NodeConnPool
	NodeCount   int         // total number of nodes (including server)
	VertexBegin int         // start of vertex range
	VertexEnd   int         // end of vertex range
	stored      map[int]int // stored neighbor vertex values
	ConnPool    graphnet.NodeConnPool
}
