// This implementation groups many nodes into a single goroutine
package parallel

import (
	"graph"
	"runtime"
	"sync"
)

// colorNodeParallel speculatively colors a single node, not paying attention
// to data consistency (this will be detected in conflict resolution)
func colorNodeParallel2(u []*graph.Node, start, end, maxColor int,
	wg *sync.WaitGroup) {

	defer wg.Done()

	neighborColors := make([]bool, maxColor)
	neighborColorsZeros := make([]bool, maxColor)

	for i := start; i < end; i++ {
		copy(neighborColors[:], neighborColorsZeros[:])

		for _, neighbor := range u[i].Adj {
			neighborColors[neighbor.Value] = true
		}

		for j := 0; j < maxColor; j++ {
			if !neighborColors[j] {
				u[i].Value = j
				break
			}
		}
	}
}

func checkNodeConflictsParallel2(u []*graph.Node, start, end int,
	wg *sync.WaitGroup, r *[]*graph.Node, m *sync.Mutex) {

	defer wg.Done()

	for i := start; i < end; i++ {
		node := u[i]
		for _, neighbor := range node.Adj {
			if neighbor.Value == node.Value && neighbor.Index > node.Index {
				m.Lock()
				*r = append(*r, node)
				m.Unlock()
				break
			}
		}
	}
}

// ColorParallelGM is the driver for the parallel coloring scheme following the
// Gebremedhin-Manne color outlined in https://www.osti.gov/biblio/1246285
func ColorParallelGM2(g *graph.Graph, maxColor int) {
	var wg sync.WaitGroup
	var m sync.Mutex
	nThreads := 2 * runtime.NumCPU()

	// set u to be a list of all of the nodes in the graph; it has
	// to be a list of node pointers so we actually update the graph
	u := make([]*graph.Node, len(g.Nodes))
	for i := range g.Nodes {
		u[i] = &g.Nodes[i]
	}

	// create secondary buffer
	r := make([]*graph.Node, 0, len(u)/10)

	// repeat process until run out of nodes to recolor
	for len(u) > 0 {
		nNodes := len(u)

		nodesPerThread := nNodes / nThreads
		if nNodes%nThreads != 0 {
			nodesPerThread++
		}

		// speculative coloring
		wg.Add(nThreads)
		for i := 0; i < nThreads; i++ {
			start := i * nodesPerThread
			end := start + nodesPerThread
			if end >= nNodes {
				end = nNodes
			}
			go colorNodeParallel2(u, start, end, maxColor, &wg)
		}
		wg.Wait()

		// conflict resolution: generate a list of nodes to recolor
		// provide the channel with a reasonably-sized buffer (?), since we
		// don't need the values immediately
		wg.Add(nThreads)
		for i := 0; i < nThreads; i++ {
			start := i * nodesPerThread
			end := start + nodesPerThread
			if end >= nNodes {
				end = nNodes
			}
			go checkNodeConflictsParallel2(u, start, end, &wg, &r, &m)
		}
		wg.Wait()

		// avoid reallocation: reuse buffers
		tmp := u
		u = r
		r = tmp[:0]
	}
}
