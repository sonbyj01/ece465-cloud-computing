package main

import (
	"bufio"
	"graph"
	"os"
	"runtime"
)

// main is just for scratch work or for generating sample graphs
// the real drivers for Project 2 are in proj2/client and proj2/server.
func main() {
	// params for sample graph
	nVertices := 1000
	degree := float32(100)
	nThreads := 2 * runtime.NumCPU()
	outFile := "sample.graph"

	// generate some sample graphs for use as testcases
	g1 := graph.NewRandomGraphParallel(nVertices, degree, nThreads)

	// write file
	file, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0666)
	if err != nil {
		panic(err)
	}

	// dump
	err = g1.Dump(bufio.NewWriter(file))
	if err != nil {
		panic(err)
	}
}
