package graph

import (
	"math/rand"
	"sync"
	"time"
)

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

// NewRandomGraphParallel generates a random graph in parallel
func NewRandomGraphParallel(nodeCount int, bFactor float32,
	nThreads int) Graph {

	g := NewParallel(nodeCount, nThreads)

	nodesPerThread := nodeCount / nThreads
	if nodeCount % nThreads != 0 {
		nodesPerThread++
	}

	pEdge := bFactor / float32(nodeCount-1)

	var wg sync.WaitGroup
	wg.Add(nThreads)

	threadFunc := func(start int) {
		defer wg.Done()
		for i := start; i < start + nodesPerThread && i < nodeCount; i++ {
			source := rand.NewSource(time.Now().UnixNano())
			generator := rand.New(source)

			for j := 0; j < i; j++ {
				if generator.Float32() < pEdge {
					g.Nodes[i].Mutex.Lock()
					g.Nodes[j].Mutex.Lock()
					g.AddUndirectedEdge(i, j)
					g.Nodes[j].Mutex.Unlock()
					g.Nodes[i].Mutex.Unlock()
				}
			}
		}
	}

	for i := 0; i < nThreads; i++ {
		go threadFunc(i * nodesPerThread)
	}
	wg.Wait()

	return g
}
