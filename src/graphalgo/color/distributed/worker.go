package distributed

import (
	"encoding/binary"
	"graphnet"
	"log"
	"math"
)

// colorSpeculative speculatively colors one group of vertices and notifies
// other nodes when their neighbors are updated
func colorSpeculative(u []int, maxColor int, ws *WorkerState) {

	defer ws.ColorWg.Done()
	iBegin, iEnd := ws.VertexBegin, ws.VertexEnd
	var color int
	neighborColors := make([]bool, maxColor)
	neighborColorsDefault := make([]bool, maxColor)
	sg := ws.Subgraph
	buf := make([]byte, 8)

	// loop over vertices for this thread
	for _, i := range u {
		v := &sg.Vertices[i]

		copy(neighborColors, neighborColorsDefault)

		// speculatively color
		for _, j := range v.Adj {
			if j < iBegin || j >= iEnd {
				ws.StoredMutex.Lock()
				color = ws.Stored[j]
				ws.StoredMutex.Unlock()
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
						binary.LittleEndian.PutUint32(buf[:4], uint32(j))
						binary.LittleEndian.PutUint32(buf[4:], uint32(i+iBegin))
						// TODO: later work on buffering
						ws.ConnPool.Conns[1+(k/len(sg.Vertices))].
							WriteBytes(graphnet.MSG_VERTEX_INFO, buf,
								true)
					}
				}

				break
			}
		}
	}
}

// resolveConflicts simply marks nodes that have conflicts to be recolored
// in the next round; doesn't require any inter-node communication
func resolveConflicts(u []int, ws *WorkerState, r *[]int) {

	defer ws.DetectWg.Done()

	iBegin, iEnd := ws.VertexBegin, ws.VertexEnd
	var color int
	sg := ws.Subgraph

	for _, i := range u {
		v := &sg.Vertices[i]

		for _, j := range v.Adj {
			if j >= iBegin && j < iEnd {
				color = sg.Vertices[j-iBegin].Value
			} else {
				ws.StoredMutex.Lock()
				color = ws.Stored[j]
				ws.StoredMutex.Unlock()
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

	sg := ws.Subgraph
	buf := make([]byte, 8)

	// initialize U to be all of the vertices in sg
	u := make([]int, len(sg.Vertices))
	for i := 0; i < len(sg.Vertices); i++ {
		u[i] = i
	}
	r := make([]int, 0)

	// loop until u is empty (see break condition)
	for len(u) > 0 {
		// synchronizing the beginning of each step
		logger.Printf("Synchronizing start of round...\n")
		ws.ConnPool.BroadcastWorkers(graphnet.MSG_NODE_ROUND_START,
			[]byte{byte(ws.NodeIndex)})
		ws.StartWg.Wait()
		ws.StartWg.Add(ws.NodeCount - 2)

		// special care taken here due to potential race condition; see
		// documentation
		ws.ColorWgLock.Lock()
		ws.ColorWg.Add(ws.NodeCount - 2)
		ws.ColorWgLock.Unlock()

		logger.Printf("Beginning new round: %d vertices to be colored\n",
			len(u))

		// loop over vertices, assign each vertex a permissible color
		// send colors of boundary vertices to relevant nodes
		// receive color information from other nodes;
		// don't use supersteps, rather choose number of threads; channels
		// will be buffered anyways
		nVertices := len(u)
		verticesPerThread := int(math.Ceil(float64(nVertices) /
			float64(nThreads)))
		ws.ColorWg.Add(nThreads)
		for i := 0; i < nThreads; i++ {
			start := i * verticesPerThread
			end := (i + 1) * verticesPerThread
			if start > nVertices {
				start = nVertices
			}
			if end > nVertices {
				end = nVertices
			}

			go colorSpeculative(u[start:end], maxColor, ws)
		}

		// flush all write buffers
		ws.ConnPool.FlushAll()

		// when speculative coloring for this round is done, notify workers
		buf[0] = byte(ws.NodeIndex)
		ws.ConnPool.BroadcastWorkers(graphnet.MSG_NODE_ROUND_FINISHED, buf[:1])

		// wait until all incoming messages are successfully received;
		// this makes all steps work in lockstep
		ws.ColorWg.Wait()

		logger.Printf("Beginning conflict resolution stage\n")

		// for each boundary vertex, check for conflicts (in parallel)
		// add conflicting nodes to R
		ws.DetectWg.Add(nThreads)
		for i := 0; i < nThreads; i++ {
			start := i * verticesPerThread
			end := (i + 1) * verticesPerThread
			if start > nVertices {
				start = nVertices
			}
			if end > nVertices {
				end = nVertices
			}

			go resolveConflicts(u[start:end], ws, &r)
		}
		ws.DetectWg.Wait()

		// set U to R
		u = r
		r = u[:0]
	}

	// when done coloring, notify all nodes
	buf[0] = byte(ws.NodeIndex)
	ws.ConnPool.Broadcast(graphnet.MSG_NODE_FINISHED, buf[:1])
}
