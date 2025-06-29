package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// helper to pop all elements into a slice of values
func collectPop(h *LeftistHeap[int, int]) []int {
	result := make([]int, 0)
	for !h.IsEmpty() {
		val, err := h.PopValue()
		if err == nil {
			result = append(result, val)
		}
	}
	return result
}

func TestNewLeftistHeapPopOrder(t *testing.T) {
	data := []HeapNode[int, int]{
		CreateHeapNode(8, 8),
		CreateHeapNode(3, 3),
		CreateHeapNode(5, 5),
		CreateHeapNode(1, 1),
		CreateHeapNode(7, 7),
		CreateHeapNode(2, 2),
	}
	h := NewLeftistHeap(data, lt, false)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(data), h.Length())

	expected := []int{1, 2, 3, 5, 7, 8}
	actual := collectPop(h)
	assert.Equal(t, expected, actual)
	assert.True(t, h.IsEmpty())

	_, _, err := h.Pop()
	assert.NotNil(t, err)
}

func TestInsertPopPeekLenIsEmptyLeftist(t *testing.T) {
	h := NewLeftistHeap([]HeapNode[int, int]{}, lt, false)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())

	input := []int{6, 4, 9, 2, 5}
	expectedOrder := []int{2, 4, 5, 6, 9}

	for _, val := range input {
		h.Push(val, val)
	}

	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(input), h.Length())
	value, priority, _ := h.Peek()
	assert.Equal(t, 2, value)
	assert.Equal(t, 2, priority)

	for i, expected := range expectedOrder {
		value, priority, _ = h.Pop()
		assert.Equal(t, expected, value)
		assert.Equal(t, expected, priority)
		assert.Equal(t, len(input)-(i+1), h.Length())
	}

	assert.True(t, h.IsEmpty())
}

func TestClearCloneLeftist(t *testing.T) {
	data := []HeapNode[int, int]{
		CreateHeapNode(4, 4),
		CreateHeapNode(1, 1),
		CreateHeapNode(3, 3),
		CreateHeapNode(2, 2),
	}
	h := NewLeftistHeap(data, lt, false)
	clone := h.Clone()
	assert.Equal(t, h.Length(), clone.Length())

	h.Push(0, 0)
	clone.Push(5, 5)
	assert.Equal(t, 5, clone.Length())
	assert.Equal(t, 5, h.Length())

	h.Clear()
	assert.True(t, h.IsEmpty())
	assert.False(t, clone.IsEmpty())
}

func TestLeftistHeap_Clone(t *testing.T) {
	h := NewLeftistHeap([]HeapNode[int, int]{}, lt, false)
	h.Push(5, 5)
	h.Push(3, 3)
	h.Push(7, 7)
	h.Push(1, 1)
	h.Push(9, 9)

	clone := h.Clone()
	originalElements := collectPop(h)
	cloneElements := collectPop(clone)
	assert.Equal(t, originalElements, cloneElements)

	h = NewLeftistHeap([]HeapNode[int, int]{}, lt, false)
	h.Push(5, 5)
	h.Push(3, 3)
	clone = h.Clone()
	clone.Push(1, 1)
	assert.Equal(t, 2, h.Length())
	assert.Equal(t, 3, clone.Length())
	val, _ := clone.PopValue()
	assert.Equal(t, 1, val)
}

func TestFullLeftistHeap_Clone(t *testing.T) {
	h := NewFullLeftistHeap([]HeapNode[int, int]{}, lt, HeapConfig{UsePool: false})
	id1, _ := h.Push(5, 5)
	id2, _ := h.Push(3, 3)
	id3, _ := h.Push(7, 7)
	id4, _ := h.Push(1, 1)
	id5, _ := h.Push(9, 9)

	clone := h.Clone()
	for _, id := range []string{id1, id2, id3, id4, id5} {
		val1, _ := h.GetValue(id)
		val2, _ := clone.GetValue(id)
		assert.Equal(t, val1, val2)
	}

	h.Clear()
	h.Push(5, 5)
	h.Push(3, 3)
	clone = h.Clone()
	newID, _ := clone.Push(1, 1)
	assert.Equal(t, 2, h.Length())
	assert.Equal(t, 3, clone.Length())

	val, _ := clone.PopValue()
	assert.Equal(t, 1, val)
	_, _, err := h.Get(newID)
	assert.Error(t, err)
}

func TestPeekPopEmptyLeftist(t *testing.T) {
	h := NewLeftistHeap([]HeapNode[int, int]{}, lt, false)
	_, _, err := h.Peek()
	assert.NotNil(t, err)
	_, _, err = h.Pop()
	assert.NotNil(t, err)
	_, err = h.PopValue()
	assert.NotNil(t, err)
	_, err = h.PopPriority()
	assert.NotNil(t, err)
}

func TestLengthIsEmptyLeftist(t *testing.T) {
	h := NewLeftistHeap([]HeapNode[int, int]{}, lt, false)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())

	h.Push(10, 10)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, 1, h.Length())
}

func TestPeekValueAndPriorityLeftist(t *testing.T) {
	h := NewLeftistHeap([]HeapNode[int, int]{}, lt, false)
	_, err := h.PeekValue()
	assert.NotNil(t, err)

	h.Push(42, 10)
	val, _ := h.PeekValue()
	pri, _ := h.PeekPriority()
	assert.Equal(t, 42, val)
	assert.Equal(t, 10, pri)

	h.Push(15, 5)
	val, _ = h.PeekValue()
	pri, _ = h.PeekPriority()
	assert.Equal(t, 15, val)
	assert.Equal(t, 5, pri)

	h.Push(100, 1)
	val, _ = h.PeekValue()
	pri, _ = h.PeekPriority()
	assert.Equal(t, 100, val)
	assert.Equal(t, 1, pri)

	h.Pop()
	val, _ = h.PeekValue()
	pri, _ = h.PeekPriority()
	assert.Equal(t, 15, val)
	assert.Equal(t, 5, pri)
}

func TestPopValueAndPriorityLeftist(t *testing.T) {
	h := NewLeftistHeap([]HeapNode[int, int]{
		CreateHeapNode(42, 10),
		CreateHeapNode(15, 5),
		CreateHeapNode(100, 1),
	}, lt, false)

	val, err := h.PopValue()
	assert.Nil(t, err)
	assert.Equal(t, 100, val)
	peekVal, err := h.PeekValue()
	assert.Nil(t, err)
	assert.Equal(t, 15, peekVal)

	pri, err := h.PopPriority()
	assert.Nil(t, err)
	assert.Equal(t, 5, pri)
	peekVal, err = h.PeekValue()
	assert.Nil(t, err)
	assert.Equal(t, 42, peekVal)

	h.Clear()
	_, err = h.PopValue()
	assert.NotNil(t, err)
	_, err = h.PopPriority()
	assert.NotNil(t, err)
}

func TestNewLeftistHeapConstruction(t *testing.T) {
	data := []HeapNode[int, int]{
		CreateHeapNode(8, 8),
		CreateHeapNode(3, 3),
		CreateHeapNode(5, 5),
	}
	h := NewFullLeftistHeap(data, lt, HeapConfig{UsePool: false})
	assert.Equal(t, 3, h.Length())
	val, err := h.PeekValue()
	assert.Nil(t, err)
	assert.Equal(t, 3, val)
}

func TestLeftistHeapInsertAndPop(t *testing.T) {
	h := NewFullLeftistHeap([]HeapNode[int, int]{}, lt, HeapConfig{UsePool: false})

	h.Push(5, 5)
	h.Push(3, 3)
	h.Push(7, 7)
	val, err := h.PeekValue()
	assert.Nil(t, err)
	assert.Equal(t, 3, val)

	value, priority, err := h.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 3, value)
	assert.Equal(t, 3, priority)
	val, err = h.PeekValue()
	assert.Nil(t, err)
	assert.Equal(t, 5, val)
}

func TestLeftistHeapClearAndClone(t *testing.T) {
	data := []HeapNode[int, int]{
		CreateHeapNode(8, 8),
		CreateHeapNode(3, 3),
	}
	h := NewFullLeftistHeap(data, lt, HeapConfig{UsePool: false})

	clone := h.Clone()
	assert.Equal(t, h.Length(), clone.Length())
	value, priority, _ := h.Peek()
	cloneValue, clonePriority, _ := clone.Peek()
	assert.Equal(t, value, cloneValue)
	assert.Equal(t, priority, clonePriority)

	h.Clear()
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	_, _, err := h.Peek()
	assert.NotNil(t, err)
}

func TestLeftistHeapInsertReturnsID(t *testing.T) {
	h := NewFullLeftistHeap([]HeapNode[int, int]{}, lt, HeapConfig{UsePool: false})

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

	// Verify elements can be retrieved using IDs
	val1, _ := h.GetValue(id1)
	val2, _ := h.GetValue(id2)
	val3, _ := h.GetValue(id3)
	assert.Equal(t, 10, val1)
	assert.Equal(t, 20, val2)
	assert.Equal(t, 30, val3)

	// Test ID continues after operations
	h.Pop()
	id4, _ := h.Push(40, 40)
	assert.NotEqual(t, id1, id4)
	assert.NotEqual(t, id2, id4)
	assert.NotEqual(t, id3, id4)
}

func TestLeftistHeapInsertIDAfterClear(t *testing.T) {
	h := NewFullLeftistHeap([]HeapNode[int, int]{}, lt, HeapConfig{UsePool: false})

	id1, _ := h.Push(10, 10)
	h.Clear()
	id2, _ := h.Push(20, 20)

	// Both should be unique UUIDs
	assert.NotEqual(t, id1, id2)
	assert.Greater(t, len(id1), 0)
	assert.Greater(t, len(id2), 0)
}

func TestLeftistHeap_InsertNoID(t *testing.T) {
	h := NewLeftistHeap([]HeapNode[int, int]{}, lt, false)

	// LeftistHeap Push should not return ID
	h.Push(10, 10)
	h.Push(20, 20)

	assert.Equal(t, 2, h.Length())
	val1, _ := h.PopValue()
	val2, _ := h.PopValue()
	assert.Equal(t, 10, val1)
	assert.Equal(t, 20, val2)
}

// -------------------------------- Leftist Heap Benchmarks --------------------------------

func BenchmarkFullLeftistHeap_Insertion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewFullLeftistHeap(
		data,
		lt,
		HeapConfig{UsePool: false, IDGenerator: &IntegerIDGenerator{NextID: 0}},
	)

	insertions := generateRandomNumbersv1(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Push(insertions[i], insertions[i])
	}
}

func BenchmarkFullLeftistHeap_Deletion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewFullLeftistHeap(
		data,
		lt,
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

func BenchmarkLeftistHeap_Insertion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewLeftistHeap(data, lt, false)

	insertions := generateRandomNumbersv1(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Push(insertions[i], insertions[i])
	}
}

func BenchmarkLeftistHeap_Deletion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewLeftistHeap(data, lt, false)

	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Pop()
	}
}
