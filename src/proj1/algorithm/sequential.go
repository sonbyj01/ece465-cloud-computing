package algorithm

import "proj1/graph"

// ColorSequential performs a naive sequential.go Delta+1 coloring
// (suboptimal chromatic number, but very simple valid coloring)
func ColorSequential(g *graph.Graph, maxColor int) {
	neighborColors := make([]bool, maxColor)

	for i := range g.Nodes {
		node := &g.Nodes[i]

		for i := 0; i < maxColor; i++ {
			neighborColors[i] = false
		}

		for _, neighbor := range node.Adj {
			// TODO: remove; this was a problem at some point
			if neighbor == nil {
				panic("Nil in adjacency list")
			}
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