package routing

import "github.com/BochkovDev/kademlia-go/node"

// IKBucket defines the interface for managing a K-bucket in the Kademlia DHT routing table.
//
// The IKBucket interface provides a set of methods for interacting with a K-bucket,
// a data structure that stores a subset of nodes in the network, based on the Kademlia protocol's
// routing and proximity rules. Each K-bucket holds a limited number of nodes (defined by KSize),
// allowing efficient routing and node management by maintaining the closest nodes to a given ID.
//
// Methods:
//
//   - Nodes() []*node.INode:
//     Returns a slice of nodes currently stored in the K-bucket. This list includes
//     the nodes ordered by proximity to the bucket's range in the keyspace.
//
//   - KSize() uint8:
//     Returns the maximum capacity of nodes that can be stored in the K-bucket. This is typically
//     a fixed value, ensuring consistency across buckets and enabling efficient routing.
//
//   - Add(newNode *node.INode):
//     Adds a new node to the K-bucket. If the bucket is full, this method may replace the least
//     recently seen node depending on the protocol's eviction policy.
//
//   - Remove(id node.NodeID):
//     Removes a node from the K-bucket based on its unique identifier. This is used to
//     discard unreachable or outdated nodes, maintaining the relevance of nodes in the bucket.
//
//   - Contains(id node.NodeID) bool:
//     Checks if a node with the given identifier exists in the K-bucket. This is helpful for
//     avoiding duplicate entries and quickly locating nodes within the bucket.
//
//   - IsFull() bool:
//     Returns true if the K-bucket has reached its capacity (KSize) and can no longer
//     accept new nodes. This helps enforce the K-bucket size limitation.
//
//   - Size() uint8:
//     Returns the current number of nodes in the K-bucket, providing an overview of how
//     many active nodes are currently being tracked in this bucket.
//
//   - Clear():
//     Clears all nodes from the K-bucket, effectively resetting it. This can be useful
//     for maintenance or reinitialization purposes.
//
// By implementing the IKBucket interface, it is possible to manage nodes efficiently
// within a Kademlia-based DHT, following the protocol’s rules for storing and
// retrieving nodes based on proximity and reachability.
type IKBucket interface {
	Nodes() []*node.INode
	KSize() uint8
	Add(newNode *node.INode)
	Remove(id node.NodeID)
	Contains(id node.NodeID) bool
	IsFull() bool
	Size() uint8
	Clear()
}
