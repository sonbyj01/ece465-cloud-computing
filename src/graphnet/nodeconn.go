package graphnet

import (
	"bufio"
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
		orderedPool[conn.index] = conn
	}

	*ncp = orderedPool
}

// NodeConn is a struct to keep track of a single connection from this node
// to another node
type NodeConn struct {
	channel chan string
	reader  *bufio.Reader
	writer  *bufio.Writer
	conn    net.Conn
	logger  *log.Logger
	index   int
}

// Read listens on the connection's socket and outputs messages to the
// connection channel
func (conn *NodeConn) Read() {
	for {
		line, err := conn.reader.ReadString('\n')

		// error if end of file reached
		if err != nil {
			break
		}

		conn.channel <- line
	}

	err := conn.conn.Close()
	if err != nil {
		conn.logger.Fatal(err)
	}
}

// Write sends messages sent to the channel over the network socket
func (conn *NodeConn) Write() {
	for data := range conn.channel {
		_, err := conn.writer.WriteString(data)
		if err != nil {
			conn.logger.Fatal(err)
		}

		err = conn.writer.Flush()
		if err != nil {
			conn.logger.Fatal(err)
		}
	}
}

// Channel is a getter for a connection's channel
func (conn *NodeConn) Channel() *chan string {
	return &conn.channel
}

// SetIndex sets the connection's node index
func (conn *NodeConn) SetIndex(index int) {
	conn.index = index
}

// NewNodeConn returns a new node connection object for sending messages
// to other nodes
func NewNodeConn(conn net.Conn, logger *log.Logger) *NodeConn {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	nodeConn := &NodeConn{
		channel: make(chan string),
		conn:    conn,
		reader:  reader,
		writer:  writer,
		logger:  logger,
	}

	// begin listening for reading and writing
	go nodeConn.Read()
	go nodeConn.Write()

	return nodeConn
}
