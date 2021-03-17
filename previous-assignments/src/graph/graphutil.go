package graph

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Print prints out a list of a graph's vertices and values, as well as their
// neighbors and values
func (g *Graph) Print() {
	for i := range g.Vertices {
		fmt.Printf("%d: %d\n", i, g.Vertices[i].Value)
		for _, j := range g.Vertices[i].Adj {
			fmt.Printf("\t%d: %d\n", j, g.Vertices[j].Value)
		}
	}
}

// PrintSubgraph prints out a list of a graph's vertices (with an index offset)
// and values, as well as their neighbors
func (g *Graph) PrintSubgraph(offset int) string {
	res := ""
	for i := range g.Vertices {
		res += fmt.Sprintf("%d: %d; ", i+offset, g.Vertices[i].Value)
		for _, j := range g.Vertices[i].Adj {
			res += fmt.Sprintf("%d, ", j)
		}
		res += "\n"
	}
	return res
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
		vertexComponents := strings.Split(scanner.Text(), ";")

		// get vertex value
		value, err := strconv.Atoi(vertexComponents[0])
		if err != nil {
			return nil, err
		}
		v.Value = value

		// get vertex adjacency list, or skip if no adjacent vertices
		if len(vertexComponents[1]) == 0 {
			continue
		}
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
	// write number of vertices
	_, err := io.WriteString(writer, strconv.Itoa(len(g.Vertices))+"\n")
	if err != nil {
		return err
	}

	// write each vertex
	for i := range g.Vertices {
		v := &g.Vertices[i]
		s := strconv.Itoa(v.Value) + ";"
		for index, j := range v.Adj {
			if index > 0 {
				s += ","
			}
			s += strconv.Itoa(j)
		}
		_, err = io.WriteString(writer, s+"\n")
		if err != nil {
			return err
		}
	}

	return nil
}
