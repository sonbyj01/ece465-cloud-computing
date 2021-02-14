package distributed

import (
	"graph"
	"graphnet"
	"sync"
)

// WorkerState holds the algorithm state for a worker node
type WorkerState struct {
	Subgraph    *graph.Graph
	NodeIndex   int            // node index in NodeConnPool
	NodeCount   int            // total number of nodes (including server)
	VertexBegin int            // start of vertex range
	VertexEnd   int            // end of vertex range
	Stored      map[int]int    // received neighbor vertex values
	ColorWg     sync.WaitGroup // WaitGroup for speculative coloring
	DetectWg    sync.WaitGroup // WaitGroup for conflict detection
	ConnPool    graphnet.NodeConnPool
}

// NewWorkerState initializes a new WorkerState
func NewWorkerState() *WorkerState {
	ws := WorkerState{}

	// DetectWg is used as a lock to prevent ColorWg semaphore from being
	// decremented before it is set; default is "locked"
	ws.DetectWg.Add(1)

	return &ws
}
