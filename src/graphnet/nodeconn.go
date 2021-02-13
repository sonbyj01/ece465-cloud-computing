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
type NodeConnPool []*NodeConn

// NewNodeConnPool generates a NodeConnPool (array of NodeConn)
func NewNodeConnPool() NodeConnPool {
	return make(NodeConnPool, 0)
}

// AddUnregistered adds a NodeConn to the nodeConnPool
func (ncp *NodeConnPool) AddUnregistered(conn *NodeConn) {
	*ncp = append(*ncp, conn)
}

// Register reorganizes the NodeConn structs into their proper index
// after they have all received their indices
func (ncp *NodeConnPool) Register() {
	orderedPool := make(NodeConnPool, len(*ncp)+1)

	for _, conn := range *ncp {
		orderedPool[conn.Index] = conn
	}

	*ncp = orderedPool
}

// Dispatch tells the Read function how many bytes to read and what to do with
// the bytes
type Dispatch struct {
	numBytes int
	callback func([]byte)
}

// NewDispatch creates a new dispatch
func NewDispatch(numBytes int, callback func([]byte)) Dispatch {
	return Dispatch{numBytes, callback}
}

// NodeConn is a struct to keep track of a single connection from this node
// to another node
type NodeConn struct {
	channel     chan string
	reader      *bufio.Reader
	writer      *bufio.Writer
	conn        net.Conn
	logger      *log.Logger
	Index       int
	dispatchTab map[byte]Dispatch
}

// Read listens on the connection's socket and outputs messages to the
// connection channel
func (conn *NodeConn) Read() {
	for {
		b, err := conn.reader.ReadByte()
		if err != nil {
			conn.logger.Println(err)
			break
		}

		// look up action in dispatch table
		dispatch := conn.dispatchTab[b]
		buf := make([]byte, dispatch.numBytes)
		n, err := io.ReadFull(conn.reader, buf)
		if n != dispatch.numBytes || err != nil {
			conn.logger.Println(err)
			break
		}

		// dispatch action
		go dispatch.callback(buf)
	}

	err := conn.conn.Close()
	if err != nil {
		conn.logger.Fatal(err)
	}
}

// TODO: remove
// Write sends messages sent to the channel over the network socket
//func (conn *NodeConn) Write() {
//	for data := range conn.channel {
//		_, err := conn.writer.WriteString(data)
//		if err != nil {
//			conn.logger.Fatal(err)
//		}
//
//		err = conn.writer.Flush()
//		if err != nil {
//			conn.logger.Fatal(err)
//		}
//	}
//}

// WriteBytes allows you to write messages directly to the socket
func (conn *NodeConn) WriteBytes(messageType byte, buffer []byte) {
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

// Channel is a getter for a connection's channel
func (conn *NodeConn) Channel() *chan string {
	return &conn.channel
}

// Close closes the NodeConn's connection
func (conn *NodeConn) Close() {
	err := conn.conn.Close()
	if err != nil {
		conn.logger.Fatal(err)
	}
}

// NewNodeConn returns a new node connection object for sending messages
// to other nodes
func NewNodeConn(conn net.Conn, logger *log.Logger,
	dispatchTab map[byte]Dispatch) *NodeConn {

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	nodeConn := &NodeConn{
		channel:     make(chan string),
		conn:        conn,
		reader:      reader,
		writer:      writer,
		logger:      logger,
		dispatchTab: dispatchTab,
	}

	// begin listening for reading and writing
	go nodeConn.Read()
	//go nodeConn.Write()

	return nodeConn
}
