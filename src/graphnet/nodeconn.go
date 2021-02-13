package graphnet

import (
	"bufio"
	"io"
	"log"
	"net"
)

// NodeConnPool keeps track of all node connections in an array; they are
// initially dumped in the unregistered list and moved to the correct index in
// the registered array when their index is received
type NodeConnPool struct {
	Conns      []*NodeConn
	Index      int
	registered bool
}

// NewNodeConnPool generates a NodeConnPool (array of NodeConn)
func NewNodeConnPool() NodeConnPool {
	return NodeConnPool{}
}

// AddUnregistered adds a NodeConn to the nodeConnPool
func (ncp *NodeConnPool) AddUnregistered(conn *NodeConn) {
	ncp.Conns = append(ncp.Conns, conn)
}

// Register reorganizes the NodeConn structs into their proper index
// after they have all received their indices
func (ncp *NodeConnPool) Register() {
	orderedPool := make([]*NodeConn, len(ncp.Conns)+1)
	var index int

	for _, conn := range ncp.Conns {
		orderedPool[conn.Index] = conn
	}

	// find which index the current node is (i.e., it should be the only
	// nil connection, all other connections should be properly filled)
	for i, conn := range orderedPool {
		if conn == nil {
			index = i
			break
		}
	}

	ncp.Conns = orderedPool
	ncp.Index = index
	ncp.registered = true
}

// Broadcast sends a message to all other (active) nodes
func (ncp *NodeConnPool) Broadcast(msgType byte, buf []byte) {
	// ncp should be registered first, i.e., indices should be correct
	if !ncp.registered {
		panic("Unregistered NodeConnPool")
	}

	for _, nodeConn := range ncp.Conns {
		if nodeConn != nil {
			nodeConn.WriteBytes(msgType, buf)
		}
	}
}

// BroadcastWorkers sends a message to all (other) (active) worker nodes
func (ncp *NodeConnPool) BroadcastWorkers(msgType byte, buf []byte) {
	// ncp should be registered first, i.e., indices should be correct
	if !ncp.registered {
		panic("Unregistered NodeConnPool")
	}

	for i, nodeConn := range ncp.Conns {
		if nodeConn != nil && i > 0 {
			nodeConn.WriteBytes(msgType, buf)
		}
	}
}

// Remove removes a node
func (ncp *NodeConnPool) Remove(nodeIndex int) {
	// ncp should be registered first, i.e., indices should be correct
	if !ncp.registered {
		panic("Unregistered NodeConnPool")
	}

	// close connection
	ncp.Conns[nodeIndex].Close()
	ncp.Conns[nodeIndex] = nil
}

// Dispatch is a callback that takes a fixed-length slice of bytes, and is
// associated with a particular message type. The length of the slice of bytes
// is defined in graphnet.NUM_BYTES_MAP
type Dispatch func([]byte, *NodeConn)

// NodeConn is a struct to keep track of a single connection from this node
// to another node
type NodeConn struct {
	reader      *bufio.Reader
	writer      *bufio.Writer
	conn        net.Conn
	logger      *log.Logger
	Index       int
	dispatchTab map[byte]Dispatch
	open		bool
}

// Read listens on the connection's socket and outputs messages to the
// connection channel
func (conn *NodeConn) Read() {
	for {
		b, err := conn.reader.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			conn.logger.Fatal(err)
		}

		// look up action in dispatch table
		buf := make([]byte, NUM_BYTES_MAP[b])
		n, err := io.ReadFull(conn.reader, buf)
		if err == io.EOF {
			break
		} else if n != NUM_BYTES_MAP[b] || err != nil {
			conn.logger.Fatal(err)
		}

		// dispatch action
		go conn.dispatchTab[b](buf, conn)
	}

	conn.Close()
}

// WriteBytes allows you to write messages directly to the socket
func (conn *NodeConn) WriteBytes(messageType byte, buffer []byte) {
	if !conn.open {
		return
	}

	err := conn.writer.WriteByte(messageType)
	if err != nil {
		conn.logger.Fatal(err)
	}

	n, err := conn.writer.Write(buffer)
	if n != len(buffer) || err != nil {
		conn.logger.Fatal(err)
	}

	err = conn.writer.Flush()
	if err != nil {
		conn.logger.Fatal(err)
	}
}

// Close closes the NodeConn's connection
func (conn *NodeConn) Close() {
	if !conn.open {
		return
	}

	err := conn.conn.Close()
	if err != nil {
		conn.logger.Fatal(err)
	}
	conn.open = false
}

// NewNodeConn returns a new node connection object for sending messages
// to other nodes
func NewNodeConn(conn net.Conn, logger *log.Logger,
	dispatchTab map[byte]Dispatch) *NodeConn {

	nodeConn := &NodeConn{
		conn:        conn,
		reader:      bufio.NewReader(conn),
		writer:      bufio.NewWriter(conn),
		logger:      logger,
		dispatchTab: dispatchTab,
		open:        true,
	}

	// begin listening for reading
	go nodeConn.Read()

	return nodeConn
}
