package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func collectSimpleSkew(h *SimpleSkewHeap[int, int]) []int {
	result := make([]int, 0)
	for !h.IsEmpty() {
		val, _ := h.PopValue()
		result = append(result, val)
	}
	return result
}

func TestNewSkewHeapPopOrder(t *testing.T) {
	data := []HeapNode[int, int]{
		CreateHeapNode(9, 9),
		CreateHeapNode(4, 4),
		CreateHeapNode(6, 6),
		CreateHeapNode(1, 1),
		CreateHeapNode(7, 7),
		CreateHeapNode(3, 3),
	}
	h := NewSimpleSkewHeap(data, lt, false)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(data), h.Length())

	expected := []int{1, 3, 4, 6, 7, 9}
	actual := collectSimpleSkew(h)
	assert.Equal(t, expected, actual)
	assert.True(t, h.IsEmpty())

	_, _, err := h.Pop()
	assert.NotNil(t, err)
}

func TestInsertPopPeekLenIsEmptySkew(t *testing.T) {
	h := NewSimpleSkewHeap([]HeapNode[int, int]{}, lt, false)
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
	value, priority, _ := h.Peek()
	assert.Equal(t, 2, value)
	assert.Equal(t, 2, priority)

	for i, expected := range expectedOrder {
		value, priority, err := h.Pop()
		assert.Nil(t, err)
		assert.NotNil(t, value)
		assert.Equal(t, expected, value)
		assert.Equal(t, expected, priority)
		assert.Equal(t, len(input)-(i+1), h.Length())
	}

	assert.True(t, h.IsEmpty())
	_, _, err = h.Peek()
	assert.NotNil(t, err)
}

func TestClearCloneSkew(t *testing.T) {
	data := []HeapNode[int, int]{
		CreateHeapNode(4, 4),
		CreateHeapNode(1, 1),
		CreateHeapNode(3, 3),
		CreateHeapNode(2, 2),
	}
	h := NewSimpleSkewHeap(data, lt, false)
	assert.Equal(t, 4, h.Length())

	// Test basic cloning
	clone := h.Clone()
	assert.Equal(t, h.Length(), clone.Length())
	value, priority, _ := h.Peek()
	cloneValue, clonePriority, _ := clone.Peek()
	assert.Equal(t, value, cloneValue)
	assert.Equal(t, priority, clonePriority)

	// Test independence of clone
	h.Push(0, 0)
	value, _, _ = h.Peek()
	assert.Equal(t, 0, value)
	cloneValue, _, _ = clone.Peek()
	assert.Equal(t, 1, cloneValue)

	// Test that clone maintains its own state
	clone.Push(5, 5)
	assert.Equal(t, 5, clone.Length())
	assert.Equal(t, 5, h.Length())

	// Test that clearing original doesn't affect clone
	h.Clear()
	assert.True(t, h.IsEmpty())
	assert.False(t, clone.IsEmpty())
	assert.Equal(t, 5, clone.Length())
}

func TestSimpleSkewHeapDeepClone(t *testing.T) {
	// Create a heap with a complex structure
	h := NewSimpleSkewHeap([]HeapNode[int, int]{}, lt, false)
	h.Push(5, 5)
	h.Push(3, 3)
	h.Push(7, 7)
	h.Push(1, 1)
	h.Push(9, 9)

	// Create a clone
	clone := h.Clone()

	// Test that all elements are in the same order
	originalElements := collectSimpleSkew(h)
	cloneElements := collectSimpleSkew(clone)
	assert.Equal(t, originalElements, cloneElements)

	// Test that modifying clone doesn't affect original
	h = NewSimpleSkewHeap([]HeapNode[int, int]{}, lt, false)
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

func TestSkewHeapDeepClone(t *testing.T) {
	// Create a heap with a complex structure
	h := NewSkewHeap([]HeapNode[int, int]{}, lt, HeapConfig{UsePool: false})
	id1, _ := h.Push(5, 5)
	id2, _ := h.Push(3, 3)
	id3, _ := h.Push(7, 7)
	id4, _ := h.Push(1, 1)
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

	newID, _ := clone.Push(1, 1)
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

func TestSkewHeapCloneWithUpdates(t *testing.T) {
	// Create a heap with a complex structure
	h := NewSkewHeap([]HeapNode[int, int]{}, lt, HeapConfig{UsePool: false})
	id1, _ := h.Push(5, 5)
	id2, _ := h.Push(3, 3)
	id3, _ := h.Push(7, 7)
	id4, _ := h.Push(1, 1)
	h.Push(9, 9)

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

func TestPeekPopEmptySkew(t *testing.T) {
	h := NewSimpleSkewHeap([]HeapNode[int, int]{}, lt, false)
	_, _, err := h.Peek()
	assert.NotNil(t, err)
	_, _, err = h.Pop()
	assert.NotNil(t, err)
	_, err = h.PopValue()
	assert.NotNil(t, err)
	_, err = h.PopPriority()
	assert.NotNil(t, err)
}

func TestLengthIsEmptySkew(t *testing.T) {
	h := NewSimpleSkewHeap([]HeapNode[int, int]{}, lt, false)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())

	h.Push(10, 10)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, 1, h.Length())
}

func TestSimpleSkewHeapInsertNoID(t *testing.T) {
	h := NewSimpleSkewHeap([]HeapNode[int, int]{}, lt, false)

	// SimpleSkewHeap's Push method should not return an ID
	// since it doesn't have node tracking
	h.Push(10, 10)
	h.Push(20, 20)
	h.Push(30, 30)

	// Verify elements were inserted correctly
	assert.Equal(t, 3, h.Length())

	// Pop elements to verify they were inserted in correct order
	val1, _ := h.PopValue()
	assert.Equal(t, 10, val1)

	val2, _ := h.PopValue()
	assert.Equal(t, 20, val2)

	val3, _ := h.PopValue()
	assert.Equal(t, 30, val3)
}

func TestPeekValueAndPrioritySkew(t *testing.T) {
	h := NewSimpleSkewHeap([]HeapNode[int, int]{}, lt, false)
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

	_, _, err := h.Pop()
	assert.Nil(t, err)
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

func TestPopValueAndPrioritySkew(t *testing.T) {
	h := NewSimpleSkewHeap([]HeapNode[int, int]{
		CreateHeapNode(42, 10),
		CreateHeapNode(15, 5),
		CreateHeapNode(100, 1),
	}, lt, false)

	val, _ := h.PopValue()
	assert.Equal(t, 100, val)
	peekValue15, _ := h.PeekValue()
	assert.Equal(t, 15, peekValue15)

	pri, _ := h.PopPriority()
	assert.Equal(t, 5, pri)
	peekValue42, _ := h.PeekValue()
	assert.Equal(t, 42, peekValue42)

	h.Clear()
	popValueAfterClear, _ := h.PopValue()
	assert.Equal(t, 0, popValueAfterClear)
	popPriorityAfterClear, _ := h.PopPriority()
	assert.Equal(t, 0, popPriorityAfterClear)
}

func TestSkewHeapInsertReturnsID(t *testing.T) {
	h := NewSkewHeap([]HeapNode[int, int]{}, lt, HeapConfig{UsePool: false})

	// Test UUID-based ID assignment
	id1, _ := h.Push(10, 10)
	id2, _ := h.Push(20, 20)
	id3, _ := h.Push(30, 30)

	// Verify IDs are unique strings (UUIDs)
	assert.NotEqual(t, id1, id2)
	assert.NotEqual(t, id2, id3)
	assert.NotEqual(t, id1, id3)
	assert.Greater(t, len(id1), 0)
	assert.Greater(t, len(id2), 0)
	assert.Greater(t, len(id3), 0)

	// Verify we can retrieve the inserted elements using the returned IDs
	val1, err := h.GetValue(id1)
	assert.Nil(t, err)
	assert.Equal(t, 10, val1)

	val2, err := h.GetValue(id2)
	assert.Nil(t, err)
	assert.Equal(t, 20, val2)

	val3, err := h.GetValue(id3)
	assert.Nil(t, err)
	assert.Equal(t, 30, val3)

	// Test that IDs continue after operations
	h.Pop() // Remove one element
	id4, _ := h.Push(40, 40)
	assert.NotEqual(t, id1, id4)
	assert.NotEqual(t, id2, id4)
	assert.NotEqual(t, id3, id4)

	// Verify the new element can be retrieved
	val4, err := h.GetValue(id4)
	assert.Nil(t, err)
	assert.Equal(t, 40, val4)
}

func TestSkewHeapInsertIDAfterClear(t *testing.T) {
	h := NewSkewHeap([]HeapNode[int, int]{}, lt, HeapConfig{UsePool: false})

	// Push some elements
	id1, _ := h.Push(10, 10)
	id2, _ := h.Push(20, 20)
	assert.NotEqual(t, id1, id2)
	assert.Greater(t, len(id1), 0)
	assert.Greater(t, len(id2), 0)

	// Clear the heap
	h.Clear()

	// Push after clear should get a new unique UUID
	id3, _ := h.Push(30, 30)
	assert.NotEqual(t, id1, id3)
	assert.NotEqual(t, id2, id3)
	assert.Greater(t, len(id3), 0)

	// Verify the element can be retrieved
	val3, err := h.GetValue(id3)
	assert.Nil(t, err)
	assert.Equal(t, 30, val3)
}

// -------------------------------- Skew Heap Benchmarks --------------------------------

func BenchmarkSkewHeapInsertion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewSkewHeap(
		data,
		func(a, b int) bool { return a < b },
		HeapConfig{UsePool: false, IDGenerator: &IntegerIDGenerator{NextID: 0}},
	)
	b.ReportAllocs()

	insertions := generateRandomNumbersv1(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Push(insertions[i], insertions[i])
	}
}

func BenchmarkSkewHeapDeletion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewSkewHeap(
		data,
		func(a, b int) bool { return a < b },
		HeapConfig{UsePool: false, IDGenerator: &IntegerIDGenerator{NextID: 0}},
	)

	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Pop()
	}
}

func BenchmarkSimpleSkewHeapInsertion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewSimpleSkewHeap(data, func(a, b int) bool { return a < b }, false)
	b.ReportAllocs()

	insertions := generateRandomNumbersv1(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Push(insertions[i], insertions[i])
	}
}

func BenchmarkSimpleSkewHeapDeletion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewSimpleSkewHeap(data, func(a, b int) bool { return a < b }, false)

	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Pop()
	}
}
