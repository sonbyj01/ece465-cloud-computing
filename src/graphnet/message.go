package graphnet

// this file the description of messages for internode communication
// over sockets

const (
	MSG_VERTEX_INFO         = byte(iota)
	MSG_NODE_FINISHED       = byte(iota)
	MSG_NODE_ROUND_FINISHED = byte(iota)


	// the following message types are for server-worker handshake
	MSG_NODE_INDEX_COUNT	= byte(iota)	// server notifies a worker node of
											// total node count and its index
	MSG_NODE_ADDRESS		= byte(iota)	// server notifies a worker node of
											// another work nodes' address
)

// mapping each message type to its number of bytes
var NUM_BYTES_MAP = map[byte]int{
	MSG_VERTEX_INFO: 8,				// 0-3: color, 4-7: index
	MSG_NODE_FINISHED: 1,			// 0: node index
	MSG_NODE_ROUND_FINISHED: 1,		// 0: node index
	MSG_NODE_INDEX_COUNT: 2,		// 0: node index, 1: total nodes
									// (including server)
	MSG_NODE_ADDRESS: 7,			// 0: node index, 1-4: ipv4 address,
									// 5-6: port
}

type VertexData struct {
	Vertices []int
	Colors   []int16
}

type VertexMessage struct {
	Type byte
	Data []VertexData
}
