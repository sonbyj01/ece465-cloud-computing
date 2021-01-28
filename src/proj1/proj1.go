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

// colorParallel is the driver for the parallel coloring scheme
func colorParallel(g *graph.Graph, maxColor int) {
	var wg sync.WaitGroup

	// repeat process until run out of nodes to recolor
	for u := g.Nodes; len(u) > 0; {
		// speculative coloring
		wg.Add(len(u))
		for i := range u {
			go colorNodeParallel(&u[i], &wg, maxColor)
		}
		wg.Wait()

		// conflict resolution: generate a list of nodes to recolor
		// TODO: working here
		break
	}
}

func main() {
	N := 100
	completeGraph := graph.NewCompleteGraph(N)

	// maxColor for a very simple coloring algorithm
	maxColor := 100

	completeGraph.Print()
	colorParallel(&completeGraph, maxColor)
	completeGraph.Print()

	fmt.Printf("isColored: %t", completeGraph.CheckValidColoring())
}
