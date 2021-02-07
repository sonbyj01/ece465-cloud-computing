package graph

import "fmt"

// Print prints out a list of a graph's nodes and values, as well as their
// neighbors and values
func (g *Graph) Print() {
	for i := range g.Vertices {
		fmt.Printf("%d: %d\n", i, g.Vertices[i].Value)
		for _, j := range g.Vertices[i].Adj {
			fmt.Printf("\t%d: %d\n", j, g.Vertices[j].Value)
		}
	}
}

// CheckValidColoring checks whether a graph is appropriately colored
func (g *Graph) CheckValidColoring() bool {
	for i := range g.Vertices {
		for _, j := range g.Vertices[i].Adj {
			if g.Vertices[i].Value == g.Vertices[j].Value {
				return false
			}
		}
	}
	return true
}

// Read reads a grpah from file
func Read() Graph {
	// TODO
	return Graph{};
}

// Write writes a graph to file
func (g *Graph) Write() {
	// TODO
}
