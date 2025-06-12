package heapcraft

import (
	"math/rand"
	"testing"
	"time"

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

func collectSkewHeap(h *SkewHeap[int, int]) []int {
	result := make([]int, 0)
	for !h.IsEmpty() {
		value, _ := h.PopValue()
		result = append(result, value)
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
	h := NewSimpleSkewHeap(data, lt)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(data), h.Length())

	expected := []int{1, 3, 4, 6, 7, 9}
	actual := collectSimpleSkew(h)
	assert.Equal(t, expected, actual)
	assert.True(t, h.IsEmpty())

	_, err := h.Pop()
	assert.NotNil(t, err)
}

func TestInsertPopPeekLenIsEmptySkew(t *testing.T) {
	h := NewSimpleSkewHeap([]HeapNode[int, int]{}, lt)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	_, err := h.Peek()
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
		h.Push(pair.Value(), pair.Priority())
	}

	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(input), h.Length())
	peekNode, _ := h.Peek()
	assert.Equal(t, 2, peekNode.Value())

	for i, expected := range expectedOrder {
		popped, err := h.Pop()
		assert.Nil(t, err)
		assert.NotNil(t, popped)
		assert.Equal(t, expected, popped.Value())
		assert.Equal(t, expected, popped.Priority())
		assert.Equal(t, len(input)-(i+1), h.Length())
	}

	assert.True(t, h.IsEmpty())
	_, err = h.Peek()
	assert.NotNil(t, err)
}

func TestClearCloneSkew(t *testing.T) {
	data := []HeapNode[int, int]{
		CreateHeapNode(4, 4),
		CreateHeapNode(1, 1),
		CreateHeapNode(3, 3),
		CreateHeapNode(2, 2),
	}
	h := NewSimpleSkewHeap(data, lt)
	assert.Equal(t, 4, h.Length())

	// Test basic cloning
	clone := h.Clone()
	assert.Equal(t, h.Length(), clone.Length())
	hPeekNode, _ := h.Peek()
	clonePeekNode, _ := clone.Peek()
	assert.Equal(t, hPeekNode.Value(), clonePeekNode.Value())

	// Test independence of clone
	h.Push(0, 0)
	hPeekNodeAfterInsert, _ := h.Peek()
	assert.Equal(t, 0, hPeekNodeAfterInsert.Value())
	clonePeekNodeAfterInsert, _ := clone.Peek()
	assert.Equal(t, 1, clonePeekNodeAfterInsert.Value())

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
	h := NewSimpleSkewHeap([]HeapNode[int, int]{}, lt)
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
	h = NewSimpleSkewHeap([]HeapNode[int, int]{}, lt)
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
	h := NewSkewHeap([]HeapNode[int, int]{}, lt)
	id1 := h.Push(5, 5)
	id2 := h.Push(3, 3)
	id3 := h.Push(7, 7)
	id4 := h.Push(1, 1)
	h.Push(9, 9)

	// Create a clone
	clone := h.Clone()

	// Test that all elements are preserved with their IDs
	for _, id := range []uint{id1, id2, id3, id4} {
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
	_, err := h.Get(newID)
	assert.Error(t, err)
	_, err = clone.Get(newID)
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
	h := NewSkewHeap([]HeapNode[int, int]{}, lt)
	id1 := h.Push(5, 5)
	id2 := h.Push(3, 3)
	id3 := h.Push(7, 7)
	id4 := h.Push(1, 1)
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
	h := NewSimpleSkewHeap([]HeapNode[int, int]{}, lt)
	_, err := h.Peek()
	assert.NotNil(t, err)
	_, err = h.Pop()
	assert.NotNil(t, err)
	_, err = h.PopValue()
	assert.NotNil(t, err)
	_, err = h.PopPriority()
	assert.NotNil(t, err)
}

func TestLengthIsEmptySkew(t *testing.T) {
	h := NewSimpleSkewHeap([]HeapNode[int, int]{}, lt)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())

	h.Push(10, 10)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, 1, h.Length())
}

func TestSimpleSkewHeapInsertNoID(t *testing.T) {
	h := NewSimpleSkewHeap([]HeapNode[int, int]{}, lt)

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
	h := NewSimpleSkewHeap([]HeapNode[int, int]{}, lt)
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

func TestPopValueAndPrioritySkew(t *testing.T) {
	h := NewSimpleSkewHeap([]HeapNode[int, int]{
		CreateHeapNode(42, 10),
		CreateHeapNode(15, 5),
		CreateHeapNode(100, 1),
	}, lt)

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

func TestSkewHeapGetOperations(t *testing.T) {
	h := NewSkewHeap([]HeapNode[int, int]{
		CreateHeapNode(42, 10),
		CreateHeapNode(15, 5),
		CreateHeapNode(100, 1),
	}, lt)

	val, err := h.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 42, val)

	pri, err := h.GetPriority(2)
	assert.Nil(t, err)
	assert.Equal(t, 5, pri)

	pair, err := h.Get(3)
	assert.Nil(t, err)
	assert.Equal(t, 100, pair.Value())
	assert.Equal(t, 1, pair.Priority())

	_, err = h.GetValue(4)
	assert.Error(t, err)
	_, err = h.GetPriority(0)
	assert.Error(t, err)
	_, err = h.Get(999)
	assert.Error(t, err)
}

func TestSkewHeapUpdateOperations(t *testing.T) {
	h := NewSkewHeap([]HeapNode[int, int]{
		CreateHeapNode(42, 10),
		CreateHeapNode(15, 5),
		CreateHeapNode(100, 1),
	}, lt)

	err := h.UpdateValue(2, 25)
	assert.Nil(t, err)
	val, _ := h.GetValue(2)
	assert.Equal(t, 25, val)

	err = h.UpdatePriority(1, 2)
	assert.Nil(t, err)
	pri, _ := h.GetPriority(1)
	assert.Equal(t, 2, pri)

	err = h.UpdateValue(999, 0)
	assert.Error(t, err)
	err = h.UpdatePriority(0, 0)
	assert.Error(t, err)
}

func TestSkewHeapUpdatePriorityPositions(t *testing.T) {
	h := NewSkewHeap([]HeapNode[int, int]{
		CreateHeapNode(1, 1),
		CreateHeapNode(2, 2),
		CreateHeapNode(3, 3),
		CreateHeapNode(4, 4),
		CreateHeapNode(5, 5),
		CreateHeapNode(6, 6),
	}, lt)

	err := h.UpdatePriority(1, 7)
	assert.Nil(t, err)
	peekNode1, err := h.Peek()
	assert.Nil(t, err)
	assert.Equal(t, 2, peekNode1.Value())

	err = h.UpdatePriority(4, 0)
	assert.Nil(t, err)
	peekNode2, err := h.Peek()
	assert.Nil(t, err)
	assert.Equal(t, 4, peekNode2.Value())

	err = h.UpdatePriority(2, 8)
	assert.Nil(t, err)
	val, _ := h.GetValue(2)
	assert.Equal(t, 2, val)

	err = h.UpdatePriority(3, 9)
	assert.Nil(t, err)
	val, _ = h.GetValue(3)
	assert.Equal(t, 3, val)

	expected := []int{4, 5, 6, 1, 2, 3}
	actual := collectSkewHeap(h)
	assert.Equal(t, expected, actual)
}

func TestSkewHeapParentPointers(t *testing.T) {
	h := NewSkewHeap([]HeapNode[int, int]{
		CreateHeapNode(1, 1),
		CreateHeapNode(2, 2),
		CreateHeapNode(3, 3),
	}, lt)

	assert.Nil(t, h.root.parent)
	assert.Equal(t, h.root, h.root.left.parent)
	assert.Equal(t, h.root, h.root.right.parent)

	err := h.UpdatePriority(2, 0)
	assert.Nil(t, err)
	assert.Nil(t, h.root.parent)

	assert.Equal(t, h.root, h.root.left.parent)
	assert.Equal(t, h.root.left, h.root.left.left.parent)

	h.Pop()
	assert.Nil(t, h.root.parent)
	if h.root.left != nil {
		assert.Equal(t, h.root, h.root.left.parent)
	}
	if h.root.right != nil {
		assert.Equal(t, h.root, h.root.right.parent)
	}
}

// Skew Heap Benchmarks
func BenchmarkSkewHeapInsertion(b *testing.B) {
	N := 10_000
	data := make([]HeapNode[int, int], 0)
	heap := NewSkewHeap(data, func(a, b int) bool { return a < b })
	b.ReportAllocs()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var num int
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			num = r.Intn(N)
			heap.Push(num, num)
		}
	})
}

func BenchmarkSkewHeapDeletion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewSkewHeap(data, func(a, b int) bool { return a < b })

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
	N := 10_000
	data := make([]HeapNode[int, int], 0)
	heap := NewSimpleSkewHeap(data, func(a, b int) bool { return a < b })
	b.ReportAllocs()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var num int
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			num = r.Intn(N)
			heap.Push(num, num)
		}
	})
}

func BenchmarkSimpleSkewHeapDeletion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewSimpleSkewHeap(data, func(a, b int) bool { return a < b })

	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Pop()
	}
}

func TestSkewHeapInsertReturnsID(t *testing.T) {
	h := NewSkewHeap([]HeapNode[int, int]{}, lt)

	// Test that Push returns sequential IDs starting from 1
	id1 := h.Push(10, 10)
	assert.Equal(t, uint(1), id1)

	id2 := h.Push(20, 20)
	assert.Equal(t, uint(2), id2)

	id3 := h.Push(30, 30)
	assert.Equal(t, uint(3), id3)

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

	// Test that IDs continue incrementing after operations
	h.Pop() // Remove one element
	id4 := h.Push(40, 40)
	assert.Equal(t, uint(4), id4)

	// Verify the new element can be retrieved
	val4, err := h.GetValue(id4)
	assert.Nil(t, err)
	assert.Equal(t, 40, val4)
}

func TestSkewHeapInsertIDAfterClear(t *testing.T) {
	h := NewSkewHeap([]HeapNode[int, int]{}, lt)

	// Push some elements
	id1 := h.Push(10, 10)
	id2 := h.Push(20, 20)
	assert.Equal(t, uint(1), id1)
	assert.Equal(t, uint(2), id2)

	// Clear the heap
	h.Clear()

	// Push after clear should start from ID 1 again
	id3 := h.Push(30, 30)
	assert.Equal(t, uint(1), id3)

	// Verify the element can be retrieved
	val3, err := h.GetValue(id3)
	assert.Nil(t, err)
	assert.Equal(t, 30, val3)
}

func TestSkewHeapInsertIDWithInitialData(t *testing.T) {
	data := []HeapNode[int, int]{
		CreateHeapNode(42, 10),
		CreateHeapNode(15, 5),
		CreateHeapNode(100, 1),
	}

	h := NewSkewHeap(data, lt)

	// The constructor should have assigned IDs 1, 2, 3
	val1, err := h.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 42, val1)

	val2, err := h.GetValue(2)
	assert.Nil(t, err)
	assert.Equal(t, 15, val2)

	val3, err := h.GetValue(3)
	assert.Nil(t, err)
	assert.Equal(t, 100, val3)

	// Next push should get ID 4
	id4 := h.Push(200, 200)
	assert.Equal(t, uint(4), id4)

	val4, err := h.GetValue(id4)
	assert.Nil(t, err)
	assert.Equal(t, 200, val4)
}
