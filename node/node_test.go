package node_test

import (
	"net"
	"testing"

	"github.com/BochkovDev/kademlia-go/node"
)

// TestNewNode checks that NewNode creates a Node with the expected NodeID, address, and port.
func TestNewNode(t *testing.T) {
	data := []byte("test_node_data")
	address := net.ParseIP("192.168.1.1")
	port := uint16(8080)

	testNode := node.NewNode(data, address, port)

	expectedID := node.NewNodeID(data)
	if testNode.ID != expectedID {
		t.Errorf("node.NewNode failed, expected ID %x, got %x", expectedID, testNode.ID)
	}

	if !testNode.Address.Equal(address) {
		t.Errorf("node.NewNode failed, expected Address %s, got %s", address, testNode.Address)
	}

	if testNode.Port != port {
		t.Errorf("node.NewNode failed, expected Port %d, got %d", port, testNode.Port)
	}
}

// TestDistance checks the Distance method by verifying the XOR calculation between two NodeIDs.
func TestDistance(t *testing.T) {
	data1 := []byte("node_1")
	data2 := []byte("node_2")
	address1 := net.ParseIP("192.168.1.1")
	address2 := net.ParseIP("192.168.1.2")
	port := uint16(8080)

	node1 := node.NewNode(data1, address1, port)
	node2 := node.NewNode(data2, address2, port)

	expectedDistance := node1.ID.XOR(node2.ID)
	calculatedDistance := node1.Distance(node2)

	if calculatedDistance != expectedDistance {
		t.Errorf("Distance failed, expected %x, got %x", expectedDistance, calculatedDistance)
	}
}

// TestSameNodeDistance checks that the distance between a node and itself is zero.
func TestSameNodeDistance(t *testing.T) {
	data := []byte("same_node_data")
	address := net.ParseIP("192.168.1.1")
	port := uint16(8080)

	testNode := node.NewNode(data, address, port)
	distance := testNode.Distance(testNode)

	var zeroDistance [20]byte
	if distance != zeroDistance {
		t.Errorf("Distance failed, expected zero distance, got %x", distance)
	}
}
