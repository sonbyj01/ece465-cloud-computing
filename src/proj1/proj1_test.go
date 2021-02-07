package main

import (
	"graph"
	"graphalgo/color/parallel"
	"graphalgo/color/sequential"
	"math"
	"math/rand"
	"testing"
)

// countEdges is a helper for TestBranchingFactor
func countEdges(g graph.Graph) int {
	edges := 0

	for i := 0; i < len(g.Vertices); i++ {
		edges += len(g.Vertices[i].Adj)
	}

	return edges
}

// TestAverageDegree verifies that the average degree is within tolerance
func TestAverageDegree(t *testing.T) {
	const maxGraphSize, maxBfRatio, tolerance = 1000, 0.75, 0.20

	// try some random large graph sizes and branching factors
	for i := 0; i < 20; i++ {
		N := maxGraphSize/2 + rand.Int()%(maxGraphSize/2)
		desiredBf := rand.Float64() * float64(N) * maxBfRatio
		actualBf := float64(countEdges(
			graph.NewRandomGraphParallel(N, float32(desiredBf), 50))) /
			float64(N)

		t.Logf("Test: NewRandomGraph(%d, %f)", N, desiredBf)

		if math.Abs(desiredBf-actualBf) > tolerance*desiredBf {
			t.Errorf("Branching factor error. Got %f, desired %f",
				actualBf, desiredBf)
		}
	}
}

// TestSequential checks that the sequential coloring works
func TestSequential(t *testing.T) {
	N := 1000
	deg := float32(30)
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

	t.Logf("Test: NewRandomGraph(%d, %f)", N, deg)
	g = graph.NewRandomGraph(N, deg)
	sequential.ColorSequential(&g, maxColor)
	if !g.CheckValidColoring() {
		t.Errorf("NewRandomGraph is improperly colored")
	}
}

// TestParallel checks that the parallel coloring works
func TestParallel(t *testing.T) {
	N := 1000
	deg := float32(30)
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

	t.Logf("Test: NewRandomGraph(%d, %f)", N, deg)
	g = graph.NewRandomGraph(N, deg)
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

	t.Logf("Test: NewRingGraph(%d)", N)
	g = graph.NewRingGraph(N)
	parallel.ColorParallelGM2(&g, maxColor)
	if !g.CheckValidColoring() {
		t.Errorf("NewRingGraph is improperly colored")
	}

	t.Logf("Test: NewRandomGraph(%d, %f)", N, deg)
	g = graph.NewRandomGraph(N, deg)
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

// BenchmarkNewRandomGraph benches the time it takes to generate
// a new random graph
func BenchmarkNewRandomGraph(b *testing.B) {
	N := 10000
	deg := float32(1000)

	for i := 0; i < b.N; i++ {
		graph.NewRandomGraph(N, deg)
	}
}

// BenchmarkNewRandomGraphParallel benches the time it takes to generate
// a new random graph in parallel
func BenchmarkNewRandomGraphParallel(b *testing.B) {
	N := 10000
	deg := float32(1000)

	for i := 0; i < b.N; i++ {
		graph.NewRandomGraphParallel(N, deg, 50)
	}
}

type coloringAlgorithm = func(*graph.Graph, int)

// benchmarkColoring is a helper for the BenchmarkColor* benchmarks
func benchmarkColoring(b *testing.B, N int, deg float32, ca coloringAlgorithm) {
	maxColor := 3 * int(deg) / 2

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		g := graph.NewRandomGraphParallel(N, deg, 50)
		b.StartTimer()

		ca(&g, maxColor)
	}
}

// BenchmarkColorSequentialV100Bf10 benchmarks parallel coloring with 100
// nodes and average branching factor of 10
func BenchmarkColorSequentialV100Bf10(b *testing.B) {
	benchmarkColoring(b, 100, 10, sequential.ColorSequential)
}

// BenchmarkColorSequentialV1000Bf100 benchmarks parallel coloring with
// 1000 nodes and average branching factor of 100
func BenchmarkColorSequentialV1000Bf100(b *testing.B) {
	benchmarkColoring(b, 1000, 100, sequential.ColorSequential)
}

// BenchmarkColorSequentialV10000Bf1000 benchmarks parallel coloring with
// 10000 nodes and average branching factor of 1000
func BenchmarkColorSequentialV10000Bf1000(b *testing.B) {
	benchmarkColoring(b, 10000, 1000, sequential.ColorSequential)
}

// BenchmarkColorSequentialV50000Bf5000 benchmarks parallel coloring with
// 50000 nodes and average branching factor of 5000
func BenchmarkColorSequentialV50000Bf5000(b *testing.B) {
	benchmarkColoring(b, 50000, 5000, sequential.ColorSequential)
}

// BenchmarkColorParallelGMV100Bf10 benchmarks parallel coloring with 100
// nodes and average branching factor of 10
func BenchmarkColorParallelGMV100Bf10(b *testing.B) {
	benchmarkColoring(b, 100, 10, parallel.ColorParallelGM)
}

// BenchmarkColorParallelGMV1000Bf100 benchmarks parallel coloring with
// 1000 nodes and average branching factor of 100
func BenchmarkColorParallelGMV1000Bf100(b *testing.B) {
	benchmarkColoring(b, 1000, 100, parallel.ColorParallelGM)
}

// BenchmarkColorParallelGMV10000Bf1000 benchmarks parallel coloring with
// 10000 nodes and average branching factor of 1000
func BenchmarkColorParallelGMV10000Bf1000(b *testing.B) {
	benchmarkColoring(b, 10000, 1000, parallel.ColorParallelGM)
}

// BenchmarkColorParallelGMV50000Bf5000 benchmarks parallel coloring with
// 50000 nodes and average branching factor of 5000
func BenchmarkColorParallelGMV50000Bf5000(b *testing.B) {
	benchmarkColoring(b, 50000, 5000, parallel.ColorParallelGM)
}

// BenchmarkColorParallelGM2V100Bf10 benchmarks parallel coloring with 100
// nodes and average branching factor of 10
func BenchmarkColorParallelGM2V100Bf10(b *testing.B) {
	benchmarkColoring(b, 100, 10, parallel.ColorParallelGM2)
}

// BenchmarkColorParallelGM2V1000Bf100 benchmarks parallel coloring with
// 1000 nodes and average branching factor of 100
func BenchmarkColorParallelGM2V1000Bf100(b *testing.B) {
	benchmarkColoring(b, 1000, 100, parallel.ColorParallelGM2)
}

// BenchmarkColorParallelGM2V10000Bf1000 benchmarks parallel coloring with
// 10000 nodes and average branching factor of 1000
func BenchmarkColorParallelGM2V10000Bf1000(b *testing.B) {
	benchmarkColoring(b, 10000, 1000, parallel.ColorParallelGM2)
}

// BenchmarkColorParallelGM2V50000Bf5000 benchmarks parallel coloring with
// 50000 nodes and average branching factor of 5000
func BenchmarkColorParallelGM2V50000Bf5000(b *testing.B) {
	benchmarkColoring(b, 50000, 5000, parallel.ColorParallelGM2)
}
