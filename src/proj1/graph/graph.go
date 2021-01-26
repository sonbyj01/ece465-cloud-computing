package graph

// Node represents a node of a Graph object
type Node struct {
	value int
	edges []int
}

// Graph represents a generic graph data structure
type Graph struct {
	nodes []Node
}

// AddNode adds a node to a graph
func (g *Graph) AddNode(value int) {
	g.nodes = append(g.nodes, Node{value, make([]int, 0)})
}

// AddUndirectedEdge adds an undirected edge between two nodes in a graph
func (g *Graph) AddUndirectedEdge(n1, n2 int) {
	g.nodes[n1].edges = append(g.nodes[n1].edges, n2)
	g.nodes[n2].edges = append(g.nodes[n2].edges, n1)
}
