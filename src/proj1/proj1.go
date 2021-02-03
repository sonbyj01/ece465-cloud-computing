package main

import (
	"fmt"
	"proj1/algorithm"
	"proj1/graph"
)

// main is a sample entrypoint to show how to generate a graph and use the
// graph coloring functions, but you can see that most of our tests and
// benchmarks are in proj1_test.go
func main() {
	N := 12000

	fmt.Printf("Generating complete graph...\n")
	completeGraph := graph.NewCompleteGraph(N)

	// maxColor for a very simple coloring algorithm
	maxColor := 3 * N / 2

	// perform coloring
	fmt.Printf("Graph coloring...\n")
	//algorithm.ColorParallelGM(&completeGraph, maxColor)
	algorithm.ColorSequential(&completeGraph, maxColor)

	// check that the graph coloring worked
	fmt.Printf("isColored: %t", completeGraph.CheckValidColoring())
}
