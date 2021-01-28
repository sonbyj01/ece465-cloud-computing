package main

import (
	"fmt"
	"proj1/graph"
	"runtime"
)

func f(from string) {
	for i := 0; i < 3; i++ {
		fmt.Println(from, ":", i)
	}
}

func printNodes(g graph.Graph) {
	fmt.Println("===============")
	for index, node := range g.Nodes {
		fmt.Println("-- DEBUGGING-- ", "Index: ", index,
			"Node Value: ", node.Value, "Edge Array: ", node.Edges)
	}
	fmt.Println("===============")
}

// speculative coloring
func coloring(g graph.Graph) {
	for index, node := range g.Nodes {
		for index2, _ := range node.Edges {
			if g.Nodes[index].Value == g.Nodes[index2].Value {
				g.Nodes[index].Value++
			}
		}
	}
}

// conflict detection and resolution
func conflictResolution(g graph.Graph) {
	for index, node := range g.Nodes {
		for index2, _ := range node.Edges {
			if g.Nodes[index].Value == g.Nodes[index2].Value {
				if index > index2 {
					g.Nodes[index].Value++
				} else {
					g.Nodes[index2].Value++
				}
			}
		}
	}
}

// showNo - Shows no from 0 to 99
func showNo() {
	for i := 0; i < 100; i++ {
		fmt.Println("value of i=", i)
	}
} // showAlphabets - shows alphabets from a-z
func showAlphabets() {
	for j := 'a'; j <= 'z'; j++ {
		fmt.Println("value of j=", string(j))
	}
}

// Verify number of logical cores available
func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func execute(id int) {
	fmt.Printf("id: %d\n", id)
}

func main() {
	//generatedGraph := graph.CompleteGraph(10)
	//printNodes(generatedGraph)
	//go coloring(generatedGraph)
	//go conflictResolution(generatedGraph)
	//printNodes(generatedGraph)

	//for i := 0; i < 5; i++ {
	//	go showNo()
	//	go showAlphabets()
	//}

	fmt.Println(runtime.NumCPU())

	//fmt.Println("Started")
	//for i := 0; i < 10; i++ {
	//	go execute(i)
	//}
	//time.Sleep(time.Second * 2)
	//fmt.Println("Finished")
}
