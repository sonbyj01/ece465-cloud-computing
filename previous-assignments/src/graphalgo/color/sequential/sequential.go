package sequential

import "graph"

// ColorSequential performs a naive sequential.go Delta+1 coloring
// (suboptimal chromatic number, but very simple valid coloring)
func ColorSequential(g *graph.Graph, maxColor int) {
	neighborColors := make([]bool, maxColor)
	neighborColorsDefault := make([]bool, maxColor)

	for i := range g.Vertices {
		v := &g.Vertices[i]

		copy(neighborColors, neighborColorsDefault)

		for _, j := range v.Adj {
			neighborColors[g.Vertices[j].Value] = true
		}

		colorFound := false
		for j := 0; j < maxColor; j++ {
			if !neighborColors[j] {
				g.Vertices[i].Value = j
				colorFound = true
				break
			}
		}

		if !colorFound {
			panic("maxColor exceeded")
		}
	}
}
