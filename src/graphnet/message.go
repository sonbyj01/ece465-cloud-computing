package graphnet

// this file the description of messages for internode communication
// over sockets

const (
	// communication during coloring

	// MSG_VERTEX_INFO for exchanging neighbor vertex info
	MSG_VERTEX_INFO         = byte(iota)
	// MSG_NODE_FINISHED when node completely finished coloring
	MSG_NODE_FINISHED       = byte(iota)
	// MSG_NODE_ROUND_FINISHED when node finished one coloring round
	MSG_NODE_ROUND_FINISHED = byte(iota)

	// the following message types are for server-worker handshake

	// MSG_NODE_INDEX_COUNT server gives worker index & total nodes
	MSG_NODE_INDEX_COUNT = byte(iota)
	// MSG_NODE_ADDRESS server gives worker connection info
	MSG_NODE_ADDRESS     = byte(iota)
	// MSG_DIALER_INDEX dialer worker tells dialee worker its node index
	MSG_DIALER_INDEX     = byte(iota)
)

// NUM_BYTES_MAP maps each message type to its number of bytes
var NUM_BYTES_MAP = map[byte]int{
	MSG_VERTEX_INFO:         8, // 0-3: color, 4-7: index
	MSG_NODE_FINISHED:       1, // 0: node index
	MSG_NODE_ROUND_FINISHED: 1, // 0: node index
	MSG_NODE_INDEX_COUNT:    2, // 0: node index, 1: total nodes including serv
	MSG_NODE_ADDRESS:        7, // 0: node index, 1-4: ipv4 address, 5-6 port
	MSG_DIALER_INDEX:        1, // 0: incoming node index
}

// VertexData is used to send a buffer of vertex data
// TODO: use this and buffer vertex info; currently only one vertex is
// 		sent at a time
// TODO: make this send a fixed-size array rather than a slice
// 		(is this more efficient?)
type VertexData struct {
	Vertices []int
	Colors   []int16
}
