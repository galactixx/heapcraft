package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNode is a simple struct for testing the pool functionality
type TestNode struct{ Value int }

// TestDefaultPool tests the default pool functionality
func TestDefaultPool(t *testing.T) {
	constructor := func() TestNode {
		return TestNode{Value: 42}
	}

	pool := newDefaultPool(constructor)

	node1 := pool.Get()
	assert.Equal(t, 42, node1.Value)

	pool.Put(node1)
	node2 := pool.Get()
	assert.Equal(t, 42, node2.Value)

	assert.NotSame(t, &node1, &node2)
}

// TestSyncPool tests the sync pool functionality
func TestSyncPool(t *testing.T) {
	constructor := func() TestNode {
		return TestNode{Value: 100}
	}

	pool := newSyncPool(constructor)
	node1 := pool.Get()
	assert.Equal(t, 100, node1.Value)
	pool.Put(node1)
	node2 := pool.Get()
	assert.Equal(t, 100, node2.Value)
}

// TestNewPool tests the newPool function with usePool=true
func TestNewPoolWithSyncPool(t *testing.T) {
	constructor := func() TestNode {
		return TestNode{Value: 200}
	}

	pool := newPool(true, constructor)
	node := pool.Get()
	assert.Equal(t, 200, node.Value)
	pool.Put(node)
}

// TestNewPool tests the newPool function with usePool=false
func TestNewPoolWithDefaultPool(t *testing.T) {
	constructor := func() TestNode {
		return TestNode{Value: 300}
	}

	pool := newPool(false, constructor)
	node1 := pool.Get()
	assert.Equal(t, 300, node1.Value)
	pool.Put(node1)
	node2 := pool.Get()
	assert.Equal(t, 300, node2.Value)
	assert.NotSame(t, &node1, &node2)
}

// TestPoolInterface tests that both pool types implement the interface correctly
func TestPoolInterface(t *testing.T) {
	constructor := func() TestNode {
		return TestNode{Value: 500}
	}

	defaultPool := newDefaultPool(constructor)
	testPoolInterface(t, defaultPool, "defaultPool")
	syncPool := newSyncPool(constructor)
	testPoolInterface(t, syncPool, "syncPool")
}

// testPoolInterface is a helper function to test pool interface methods
func testPoolInterface(t *testing.T, p pool[TestNode], poolType string) {
	node := p.Get()
	assert.Equal(t, 500, node.Value, poolType)

	p.Put(node)
	node2 := p.Get()
	assert.Equal(t, 500, node2.Value, poolType)
}

// TestPoolConstructorFunctions tests the constructor functions
func TestPoolConstructorFunctions(t *testing.T) {
	constructor := func() TestNode {
		return TestNode{Value: 999}
	}

	defaultPool := newDefaultPool(constructor)
	assert.NotNil(t, defaultPool)
	syncPool := newSyncPool(constructor)
	assert.NotNil(t, syncPool)
	pool1 := newPool(true, constructor)
	assert.NotNil(t, pool1)
	pool2 := newPool(false, constructor)
	assert.NotNil(t, pool2)
}
