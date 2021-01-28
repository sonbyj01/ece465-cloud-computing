package graph

import (
	"math/rand"
)

// NewCompleteGraph generates a complete graph with nodeCount nodes
func NewCompleteGraph(nodeCount int) Graph {
	g := New(nodeCount)

	for i := 1; i < nodeCount; i++ {
		for j := 0; j < i; j++ {
			g.AddUndirectedEdge(i, j)
		}
	}

	return g
}

// NewRingGraph generates a graph in which each node has exactly two neighbors
func NewRingGraph(nodeCount int) Graph {
	g := New(nodeCount)

	for i := 0; i < nodeCount-1; i++ {
		g.AddUndirectedEdge(i, i+1)
	}

	// if n=2, this edge is already created
	if nodeCount > 2 {
		g.AddUndirectedEdge(0, nodeCount-1)
	}

	return g
}

// NewRandomGraph generates a graph with nodeCount nodes and average
// branching factor bFactor
func NewRandomGraph(nodeCount int, bFactor float32) Graph {
	g := New(nodeCount)

	// if branching factor is bFactor, then given n1, n2 nodes in g, then the
	// probability of an undirected edge is bFactor / (nodeCount - 1)
	pEdge := bFactor / float32(nodeCount-1)

	for i := 1; i < nodeCount; i++ {
		for j := 0; j < i; j++ {
			if rand.Float32() < pEdge {
				g.AddUndirectedEdge(i, j)
			}
		}
	}

	return g
}
