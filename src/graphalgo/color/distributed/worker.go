package distributed

import (
	"graphnet"
	"log"
	"math"
	"sync"
)

// colorSpeculative speculatively colors one group of vertices and notifies
// other nodes when their neighbors are updated
func colorSpeculative(u []int, maxColor int, ws *WorkerState,
	wg *sync.WaitGroup) {

	defer wg.Done()
	iBegin, iEnd := ws.VertexBegin, ws.VertexEnd
	var color int
	neighborColors := make([]bool, maxColor)
	neighborColorsDefault := make([]bool, maxColor)
	sg := ws.Subgraph

	// loop over vertices for this thread
	for _, i := range u {
		v := &sg.Vertices[i]

		copy(neighborColors, neighborColorsDefault)

		// speculatively color
		for _, j := range v.Adj {
			if j < iBegin || j >= iEnd {
				color = ws.stored[j]
			} else {
				color = sg.Vertices[j-iBegin].Value
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
						// TODO: send vertex data to appropriate node
					}
				}

				break
			}
		}
	}
}

// resolveConflicts simply marks nodes that have conflicts to be recolored
// in the next round; doesn't require any inter-node communication
func resolveConflicts(u []int, ws *WorkerState, r *[]int, wg *sync.WaitGroup) {

	defer wg.Done()

	iBegin, iEnd := ws.VertexBegin, ws.VertexEnd
	var color int
	sg := ws.Subgraph

	for _, i := range u {
		v := &sg.Vertices[i]

		for _, j := range v.Adj {
			if j >= iBegin && j < iEnd {
				color = sg.Vertices[j-iBegin].Value
			} else {
				color = ws.stored[j]
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
func ColorDistributed(ws *WorkerState, maxColor, nThreads int,
	logger *log.Logger) {

	var wg sync.WaitGroup
	sg := ws.Subgraph
	//nNodes := len(sg.Vertices)
	buf := make([]byte, 8)

	// initialize U to be all of the vertices in sg
	u := make([]int, len(sg.Vertices))
	for i := 0; i < len(sg.Vertices); i++ {
		u[i] = i
	}

	r := make([]int, 0)

	// TODO: remove
	// keep track of how many other nodes are still computing
	//totalNodes := nNodes - 1

	// loop until u is empty
	for len(u) > 0 {
		logger.Printf("Beginning new round: %d nodes to be colored\n",
			len(u))

		// listen on socket (async) until all vertex information received
		wg.Add(1)
		go func() {
			defer wg.Done()

			// TODO: use handlers
			//nodesWaiting := totalNodes
			//for vertexData := range nodes[0].GetVertexChannel() {
			//	if vertexData.Type == graphnet.MSG_NODE_ROUND_FINISHED {
			//		// one node finished its round
			//		nodesWaiting--
			//	} else if vertexData.Type == graphnet.MSG_NODE_FINISHED {
			//		// one node finished all its work
			//		totalNodes--
			//		nodesWaiting--
			//	} else {
			//		// a node sent updated color information
			//		colors := vertexData.Data.Colors
			//		for i, index := range vertexData.Data.Vertices {
			//			sg.stored[index] = int(colors[i])
			//		}
			//	}
			//
			//	// loop break condition
			//	if nodesWaiting == 0 {
			//		break
			//	}
			//}
		}()

		// loop over vertices, assign each vertex a permissible color
		// send colors of boundary vertices to relevant nodes
		// receive color information from other nodes;
		// don't use supersteps, rather choose number of threads; channels
		// will be buffered anyways
		nVertices := len(u)
		verticesPerThread := int(math.Ceil(float64(nVertices / nThreads)))
		wg.Add(nThreads)
		for i := 0; i < nThreads; i++ {
			start := i * verticesPerThread
			end := (i + 1) * verticesPerThread
			if end > nVertices {
				end = nVertices
			}

			go colorSpeculative(u[start:end], maxColor, ws, &wg)
		}

		// when speculative coloring for this round is done, notify workers
		buf[0] = byte(ws.NodeIndex)
		ws.ConnPool.BroadcastWorkers(graphnet.MSG_NODE_ROUND_FINISHED, buf[:1])

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

			go resolveConflicts(u[start:end], ws, &r, &wg)
		}
		wg.Wait()

		// set U to R
		u = r
		r = u[:0]
	}

	// when done coloring, notify all nodes
	buf[0] = byte(ws.NodeIndex)
	ws.ConnPool.Broadcast(graphnet.MSG_NODE_FINISHED, buf[:1])
}
