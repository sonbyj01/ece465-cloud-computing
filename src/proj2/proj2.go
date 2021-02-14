package main

import (
	"bufio"
	"fmt"
	"graph"
	"log"
	"os"
	"runtime"
)

// main is just for scratch work or for generating sample graphs
// the real drivers for Project 2 are in proj2/client and proj2/server.
func main() {
	// params for sample graph
	nVertices := 100
	degree := float32(10)
	nThreads := 2 * runtime.NumCPU()
	outFile := fmt.Sprintf("res/sample%d.graph", nVertices)

	// generate some sample graphs for use as testcases
	g := graph.NewRandomGraphParallel(nVertices, degree, nThreads)

	// write file
	log.Printf("Creating graph file %s...\n", outFile)
	file, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0666)
	if err != nil {
		log.Panic(err)
	}

	// dump
	log.Printf("Writing to file...")
	writer := bufio.NewWriter(file)
	err = g.Dump(writer)
	if err != nil {
		log.Panic(err)
	}
	err = writer.Flush()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Done\n")
}
