// Package graph includes a graph data structure and some graph utilities,
// including graph generation and partitioning
package graph

import (
	"sync"
)

// Vertex represents a vertex (node) of a Graph object, with an adjacency list
// of indices of other vertices
type Vertex struct {
	Value int
	Adj   []int
	Mutex sync.Mutex
}

// Graph represents a very simple graph data structure
type Graph struct {
	Vertices []Vertex
}
