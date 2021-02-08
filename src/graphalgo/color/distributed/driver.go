package distributed

import "graphnet"

// FindEdgeVertices returns a boolean array indicating which nodes are on
// the partition edge
// TODO: can easily parallelize this over multiple cores on a single node
func FindEdgeVertices(sg *Subgraph) []bool {
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

// ColorDistributed is the main driver for the distributed coloring algorithm
// on the slave node, and is called after all the connections are set up
func ColorDistributed(sg *Subgraph, nodes []graphnet.Node, s int) {
	edgeVertexMap := FindEdgeVertices(sg)

	// initialize U to be all of the vertices in sg
	u := make([]int, len(sg.Vertices))
	for i := 0; i < len(sg.Vertices); i++ {
		u[i] = i
	}

	// loop until u is empty
	for len(u) > 0 {
		// partition u into subsets "supersteps" of size s
		nVertices := len(u)
		verticesPerSuperstep := nVertices / s
		if nVertices % s != 0 {
			verticesPerSuperstep++
		}

		// for each superset (in parallel): {
		//		loop over vertices, assign each vertex a permissible color
		//		send colors of boundary vertices to relevant nodes
		//		receive color information from other nodes
		// }

		// wait until all incoming messages are successfully received
		// for each boundary vertex, check for conflicts (in parallel)
		// add conflicting nodes to R

		// set U to R
	}
}