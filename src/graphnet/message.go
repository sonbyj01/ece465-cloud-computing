package graphnet

// this file the description of messages for internode communication
// over sockets

const (
	// MSG_ACK is a generic message used for synchronization
	MSG_ACK = byte(iota)

	// messages for communication during coloring

	// MSG_VERTEX_INFO for exchanging neighbor vertex info
	MSG_VERTEX_INFO = byte(iota)
	// MSG_NODE_FINISHED when node completely finished coloring
	MSG_NODE_FINISHED = byte(iota)
	// MSG_NODE_ROUND_FINISHED when node finished one coloring round
	MSG_NODE_ROUND_FINISHED = byte(iota)

	// messages for server-worker handshake

	// MSG_NODE_INDEX_COUNT server gives worker index & total nodes
	MSG_NODE_INDEX_COUNT = byte(iota)
	// MSG_NODE_ADDRESS server gives worker connection info
	MSG_NODE_ADDRESS = byte(iota)
	// MSG_DIALER_INDEX dialer worker tells dialee worker its node index
	MSG_DIALER_INDEX = byte(iota)

	// messages for sending subgraph

	// MSG_SUBGRAPH is for sending subgraph
	MSG_SUBGRAPH = byte(iota)

	// MSG_CONT indicates not to send a message type, this buffer is a
	// continuation of the last byte buffer
	MSG_CONT = byte(255)
)

// NUM_BYTES_MAP maps each message type to its number of bytes; -1 indicates
// reading arbitrary-length data as string until DELIM_EOF is found
var NUM_BYTES_MAP = map[byte]int{
	MSG_ACK:                 1,  // 0: node index
	MSG_VERTEX_INFO:         8,  // 0-3: color, 4-7: index
	MSG_NODE_FINISHED:       1,  // 0: node index
	MSG_NODE_ROUND_FINISHED: 1,  // 0: node index
	MSG_NODE_INDEX_COUNT:    2,  // 0: node index, 1: total nodes including serv
	MSG_NODE_ADDRESS:        7,  // 0: node index, 1-4: ipv4 address, 5-6 port
	MSG_DIALER_INDEX:        1,  // 0: incoming node index
	MSG_SUBGRAPH:            -1, // variable length string until DELIM_EOF
}

// DELIM_EOF is used to indicate end of string; null byte is used arbitrarily
const DELIM_EOF = byte(0)

// VertexData is used to send a buffer of vertex data
// TODO: use this and buffer vertex info; currently only one vertex is
// 		sent at a time
// TODO: make this send a fixed-size array rather than a slice
// 		(is this more efficient?)
type VertexData struct {
	Vertices []int
	Colors   []int16
}
