package routing

import (
	"sync"

	"github.com/BochkovDev/kademlia-go/node"
)

// KBucket represents a container for nodes within the Kademlia network.
//
// A KBucket is a segment of Kademlia's routing table that maintains a list of nodes
// within a specific distance range from the current node. Each KBucket has a limited
// size and follows a Least Recently Seen (LRS) eviction policy, where the oldest nodes
// are removed to make space for newly added or recently active nodes.
//
// References:
//   - [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"] [Section 2.2, "Node State"]
//     https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf
type KBucket struct {
	// nodes is a slice of nodes stored in this KBucket. These nodes represent peers at a specific
	// distance range from the current node. The slice maintains nodes in order of
	// activity, with the most recently seen node positioned at the end.
	nodes []node.INode

	// ksize is the maximum number of nodes that the KBucket can contain. If this limit is reached
	// when adding a new node, the oldest node is evicted to make room for the new node.
	ksize uint8

	// mu is a mutex used to synchronize access to the nodes slice, ensuring that all operations
	// on the KBucket are thread-safe in concurrent environments.
	mu sync.Mutex
}

// NewKBucket creates and returns a new KBucket instance with a specified capacity for storing nodes.
//
// Parameters:
//   - ksize uint8: The maximum number of nodes (K) that this KBucket can hold, based on
//     the Kademlia protocol's specifications.
//
// Returns:
//   - *KBucket: A pointer to a newly created KBucket, initialized with an empty node list
//     and a mutex for thread safety.
func NewKBucket(ksize uint8) *KBucket {
	return &KBucket{
		nodes: make([]node.INode, 0, ksize),
		ksize: ksize,
		mu:    sync.Mutex{},
	}
}

// Nodes returns a slice of nodes stored in the KBucket.
//
// This method provides access to the nodes contained within the KBucket, representing peers
// at a specific distance from the current node. The nodes are ordered by their last-seen time.
//
// Returns:
//   - []node.INode: A slice of nodes currently stored in the KBucket.
func (kb *KBucket) Nodes() []node.INode {
	return kb.nodes
}

// KSize returns the maximum number of nodes that the KBucket can hold.
//
// This method provides the maximum capacity of the KBucket, which is fixed and determined
// during initialization. The capacity is used to manage node evictions when the KBucket is full.
//
// Returns:
//   - uint8: The maximum number of nodes that can be stored in the KBucket.
func (kb *KBucket) KSize() uint8 {
	return kb.ksize
}

// Add inserts a new node into the KBucket.
//
// If the node already exists, it is removed from its current position and re-added to the end
// of the list to reflect its recent activity. If the KBucket is full and does not contain the new node,
// the oldest node (at the beginning) is removed to make space.
//
// Parameters:
//   - newNode node.INode: The node to be added to the KBucket.
//
// Notes:
//   - This method uses a mutex to ensure thread safety while modifying the list of nodes.
//
// References:
//   - [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"] [Section 2.2, "Node State"]//
//
// [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"]: https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf
func (kb *KBucket) Add(newNode node.INode) {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	for i, n := range kb.nodes {
		if n.ID().Equals(newNode.ID()) {
			kb.nodes = append(kb.nodes[:i], kb.nodes[i+1:]...)
			kb.nodes = append(kb.nodes, newNode)
			return
		}
	}

	if len(kb.nodes) >= int(kb.ksize) {
		kb.nodes = kb.nodes[1:]
	}

	kb.nodes = append(kb.nodes, newNode)
}

// Remove deletes a node from the KBucket based on its NodeID.
//
// This method searches for the node with the specified NodeID in the KBucket. If found,
// it removes the node from the list, maintaining the order of remaining nodes.
//
// Parameters:
//   - id node.NodeID: The NodeID of the node to be removed from the KBucket.
//
// Notes:
//   - This method uses a mutex to ensure thread safety while modifying the list of nodes.
func (kb *KBucket) Remove(id node.NodeID) {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	for i, n := range kb.nodes {
		if n.ID().Equals(id) {
			kb.nodes = append(kb.nodes[:i], kb.nodes[i+1:]...)
			return
		}
	}
}

// Contains checks if a node with a specific NodeID is present in the KBucket.
//
// This method iterates over the nodes stored in the KBucket and compares the NodeID of each node with
// the provided NodeID. If a match is found, it returns true, indicating that the node is present in the
// KBucket. Otherwise, it returns false.
//
// Parameters:
//   - id node.NodeID: The NodeID of the node to be checked for presence in the KBucket.
//
// Returns:
//   - bool: Returns true if the node with the specified NodeID exists in the KBucket, false otherwise.
//
// Notes:
//   - This method uses a mutex to ensure thread safety while modifying the list of nodes.
func (kb *KBucket) Contains(id node.NodeID) bool {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	for _, n := range kb.nodes {
		if n.ID().Equals(id) {
			return true
		}
	}
	return false
}

// IsFull checks whether the KBucket has reached its maximum capacity.
//
// This method compares the current number of nodes in the KBucket with the MaxSize and returns true if the
// KBucket has reached or exceeded its limit, indicating that no more nodes can be added without eviction.
// If the KBucket has space for more nodes, it returns false.
//
// Returns:
//   - bool: True if the KBucket is full, false otherwise.
//
// Notes:
//   - This method uses a mutex to ensure thread safety while checking the size of the node list.
func (kb *KBucket) IsFull() bool {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	return len(kb.nodes) >= int(kb.ksize)
}

// Size returns the current number of nodes in the KBucket.
//
// This method calculates and returns the number of nodes currently stored in the KBucket. It is useful for
// monitoring the number of nodes in the bucket and making decisions based on the current size. The size is
// returned as an unsigned 8-bit integer.
//
// Returns:
//   - uint8: The number of nodes currently stored in the KBucket.
//
// Notes:
//   - This method uses a mutex to ensure thread safety while accessing the node list.
func (kb *KBucket) Size() uint8 {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	return uint8(len(kb.nodes))
}

// Clear empties the KBucket by removing all nodes.
//
// This method sets the list of nodes in the KBucket to nil, effectively clearing the KBucket. It ensures that
// all nodes are removed from the KBucket, and the list is reset to its initial state. This is useful when
// a fresh set of nodes needs to be loaded or when the KBucket needs to be reset.
//
// Notes:
//   - This method uses a mutex to ensure thread safety during the operation.
func (kb *KBucket) Clear() {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	kb.nodes = nil
}
