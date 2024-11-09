package node

import (
	"net"
)

// Node represents a node in the Kademlia DHT network.
//
// Each node is identified by a unique identifier (ID) and can have
// an IP address through which it is accessible, along with a port for establishing connections.
// It also stores the last seen time of the node.
// This struct is used to store and share information about other nodes in the Kademlia network.
//
// Fields:
//
//   - ID NodeID:
//     Unique identifier of the node in the Kademlia network. This ID is
//     computed based on a hash, such as the IP address and other data.
//     It is used for sorting nodes and finding the closest nodes to the current node.
//
//   - Address net.IP:
//     The IP address of the node, which can be used for establishing connections.
//     It can be either an IPv4 or IPv6 address, depending on the network configuration.
//
//   - Port uint16:
//     The port the node is listening on for incoming connections. The port must be in the range
//     0-65535. It is used for connections over TCP or UDP for data exchange in the Kademlia network.
//
// References:
//   - [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"] [Section 2.2, "Node State"]
//
// [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"]: https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf
type Node struct {
	id      NodeID
	address net.IP
	port    uint16
}

// NewNode creates and returns a new Node instance with a unique NodeID,
// based on the provided data, IP address, and port.
//
// Parameters:
//   - data []byte: Byte slice input used to generate the NodeID, typically based on unique information like IP and port.
//   - address net.IP: The IP address of the node, which can be IPv4 or IPv6, specifying its network location.
//   - port uint16: The port number the node listens on, used to facilitate network communication.
//
// Returns:
//   - *Node: A pointer to a newly created Node, with its ID, address, and port initialized.
func NewNode(data []byte, address net.IP, port uint16) *Node {
	return &Node{
		id:      NewNodeID(data),
		address: address,
		port:    port,
	}
}

// ID returns the NodeID of the current node.
//
// The NodeID is a unique identifier generated from the node's relevant data,
// such as its IP address and other parameters, used for sorting and determining
// proximity to other nodes in the Kademlia network.
//
// Returns:
//   - NodeID: The unique identifier of the node.
func (n *Node) ID() NodeID {
	return n.id
}

// Address returns the IP address of the current node.
//
// The IP address is used for network communication and can be either IPv4 or IPv6.
// This address is necessary for establishing connections with other nodes in the network.
//
// Returns:
//   - net.IP: The IP address of the node.
func (n *Node) Address() net.IP {
	return n.address
}

// Port returns the port number the current node is listening on.
//
// The port is used for establishing network connections, either over TCP or UDP,
// and must be within the valid range (0-65535) to ensure proper communication.
//
// Returns:
//   - uint16: The port number on which the node is listening.
func (n *Node) Port() uint16 {
	return n.port
}

// Distance calculates the distance between the current node and another node in the Kademlia DHT.
//
// The distance is determined using the XOR metric, which is applied between the NodeIDs of
// the current node and the other node. The result is a 160-bit value that represents the
// proximity or distance between the nodes in the Kademlia keyspace.
//
// The smaller the result, the closer the nodes are in the network.
//
// References:
//   - [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"] [Section 2.1, "XOR Metric"]
//
// [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"]: https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf
func (n *Node) Distance(other *Node) [20]byte {
	return n.ID().XOR(other.ID())
}
