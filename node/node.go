package node

import (
	"net"
)

// Node represents a node in the Kademlia DHT network.
//
// Each Node is identified by a unique identifier (ID), and it is associated
// with an IP address and port for establishing network connections. The Node
// struct is fundamental in the Kademlia protocol, storing the necessary information
// for routing, communication, and maintaining a decentralized distributed hash table (DHT).
//
// References:
//   - [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"] [Section 2.2, "Node State"]
//
// [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"]: https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf
type Node struct {
	id      NodeID // Unique identifier for the node.
	address net.IP // IP address for network communication (IPv4 or IPv6).
	port    uint16 // Port for listening to incoming connections (range 0-65535).
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
