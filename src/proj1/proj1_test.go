package main

import (
	"graph"
	"graphalgo/color/parallel"
	"graphalgo/color/sequential"
	"math"
	"math/rand"
	"runtime"
	"testing"
)

// countEdges is a helper for TestBranchingFactor
func countEdges(g graph.Graph) int {
	edges := 0

	for i := 0; i < len(g.Nodes); i++ {
		edges += len(g.Nodes[i].Adj)
	}

	return edges
}

// TestBranchingFactor verifies that the branching factor is
// approximately correct
func TestBranchingFactor(t *testing.T) {
	const maxGraphSize, maxBfRatio, threshold = 1000, 0.75, 0.20

	// try some random large graph sizes and branching factors
	for i := 0; i < 20; i++ {
		N := maxGraphSize/2 + rand.Int()%(maxGraphSize/2)
		desiredBf := rand.Float64() * float64(N) * maxBfRatio
		actualBf := float64(countEdges(
			graph.NewRandomGraphParallel(N, float32(desiredBf), 50))) /
				float64(N)

		t.Logf("Test: NewRandomGraph(%d, %f)", N, desiredBf)

		if math.Abs(desiredBf-actualBf) > threshold*desiredBf {
			t.Errorf("Branching factor error. Got %f, desired %f",
				actualBf, desiredBf)
		}
	}
}

// checkIndices is a helper for TestIndices
func checkIndices(g *graph.Graph, n int) bool {
	if n != len(g.Nodes) {
		return false
	}
	for i := range g.Nodes {
		if i != g.Nodes[i].Index {
			return false
		}
	}
	return true
}

// TestIndices verifies that the indices are correct
func TestIndices(t *testing.T) {
	N := 1000
	bf := float32(30)

	if g := graph.NewCompleteGraph(N); !checkIndices(&g, N) {
		t.Errorf("NewCompleteGraph creates wrong indices")
	}

	if g := graph.NewRingGraph(N); !checkIndices(&g, N) {
		t.Errorf("NewRingGraph creates wrong indices")
	}

	if g := graph.NewRandomGraph(N, bf); !checkIndices(&g, N) {
		t.Errorf("NewRandomGraph creates wrong indices")
	}

	if g := graph.NewRandomGraphParallel(N, bf, 50);
		!checkIndices(&g, N) {
		t.Errorf("NewRandomGraphParallel creates wrong indices")
	}
}

// TestSequential checks that the sequential coloring works
func TestSequential(t *testing.T) {
	N := 1000
	bf := float32(30)
	maxColor := 1000

	t.Logf("Test: NewCompleteGraph(%d)", N)
	g := graph.NewCompleteGraph(N)
	sequential.ColorSequential(&g, maxColor)
	if !g.CheckValidColoring() {
		t.Errorf("NewCompleteGraph is improperly colored")
	}

	t.Logf("Test: NewCompleteGraph(%d)", N)
	g = graph.NewRingGraph(N)
	sequential.ColorSequential(&g, maxColor)
	if !g.CheckValidColoring() {
		t.Errorf("NewRingGraph is improperly colored")
	}

	t.Logf("Test: NewRandomGraph(%d, %f)", N, bf)
	g = graph.NewRandomGraph(N, bf)
	sequential.ColorSequential(&g, maxColor)
	if !g.CheckValidColoring() {
		t.Errorf("NewRandomGraph is improperly colored")
	}
}

// TestParallel checks that the parallel coloring works
func TestParallel(t *testing.T) {
	N := 1000
	bf := float32(30)
	maxColor := 1000

	t.Logf("Test: NewCompleteGraph(%d)", N)
	g := graph.NewCompleteGraph(N)
	parallel.ColorParallelGM(&g, maxColor)
	if !g.CheckValidColoring() {
		t.Errorf("NewCompleteGraph is improperly colored")
	}

	t.Logf("Test: NewCompleteGraph(%d)", N)
	g = graph.NewRingGraph(N)
	parallel.ColorParallelGM(&g, maxColor)
	if !g.CheckValidColoring() {
		t.Errorf("NewRingGraph is improperly colored")
	}

	t.Logf("Test: NewRandomGraph(%d, %f)", N, bf)
	g = graph.NewRandomGraph(N, bf)
	parallel.ColorParallelGM(&g, maxColor)
	if !g.CheckValidColoring() {
		t.Errorf("NewRandomGraph is improperly colored")
	}

	t.Logf("Test: NewCompleteGraph(%d)", N)
	g = graph.NewCompleteGraph(N)
	parallel.ColorParallelGM2(&g, maxColor)
	if !g.CheckValidColoring() {
		t.Errorf("NewCompleteGraph is improperly colored")
	}

	t.Logf("Test: NewCompleteGraph(%d)", N)
	g = graph.NewRingGraph(N)
	parallel.ColorParallelGM2(&g, maxColor)
	if !g.CheckValidColoring() {
		t.Errorf("NewRingGraph is improperly colored")
	}

	t.Logf("Test: NewRandomGraph(%d, %f)", N, bf)
	g = graph.NewRandomGraph(N, bf)
	parallel.ColorParallelGM2(&g, maxColor)
	if !g.CheckValidColoring() {
		t.Errorf("NewRandomGraph is improperly colored")
	}
}

// BenchmarkNewGraph benches the time to generate a new graph
// and number its nodes
func BenchmarkNewGraph(b *testing.B) {
	N := 50000

	for i := 0; i < b.N; i++ {
		graph.New(N)
	}
}

// BenchmarkNewGraphParallel benches the time to generate a new graph
// and number its nodes in parallel
func BenchmarkNewGraphParallel(b *testing.B) {
	N := 50000

	for i := 0; i < b.N; i++ {
		graph.NewParallel(N, runtime.NumCPU())
	}
}

// BenchmarkNewRandomGraph benches the time it takes to generate
// a new random graph
func BenchmarkNewRandomGraph(b *testing.B) {
	N := 10000
	bf := float32(1000)

	for i := 0; i < b.N; i++ {
		graph.NewRandomGraph(N, bf)
	}
}

// BenchmarkNewRandomGraphParallel benches the time it takes to generate
// a new random graph in parallel
func BenchmarkNewRandomGraphParallel(b *testing.B) {
	N := 10000
	bf := float32(1000)

	for i := 0; i < b.N; i++ {
		graph.NewRandomGraphParallel(N, bf, 50)
	}
}

// checkGraph checks the a graph to make sure that there are no nils
// in the adjacency lists (this was a problem at some point)
func checkAdjacencyLists(g *graph.Graph) {
	for i := range g.Nodes {
		for j := range g.Nodes[i].Adj {
			if g.Nodes[i].Adj[j] == nil {
				panic("Nil in adjacency list")
			}
		}
	}
}

type coloringAlgorithm = func(*graph.Graph, int)
var algorithms = []coloringAlgorithm {
	sequential.ColorSequential,
	parallel.ColorParallelGM,
	parallel.ColorParallelGM2,
}

// benchmarkColoring is a helper for the BenchmarkColor* benchmarks
func benchmarkColoring(b *testing.B, N int, bf float32, ca coloringAlgorithm) {
	maxColor := 3 * int(bf) / 2

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		g := graph.NewRandomGraphParallel(N, bf, 50)
		b.StartTimer()

		ca(&g, maxColor)
	}
}

// BenchmarkColorSequentialV100Bf10 benchmarks sequential.go coloring with 100
// nodes and average branching factor of 10
func BenchmarkColorSequentialV100Bf10(b *testing.B) {
	benchmarkColoring(b, 100, 10, sequential.ColorSequential)
}

// BenchmarkColorSequentialV1000Bf100 benchmarks sequential.go coloring with
// 1000 nodes and average branching factor of 100
func BenchmarkColorSequentialV1000Bf100(b *testing.B) {
	benchmarkColoring(b, 1000, 100, sequential.ColorSequential)
}

// BenchmarkColorSequentialV10000Bf1000 benchmarks sequential.go coloring with
// 10000 nodes and average branching factor of 1000
func BenchmarkColorSequentialV10000Bf1000(b *testing.B) {
	benchmarkColoring(b, 10000, 1000, sequential.ColorSequential)
}

// BenchmarkColorSequentialV50000Bf5000 benchmarks sequential.go coloring with
// 50000 nodes and average branching factor of 5000
func BenchmarkColorSequentialV50000Bf5000(b *testing.B) {
	benchmarkColoring(b, 50000, 5000, sequential.ColorSequential)
}

// BenchmarkColorSequentialV100Bf10 benchmarks sequential.go coloring with 100
// nodes and average branching factor of 10
func BenchmarkColorParallelGMV100Bf10(b *testing.B) {
	benchmarkColoring(b, 100, 10, parallel.ColorParallelGM)
}

// BenchmarkColorParallelGMV1000Bf100 benchmarks sequential.go coloring with
// 1000 nodes and average branching factor of 100
func BenchmarkColorParallelGMV1000Bf100(b *testing.B) {
	benchmarkColoring(b, 1000, 100, parallel.ColorParallelGM)
}

// BenchmarkColorParallelGMV10000Bf1000 benchmarks sequential.go coloring with
// 10000 nodes and average branching factor of 1000
func BenchmarkColorParallelGMV10000Bf1000(b *testing.B) {
	benchmarkColoring(b, 10000, 1000, parallel.ColorParallelGM)
}

// BenchmarkColorParallelGMV50000Bf5000 benchmarks sequential.go coloring with
// 50000 nodes and average branching factor of 5000
func BenchmarkColorParallelGMV50000Bf5000(b *testing.B) {
	benchmarkColoring(b, 50000, 5000, parallel.ColorParallelGM)
}

// BenchmarkColorParallelGMV50000Bf5000 benchmarks sequential.go coloring with
// 50000 nodes and average branching factor of 5000
func BenchmarkColorParallelGM2V50000Bf5000(b *testing.B) {
	benchmarkColoring(b, 50000, 5000, parallel.ColorParallelGM2)
}
