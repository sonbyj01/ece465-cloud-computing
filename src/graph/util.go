package graph

import "fmt"

// Print prints out a list of a graph's nodes and values, as well as their
// neighbors and values
func (g *Graph) Print() {
	for i := range g.Nodes {
		fmt.Printf("%d: %d\n", i, g.Nodes[i].Value)
		for _, neighbor := range g.Nodes[i].Adj {
			fmt.Printf("\t%d: %d\n", neighbor.Index, neighbor.Value)
		}
	}
}

// CheckValidColoring checks whether a graph is appropriately colored
func (g *Graph) CheckValidColoring() bool {
	for i := range g.Nodes {
		for _, neighbor := range g.Nodes[i].Adj {
			if g.Nodes[i].Value == neighbor.Value {
				return false
			}
		}
	}
	return true
}
