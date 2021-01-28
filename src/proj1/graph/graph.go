package graph

// Node represents a node of a Graph object
type Node struct {
	value int
	edges []*Node
}

// Graph represents a generic graph data structure
type Graph struct {
	nodes []Node
}

// AddNode adds a node to a graph
func (g *Graph) AddNode(value int) {
	g.nodes = append(g.nodes, Node{value, make([]*Node, 0)})
}

// AddUndirectedEdge adds an undirected edge between two nodes in a graph
func (g *Graph) AddUndirectedEdge(n1, n2 int) {
	if n1 >= len(g.nodes) || n2 >= len(g.nodes) {
		panic("Invalid node indices")
	}

	g.nodes[n1].edges = append(g.nodes[n1].edges, &g.nodes[n2])
	g.nodes[n2].edges = append(g.nodes[n2].edges, &g.nodes[n1])
}

// AddUndirectedEdge adds an undirected edge to another node pointer
// note that this doesn't check if the second node is within the same graph
func (n1 *Node) AddUndirectedEdge(n2 *Node) {
	n2.edges = append(n2.edges, n1)
	n1.edges = append(n1.edges, n2)
}