package node

import (
	"net"
)

// INode defines the interface for interacting with nodes in a Kademlia DHT network.
//
// The INode interface provides a standard set of methods for interacting with nodes,
// enabling the abstraction of node properties such as identifiers, network addresses,
// and distance calculations. Implementations of this interface can include various
// node types, which may extend the functionality with additional attributes or behaviors,
// such as custom metadata or specialized communication capabilities.
//
// By implementing the INode interface, various node structures can be created to handle
// different use cases or extended behaviors, while maintaining compatibility with the
// Kademlia routing logic. This allows for flexibility and extensibility in how nodes
// are represented and interacted with in the distributed network.
type INode interface {
	// ID returns the unique identifier of the node, used for routing and distance calculations
	// within the Kademlia network.
	ID() NodeID

	// Address returns the IP address of the node, either IPv4 or IPv6, which is used for
	// establishing network communication.
	Address() net.IP

	// Port returns the port number the node listens on, used in conjunction with the IP address
	// to facilitate network connections.
	Port() uint16

	// Distance calculates and returns the XOR-based distance between the current node and another
	// node, which is essential for determining proximity and routing decisions in the Kademlia protocol.
	Distance(other *Node) [20]byte
}
