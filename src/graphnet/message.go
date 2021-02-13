package graphnet

// this file the description of messages for internode communication
// over sockets

const (
	MSG_VERTEX_INFO         = byte(iota)
	MSG_NODE_FINISHED       = byte(iota)
	MSG_NODE_ROUND_FINISHED = byte(iota)
)

type VertexData struct {
	Vertices []int
	Colors   []int16
}

type VertexMessage struct {
	Type byte
	Data []VertexData
}
