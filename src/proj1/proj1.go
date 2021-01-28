package main

import (
	"proj1/graph"
)

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

func main() {
	N := 100
	completeGraph := graph.NewCompleteGraph(N)

	// maxColor for a very simple coloring algorithm
	maxColor := 100

	completeGraph.Print()
	//for u := completeGraph; len(u.Nodes) > 0; {
	//	colorSequential(&u, maxColor)
	//}
	colorSequential(&completeGraph, maxColor)
	completeGraph.Print()
}
