package graphnet

// this file the description of messages for internode communication
// over sockets

type VertexMessageType int8

const (
	MSG_VERTEX_INFO         VertexMessageType = iota
	MSG_NODE_FINISHED       VertexMessageType = iota
	MSG_NODE_ROUND_FINISHED VertexMessageType = iota
)

type VertexData struct {
	Vertices []int
	Colors   []int16
}

type VertexMessage struct {
	Type VertexMessageType
	Data VertexData
}
