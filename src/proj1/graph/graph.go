package graph

import (
	"fmt"
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

// New returns a new graph
// Since this is namespaced under the graph package, can be used from the
// outside as graph.New(...)
func New(nodeCount int) Graph {
	g := Graph{make([]Node, nodeCount)}

	for i := 0; i < nodeCount; i++ {
		g.Nodes[i].Index = i
	}

	return g
}

// NewParallel returns a new graph numbered in parallel
func NewParallel(nodeCount int, nThreads int) Graph {
	g := Graph{make([]Node, nodeCount)}

	nodesPerThread := nodeCount / nThreads
	if nodeCount % nThreads != 0 {
		nodesPerThread++
	}

	var wg sync.WaitGroup
	wg.Add(nThreads)
	for i := 0; i < nThreads; i++ {
		go func(start int) {
			defer wg.Done()
			for j := start; j < start + nodesPerThread && j < nodeCount; j++ {
				g.Nodes[j].Index = j
			}
		}(i * nodesPerThread)
	}
	wg.Wait()

	return g
}

// AddNode adds a node to a graph
func (g *Graph) AddNode(value int) {
	g.Nodes = append(g.Nodes, Node{
		value,
		len(g.Nodes),
		make([]*Node, 0),
		sync.Mutex{},
	})
}

// AddUndirectedEdge adds an undirected edge between two nodes in a graph
// Note that this doesn't check for duplicate edges
func (g *Graph) AddUndirectedEdge(n1, n2 int) {
	if n1 >= len(g.Nodes) || n2 >= len(g.Nodes) {
		panic("Invalid node indices")
	}

	g.Nodes[n1].Adj = append(g.Nodes[n1].Adj, &g.Nodes[n2])
	g.Nodes[n2].Adj = append(g.Nodes[n2].Adj, &g.Nodes[n1])
}

// AddUndirectedEdge adds an undirected edge to another node pointer
// Note that this doesn't check if the second node is within the same graph,
// and this doesn't check for duplicate edges
func (n1 *Node) AddUndirectedEdge(n2 *Node) {
	n2.Adj = append(n2.Adj, n1)
	n1.Adj = append(n1.Adj, n2)
}

// Print prints out a list of a graph's nodes and values, as well as their
// neighbors and values
func (g *Graph) Print() {
	for i := range g.Nodes {
		fmt.Printf("%d: %d\n", i, g.Nodes[i].Value)
		for _, neighbor := range g.Nodes[i].Adj {
			fmt.Printf("\t%d: %d\n", neighbor.Index, neighbor.Value)
		}
	}
}

// CheckValidColoring checks whether a graph is appropriately colored
func (g *Graph) CheckValidColoring() bool {
	for i := range g.Nodes {
		for _, neighbor := range g.Nodes[i].Adj {
			if g.Nodes[i].Value == neighbor.Value {
				return false
			}
		}
	}
	return true
}
