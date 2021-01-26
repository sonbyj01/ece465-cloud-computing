package graph

// CompleteGraph generates a complete graph with nodeCount nodes
func CompleteGraph(nodeCount int) Graph {
	graph := Graph{make([]Node, nodeCount)}

	for i := 0; i < nodeCount; i++ {
		for j := 0; j < i; j++ {
			graph.AddUndirectedEdge(i, j)
		}
	}

	return graph
}

// RingGraph generates a graph in which each node has exactly two neighbors
func RingGraph(nodeCount int) Graph {
	// TODO: implement this
	return Graph{}
}

// RandomGraph generates a graph with nodeCount nodes and average
// branching factor bFactor
func RandomGraph(nodeCount int, bFactor float32) Graph {
	// TODO: implement this
	return Graph{}
}
