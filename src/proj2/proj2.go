package main

import (
	"fmt"
	"graph"
)

func main() {
	g := graph.NewRandomGraph(1000, 100)

	fmt.Printf("Hello, world! %d\n", len(g.Vertices))
}
