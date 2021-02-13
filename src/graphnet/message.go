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

type VertexData struct {
	Vertices []int
	Colors   []int16
}

type VertexMessage struct {
	Type byte
	Data []VertexData
}
