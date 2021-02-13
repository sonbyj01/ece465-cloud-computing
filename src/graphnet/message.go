package graphnet

// this file the description of messages for internode communication
// over sockets

const (
	MSG_VERTEX_INFO         = byte(iota)
	MSG_NODE_FINISHED       = byte(iota)
	MSG_NODE_ROUND_FINISHED = byte(iota)


	// Following message types are for server-worker handshake

	// server notifies a worker node of total number of nodes and its index
	MSG_NODE_INDEX_COUNT	= byte(iota)

	// server notifies a worker node of another work nodes' address
	MSG_NODE_IP				= byte(iota)

	// server is finished with handshake
	MSG_HANDSHAKE_DONE		= byte(iota)
)

// mapping each message type to its number of bytes
var NUM_BYTES_MAP = map[byte]int{
	MSG_VERTEX_INFO: 0,
	MSG_NODE_FINISHED: 0,
	MSG_NODE_ROUND_FINISHED: 0,
	MSG_NODE_INDEX_COUNT: 0,
	MSG_NODE_IP: 0,
	MSG_HANDSHAKE_DONE: 0,
}

type VertexData struct {
	Vertices []int
	Colors   []int16
}

type VertexMessage struct {
	Type byte
	Data []VertexData
}
