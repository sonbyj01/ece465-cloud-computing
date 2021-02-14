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
	StoredMutex sync.Mutex     // mutex for the above (TODO: make R/W lock?)
	StartWg     sync.WaitGroup // WaitGroup for starting the round
	ColorWg     sync.WaitGroup // WaitGroup for speculative coloring
	DetectWg    sync.WaitGroup // WaitGroup for conflict detection
	AlgoStarted bool
	ConnPool    graphnet.NodeConnPool
}

// NewWorkerState initializes a new WorkerState
func NewWorkerState() *WorkerState {
	ws := WorkerState{
		Stored: make(map[int]int),
	}

	return &ws
}
