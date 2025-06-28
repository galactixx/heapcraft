package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeapNodeCreation(t *testing.T) {
	// Test CreateHeapNode
	heapNode := CreateHeapNode("test", 42)
	assert.Equal(t, "test", heapNode.value)
	assert.Equal(t, 42, heapNode.priority)

	// Test CreateHeapNode
	heapNodePtr := CreateHeapNode(123, 456.78)
	assert.Equal(t, 123, heapNodePtr.value)
	assert.Equal(t, 456.78, heapNodePtr.priority)
	assert.NotNil(t, heapNodePtr)
}

func TestHeapNodeMethods(t *testing.T) {
	// Test direct struct creation
	node := HeapNode[int, string]{
		value:    100,
		priority: "high",
	}

	assert.Equal(t, 100, node.value)
	assert.Equal(t, "high", node.priority)
}

func TestRadixPairCreation(t *testing.T) {
	// Test CreateRadixPair
	radixPair := CreateHeapNode(true, 3.14)
	assert.Equal(t, true, radixPair.value)
	assert.Equal(t, 3.14, radixPair.priority)
	assert.NotNil(t, radixPair)
}

func TestRadixPairMethods(t *testing.T) {
	// Test direct struct creation
	pair := HeapNode[[]int, uint]{
		value:    []int{1, 2, 3},
		priority: 42,
	}

	assert.Equal(t, []int{1, 2, 3}, pair.value)
	assert.Equal(t, uint(42), pair.priority)
}

func TestGenericNodeTypes(t *testing.T) {
	// Test with different type combinations
	stringIntNode := CreateHeapNode("string", 42)
	assert.Equal(t, "string", stringIntNode.value)
	assert.Equal(t, 42, stringIntNode.priority)

	boolFloatNode := CreateHeapNode(true, 3.14159)
	assert.Equal(t, true, boolFloatNode.value)
	assert.Equal(t, 3.14159, boolFloatNode.priority)

	sliceNode := CreateHeapNode([]int{1, 2, 3}, "priority")
	assert.Equal(t, []int{1, 2, 3}, sliceNode.value)
	assert.Equal(t, "priority", sliceNode.priority)
}

func TestRadixPairGenericTypes(t *testing.T) {
	// Test HeapNode with different type combinations
	complexRadix := CreateHeapNode(map[string]int{"a": 1, "b": 2}, 99.9)
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, complexRadix.value)
	assert.Equal(t, 99.9, complexRadix.priority)

	pointerRadix := CreateHeapNode(&[]int{1, 2, 3}, uint(123))
	assert.Equal(t, &[]int{1, 2, 3}, pointerRadix.value)
	assert.Equal(t, uint(123), pointerRadix.priority)
}

func TestNodeEquality(t *testing.T) {
	// Test that nodes with same values are equal
	node1 := CreateHeapNode("test", 42)
	node2 := CreateHeapNode("test", 42)
	assert.Equal(t, node1.value, node2.value)
	assert.Equal(t, node1.priority, node2.priority)

	// Test that nodes with different values are not equal
	node3 := CreateHeapNode("different", 42)
	assert.NotEqual(t, node1.value, node3.value)
	assert.Equal(t, node1.priority, node3.priority)
}

func TestPointerVsValue(t *testing.T) {
	// Test that pointer and value versions work correctly
	valueNode := CreateHeapNode("test", 42)
	ptrNode := CreateHeapNode("test", 42)

	assert.Equal(t, valueNode.value, ptrNode.value)
	assert.Equal(t, valueNode.priority, ptrNode.priority)

	// Test that pointer can be modified
	ptrNode.value = "modified"
	assert.Equal(t, "modified", ptrNode.value)
	assert.Equal(t, "test", valueNode.value)
}
