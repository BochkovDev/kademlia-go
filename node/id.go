package node

import (
	"crypto/sha1"
	"encoding/hex"
)

// NodeID represents a unique identifier for a node in the Kademlia DHT network.
// Each NodeID consists of a 160-bit (20-byte) value.
//
// References:
//   - Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric" [Section 1, "Introduction"; 2.1, "XOR Metric"].
//     Retrieved from: https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf
type NodeID [20]byte

// NewNodeID generates a unique NodeID from a given input byte slice by applying the SHA-1 hashing algorithm.
//
// This function returns a 160-bit hash, which matches the required NodeID size for Kademlia's DHT.
// SHA-1 was selected for its ability to produce uniformly distributed identifiers, a property essential
// for maintaining balanced distribution and efficient lookup performance in Kademlia networks.
//
// Parameters:
//   - data []byte: The input byte slice used to generate the NodeID. This input should contain unique
//     information about the node, such as its IP address and port, to ensure a unique identifier.
//
// References:
//   - Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric" [Section 1, "Introduction"].
//     Retrieved from: https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf
func NewNodeID(data []byte) NodeID {
	hash := sha1.Sum(data)
	return hash
}

// String converts a NodeID into its hexadecimal string representation for easy human-readable display.
func (id NodeID) String() string {
	return hex.EncodeToString(id[:])
}

// Equals compares the current NodeID with another NodeID for equality.
func (id NodeID) Equals(other NodeID) bool {
	return id == other
}

// XOR performs a bitwise XOR operation between the current NodeID and another NodeID.
//
// This operation is fundamental to calculating the XOR distance between nodes, a metric used in Kademlia
// to determine routing paths and proximity of nodes in the network. The XOR metric ensures efficient
// lookups by enabling distance-based routing.
//
// References:
//   - Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric" [Section 2.1, "XOR Metric"].
//     Retrieved from: https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf
func (id NodeID) XOR(other NodeID) [20]byte {
	var result NodeID
	for i := 0; i < len(id); i++ {
		result[i] = id[i] ^ other[i]
	}
	return result
}
