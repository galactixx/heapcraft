package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSimplePairingHeapPopOrder(t *testing.T) {
	// Create HeapPairs with values and priorities
	data := []*HeapPair[int, int]{
		CreateHeapPair(9, 9),
		CreateHeapPair(4, 4),
		CreateHeapPair(6, 6),
		CreateHeapPair(1, 1),
		CreateHeapPair(7, 7),
		CreateHeapPair(3, 3),
	}

	// Use a comparison function that compares priorities
	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap(data, cmp)

	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(data), h.Length())

	// Collect values in order
	var values []int
	for !h.IsEmpty() {
		pair := h.Pop()
		if pair != nil {
			values = append(values, pair.Value())
		}
	}

	expected := []int{1, 3, 4, 6, 7, 9}
	assert.Equal(t, expected, values)
	assert.True(t, h.IsEmpty())
	assert.Nil(t, h.Pop())
}

func TestInsertPopPeekLenIsEmptySimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap([]*HeapPair[int, int]{}, cmp)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	assert.Nil(t, h.Peek())

	// Create test data with values and priorities
	input := []*HeapPair[int, int]{
		CreateHeapPair(5, 5),
		CreateHeapPair(2, 2),
		CreateHeapPair(8, 8),
		CreateHeapPair(3, 3),
		CreateHeapPair(6, 6),
	}
	expectedOrder := []int{2, 3, 5, 6, 8}

	for _, pair := range input {
		h.Insert(pair.Value(), pair.Priority())
	}

	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(input), h.Length())
	assert.Equal(t, 2, *h.PeekValue())

	for i, expected := range expectedOrder {
		popped := h.Pop()
		assert.NotNil(t, popped)
		assert.Equal(t, expected, popped.Value())
		assert.Equal(t, len(input)-(i+1), h.Length())
	}

	assert.True(t, h.IsEmpty())
	assert.Nil(t, h.Peek())
}

func TestClearCloneSimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	data := []*HeapPair[int, int]{
		CreateHeapPair(4, 4),
		CreateHeapPair(1, 1),
		CreateHeapPair(3, 3),
		CreateHeapPair(2, 2),
	}
	h := NewSimplePairingHeap(data, cmp)
	assert.Equal(t, 4, h.Length())

	clone := h.Clone()
	assert.Equal(t, h.Length(), clone.Length())
	assert.Equal(t, *h.PeekValue(), *clone.PeekValue())

	h.Insert(0, 0)
	assert.Equal(t, 0, *h.PeekValue())
	assert.Equal(t, 1, *clone.PeekValue())

	h2 := NewSimplePairingHeap([]*HeapPair[int, int]{
		CreateHeapPair(7, 7),
		CreateHeapPair(5, 5),
		CreateHeapPair(9, 9),
	}, cmp)

	h2.Clear()
	assert.True(t, h2.IsEmpty())
}

func TestPeekPopEmptySimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap([]*HeapPair[int, int]{}, cmp)
	assert.Nil(t, h.Peek())
	assert.Nil(t, h.Pop())
	assert.Nil(t, h.PopValue())
	assert.Nil(t, h.PopPriority())
}

func TestLengthIsEmptySimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap([]*HeapPair[int, int]{}, cmp)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())

	h.Insert(10, 10)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, 1, h.Length())
}

func TestMergeWithSimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h1 := NewSimplePairingHeap([]*HeapPair[int, int]{
		CreateHeapPair(1, 1),
		CreateHeapPair(4, 4),
		CreateHeapPair(7, 7),
	}, cmp)

	h2 := NewSimplePairingHeap([]*HeapPair[int, int]{
		CreateHeapPair(2, 2),
		CreateHeapPair(3, 3),
		CreateHeapPair(5, 5),
		CreateHeapPair(6, 6),
	}, cmp)

	h1.MergeWith(h2)

	// Collect values in order
	var values []int
	for !h1.IsEmpty() {
		pair := h1.Pop()
		if pair != nil {
			values = append(values, pair.Value())
		}
	}

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, values)
}

func TestPeekValueAndPrioritySimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }

	// Test empty heap
	h := NewSimplePairingHeap([]*HeapPair[int, int]{}, cmp)
	assert.Nil(t, h.PeekValue())
	assert.Nil(t, h.PeekPriority())

	// Test single element
	h.Insert(42, 10)
	assert.Equal(t, 42, *h.PeekValue())
	assert.Equal(t, 10, *h.PeekPriority())

	// Test multiple elements
	h.Insert(15, 5)
	assert.Equal(t, 15, *h.PeekValue())
	assert.Equal(t, 5, *h.PeekPriority())

	h.Insert(100, 1)
	assert.Equal(t, 100, *h.PeekValue())
	assert.Equal(t, 1, *h.PeekPriority())

	// Test after popping
	h.Pop()
	assert.Equal(t, 15, *h.PeekValue())
	assert.Equal(t, 5, *h.PeekPriority())

	// Test after clearing
	h.Clear()
	assert.Nil(t, h.PeekValue())
	assert.Nil(t, h.PeekPriority())
}

func TestPopValueAndPrioritySimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap([]*HeapPair[int, int]{
		CreateHeapPair(42, 10),
		CreateHeapPair(15, 5),
		CreateHeapPair(100, 1),
	}, cmp)

	// Test PopValue
	val := h.PopValue()
	assert.Equal(t, 100, *val)
	assert.Equal(t, 15, *h.PeekValue())

	// Test PopPriority
	pri := h.PopPriority()
	assert.Equal(t, 5, *pri)
	assert.Equal(t, 42, *h.PeekValue())

	// Test empty heap
	h.Clear()
	assert.Nil(t, h.PopValue())
	assert.Nil(t, h.PopPriority())
}

func TestPairingHeapIDTracking(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewPairingHeap([]*HeapPair[int, int]{}, cmp)
	assert.NotNil(t, h.elements)
	assert.Equal(t, 0, len(h.elements))

	// Test ID assignment and tracking
	h.Insert(1, 10)
	h.Insert(2, 20)
	h.Insert(3, 30)

	assert.Equal(t, 3, len(h.elements))
	assert.Equal(t, uint(1), h.curID-3)

	// Verify elements are tracked
	for i := uint(1); i < h.curID; i++ {
		node, exists := h.elements[i]
		assert.True(t, exists)
		assert.Equal(t, i, node.ID())
	}

	// Test ID tracking after pop
	popped := h.Pop()
	assert.NotNil(t, popped)
	assert.Equal(t, 2, len(h.elements))
	assert.Equal(t, 1, popped.Value())

	// Test ID tracking after clear
	h.Clear()
	assert.Equal(t, 0, len(h.elements))
	assert.Equal(t, uint(1), h.curID)
}

func TestPairingHeapUpdateValue(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewPairingHeap([]*HeapPair[int, int]{}, cmp)

	// Insert some elements
	h.Insert(1, 10)
	h.Insert(2, 20)
	h.Insert(3, 30)

	// Test updating existing value
	err := h.UpdateValue(1, 100)
	assert.Nil(t, err)
	node, exists := h.elements[1]
	assert.True(t, exists)
	assert.Equal(t, 100, node.Value())

	// Test updating non-existent ID
	err = h.UpdateValue(999, 100)
	assert.NotNil(t, err)
	assert.Equal(t, "id does not link to existing node", err.Error())

	// Verify heap order is maintained
	popped := h.Pop()
	assert.NotNil(t, popped)
	assert.Equal(t, 100, popped.Value())
}

func TestPairingHeapUpdatePriority(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewPairingHeap([]*HeapPair[int, int]{}, cmp)

	// Insert some elements
	h.Insert(1, 10)
	h.Insert(2, 20)
	h.Insert(3, 30)

	// Test updating priority of non-root node
	err := h.UpdatePriority(2, 5)
	assert.Nil(t, err)

	// Verify new heap order
	popped := h.Pop()

	assert.NotNil(t, popped)
	assert.Equal(t, 2, popped.Value())
	assert.Equal(t, 5, popped.Priority())

	// Test updating root node's priority
	err = h.UpdatePriority(1, 15)
	assert.Nil(t, err)

	popped = h.Pop()
	assert.NotNil(t, popped)
	assert.Equal(t, 1, popped.Value())
}

func TestPairingHeapUpdatePriorityEdgeCases(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewPairingHeap([]*HeapPair[int, int]{}, cmp)

	// Test updating priority of single node
	h.Insert(1, 10)
	err := h.UpdatePriority(1, 20)
	assert.Nil(t, err)
	popped := h.Pop()
	assert.NotNil(t, popped)
	assert.Equal(t, 1, popped.Value())
	assert.Equal(t, 20, popped.Priority())

	// Test updating priority of first child
	h.Insert(1, 10)
	h.Insert(2, 20)
	h.Insert(3, 30)
	err = h.UpdatePriority(2, 5)
	assert.Nil(t, err)
	popped = h.Pop()
	assert.NotNil(t, popped)
	assert.Equal(t, 1, popped.Value())
	assert.Equal(t, 5, popped.Priority())

	// Test updating priority of last child
	h.Clear()
	h.Insert(1, 10)
	h.Insert(2, 20)
	h.Insert(3, 30)
	err = h.UpdatePriority(3, 5)
	assert.Nil(t, err)
	popped = h.Pop()
	assert.NotNil(t, popped)
	assert.Equal(t, 3, popped.Value())
	assert.Equal(t, 5, popped.Priority())
}

func TestPairingHeapClone(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewPairingHeap([]*HeapPair[int, int]{}, cmp)

	// Insert some elements
	h.Insert(1, 10)
	h.Insert(2, 20)
	h.Insert(3, 30)

	// Clone the heap
	clone := h.Clone()
	assert.Equal(t, h.size, clone.size)
	assert.Equal(t, h.curID, clone.curID)
	assert.Equal(t, len(h.elements), len(clone.elements))

	// Verify elements are properly cloned
	for id, node := range h.elements {
		cloneNode, exists := clone.elements[id]
		assert.True(t, exists)
		assert.Equal(t, node.Value(), cloneNode.Value())
		assert.Equal(t, node.Priority(), cloneNode.Priority())
	}
}

// TestComplexHeapStructure tests the creation and verification of a complex heap structure
// with multiple levels and siblings. It verifies that the initial structure is correct
// and that all nodes are properly inserted.
func TestComplexHeapStructure(t *testing.T) {
	// Create a heap with a complex initial structure
	// Initial structure (priority in parentheses):
	//        1(1)
	//      /     \
	//    2(2)    3(3)
	//   /   \    /  \
	// 4(4) 5(5) 6(6) 7(7)
	h := NewPairingHeap[int](nil, func(a, b int) bool { return a < b })

	// Insert nodes in a specific order to create the desired structure
	h.Insert(1, 1)
	h.Insert(2, 2)
	h.Insert(3, 3)
	h.Insert(4, 4)
	h.Insert(5, 5)
	h.Insert(6, 6)
	h.Insert(7, 7)

	// Get node IDs by finding nodes with matching values in the elements map
	nodeIDs := make(map[int]uint)
	for id, node := range h.elements {
		nodeIDs[node.value] = id
	}

	// Verify initial structure
	assert.Equal(t, 7, h.Length())
	assert.Equal(t, 1, *h.PeekValue())
	assert.Equal(t, 1, *h.PeekPriority())

	// Verify all nodes are present and in correct order
	values := make([]int, 0)
	for !h.IsEmpty() {
		val := *h.PopValue()
		values = append(values, val)
	}
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, values)
}

// TestLeafNodeUpdate tests updating a leaf node to become the new root
// by giving it the highest priority.
func TestLeafNodeUpdate(t *testing.T) {
	h := NewPairingHeap[int](nil, func(a, b int) bool { return a < b })

	// Create the initial structure
	h.Insert(1, 1)
	h.Insert(2, 2)
	h.Insert(3, 3)
	h.Insert(4, 4)
	h.Insert(5, 5)
	h.Insert(6, 6)
	h.Insert(7, 7)

	// Get node IDs
	nodeIDs := make(map[int]uint)
	for id, node := range h.elements {
		nodeIDs[node.value] = id
	}

	// Update leaf node (7) to become root
	err := h.UpdatePriority(nodeIDs[7], 0)
	assert.Nil(t, err)
	assert.Equal(t, 7, *h.PeekValue())
	assert.Equal(t, 0, *h.PeekPriority())

	// Verify all nodes are still in the heap and in correct order
	values := make([]int, 0)
	for !h.IsEmpty() {
		val := *h.PopValue()
		values = append(values, val)
	}
	assert.Equal(t, []int{7, 1, 2, 3, 4, 5, 6}, values)
}

// TestMiddleNodeUpdate tests updating a middle node to become the new root
// and verifies that the heap property is maintained.
func TestMiddleNodeUpdate(t *testing.T) {
	h := NewPairingHeap[int](nil, func(a, b int) bool { return a < b })

	// Create the initial structure
	h.Insert(1, 1)
	h.Insert(2, 2)
	h.Insert(3, 3)
	h.Insert(4, 4)
	h.Insert(5, 5)
	h.Insert(6, 6)
	h.Insert(7, 7)

	// Get node IDs
	nodeIDs := make(map[int]uint)
	for id, node := range h.elements {
		nodeIDs[node.value] = id
	}

	// Update middle node (3) to become root
	err := h.UpdatePriority(nodeIDs[3], 0)
	assert.Nil(t, err)
	assert.Equal(t, 3, *h.PeekValue())
	assert.Equal(t, 0, *h.PeekPriority())

	// Verify the heap property is maintained
	values := make([]int, 0)
	for !h.IsEmpty() {
		val := *h.PopValue()
		values = append(values, val)
	}
	assert.Equal(t, []int{3, 1, 2, 4, 5, 6, 7}, values)
}

// TestMultipleNodeUpdates tests a sequence of priority updates that
// change the heap structure multiple times, verifying that the heap
// property is maintained throughout.
func TestMultipleNodeUpdates(t *testing.T) {
	h := NewPairingHeap[int](nil, func(a, b int) bool { return a < b })

	// Create the initial structure
	h.Insert(1, 1)
	h.Insert(2, 2)
	h.Insert(3, 3)
	h.Insert(4, 4)
	h.Insert(5, 5)
	h.Insert(6, 6)
	h.Insert(7, 7)

	// Get node IDs
	nodeIDs := make(map[int]uint)
	for id, node := range h.elements {
		nodeIDs[node.value] = id
	}

	// First update: Make node 4 the root
	err := h.UpdatePriority(nodeIDs[4], 0)
	assert.Nil(t, err)
	assert.Equal(t, 4, *h.PeekValue())
	assert.Equal(t, 0, *h.PeekPriority())

	// Second update: Move node 2 up (but not higher than 4)
	err = h.UpdatePriority(nodeIDs[2], 1)
	assert.Nil(t, err)
	assert.Equal(t, 4, *h.PeekValue()) // 4 should still be root
	assert.Equal(t, 0, *h.PeekPriority())

	// Third update: Make node 6 the root
	err = h.UpdatePriority(nodeIDs[6], -1)
	assert.Nil(t, err)
	assert.Equal(t, 6, *h.PeekValue())
	assert.Equal(t, -1, *h.PeekPriority())

	// Verify final heap structure
	values := make([]int, 0)
	for !h.IsEmpty() {
		val := *h.PopValue()
		values = append(values, val)
	}
	assert.Equal(t, []int{6, 4, 1, 2, 3, 5, 7}, values)
}

// TestReversePriorityUpdates tests updating all node priorities in reverse order,
// effectively creating a heap where the nodes are ordered in reverse of their
// original insertion order.
func TestReversePriorityUpdates(t *testing.T) {
	h := NewPairingHeap[int](nil, func(a, b int) bool { return a < b })

	// Insert nodes with higher initial priorities
	h.Insert(1, 10)
	h.Insert(2, 20)
	h.Insert(3, 30)
	h.Insert(4, 40)
	h.Insert(5, 50)
	h.Insert(6, 60)
	h.Insert(7, 70)

	// Get node IDs
	nodeIDs := make(map[int]uint)
	for id, node := range h.elements {
		nodeIDs[node.value] = id
	}

	// Update priorities in reverse order
	err := h.UpdatePriority(nodeIDs[7], 1)
	assert.Nil(t, err)
	err = h.UpdatePriority(nodeIDs[6], 2)
	assert.Nil(t, err)
	err = h.UpdatePriority(nodeIDs[5], 3)
	assert.Nil(t, err)
	err = h.UpdatePriority(nodeIDs[4], 4)
	assert.Nil(t, err)
	err = h.UpdatePriority(nodeIDs[3], 5)
	assert.Nil(t, err)
	err = h.UpdatePriority(nodeIDs[2], 6)
	assert.Nil(t, err)
	err = h.UpdatePriority(nodeIDs[1], 7)
	assert.Nil(t, err)

	// Verify the heap is now in reverse order
	values := make([]int, 0)
	for !h.IsEmpty() {
		val := *h.PopValue()
		values = append(values, val)
	}
	assert.Equal(t, []int{7, 6, 5, 4, 3, 2, 1}, values)
}
