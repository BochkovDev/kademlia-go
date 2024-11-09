package routing

import "github.com/BochkovDev/kademlia-go/node"

// IKBucket defines the interface for managing a K-bucket in the Kademlia DHT routing table.
//
// The IKBucket interface provides methods for interacting with a K-bucket, which is a data structure
// used in the Kademlia protocol to store a subset of nodes based on proximity in the keyspace. A K-bucket
// stores nodes within a specific range and ensures efficient node management by keeping track of the closest
// nodes to a given ID. Each K-bucket has a fixed capacity (KSize) and follows the Kademlia protocol's rules
// for adding, removing, and evicting nodes.
//
// By implementing the IKBucket interface, a K-bucket can efficiently manage nodes in a Kademlia-based DHT,
// adhering to the protocol's requirements for proximity-based routing and node management.
type IKBucket interface {
	// Nodes returns a slice of nodes currently stored in the K-bucket, ordered by proximity
	// to the bucket's range in the keyspace. The most recently active nodes are placed at the end.
	Nodes() []*node.INode

	// KSize returns the maximum number of nodes that the K-bucket can hold.
	// This value is typically a fixed constant.
	KSize() uint8

	// Add inserts a new node into the K-bucket. If the K-bucket is full,
	// the least recently seen node may be evicted to make room for the new node.
	Add(newNode *node.INode)

	// Remove deletes a node from the K-bucket using its NodeID.
	// This is typically used to remove unreachable or outdated nodes.
	Remove(id node.NodeID)

	// Contains checks if a node with the given NodeID is present in the K-bucket.
	// Returns true if the node exists, otherwise false.
	Contains(id node.NodeID) bool

	// IsFull returns true if the K-bucket has reached its maximum capacity (KSize).
	// If true, no additional nodes can be added until space is freed.
	IsFull() bool

	// Size returns the current number of nodes in the K-bucket.
	// This helps monitor how many nodes are actively being tracked.
	Size() uint8

	// Clear removes all nodes from the K-bucket, effectively resetting it.
	// This can be useful for maintenance or reinitialization purposes.
	Clear()
}
