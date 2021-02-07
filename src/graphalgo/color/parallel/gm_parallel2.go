// This implementation groups many nodes into a single goroutine
package parallel

import (
	"graph"
	"runtime"
	"sync"
)

// colorNodeParallel speculatively colors a single node, not paying attention
// to data consistency (this will be detected in conflict resolution)
func colorNodeParallel2(g *graph.Graph, u []int, maxColor int,
	wg *sync.WaitGroup) {

	defer wg.Done()

	neighborColors := make([]bool, maxColor)
	neighborColorsZeros := make([]bool, maxColor)

	for _, i := range u {
		copy(neighborColors[:], neighborColorsZeros[:])

		v := &g.Vertices[i]
		for _, j := range v.Adj {
			neighborColors[g.Vertices[j].Value] = true
		}

		for j := 0; j < maxColor; j++ {
			if !neighborColors[j] {
				v.Value = j
				break
			}
		}
	}
}

func checkNodeConflictsParallel2(g *graph.Graph, u []int, wg *sync.WaitGroup,
	r *[]int, m *sync.Mutex) {

	defer wg.Done()

	for _, i := range u {
		v := &g.Vertices[i]
		for _, j := range v.Adj {
			if g.Vertices[j].Value == v.Value && j > i {
				m.Lock()
				*r = append(*r, i)
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
	u := make([]int, len(g.Vertices))
	for i := range g.Vertices {
		u[i] = i
	}

	// create secondary buffer
	r := make([]int, 0, len(u)/10)

	// helper function
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	// repeat process until run out of nodes to recolor
	for len(u) > 0 {
		nVertices := len(u)

		nodesPerThread := nVertices / nThreads
		if nVertices%nThreads != 0 {
			nodesPerThread++
		}

		// speculative coloring
		wg.Add(nThreads)
		for i := 0; i < nThreads; i++ {
			start := min(i * nodesPerThread, nVertices)
			end := min(start + nodesPerThread, nVertices)
			go colorNodeParallel2(g, u[start:end], maxColor, &wg)
		}
		wg.Wait()

		// conflict resolution: generate a list of nodes to recolor
		// provide the channel with a reasonably-sized buffer (?), since we
		// don't need the values immediately
		wg.Add(nThreads)
		for i := 0; i < nThreads; i++ {
			start := min(i * nodesPerThread, nVertices)
			end := min(start + nodesPerThread, nVertices)
			go checkNodeConflictsParallel2(g, u[start:end], &wg, &r, &m)
		}
		wg.Wait()

		// avoid reallocation: reuse buffers
		tmp := u
		u = r
		r = tmp[:0]
	}
}
