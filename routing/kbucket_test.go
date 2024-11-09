package routing_test

import (
	"sync"
	"testing"

	"github.com/BochkovDev/kademlia-go/node"
	"github.com/BochkovDev/kademlia-go/routing"

	"github.com/stretchr/testify/suite"
)

// KBucketTestSuite defines the test structure for KBucket tests
type KBucketTestSuite struct {
	suite.Suite
	kb    *routing.KBucket
	node1 *node.Node
	node2 *node.Node
	node3 *node.Node
	node4 *node.Node
}

// SetupTest initializes the necessary data before each test
func (suite *KBucketTestSuite) SetupTest() {
	suite.node1 = &node.Node{ID: node.NewNodeID([]byte("node_1"))}
	suite.node2 = &node.Node{ID: node.NewNodeID([]byte("node_2"))}
	suite.node3 = &node.Node{ID: node.NewNodeID([]byte("node_3"))}
	suite.node4 = &node.Node{ID: node.NewNodeID([]byte("node_4"))}
	suite.kb = &routing.KBucket{MaxSize: 3}
}

// TearDownTest clears the KBucket after each test
func (suite *KBucketTestSuite) TearDownTest() {
	suite.kb.Clear()
}

// TestAddNode tests adding nodes to the KBucket
func (suite *KBucketTestSuite) TestAddNode() {
	suite.kb.Add(suite.node1)
	suite.kb.Add(suite.node2)

	suite.Equal(uint8(2), suite.kb.Size(), "Expected KBucket size to be 2")
	suite.True(suite.kb.Contains(suite.node1.ID), "KBucket should contain node1")
	suite.True(suite.kb.Contains(suite.node2.ID), "KBucket should contain node2")
}

// TestAddNodeEviction tests adding a new node when the KBucket is full
func (suite *KBucketTestSuite) TestAddNodeEviction() {
	suite.kb.Add(suite.node1)
	suite.kb.Add(suite.node2)
	suite.kb.Add(suite.node3)

	// The KBucket is now full with 3 nodes.
	suite.Equal(uint8(3), suite.kb.Size(), "Expected KBucket size to be 3")

	// Adding a new node should evict the oldest one (node1).
	suite.kb.Add(suite.node4)
	suite.Equal(uint8(3), suite.kb.Size(), "Expected KBucket size to be 3 after eviction")
	suite.False(suite.kb.Contains(suite.node1.ID), "KBucket should not contain node1 after eviction")
	suite.True(suite.kb.Contains(suite.node2.ID), "KBucket should contain node2")
	suite.True(suite.kb.Contains(suite.node3.ID), "KBucket should contain node3")
	suite.True(suite.kb.Contains(suite.node4.ID), "KBucket should contain node4")
}

// TestRemoveNode tests removing nodes from the KBucket
func (suite *KBucketTestSuite) TestRemoveNode() {
	suite.kb.Add(suite.node1)
	suite.kb.Add(suite.node2)

	suite.True(suite.kb.Contains(suite.node1.ID), "KBucket should contain node1")
	suite.kb.Remove(suite.node1.ID)
	suite.False(suite.kb.Contains(suite.node1.ID), "KBucket should not contain node1 after removal")
}

// TestContainsNode tests if a node is contained in the KBucket
func (suite *KBucketTestSuite) TestContainsNode() {
	suite.kb.Add(suite.node1)

	suite.True(suite.kb.Contains(suite.node1.ID), "KBucket should contain node1")
	suite.False(suite.kb.Contains(suite.node2.ID), "KBucket should not contain node2")
}

// TestIsFull tests checking if the KBucket is full
func (suite *KBucketTestSuite) TestIsFull() {
	suite.kb.Add(suite.node1)
	suite.kb.Add(suite.node2)

	suite.False(suite.kb.IsFull(), "KBucket should not be full")

	suite.kb.Add(suite.node3)
	suite.True(suite.kb.IsFull(), "KBucket should be full after adding 3 nodes")
}

// TestSize tests getting the size of the KBucket
func (suite *KBucketTestSuite) TestSize() {
	suite.kb.Add(suite.node1)
	suite.kb.Add(suite.node2)

	suite.Equal(uint8(2), suite.kb.Size(), "Expected KBucket size to be 2")
}

// TestClear tests clearing the KBucket
func (suite *KBucketTestSuite) TestClear() {
	suite.kb.Add(suite.node1)
	suite.kb.Add(suite.node2)

	suite.Equal(uint8(2), suite.kb.Size(), "Expected KBucket size to be 2 before clearing")

	suite.kb.Clear()

	suite.Equal(uint8(0), suite.kb.Size(), "Expected KBucket size to be 0 after clearing")
	suite.False(suite.kb.Contains(suite.node1.ID), "KBucket should not contain node1 after clearing")
	suite.False(suite.kb.Contains(suite.node2.ID), "KBucket should not contain node2 after clearing")
}

// TestKBucketTestSuite runs the test suite
func TestKBucketTestSuite(t *testing.T) {
	suite.Run(t, new(KBucketTestSuite))
}

func TestKBucketAddConcurrency(t *testing.T) {
	kb := &routing.KBucket{
		Nodes:   []*node.Node{},
		MaxSize: 10,
	}

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			kb.Add(&node.Node{ID: node.NewNodeID([]byte{byte(i)})})
		}(i)
	}

	wg.Wait()

	if kb.Size() > kb.MaxSize {
		t.Errorf("KBucket contains more than MaxSize nodes. Size: %d", kb.Size())
	}
}

// TestKBucketAddRemoveConcurrency verifies the KBucket's behavior when adding and removing nodes concurrently.
func TestKBucketAddRemoveConcurrency(t *testing.T) {
	kb := &routing.KBucket{
		Nodes:   []*node.Node{},
		MaxSize: 10,
	}

	var wg sync.WaitGroup
	numOperations := 100

	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			nodeID := node.NewNodeID([]byte{byte(i % 256)})
			kb.Add(&node.Node{ID: nodeID})
			kb.Remove(nodeID)
		}(i)
	}

	wg.Wait()

	if kb.Size() > uint8(kb.MaxSize) {
		t.Errorf("KBucket size exceeds MaxSize: got %d, expected <= %d", kb.Size(), kb.MaxSize)
	}
}

func TestKBucketConcurrentClear(t *testing.T) {
	kb := &routing.KBucket{
		Nodes:   []*node.Node{},
		MaxSize: 10,
	}

	for i := 0; i < 10; i++ {
		kb.Add(&node.Node{ID: node.NewNodeID([]byte{byte(i)})})
	}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			kb.Clear()
		}()
	}

	wg.Wait()

	if kb.Size() != 0 {
		t.Errorf("KBucket should be empty after concurrent clears, but it contains %d nodes", kb.Size())
	}
}

// TestKBucketConcurrentContains tests the Contains method with concurrent access.
func TestKBucketConcurrentContains(t *testing.T) {
	kb := &routing.KBucket{
		Nodes:   []*node.Node{},
		MaxSize: 10,
	}

	for i := 0; i < 10; i++ {
		kb.Add(&node.Node{ID: node.NewNodeID([]byte{byte(i)})})
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			nodeID := node.NewNodeID([]byte{byte(i % 256)})
			_ = kb.Contains(nodeID)
		}(i)
	}

	wg.Wait()
}
