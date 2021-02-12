package distributed

import (
	"graphnet"
	"log"
	"math"
	"sync"
)

// findEdgeVertices returns a boolean array indicating which nodes are on
// the partition edge
// TODO: can easily parallelize this over multiple cores on a single node
// TODO: is this even needed? this is kind of handled automatically
//func findEdgeVertices(sg *Subgraph) []bool {
//	edgeVertexMap := make([]bool, len(sg.Vertices))
//
//	for i := range sg.Vertices {
//		for _, j := range sg.Vertices[i].Adj {
//			if j < sg.iBegin || j >= sg.iEnd {
//				edgeVertexMap[j] = true
//				break
//			}
//		}
//	}
//
//	return edgeVertexMap
//}

// colorSpeculative speculatively colors one group of vertices and notifies
// other nodes when their neighbors are updated
func colorSpeculative(u []int, maxColor int, sg *Subgraph,
	nodes []graphnet.Node, wg *sync.WaitGroup) {

	defer wg.Done()
	iBegin, iEnd := sg.iBegin, sg.iEnd
	var color int
	neighborColors := make([]bool, maxColor)
	neighborColorsDefault := make([]bool, maxColor)

	// loop over vertices for this thread
	for _, i := range u {
		v := &sg.Vertices[i]

		copy(neighborColors, neighborColorsDefault)

		// speculatively color
		for _, j := range v.Adj {
			if j < iBegin || j >= iEnd {
				color = sg.stored[j]
			} else {
				color = sg.Vertices[j-sg.iBegin].Value
			}

			neighborColors[color] = true
		}

		// find first valid color
		for j := range neighborColors {
			if !neighborColors[j] {
				v.Value = j

				// notify all larger neighbors in different subgraphs
				for _, k := range v.Adj {
					if k >= iEnd {
						// TODO: fix this
						sg.sendVertexData(&nodes[0], graphnet.VertexData{
							Vertices: make([]int, 0),
							Colors:   make([]int16, 0),
						})
					}
				}

				break
			}
		}
	}
}

// resolveConflicts simply marks nodes that have conflicts to be recolored
// in the next round; doesn't require any inter-node communication
func resolveConflicts(u []int, sg *Subgraph, r *[]int, wg *sync.WaitGroup) {

	defer wg.Done()

	iBegin, iEnd := sg.iBegin, sg.iEnd
	var color int

	for _, i := range u {
		v := &sg.Vertices[i]

		for _, j := range v.Adj {
			if j >= iBegin && j < iEnd {
				color = sg.Vertices[j-iBegin].Value
			} else {
				color = sg.stored[j]
			}

			// if conflict detected, set larger-indexed node to be recolored
			if color == v.Value && i+iBegin > j {
				*r = append(*r, i)
			}
		}
	}
}

// ColorDistributed is the main driver for the distributed coloring algorithm
// on the slave node, and is called after all the connections are set up
func ColorDistributed(sg *Subgraph, nodes []graphnet.Node,
	maxColor, nThreads int, logger *log.Logger) {

	//edgeVertexMap := findEdgeVertices(sg)
	var wg sync.WaitGroup
	nNodes := len(nodes)

	// initialize U to be all of the vertices in sg
	u := make([]int, len(sg.Vertices))
	for i := 0; i < len(sg.Vertices); i++ {
		u[i] = i
	}

	r := make([]int, 0)

	// keep track of how many other nodes are still computing
	totalNodes := nNodes - 1

	// loop until u is empty
	for len(u) > 0 {
		logger.Printf("Beginning new round: %d nodes to be colored\n",
			len(u))

		// listen on socket (async) until all vertex information received
		wg.Add(1)
		go func() {
			defer wg.Done()

			nodesWaiting := totalNodes
			for vertexData := range nodes[0].GetVertexChannel() {
				if vertexData.Type == graphnet.MSG_NODE_ROUND_FINISHED {
					// one node finished its round
					nodesWaiting--
				} else if vertexData.Type == graphnet.MSG_NODE_FINISHED {
					// one node finished all its work
					totalNodes--
					nodesWaiting--
				} else {
					// a node sent updated color information
					colors := vertexData.Data.Colors
					for i, index := range vertexData.Data.Vertices {
						sg.stored[index] = int(colors[i])
					}
				}

				// loop break condition
				if nodesWaiting == 0 {
					break
				}
			}
		}()

		// loop over vertices, assign each vertex a permissible color
		// send colors of boundary vertices to relevant nodes
		// receive color information from other nodes;
		// don't use supersteps, rather choose number of threads; channels
		// will be buffered anyways
		// TODO: run these in parallel
		nVertices := len(u)
		verticesPerThread := int(math.Ceil(float64(nVertices / nThreads)))
		wg.Add(nThreads)
		for i := 0; i < nThreads; i++ {
			start := i * verticesPerThread
			end := (i + 1) * verticesPerThread
			if end > nVertices {
				end = nVertices
			}

			go colorSpeculative(u[start:end], maxColor, sg, nodes, &wg)
		}

		// when speculative coloring for this round is done, notify all nodes
		// TODO: broadcast this to all nodes
		sg.sendControlMessage(&nodes[0], graphnet.MSG_NODE_ROUND_FINISHED)

		// wait until all incoming messages are successfully received;
		// this makes all steps work in lockstep
		wg.Wait()

		logger.Printf("Beginning conflict resolution stage\n")

		// for each boundary vertex, check for conflicts (in parallel)
		// add conflicting nodes to R
		wg.Add(nThreads)
		for i := 0; i < nThreads; i++ {
			start := i * verticesPerThread
			end := (i + 1) * verticesPerThread
			if end > nVertices {
				end = nVertices
			}

			go resolveConflicts(u[start:end], sg, &r, &wg)
		}
		wg.Wait()

		// set U to R
		u = r
		r = u[:0]
	}

	// TODO: when done coloring, notify all worker nodes and server
	sg.sendControlMessage(&nodes[0], graphnet.MSG_NODE_FINISHED)
}
