package graph

// Partition simply partitions a graph into nNodes subgraphs comprising
// consecutive vertices
func (g *Graph) Partition(nNodes int) {

	nVertices := len(g.Vertices)
	verticesPerNode := nVertices / nNodes
	if nVertices%nNodes != 0 {
		verticesPerNode++
	}

}
