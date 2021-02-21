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
	ColorWgLock sync.Mutex     // to protect the consistency of colorWg
	ConnPool    graphnet.NodeConnPool
	State       AlgoState
}

// NewWorkerState initializes a new WorkerState
func NewWorkerState() *WorkerState {
	ws := WorkerState{
		Stored: make(map[int]int),
		State:  STATE_INIT,
	}

	return &ws
}

// AlgoState is used to determine the current state of the algorithm (e.g.,
// for heartbeat purposes and to have clean cleanup procedures)
type AlgoState int

const (
	// STATE_INIT means startup and/or handshake
	STATE_INIT AlgoState = iota

	// STATE_RUNNING means algo is running
	STATE_RUNNING AlgoState = iota

	// STATE_FINISHED means algo is done
	STATE_FINISHED AlgoState = iota
)
