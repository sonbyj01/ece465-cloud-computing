package distributed

import "graph"

type Subgraph struct {
	graph.Graph

	// start and end indices
	iBegin, iEnd	int
}