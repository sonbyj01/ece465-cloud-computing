package distributed

import (
	"graphnet"
	"math"
	"sync"
)

// findEdgeVertices returns a boolean array indicating which nodes are on
// the partition edge
// TODO: can easily parallelize this over multiple cores on a single node
func findEdgeVertices(sg *Subgraph) []bool {
	edgeVertexMap := make([]bool, len(sg.Vertices))

	for i := range sg.Vertices {
		for _, j := range sg.Vertices[i].Adj {
			if j < sg.iBegin || j >= sg.iEnd {
				edgeVertexMap[j] = true
				break
			}
		}
	}

	return edgeVertexMap
}

// coloring for one vertex superset
func colorSpeculative(u []int, maxColor int, sg *Subgraph,
	nodes []graphnet.Node, wg *sync.WaitGroup) {

	defer wg.Done()
	iBegin, iEnd := sg.iBegin, sg.iEnd
	var color int
	neighborColors := make([]bool, maxColor)
	neighborColorsDefault := make([]bool, maxColor)

	// loop over vertices in superset
	for _, i := range u {
		v := &sg.Vertices[i]

		copy(neighborColors, neighborColorsDefault)

		// speculatively color
		for _, j := range v.Adj {
			if j < iBegin || j >= iEnd {
				color = sg.stored[j]
			} else {
				color = sg.Vertices[j - sg.iBegin].Value
			}

			neighborColors[color] = true
		}

		// find first valid color
		for j := range neighborColors {
			if !neighborColors[j] {
				v.Value = j
				// TODO: if need to notify another node, notify them
				break
			}
		}
	}
}

// ColorDistributed is the main driver for the distributed coloring algorithm
// on the slave node, and is called after all the connections are set up
func ColorDistributed(sg *Subgraph, nodes []graphnet.Node,
	maxColor, nThreads int) {

	edgeVertexMap := findEdgeVertices(sg)
	var wg sync.WaitGroup
	nNodes := len(nodes)

	// initialize U to be all of the vertices in sg
	u := make([]int, len(sg.Vertices))
	for i := 0; i < len(sg.Vertices); i++ {
		u[i] = i
	}

	// loop until u is empty
	for len(u) > 0 {
		// listen on socket until all vertex information received
		// TODO: do this asynchronously
		//for nodesCompleted := 0; nodesCompleted < nNodes-1; nodesCompleted++ {
			// TODO: get node info until receive message that node is completed
		//}

		// loop over vertices, assign each vertex a permissible color
		// send colors of boundary vertices to relevant nodes
		// receive color information from other nodes;
		// don't use supersteps, rather choose number of threads; channels
		// will be buffered anyways
		// TODO: run these in parallel
		nVertices := len(u)
		verticesPerThread := int(math.Ceil(float64(nVertices / nThreads)))
		for i := 0; i < nThreads; i++ {
			start := i * verticesPerThread
			end := (i + 1) * verticesPerThread
			if end > nVertices {
				end = nVertices
			}

			go colorSpeculative(u[start:end], maxColor, sg, nodes, &wg)
		}

		// TODO: when all speculative coloring is done, notify all nodes

		// for each superset (in parallel): {
		// }

		// wait until all incoming messages are successfully received
		// for each boundary vertex, check for conflicts (in parallel)
		// add conflicting nodes to R
		// TODO: implement conflict resolution

		// set U to R
		// TODO: implement this
	}

	// TODO: when done coloring, notify all nodes
}
