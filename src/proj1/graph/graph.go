package graph

// Node represents a node of a Graph object
type Node struct {
	Value int
	Edges []int
}

// Graph represents a generic graph data structure
type Graph struct {
	Nodes []Node
}

// AddNode adds a node to a graph
func (g *Graph) AddNode(value int) {
	g.Nodes = append(g.Nodes, Node{value, make([]int, 0)})
}

// AddUndirectedEdge adds an undirected edge between two nodes in a graph
func (g *Graph) AddUndirectedEdge(n1, n2 int) {
	g.Nodes[n1].Edges = append(g.Nodes[n1].Edges, n2)
	g.Nodes[n2].Edges = append(g.Nodes[n2].Edges, n1)
}
