package node

import (
	"net"
)

// INode defines the interface for interacting with nodes in a Kademlia DHT network.
//
// The INode interface provides a standardized way for interacting with nodes,
// allowing for the customization and extension of node properties. This can
// include various implementations that may add additional attributes or behaviors
// to nodes, such as custom metadata or extended communication capabilities.
//
// Methods:
//
//   - ID() NodeID:
//     Returns the unique identifier of the node. This ID is used for routing and
//     distance calculations within the Kademlia network.
//
//   - Address() net.IP:
//     Returns the IP address of the node. The address is used for network communication
//     and can be either an IPv4 or IPv6 address.
//
//   - Port() uint16:
//     Returns the port number the node listens on. The port is used in conjunction with
//     the IP address to establish network connections.
//
//   - Distance(other *Node) [20]byte:
//     Calculates and returns the distance between the current node and another node using
//     the XOR metric. The distance is used to determine the proximity of nodes in the
//     keyspace, which is essential for routing and lookup operations in the Kademlia protocol.
//
// By implementing this interface, it is possible to create alternative versions of Node
// structures that include additional information or implement enhanced behaviors, while
// maintaining compatibility with the Kademlia routing and lookup logic.
type INode interface {
	ID() NodeID
	Address() net.IP
	Port() uint16
	Distance(other *Node) [20]byte
}
