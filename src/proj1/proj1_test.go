package main

import (
	"math"
	"math/rand"
	"proj1/graph"
	"testing"
)

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
	for i := 0; i < 100; i++ {
		N := maxGraphSize / 2 + rand.Int() % (maxGraphSize / 2)
		desiredBf := rand.Float64() * float64(N) * maxBfRatio
		actualBf := float64(countEdges(
			graph.NewRandomGraph(N, float32(desiredBf)))) / float64(N)

		t.Logf("Test: NewRandomGraph(%d, %f)", N, desiredBf)

		if math.Abs(desiredBf-actualBf) > threshold*desiredBf {
			t.Errorf("Branching factor error. Got %f, desired %f",
				actualBf, desiredBf)
		}
	}
}
