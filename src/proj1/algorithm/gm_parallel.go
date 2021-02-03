package algorithm

import (
	"proj1/graph"
	"sync"
)

// colorNodeParallel speculatively colors a single node, not paying attention
// to data consistency (this will be detected in conflict resolution)
func colorNodeParallel(n *graph.Node, wg *sync.WaitGroup, maxColor int) {
	defer wg.Done()

	neighborColors := make([]bool, maxColor)

	for _, neighbor := range n.Adj {
		neighborColors[neighbor.Value] = true
	}

	for i := 0; i < maxColor; i++ {
		if !neighborColors[i] {
			n.Value = i
			return
		}
	}
	panic("maxColor exceeded")
}

func checkNodeConflictsParallel(n *graph.Node, wg *sync.WaitGroup,
	ch chan *graph.Node) {
	defer wg.Done()

	for _, neighbor := range n.Adj {
		if neighbor.Value == n.Value && neighbor.Index > n.Index {
			ch <- n
			return
		}
	}
}

// ColorParallelGM is the driver for the parallel coloring scheme following the
// Gebremedhin-Manne algorithm outlined in https://www.osti.gov/biblio/1246285
func ColorParallelGM(g *graph.Graph, maxColor int) {
	var wg sync.WaitGroup

	// set u to be a list of all of the nodes in the graph; it has
	// to be a list of node pointers so we actually update the graph
	u := make([]*graph.Node, len(g.Nodes))
	for i := range g.Nodes {
		u[i] = &g.Nodes[i]
	}

	// repeat process until run out of nodes to recolor
	for len(u) > 0 {
		// speculative coloring
		wg.Add(len(u))
		for i := range u {
			go colorNodeParallel(u[i], &wg, maxColor)
		}
		wg.Wait()

		// conflict resolution: generate a list of nodes to recolor
		// provide the channel with a reasonably-sized buffer (?), since we
		// don't need the values immediately
		wg.Add(len(u))
		ch := make(chan *graph.Node, 64)
		for i := range u {
			go checkNodeConflictsParallel(u[i], &wg, ch)
		}

		// monitor to watch for the parallel routines to finish and close the
		// channel so the range loop below this knows when to finish; this
		// doesn't really have to happen in parallel if we make the channel
		// large enough, but this allows us to keep the channel buffer small
		go func() {
			wg.Wait()
			close(ch)
		}()

		u = make([]*graph.Node, 0)
		for node := range ch {
			u = append(u, node)
		}
	}
}