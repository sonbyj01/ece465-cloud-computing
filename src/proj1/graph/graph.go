package graph

// Node represents a node of a Graph object, with an adjacency list
// of pointers to other nodes
type Node struct {
	Value int
	Adj   []*Node
}

// Graph represents a very simple graph data structure
type Graph struct {
	Nodes []Node
}

// New returns a new graph
// Since this is namespaced under the graph package, can be used from the
// outside as graph.New(...)
func New(nodeCount int) Graph {
	return Graph{make([]Node, nodeCount)}
}

// AddNode adds a node to a graph
func (g *Graph) AddNode(value int) {
	g.Nodes = append(g.Nodes, Node{value, make([]*Node, 0)})
}

// AddUndirectedEdge adds an undirected edge between two nodes in a graph
// Note that this doesn't check for duplicate edges
func (g *Graph) AddUndirectedEdge(n1, n2 int) {
	if n1 >= len(g.Nodes) || n2 >= len(g.Nodes) {
		panic("Invalid node indices")
	}

	g.Nodes[n1].Adj = append(g.Nodes[n1].Adj, &g.Nodes[n2])
	g.Nodes[n2].Adj = append(g.Nodes[n2].Adj, &g.Nodes[n1])
}

// AddUndirectedEdge adds an undirected edge to another node pointer
// Note that this doesn't check if the second node is within the same graph,
// and this doesn't check for duplicate edges
func (n1 *Node) AddUndirectedEdge(n2 *Node) {
	n2.Adj = append(n2.Adj, n1)
	n1.Adj = append(n1.Adj, n2)
}
