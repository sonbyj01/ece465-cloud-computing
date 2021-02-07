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

// NewCompleteGraph generates a complete graph with nVertices nodes
func NewCompleteGraph(nVertices int) Graph {
	g := New(nVertices)

	for i := 1; i < nVertices; i++ {
		for j := 0; j < i; j++ {
			g.AddUndirectedEdge(i, j)
		}
	}

	return g
}

// NewRingGraph generates a graph in which each node has exactly two neighbors
func NewRingGraph(nVertices int) Graph {
	g := New(nVertices)

	for i := 0; i < nVertices-1; i++ {
		g.AddUndirectedEdge(i, i+1)
	}

	// if n=2, this edge is already created
	if nVertices > 2 {
		g.AddUndirectedEdge(0, nVertices-1)
	}

	return g
}

// NewRandomGraph generates a graph with nVertices nodes and average
// branching factor bFactor
func NewRandomGraph(nVertices int, degree float32) Graph {
	g := New(nVertices)

	// if branching factor is bFactor, then given n1, n2 nodes in g, then the
	// probability of an undirected edge is bFactor / (nVertices - 1)
	pEdge := degree / float32(nVertices-1)

	for i := 1; i < nVertices; i++ {
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

	//g := NewParallel(nVertices, nThreads)
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
