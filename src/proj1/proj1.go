package main

import (
	"fmt"
	"proj1/graph"
	"sync"
)

// colorSequential performs a naive sequential Delta+1 coloring
// (suboptimal chromatic number, but very simple valid coloring)
func colorSequential(g *graph.Graph, maxColor int) {
	neighborColors := make([]bool, maxColor)

	for _, node := range g.Nodes {
		for i := 0; i < maxColor; i++ {
			neighborColors[i] = false
		}

		for _, neighbor := range node.Adj {
			neighborColors[neighbor.Value] = true
		}

		colorFound := false
		for i := 0; i < maxColor; i++ {
			if !neighborColors[i] {
				g.Nodes[node.Index].Value = i
				colorFound = true
				break
			}
		}

		if !colorFound {
			panic("maxColor exceeded")
		}
	}
}

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

// colorParallel is the driver for the parallel coloring scheme
func colorParallel(g *graph.Graph, maxColor int) {
	var wg sync.WaitGroup

	fmt.Printf("Beginning parallel coloring\n")

	// set u to be a list of all of the nodes in the graph; it has
	// to be a list of node pointers so we actually update the graph
	u := make([]*graph.Node, len(g.Nodes))
	for i, _ := range g.Nodes {
		u[i] = &g.Nodes[i]
	}

	// repeat process until run out of nodes to recolor
	for len(u) > 0 {
		fmt.Printf("New speculative coloring round, |u|=%d\n", len(u))

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

func main() {
	N := 12000

	fmt.Printf("Generating complete graph...\n")
	completeGraph := graph.NewCompleteGraph(N)

	// maxColor for a very simple coloring algorithm
	maxColor := 3 * N / 2

	// perform coloring
	fmt.Printf("Graph coloring...\n")
	//colorParallel(&completeGraph, maxColor)
	colorSequential(&completeGraph, maxColor)

	// check that the graph coloring worked
	fmt.Printf("isColored: %t", completeGraph.CheckValidColoring())
}
