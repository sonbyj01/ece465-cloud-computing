package main

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
	g.nodes = append(g, Node{value, make([]nodes, 0)})
}

// AddUndirectedEdge adds an undirected edge between two nodes in a graph
func (g *Graph) AddUndirectedEdge(n1, n2 *Node) {
	g.nodes[n1] = append(g.nodes[n1], n2)
	g.nodes[n2] = append(g.nodes[n2], n1)
}
