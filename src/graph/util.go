package graph

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

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

// Load reads a graph from file
func Load(reader io.Reader) (*Graph, error) {
	scanner := bufio.NewScanner(reader)

	// get number of vertices and create graph
	var nVertices int
	scanner.Scan()
	nVertices, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, err
	}
	g := Graph{make([]Vertex, nVertices)}

	// read in all vertices
	for i := 0; scanner.Scan(); i++ {
		v := &g.Vertices[i]

		// split into value and adj list
		vertexComponents := strings.Split(scanner.Text(), ":")

		// get vertex value
		value, err := strconv.Atoi(vertexComponents[0])
		if err != nil {
			return nil, err
		}
		v.Value = value

		// get vertex adjacency list
		for _, adj := range strings.Split(vertexComponents[1], ",") {
			adjInt, err := strconv.Atoi(adj)
			if err != nil {
				return nil, err
			}
			v.Adj = append(v.Adj, adjInt)
		}
	}

	return &g, nil
}

// Dump writes a graph to file
func (g *Graph) Dump(writer io.Writer) error {
	_, err := io.WriteString(writer, strconv.Itoa(len(g.Vertices)) + "\n")
	if err != nil {
		return err
	}

	for i := range g.Vertices {
		v := &g.Vertices[i]
		s := strconv.Itoa(v.Value) + ";"
		for j := range v.Adj {
			s += strconv.Itoa(j) + ","
		}
		_, err = io.WriteString(writer, s[:len(s)-1] + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
