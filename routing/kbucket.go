package routing

import (
	"sync"

	"github.com/BochkovDev/kademlia-go/node"
)

// KBucket represents a container for nodes within the Kademlia network.
//
// A KBucket is part of Kademlia's routing table that maintains a list of nodes at a specific
// distance range from the current node. It is used to keep track of neighboring nodes and organize
// searches for the closest nodes in the network. Each KBucket has a limited size and may evict the
// oldest nodes when the maximum capacity is reached.
//
// Fields:
//   - Nodes []*node.Node:
//     A slice of pointers to nodes stored in this KBucket. The nodes represent peers at a particular
//     distance from the current node. Nodes are maintained in order of insertion, with newly added
//     nodes appended to the end, and active nodes moved to the end upon interaction.
//   - MaxSize uint8:
//     The maximum number of nodes that the KBucket can contain. If this limit is reached when adding a new
//     node, the oldest node is removed to make space.
//   - mu sync.Mutex:
//     A mutex used to protect access to the Nodes slice and other fields, ensuring that modifications to
//     the KBucket are thread-safe. This is crucial in a distributed network environment where multiple
//     goroutines may attempt to add, remove, or access nodes concurrently.
//
// References:
//   - [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"] [Section 2.2, "Node State"]
//
// [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"]: https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf
type KBucket struct {
	Nodes   []*node.Node
	MaxSize uint8
	mu      sync.Mutex
}

// Add inserts a new node into the KBucket.
//
// If the node already exists, it is removed from its current
// position and added to the end of the list to reflect its recent activity. If the KBucket is full and
// does not contain the new node, the oldest node (at the beginning) is removed to make space.
//
// This method ensures that the KBucket maintains a list of nodes sorted by their last-seen time, with the
// most recently seen node at the end of the list. Such behavior is crucial for maintaining efficient lookups
// and evicting inactive nodes.
//
// Parameters:
//   - newNode *node.Node: The node to be added to the KBucket.
//
// Notes:
//   - This method uses a mutex to ensure thread safety while modifying the list of nodes.
//
// References:
//   - [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"] [Section 2.2, "Node State"]
//
// [Maymounkov, Petar; Mazieres, David. "Kademlia: A Peer-to-peer Information System Based on the XOR Metric"]: https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf
func (kb *KBucket) Add(newNode *node.Node) {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	for i, n := range kb.Nodes {
		if n.ID.Equals(newNode.ID) {
			kb.Nodes = append(kb.Nodes[:i], kb.Nodes[i+1:]...)
			kb.Nodes = append(kb.Nodes, newNode)
			return
		}
	}

	if len(kb.Nodes) >= int(kb.MaxSize) {
		kb.Nodes = kb.Nodes[1:]
	}

	kb.Nodes = append(kb.Nodes, newNode)
}

// Remove deletes a node from the KBucket based on its NodeID.
//
// This method searches for the node with the specified NodeID in the KBucket. If the node is found,
// it is removed from the list of nodes. The removal operation maintains the integrity of the list, ensuring
// the order of nodes is preserved.
//
// Parameters:
//   - id node.NodeID: The NodeID of the node to be removed from the KBucket.
//
// Notes:
//   - This method uses a mutex to ensure thread safety while modifying the list of nodes.
func (kb *KBucket) Remove(id node.NodeID) {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	for i, n := range kb.Nodes {
		if n.ID.Equals(id) {
			kb.Nodes = append(kb.Nodes[:i], kb.Nodes[i+1:]...)
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

	for _, n := range kb.Nodes {
		if n.ID.Equals(id) {
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

	return len(kb.Nodes) >= int(kb.MaxSize)
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

	return uint8(len(kb.Nodes))
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

	kb.Nodes = nil
}
