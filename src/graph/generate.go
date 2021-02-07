package graph

import (
	"math/rand"
	"sync"
	"time"
)

// New returns a new graph
func New(nVertices int) Graph {
	return Graph{make([]Vertex, nVertices)}
}

//// NewParallel returns a new graph numbered in parallel
//func NewParallel(nVertices int, nThreads int) Graph {
//	g := Graph{make([]Node, nodeCount)}
//
//	nodesPerThread := nodeCount / nThreads
//	if nodeCount % nThreads != 0 {
//		nodesPerThread++
//	}
//
//	var wg sync.WaitGroup
//	wg.Add(nThreads)
//	for i := 0; i < nThreads; i++ {
//		go func(start int) {
//			defer wg.Done()
//			for j := start; j < start + nodesPerThread && j < nodeCount; j++ {
//				g.Nodes[j].Index = j
//			}
//		}(i * nodesPerThread)
//	}
//	wg.Wait()
//
//	return g
//}

// AddNode adds a node to a graph
func (g *Graph) AddNode(value int) {
	g.Vertices = append(g.Vertices, Vertex{
		value,
		make([]int, 0),
		sync.Mutex{},
	})
}

// AddUndirectedEdge adds an undirected edge between two nodes in a graph
// Note that this doesn't check for duplicate edges
func (g *Graph) AddUndirectedEdge(n1, n2 int) {
	if n1 >= len(g.Vertices) || n2 >= len(g.Vertices) {
		panic("Invalid node indices")
	}

	g.Vertices[n1].Adj = append(g.Vertices[n1].Adj, n2)
	g.Vertices[n2].Adj = append(g.Vertices[n2].Adj, n1)
}

//// AddUndirectedEdge adds an undirected edge to another node pointer
//// Note that this doesn't check if the second node is within the same graph,
//// and this doesn't check for duplicate edges
//func (n1 *Node) AddUndirectedEdge(n2 *Node) {
//	n2.Adj = append(n2.Adj, n1)
//	n1.Adj = append(n1.Adj, n2)
//}

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
func NewRandomGraph(nodeCount int, degree float32) Graph {
	g := New(nodeCount)

	// if branching factor is bFactor, then given n1, n2 nodes in g, then the
	// probability of an undirected edge is bFactor / (nodeCount - 1)
	pEdge := degree / float32(nodeCount-1)

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
func NewRandomGraphParallel(nVertices int, degree float32,
	nThreads int) Graph {

	//g := NewParallel(nodeCount, nThreads)
	g := New(nVertices)

	nodesPerThread := nVertices / nThreads
	if nVertices%nThreads != 0 {
		nodesPerThread++
	}

	pEdge := degree / float32(nVertices-1)

	var wg sync.WaitGroup
	wg.Add(nThreads)

	threadFunc := func(start int) {
		defer wg.Done()
		for i := start; i < start+nodesPerThread && i < nVertices; i++ {
			source := rand.NewSource(time.Now().UnixNano())
			generator := rand.New(source)

			for j := 0; j < i; j++ {
				if generator.Float32() < pEdge {
					g.Vertices[i].Mutex.Lock()
					g.Vertices[j].Mutex.Lock()
					g.AddUndirectedEdge(i, j)
					g.Vertices[j].Mutex.Unlock()
					g.Vertices[i].Mutex.Unlock()
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
