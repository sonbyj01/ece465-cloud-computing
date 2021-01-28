package main

import (
	"fmt"
	"proj1/graph"
	"sync"
)

func printNodes(g graph.Graph) {
	fmt.Println("===============")
	for index, node := range g.Nodes {
		fmt.Println("-- DEBUGGING-- ", "Index: ", index,
			"Node Value: ", node.Value, "Edge Array: ", node.Adj)
	}
	fmt.Println("===============")
}

// speculative coloring
func coloring(g graph.Graph) {
	var wg sync.WaitGroup

	for index, node := range g.Nodes {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for _, neighborNode := range node.Adj {
				if g.Nodes[index].Value == neighborNode.Value {
					g.Nodes[index].Value++
				}
			}
		}()
	}
	wg.Wait()
}

// conflict detection and resolution
func conflictResolution(g graph.Graph) {
	var wg sync.WaitGroup
	var temp graph.Graph

	for index, node := range g.Nodes {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for _, neighborNode := range node.Adj {
				if g.Nodes[index].Value == neighborNode.Value {
					temp.Nodes = append(temp.Nodes, g.Nodes[index])
				}
			}
		}()
	}
	wg.Wait()
	g.Nodes = temp.Nodes
}

func main() {
	generatedGraph := graph.CompleteGraph(4)

	count := 1

	printNodes(generatedGraph)

	for len(generatedGraph.Nodes) > 0 {
		coloring(generatedGraph)
		conflictResolution(generatedGraph)

		count++

		if count == 500 {
			break
		}
	}

	printNodes(generatedGraph)

	//printNodes(generatedGraph)
	//coloring(generatedGraph)
	//printNodes(generatedGraph)
	//conflictResolution(generatedGraph)
	//printNodes(generatedGraph)

	//fmt.Println(runtime.NumCPU())
}
