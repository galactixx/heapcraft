package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock implementations to test interfaces
type mockSimpleNode struct {
	value    int
	priority int
}

func (m mockSimpleNode) Value() int    { return m.value }
func (m mockSimpleNode) Priority() int { return m.priority }

type mockNode struct {
	id       string
	value    string
	priority float64
}

func (m mockNode) ID() string        { return m.id }
func (m mockNode) Value() string     { return m.value }
func (m mockNode) Priority() float64 { return m.priority }

func TestHeapNodeCreation(t *testing.T) {
	// Test CreateHeapNode
	heapNode := CreateHeapNode("test", 42)
	assert.Equal(t, "test", heapNode.Value())
	assert.Equal(t, 42, heapNode.Priority())

	// Test CreateHeapNode
	heapNodePtr := CreateHeapNode(123, 456.78)
	assert.Equal(t, 123, heapNodePtr.Value())
	assert.Equal(t, 456.78, heapNodePtr.Priority())
	assert.NotNil(t, heapNodePtr)
}

func TestHeapNodeMethods(t *testing.T) {
	// Test direct struct creation
	node := HeapNode[int, string]{
		value:    100,
		priority: "high",
	}

	assert.Equal(t, 100, node.Value())
	assert.Equal(t, "high", node.Priority())
}

func TestRadixPairCreation(t *testing.T) {
	// Test CreateRadixPair
	radixPair := CreateHeapNode(true, 3.14)
	assert.Equal(t, true, radixPair.Value())
	assert.Equal(t, 3.14, radixPair.Priority())
	assert.NotNil(t, radixPair)
}

func TestRadixPairMethods(t *testing.T) {
	// Test direct struct creation
	pair := HeapNode[[]int, uint]{
		value:    []int{1, 2, 3},
		priority: 42,
	}

	assert.Equal(t, []int{1, 2, 3}, pair.Value())
	assert.Equal(t, uint(42), pair.Priority())
}

func TestSimpleNodeInterface(t *testing.T) {
	// Test with mock implementation
	mock := mockSimpleNode{
		value:    999,
		priority: 888,
	}

	var simpleNode SimpleNode[int, int] = mock
	assert.Equal(t, 999, simpleNode.Value())
	assert.Equal(t, 888, simpleNode.Priority())

	// Test with HeapNode
	heapNode := CreateHeapNode("hello", 123)
	var simpleNode2 SimpleNode[string, int] = heapNode
	assert.Equal(t, "hello", simpleNode2.Value())
	assert.Equal(t, 123, simpleNode2.Priority())
}

func TestNodeInterface(t *testing.T) {
	// Test with mock implementation
	mock := mockNode{
		id:       "123",
		value:    "test",
		priority: 45.67,
	}

	var node Node[string, float64] = mock
	assert.Equal(t, "123", node.ID())
	assert.Equal(t, "test", node.Value())
	assert.Equal(t, 45.67, node.Priority())
}

func TestGenericNodeTypes(t *testing.T) {
	// Test with different type combinations
	stringIntNode := CreateHeapNode("string", 42)
	assert.Equal(t, "string", stringIntNode.Value())
	assert.Equal(t, 42, stringIntNode.Priority())

	boolFloatNode := CreateHeapNode(true, 3.14159)
	assert.Equal(t, true, boolFloatNode.Value())
	assert.Equal(t, 3.14159, boolFloatNode.Priority())

	sliceNode := CreateHeapNode([]int{1, 2, 3}, "priority")
	assert.Equal(t, []int{1, 2, 3}, sliceNode.Value())
	assert.Equal(t, "priority", sliceNode.Priority())
}

func TestRadixPairGenericTypes(t *testing.T) {
	// Test HeapNode with different type combinations
	complexRadix := CreateHeapNode(map[string]int{"a": 1, "b": 2}, 99.9)
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, complexRadix.Value())
	assert.Equal(t, 99.9, complexRadix.Priority())

	pointerRadix := CreateHeapNode(&[]int{1, 2, 3}, uint(123))
	assert.Equal(t, &[]int{1, 2, 3}, pointerRadix.Value())
	assert.Equal(t, uint(123), pointerRadix.Priority())
}

func TestNodeEquality(t *testing.T) {
	// Test that nodes with same values are equal
	node1 := CreateHeapNode("test", 42)
	node2 := CreateHeapNode("test", 42)
	assert.Equal(t, node1.Value(), node2.Value())
	assert.Equal(t, node1.Priority(), node2.Priority())

	// Test that nodes with different values are not equal
	node3 := CreateHeapNode("different", 42)
	assert.NotEqual(t, node1.Value(), node3.Value())
	assert.Equal(t, node1.Priority(), node3.Priority())
}

func TestPointerVsValue(t *testing.T) {
	// Test that pointer and value versions work correctly
	valueNode := CreateHeapNode("test", 42)
	ptrNode := CreateHeapNode("test", 42)

	assert.Equal(t, valueNode.Value(), ptrNode.Value())
	assert.Equal(t, valueNode.Priority(), ptrNode.Priority())

	// Test that pointer can be modified
	ptrNode.value = "modified"
	assert.Equal(t, "modified", ptrNode.Value())
	assert.Equal(t, "test", valueNode.Value())
}
