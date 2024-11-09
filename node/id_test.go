package node_test

import (
	"crypto/sha1"
	"encoding/hex"
	"testing"

	"github.com/BochkovDev/kademlia-go/node"
)

// TestNewNodeID checks that NewNodeID generates the expected SHA-1 hash for a given input.
func TestNewNodeID(t *testing.T) {
	data := []byte("node_data")
	expectedHash := sha1.Sum(data)
	nodeID := node.NewNodeID(data)

	if nodeID != expectedHash {
		t.Errorf("NewNodeID failed, expected %x, got %x", expectedHash, nodeID)
	}
}

// TestString checks that the String method returns the correct hexadecimal representation.
func TestString(t *testing.T) {
	data := []byte("node_data")
	nodeID := node.NewNodeID(data)
	expectedStr := hex.EncodeToString(nodeID[:])

	if nodeID.String() != expectedStr {
		t.Errorf("String() failed, expected %s, got %s", expectedStr, nodeID.String())
	}
}

// TestEquals checks the Equals method by comparing identical and different NodeIDs.
func TestEquals(t *testing.T) {
	data1 := []byte("node_data_1")
	data2 := []byte("node_data_2")

	nodeID1 := node.NewNodeID(data1)
	sameNodeID := node.NewNodeID(data1)
	nodeID2 := node.NewNodeID(data2)

	if !nodeID1.Equals(sameNodeID) {
		t.Error("Equals() failed, expected nodeID1 to equal sameNodeID")
	}
	if nodeID1.Equals(nodeID2) {
		t.Error("Equals() failed, expected nodeID1 to not equal nodeID2")
	}
}

// TestXOR checks the XOR method, ensuring it performs a bitwise XOR correctly.
func TestXOR(t *testing.T) {
	data1 := []byte("node_data_1")
	data2 := []byte("node_data_2")

	nodeID1 := node.NewNodeID(data1)
	nodeID2 := node.NewNodeID(data2)

	xorResult := nodeID1.XOR(nodeID2)
	for i := 0; i < len(nodeID1); i++ {
		expected := nodeID1[i] ^ nodeID2[i]
		if xorResult[i] != expected {
			t.Errorf("XOR() failed at byte %d, expected %x, got %x", i, expected, xorResult[i])
		}
	}
}
