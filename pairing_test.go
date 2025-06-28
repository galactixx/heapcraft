package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSimplePairingHeapPopOrder(t *testing.T) {
	data := []HeapNode[int, int]{
		CreateHeapNode(9, 9),
		CreateHeapNode(4, 4),
		CreateHeapNode(6, 6),
		CreateHeapNode(1, 1),
		CreateHeapNode(7, 7),
		CreateHeapNode(3, 3),
	}

	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap(data, cmp, false)

	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(data), h.Length())

	var values []int
	for !h.IsEmpty() {
		popped, _, err := h.Pop()
		if err == nil {
			values = append(values, popped)
		}
	}

	expected := []int{1, 3, 4, 6, 7, 9}
	assert.Equal(t, expected, values)
	assert.True(t, h.IsEmpty())
	_, _, err := h.Pop()
	assert.NotNil(t, err)
}

func TestInsertPopPeekLenIsEmptySimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap([]HeapNode[int, int]{}, cmp, false)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	_, _, err := h.Peek()
	assert.NotNil(t, err)

	input := []HeapNode[int, int]{
		CreateHeapNode(5, 5),
		CreateHeapNode(2, 2),
		CreateHeapNode(8, 8),
		CreateHeapNode(3, 3),
		CreateHeapNode(6, 6),
	}
	expectedOrder := []int{2, 3, 5, 6, 8}

	for _, pair := range input {
		h.Push(pair.value, pair.priority)
	}

	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(input), h.Length())
	peekValue, _ := h.PeekValue()
	assert.Equal(t, 2, peekValue)

	for i, expected := range expectedOrder {
		popped, _, err := h.Pop()
		assert.Nil(t, err)
		assert.Equal(t, expected, popped)
		assert.Equal(t, len(input)-(i+1), h.Length())
	}

	assert.True(t, h.IsEmpty())
	_, _, err = h.Peek()
	assert.NotNil(t, err)
}

func TestClearCloneSimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap([]HeapNode[int, int]{
		CreateHeapNode(1, 1),
		CreateHeapNode(2, 2),
		CreateHeapNode(3, 3),
	}, cmp, false)

	// Test basic cloning
	clone := h.Clone()
	assert.Equal(t, h.Length(), clone.Length())
	hPeekValue, _ := h.PeekValue()
	clonePeekValue, _ := clone.PeekValue()
	assert.Equal(t, hPeekValue, clonePeekValue)
	hPeekPriority, _ := h.PeekPriority()
	clonePeekPriority, _ := clone.PeekPriority()
	assert.Equal(t, hPeekPriority, clonePeekPriority)

	// Test independence of clone
	h.Push(0, 0)
	hPeekValueAfterInsert, _ := h.PeekValue()
	assert.Equal(t, 0, hPeekValueAfterInsert)
	clonePeekValueAfterInsert, _ := clone.PeekValue()
	assert.Equal(t, 1, clonePeekValueAfterInsert)

	// Test that clone maintains its own state
	clone.Push(5, 5)
	assert.Equal(t, 4, clone.Length())
	assert.Equal(t, 4, h.Length())

	// Test that clearing original doesn't affect clone
	h.Clear()
	assert.True(t, h.IsEmpty())
	assert.False(t, clone.IsEmpty())
	assert.Equal(t, 4, clone.Length())
}

func TestSimplePairingHeapDeepClone(t *testing.T) {
	// Create a heap with a complex structure
	h := NewSimplePairingHeap([]HeapNode[int, int]{}, lt, false)
	h.Push(5, 5)
	h.Push(3, 3)
	h.Push(7, 7)
	h.Push(1, 1)
	h.Push(9, 9)

	// Create a clone
	clone := h.Clone()

	// Test that all elements are in the same order
	originalElements := make([]int, 0)
	cloneElements := make([]int, 0)

	for !h.IsEmpty() {
		val, _ := h.PopValue()
		originalElements = append(originalElements, val)
	}

	for !clone.IsEmpty() {
		val, _ := clone.PopValue()
		cloneElements = append(cloneElements, val)
	}

	assert.Equal(t, originalElements, cloneElements)

	// Test that modifying clone doesn't affect original
	h = NewSimplePairingHeap([]HeapNode[int, int]{}, lt, false)
	h.Push(5, 5)
	h.Push(3, 3)
	clone = h.Clone()

	clone.Push(1, 1)
	assert.Equal(t, 2, h.Length())
	assert.Equal(t, 3, clone.Length())

	// Test that clone maintains heap property
	val, _ := clone.PopValue()
	assert.Equal(t, 1, val)
}

func TestPairingHeapDeepClone(t *testing.T) {
	// Create a heap with a complex structure
	h := NewPairingHeap([]HeapNode[int, int]{}, lt, false)
	id1 := h.Push(5, 5)
	id2 := h.Push(3, 3)
	id3 := h.Push(7, 7)
	id4 := h.Push(1, 1)
	h.Push(9, 9)

	// Create a clone
	clone := h.Clone()

	// Test that all elements are preserved with their IDs
	for _, id := range []string{id1, id2, id3, id4} {
		val1, err1 := h.GetValue(id)
		val2, err2 := clone.GetValue(id)
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, val1, val2)
	}

	// Test that modifying clone doesn't affect original
	h.Clear()
	h.Push(5, 5)
	h.Push(3, 3)
	clone = h.Clone()

	newID := clone.Push(1, 1)
	assert.Equal(t, 2, h.Length())
	assert.Equal(t, 3, clone.Length())

	// Test that clone maintains heap property and node tracking
	val, _ := clone.PopValue()
	assert.Equal(t, 1, val)

	// Test that new nodes in clone have unique IDs
	_, _, err := h.Get(newID)
	assert.Error(t, err)
	_, _, err = clone.Get(newID)
	assert.Error(t, err)

	// Test that clone maintains independent node tracking
	h.Push(10, 10)
	clone.Push(20, 20)

	hVal, _ := h.PeekValue()
	cloneVal, _ := clone.PeekValue()
	assert.Equal(t, hVal, cloneVal)
}

func TestPairingHeapCloneWithUpdates(t *testing.T) {
	// Create a heap with a complex structure
	h := NewPairingHeap([]HeapNode[int, int]{}, lt, false)
	id1 := h.Push(5, 5)
	id2 := h.Push(3, 3)
	id3 := h.Push(7, 7)
	id4 := h.Push(1, 1)

	// Create a clone
	clone := h.Clone()

	// Update values in original
	err := h.UpdateValue(id1, 50)
	assert.NoError(t, err)
	err = h.UpdatePriority(id2, 30)
	assert.NoError(t, err)

	// Verify clone remains unchanged
	val1, _ := clone.GetValue(id1)
	val2, _ := clone.GetValue(id2)
	assert.Equal(t, 5, val1)
	assert.Equal(t, 3, val2)

	// Update values in clone
	err = clone.UpdateValue(id3, 70)
	assert.NoError(t, err)
	err = clone.UpdatePriority(id4, 10)
	assert.NoError(t, err)

	// Verify original remains unchanged
	val3, _ := h.GetValue(id3)
	val4, _ := h.GetValue(id4)
	assert.Equal(t, 7, val3)
	assert.Equal(t, 1, val4)

	// Test that both heaps maintain correct order after updates
	hVal, _ := h.PeekValue()
	cloneVal, _ := clone.PeekValue()
	assert.Equal(t, 1, hVal)
	assert.Equal(t, 3, cloneVal)
}

func TestPeekPopEmptySimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap([]HeapNode[int, int]{}, cmp, false)
	_, _, err := h.Peek()
	assert.NotNil(t, err)
	_, _, err = h.Pop()
	assert.NotNil(t, err)
	_, err = h.PopValue()
	assert.NotNil(t, err)
	_, err = h.PopPriority()
	assert.NotNil(t, err)
}

func TestLengthIsEmptySimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap([]HeapNode[int, int]{}, cmp, false)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())

	h.Push(10, 10)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, 1, h.Length())
}

func TestPeekValueAndPrioritySimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }

	h := NewSimplePairingHeap([]HeapNode[int, int]{}, cmp, false)
	peekValueEmpty, _ := h.PeekValue()
	assert.Equal(t, 0, peekValueEmpty)
	peekPriorityEmpty, _ := h.PeekPriority()
	assert.Equal(t, 0, peekPriorityEmpty)

	h.Push(42, 10)
	peekValue42, _ := h.PeekValue()
	assert.Equal(t, 42, peekValue42)
	peekPriority10, _ := h.PeekPriority()
	assert.Equal(t, 10, peekPriority10)

	h.Push(15, 5)
	peekValue15, _ := h.PeekValue()
	assert.Equal(t, 15, peekValue15)
	peekPriority5, _ := h.PeekPriority()
	assert.Equal(t, 5, peekPriority5)

	h.Push(100, 1)
	peekValue100, _ := h.PeekValue()
	assert.Equal(t, 100, peekValue100)
	peekPriority1, _ := h.PeekPriority()
	assert.Equal(t, 1, peekPriority1)

	h.Pop()
	peekValueAfterPop, _ := h.PeekValue()
	assert.Equal(t, 15, peekValueAfterPop)
	peekPriorityAfterPop, _ := h.PeekPriority()
	assert.Equal(t, 5, peekPriorityAfterPop)

	h.Clear()
	peekValueAfterClear, _ := h.PeekValue()
	assert.Equal(t, 0, peekValueAfterClear)
	peekPriorityAfterClear, _ := h.PeekPriority()
	assert.Equal(t, 0, peekPriorityAfterClear)
}

func TestPopValueAndPrioritySimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap([]HeapNode[int, int]{
		CreateHeapNode(42, 10),
		CreateHeapNode(15, 5),
		CreateHeapNode(100, 1),
	}, cmp, false)

	val, err := h.PopValue()
	assert.Nil(t, err)
	assert.Equal(t, 100, val)
	peekValue15AfterPop, _ := h.PeekValue()
	assert.Equal(t, 15, peekValue15AfterPop)

	pri, err := h.PopPriority()
	assert.Nil(t, err)
	assert.Equal(t, 5, pri)
	peekValue42AfterPop, _ := h.PeekValue()
	assert.Equal(t, 42, peekValue42AfterPop)

	h.Clear()
	_, err = h.PopValue()
	assert.NotNil(t, err)
	_, err = h.PopPriority()
	assert.NotNil(t, err)
}

func TestPairingHeapUpdateValue(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewPairingHeap([]HeapNode[int, int]{}, cmp, false)

	id1 := h.Push(1, 10)
	h.Push(2, 20)
	h.Push(3, 30)

	err := h.UpdateValue(id1, 100)
	assert.Nil(t, err)
	node, exists := h.elements[id1]
	assert.True(t, exists)
	assert.Equal(t, 100, node.value)

	err = h.UpdateValue("non-existent-id", 100)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNodeNotFound.Error(), err.Error())

	popped, _, err := h.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 100, popped)
}

func TestPairingHeapUpdatePriority(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewPairingHeap([]HeapNode[int, int]{}, cmp, false)

	h.Push(1, 10)
	id2 := h.Push(2, 20)
	h.Push(3, 30)

	err := h.UpdatePriority(id2, 5)
	assert.Nil(t, err)

	_, priority, err := h.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 5, priority)

	id1 := h.Push(1, 10)
	err = h.UpdatePriority(id1, 15)
	assert.Nil(t, err)

	_, priority, err = h.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 10, priority)
}

func TestPairingHeapUpdatePriorityEdgeCases(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewPairingHeap([]HeapNode[int, int]{}, cmp, false)

	id1 := h.Push(1, 10)
	err := h.UpdatePriority(id1, 20)
	assert.Nil(t, err)
	_, priority, err := h.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 20, priority)

	h.Push(1, 10)
	id2 := h.Push(2, 20)
	h.Push(3, 30)
	err = h.UpdatePriority(id2, 5)
	assert.Nil(t, err)
	_, priority, err = h.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 5, priority)

	h.Clear()
	h.Push(1, 10)
	h.Push(2, 20)
	id3 := h.Push(3, 30)
	err = h.UpdatePriority(id3, 5)
	assert.Nil(t, err)
	_, priority, err = h.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 5, priority)
}

func TestPairingHeapClone(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewPairingHeap([]HeapNode[int, int]{}, cmp, false)

	h.Push(1, 10)
	h.Push(2, 20)
	h.Push(3, 30)

	clone := h.Clone()
	assert.Equal(t, h.size, clone.size)
	assert.Equal(t, len(h.elements), len(clone.elements))

	for id, node := range h.elements {
		cloneNode, exists := clone.elements[id]
		assert.True(t, exists)
		assert.Equal(t, node.value, cloneNode.value)
		assert.Equal(t, node.priority, cloneNode.priority)
	}
}

func TestComplexHeapStructure(t *testing.T) {
	h := NewPairingHeap[int](nil, func(a, b int) bool { return a < b }, false)

	h.Push(1, 1)
	h.Push(2, 2)
	h.Push(3, 3)
	h.Push(4, 4)
	h.Push(5, 5)
	h.Push(6, 6)
	h.Push(7, 7)

	assert.Equal(t, 7, h.Length())
	peekValueComplex, _ := h.PeekValue()
	assert.Equal(t, 1, peekValueComplex)
	peekPriorityComplex, _ := h.PeekPriority()
	assert.Equal(t, 1, peekPriorityComplex)

	values := make([]int, 0)
	for !h.IsEmpty() {
		val, err := h.PopValue()
		assert.Nil(t, err)
		values = append(values, val)
	}
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, values)
}

func TestLeafNodeUpdate(t *testing.T) {
	h := NewPairingHeap[int](nil, func(a, b int) bool { return a < b }, false)

	h.Push(1, 1)
	h.Push(2, 2)
	h.Push(3, 3)
	h.Push(4, 4)
	h.Push(5, 5)
	h.Push(6, 6)
	id7 := h.Push(7, 7)

	err := h.UpdatePriority(id7, 0)
	assert.Nil(t, err)
	peekValueLeaf, _ := h.PeekValue()
	assert.Equal(t, 7, peekValueLeaf)
	peekPriorityLeaf, _ := h.PeekPriority()
	assert.Equal(t, 0, peekPriorityLeaf)

	values := make([]int, 0)
	for !h.IsEmpty() {
		val, err := h.PopValue()
		assert.Nil(t, err)
		values = append(values, val)
	}
	assert.Equal(t, []int{7, 1, 2, 3, 4, 5, 6}, values)
}

func TestMiddleNodeUpdate(t *testing.T) {
	h := NewPairingHeap[int](nil, func(a, b int) bool { return a < b }, false)

	h.Push(1, 1)
	h.Push(2, 2)
	id3 := h.Push(3, 3)
	h.Push(4, 4)
	h.Push(5, 5)
	h.Push(6, 6)
	h.Push(7, 7)

	err := h.UpdatePriority(id3, 0)
	assert.Nil(t, err)
	peekValueMiddle, _ := h.PeekValue()
	assert.Equal(t, 3, peekValueMiddle)
	peekPriorityMiddle, _ := h.PeekPriority()
	assert.Equal(t, 0, peekPriorityMiddle)

	values := make([]int, 0)
	for !h.IsEmpty() {
		val, err := h.PopValue()
		assert.Nil(t, err)
		values = append(values, val)
	}
	assert.Equal(t, []int{3, 1, 2, 4, 5, 6, 7}, values)
}

func TestMultipleNodeUpdates(t *testing.T) {
	h := NewPairingHeap[int](nil, func(a, b int) bool { return a < b }, false)

	h.Push(1, 1)
	id2 := h.Push(2, 2)
	h.Push(3, 3)
	id4 := h.Push(4, 4)
	h.Push(5, 5)
	id6 := h.Push(6, 6)
	h.Push(7, 7)

	err := h.UpdatePriority(id4, 0)
	assert.Nil(t, err)
	peekValueMultiple1, _ := h.PeekValue()
	assert.Equal(t, 4, peekValueMultiple1)
	peekPriorityMultiple1, _ := h.PeekPriority()
	assert.Equal(t, 0, peekPriorityMultiple1)

	err = h.UpdatePriority(id2, 1)
	assert.Nil(t, err)
	peekValueMultiple2, _ := h.PeekValue()
	assert.Equal(t, 4, peekValueMultiple2)
	peekPriorityMultiple2, _ := h.PeekPriority()
	assert.Equal(t, 0, peekPriorityMultiple2)

	err = h.UpdatePriority(id6, -1)
	assert.Nil(t, err)
	peekValueMultiple3, _ := h.PeekValue()
	assert.Equal(t, 6, peekValueMultiple3)
	peekPriorityMultiple3, _ := h.PeekPriority()
	assert.Equal(t, -1, peekPriorityMultiple3)

	values := make([]int, 0)
	for !h.IsEmpty() {
		val, err := h.PopValue()
		assert.Nil(t, err)
		values = append(values, val)
	}
	assert.Equal(t, []int{6, 4, 1, 2, 3, 5, 7}, values)
}

func TestReversePriorityUpdates(t *testing.T) {
	h := NewPairingHeap[int](nil, func(a, b int) bool { return a < b }, false)

	id1 := h.Push(1, 10)
	id2 := h.Push(2, 20)
	id3 := h.Push(3, 30)
	id4 := h.Push(4, 40)
	id5 := h.Push(5, 50)
	id6 := h.Push(6, 60)
	id7 := h.Push(7, 70)

	err := h.UpdatePriority(id7, 1)
	assert.Nil(t, err)
	err = h.UpdatePriority(id6, 2)
	assert.Nil(t, err)
	err = h.UpdatePriority(id5, 3)
	assert.Nil(t, err)
	err = h.UpdatePriority(id4, 4)
	assert.Nil(t, err)
	err = h.UpdatePriority(id3, 5)
	assert.Nil(t, err)
	err = h.UpdatePriority(id2, 6)
	assert.Nil(t, err)
	err = h.UpdatePriority(id1, 7)
	assert.Nil(t, err)

	values := make([]int, 0)
	for !h.IsEmpty() {
		val, err := h.PopValue()
		assert.Nil(t, err)
		values = append(values, val)
	}
	assert.Equal(t, []int{7, 6, 5, 4, 3, 2, 1}, values)
}

func TestPairingHeapGetters(t *testing.T) {
	h := NewPairingHeap[int, int](nil, func(a, b int) bool { return a < b }, false)
	id1 := h.Push(42, 10)
	id2 := h.Push(15, 5)
	id3 := h.Push(100, 1)

	value, priority, _ := h.Get(id1)
	assert.Equal(t, 42, value)
	assert.Equal(t, 10, priority)

	val, _ := h.GetValue(id2)
	assert.Equal(t, 15, val)

	pri, _ := h.GetPriority(id3)
	assert.Equal(t, 1, pri)

	_, _, err := h.Get("non-existent-id")
	assert.NotNil(t, err)
	_, err = h.GetValue("non-existent-id")
	assert.NotNil(t, err)
	_, err = h.GetPriority("non-existent-id")
	assert.NotNil(t, err)

	h.Pop()
	_, _, err = h.Get(id3)
	assert.NotNil(t, err)
}

func TestPairingHeapInsertReturnsID(t *testing.T) {
	h := NewPairingHeap([]HeapNode[int, int]{}, lt, false)

	// Test UUID-based ID assignment
	id1 := h.Push(10, 10)
	id2 := h.Push(20, 20)
	id3 := h.Push(30, 30)

	// Verify IDs are unique strings (UUIDs)
	assert.NotEqual(t, id1, id2)
	assert.NotEqual(t, id2, id3)
	assert.NotEqual(t, id1, id3)
	assert.Greater(t, len(id1), 0)
	assert.Greater(t, len(id2), 0)
	assert.Greater(t, len(id3), 0)

	// Verify elements can be retrieved using IDs
	val1, _ := h.GetValue(id1)
	val2, _ := h.GetValue(id2)
	val3, _ := h.GetValue(id3)
	assert.Equal(t, 10, val1)
	assert.Equal(t, 20, val2)
	assert.Equal(t, 30, val3)

	// Test ID continues after operations
	h.Pop()
	id4 := h.Push(40, 40)
	assert.NotEqual(t, id1, id4)
	assert.NotEqual(t, id2, id4)
	assert.NotEqual(t, id3, id4)
}

func TestPairingHeapInsertIDAfterClear(t *testing.T) {
	h := NewPairingHeap([]HeapNode[int, int]{}, lt, false)

	id1 := h.Push(10, 10)
	h.Clear()
	id2 := h.Push(20, 20)

	// Both should be unique UUIDs
	assert.NotEqual(t, id1, id2)
	assert.Greater(t, len(id1), 0)
	assert.Greater(t, len(id2), 0)
}

func TestSimplePairingHeapInsertNoID(t *testing.T) {
	h := NewSimplePairingHeap([]HeapNode[int, int]{}, lt, false)

	// SimplePairingHeap Push should not return ID
	h.Push(10, 10)
	h.Push(20, 20)

	assert.Equal(t, 2, h.Length())
	val1, _ := h.PopValue()
	val2, _ := h.PopValue()
	assert.Equal(t, 10, val1)
	assert.Equal(t, 20, val2)
}

// -------------------------------- Pairing Heap Benchmarks --------------------------------

func BenchmarkPairingHeapInsertion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewPairingHeap(data, func(a, b int) bool { return a < b }, false)

	insertions := generateRandomNumbersv1(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Push(insertions[i], insertions[i])
	}
}

func BenchmarkPairingHeapDeletion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewPairingHeap(data, func(a, b int) bool { return a < b }, false)

	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Pop()
	}
}

func BenchmarkSimplePairingHeapInsertion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewSimplePairingHeap(data, func(a, b int) bool { return a < b }, false)

	insertions := generateRandomNumbersv1(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Push(insertions[i], insertions[i])
	}
}

func BenchmarkSimplePairingHeapDeletion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewSimplePairingHeap(data, func(a, b int) bool { return a < b }, false)

	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Pop()
	}
}
