// Package graph includes a graph data structure and some graph utilities,
// including graph generation and partitioning
package graph

import (
	"sync"
)

// Node represents a node of a Graph object, with an adjacency list
// of pointers to other nodes
type Node struct {
	Value int
	Index int
	Adj   []*Node
	Mutex sync.Mutex
}

// Graph represents a very simple graph data structure
type Graph struct {
	Nodes []Node
}
