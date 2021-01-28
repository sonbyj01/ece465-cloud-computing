package graph

// Node represents a node of a Graph object
type Node struct {
	Value int
	Adj   []*Node
}

// Graph represents a generic graph data structure
type Graph struct {
	Nodes []Node
}

// AddNode adds a node to a graph
func (g *Graph) AddNode(value int) {
	g.Nodes = append(g.Nodes, Node{value, make([]*Node, 0)})
}

// AddUndirectedEdge adds an undirected edge between two nodes in a graph
func (g *Graph) AddUndirectedEdge(n1, n2 int) {
	if n1 >= len(g.Nodes) || n2 >= len(g.Nodes) {
		panic("Invalid node indices")
	}

	g.Nodes[n1].Adj = append(g.Nodes[n1].Adj, &g.Nodes[n2])
	g.Nodes[n2].Adj = append(g.Nodes[n2].Adj, &g.Nodes[n1])
}

// AddUndirectedEdge adds an undirected edge to another node pointer
// note that this doesn't check if the second node is within the same graph
func (n1 *Node) AddUndirectedEdge(n2 *Node) {
	n2.Adj = append(n2.Adj, n1)
	n1.Adj = append(n1.Adj, n2)
}
