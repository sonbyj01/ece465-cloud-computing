package distributed

import "graph"

type Subgraph struct {
	graph.Graph
	pos, iBegin, iEnd int 	// graph index/position, vertex start/end indices
	stored map[int]int		// stored neighbor vertex values
}
