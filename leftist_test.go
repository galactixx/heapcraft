package heapcraft

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// helper to pop all elements into a slice of values
func collectPop(h *SimpleLeftistHeap[int, int]) []int {
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
	data := []*HeapNode[int, int]{
		CreateHeapPairPtr(8, 8),
		CreateHeapPairPtr(3, 3),
		CreateHeapPairPtr(5, 5),
		CreateHeapPairPtr(1, 1),
		CreateHeapPairPtr(7, 7),
		CreateHeapPairPtr(2, 2),
	}
	h := NewSimpleLeftistHeap(data, lt)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(data), h.Length())

	expected := []int{1, 2, 3, 5, 7, 8}
	actual := collectPop(h)
	assert.Equal(t, expected, actual)
	assert.True(t, h.IsEmpty())

	_, err := h.Pop()
	assert.NotNil(t, err)
}

func TestInsertPopPeekLenIsEmptyLeftist(t *testing.T) {
	h := NewSimpleLeftistHeap([]*HeapNode[int, int]{}, lt)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	_, err := h.Peek()
	assert.NotNil(t, err)

	input := []*HeapNode[int, int]{
		CreateHeapPairPtr(6, 6),
		CreateHeapPairPtr(4, 4),
		CreateHeapPairPtr(9, 9),
		CreateHeapPairPtr(2, 2),
		CreateHeapPairPtr(5, 5),
	}
	expectedOrder := []int{2, 4, 5, 6, 9}

	for _, pair := range input {
		h.Insert(pair.Value(), pair.Priority())
	}

	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(input), h.Length())
	peekPair, err := h.Peek()
	assert.Nil(t, err)
	assert.Equal(t, 2, peekPair.Value())

	for i, expected := range expectedOrder {
		popped, err := h.Pop()
		assert.Nil(t, err)
		assert.Equal(t, expected, popped.Value())
		assert.Equal(t, expected, popped.Priority())
		assert.Equal(t, len(input)-(i+1), h.Length())
	}

	assert.True(t, h.IsEmpty())
	_, err = h.Peek()
	assert.NotNil(t, err)
}

func TestClearCloneLeftist(t *testing.T) {
	data := []*HeapNode[int, int]{
		CreateHeapPairPtr(4, 4),
		CreateHeapPairPtr(1, 1),
		CreateHeapPairPtr(3, 3),
		CreateHeapPairPtr(2, 2),
	}
	h := NewSimpleLeftistHeap(data, lt)
	assert.Equal(t, 4, h.Length())

	clone := h.Clone()
	assert.Equal(t, h.Length(), clone.Length())
	hPeek, _ := h.Peek()
	clonePeek, _ := clone.Peek()
	assert.Equal(t, hPeek.Value(), clonePeek.Value())

	h.Insert(0, 0)
	hPeek, _ = h.Peek()
	assert.Equal(t, 0, hPeek.Value())
	clonePeek, _ = clone.Peek()
	assert.Equal(t, 1, clonePeek.Value())

	h.Clear()
	assert.True(t, h.IsEmpty())
}

func TestPeekPopEmptyLeftist(t *testing.T) {
	h := NewSimpleLeftistHeap([]*HeapNode[int, int]{}, lt)
	_, err := h.Peek()
	assert.NotNil(t, err)
	_, err = h.Pop()
	assert.NotNil(t, err)
	_, err = h.PopValue()
	assert.NotNil(t, err)
	_, err = h.PopPriority()
	assert.NotNil(t, err)
}

func TestLengthIsEmptyLeftist(t *testing.T) {
	h := NewSimpleLeftistHeap([]*HeapNode[int, int]{}, lt)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())

	h.Insert(10, 10)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, 1, h.Length())
}

func TestPeekValueAndPriorityLeftist(t *testing.T) {
	h := NewSimpleLeftistHeap([]*HeapNode[int, int]{}, lt)
	_, err := h.PeekValue()
	assert.NotNil(t, err)
	_, err = h.PeekPriority()
	assert.NotNil(t, err)

	h.Insert(42, 10)
	val, err := h.PeekValue()
	assert.Nil(t, err)
	assert.Equal(t, 42, val)
	pri, err := h.PeekPriority()
	assert.Nil(t, err)
	assert.Equal(t, 10, pri)

	h.Insert(15, 5)
	val, err = h.PeekValue()
	assert.Nil(t, err)
	assert.Equal(t, 15, val)
	pri, err = h.PeekPriority()
	assert.Nil(t, err)
	assert.Equal(t, 5, pri)

	h.Insert(100, 1)
	val, err = h.PeekValue()
	assert.Nil(t, err)
	assert.Equal(t, 100, val)
	pri, err = h.PeekPriority()
	assert.Nil(t, err)
	assert.Equal(t, 1, pri)

	h.Pop()
	val, err = h.PeekValue()
	assert.Nil(t, err)
	assert.Equal(t, 15, val)
	pri, err = h.PeekPriority()
	assert.Nil(t, err)
	assert.Equal(t, 5, pri)

	h.Clear()
	_, err = h.PeekValue()
	assert.NotNil(t, err)
	_, err = h.PeekPriority()
	assert.NotNil(t, err)
}

func TestPopValueAndPriorityLeftist(t *testing.T) {
	h := NewSimpleLeftistHeap([]*HeapNode[int, int]{
		CreateHeapPairPtr(42, 10),
		CreateHeapPairPtr(15, 5),
		CreateHeapPairPtr(100, 1),
	}, lt)

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
	data := []*HeapNode[int, int]{
		CreateHeapPairPtr(8, 8),
		CreateHeapPairPtr(3, 3),
		CreateHeapPairPtr(5, 5),
	}
	h := NewLeftistHeap(data, lt)
	assert.Equal(t, 3, h.Length())
	val, err := h.PeekValue()
	assert.Nil(t, err)
	assert.Equal(t, 3, val)
}

func TestLeftistHeapGetOperations(t *testing.T) {
	data := []*HeapNode[int, int]{
		CreateHeapPairPtr(8, 8),
		CreateHeapPairPtr(3, 3),
		CreateHeapPairPtr(5, 5),
	}
	h := NewLeftistHeap(data, lt)

	element, err := h.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, 8, element.Value())
	assert.Equal(t, 8, element.Priority())

	value, err := h.GetValue(2)
	assert.NoError(t, err)
	assert.Equal(t, 3, value)

	priority, err := h.GetPriority(3)
	assert.NoError(t, err)
	assert.Equal(t, 5, priority)

	_, err = h.Get(999)
	assert.Error(t, err)
}

func TestLeftistHeapUpdateOperations(t *testing.T) {
	data := []*HeapNode[int, int]{
		CreateHeapPairPtr(8, 8),
		CreateHeapPairPtr(3, 3),
		CreateHeapPairPtr(5, 5),
	}
	h := NewLeftistHeap(data, lt)

	err := h.UpdateValue(1, 10)
	assert.NoError(t, err)
	value, _ := h.GetValue(1)
	assert.Equal(t, 10, value)

	err = h.UpdatePriority(2, 1)
	assert.NoError(t, err)
	pri, err := h.PeekPriority()
	assert.Nil(t, err)
	assert.Equal(t, 1, pri)

	err = h.UpdateValue(999, 10)
	assert.Error(t, err)
	err = h.UpdatePriority(999, 10)
	assert.Error(t, err)
}

func TestLeftistHeapInsertAndPop(t *testing.T) {
	h := NewLeftistHeap([]*HeapNode[int, int]{}, lt)

	h.Insert(5, 5)
	h.Insert(3, 3)
	h.Insert(7, 7)
	val, err := h.PeekValue()
	assert.Nil(t, err)
	assert.Equal(t, 3, val)

	popped, err := h.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 3, popped.Value())
	assert.Equal(t, 3, popped.Priority())
	val, err = h.PeekValue()
	assert.Nil(t, err)
	assert.Equal(t, 5, val)
}

func TestLeftistHeapClearAndClone(t *testing.T) {
	data := []*HeapNode[int, int]{
		CreateHeapPairPtr(8, 8),
		CreateHeapPairPtr(3, 3),
	}
	h := NewLeftistHeap(data, lt)

	clone := h.Clone()
	assert.Equal(t, h.Length(), clone.Length())
	hPeek, _ := h.Peek()
	clonePeek, _ := clone.Peek()
	assert.Equal(t, hPeek.Value(), clonePeek.Value())

	h.Clear()
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	_, err := h.Peek()
	assert.NotNil(t, err)
}

func TestLeftistHeapComplexUpdate(t *testing.T) {
	data := []*HeapNode[int, int]{
		CreateHeapPairPtr(8, 8),
		CreateHeapPairPtr(3, 3),
		CreateHeapPairPtr(5, 5),
		CreateHeapPairPtr(1, 1),
	}
	h := NewLeftistHeap(data, lt)

	err := h.UpdatePriority(2, 10)
	assert.NoError(t, err)
	val, err := h.PeekValue()
	assert.Nil(t, err)
	assert.Equal(t, 1, val)

	err = h.UpdatePriority(4, 0)
	assert.NoError(t, err)
	pri, err := h.PeekPriority()
	assert.Nil(t, err)
	assert.Equal(t, 0, pri)

	values := make([]int, 0)
	for !h.IsEmpty() {
		pri, err := h.PopPriority()
		assert.Nil(t, err)
		values = append(values, pri)
	}
	assert.True(t, sort.IntsAreSorted(values))
}

func TestLeftistHeapUpdatePriorityPositions(t *testing.T) {
	data := []*HeapNode[int, int]{
		CreateHeapPairPtr(1, 1),
		CreateHeapPairPtr(2, 2),
		CreateHeapPairPtr(3, 3),
		CreateHeapPairPtr(4, 4),
		CreateHeapPairPtr(5, 5),
		CreateHeapPairPtr(6, 6),
		CreateHeapPairPtr(7, 7),
	}
	h := NewLeftistHeap(data, lt)

	val, err := h.PeekValue()
	assert.Nil(t, err)
	assert.Equal(t, 1, val)
	rootID := uint(1)
	leftChildID := uint(2)
	rightChildID := uint(3)
	leafID := uint(4)

	err = h.UpdatePriority(rootID, 10)
	assert.NoError(t, err)
	val, err = h.PeekValue()
	assert.Nil(t, err)
	assert.Equal(t, 2, val)
	value, _ := h.GetValue(rootID)
	assert.Equal(t, 1, value)

	err = h.UpdatePriority(leafID, 0)
	assert.NoError(t, err)
	pri, err := h.PeekPriority()
	assert.Nil(t, err)
	assert.Equal(t, 0, pri)
	value, _ = h.GetValue(leafID)
	assert.Equal(t, 4, value)

	err = h.UpdatePriority(leftChildID, 8)
	assert.NoError(t, err)
	value, _ = h.GetValue(leftChildID)
	assert.Equal(t, 2, value)
	priority, _ := h.GetPriority(leftChildID)
	assert.Equal(t, 8, priority)

	err = h.UpdatePriority(rightChildID, 9)
	assert.NoError(t, err)
	value, _ = h.GetValue(rightChildID)
	assert.Equal(t, 3, value)
	priority, _ = h.GetPriority(rightChildID)
	assert.Equal(t, 9, priority)

	values := make([]int, 0)
	for !h.IsEmpty() {
		pri, err := h.PopPriority()
		assert.Nil(t, err)
		values = append(values, pri)
	}
	assert.True(t, sort.IntsAreSorted(values))
}

// Leftist Heap Benchmarks
func BenchmarkLeftistHeapInsertion(b *testing.B) {
	N := 10_000
	data := make([]*HeapNode[int, int], 0)
	heap := NewLeftistHeap(data, func(a, b int) bool { return a < b })
	b.ReportAllocs()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var num int
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			num = r.Intn(N)
			heap.Insert(num, num)
		}
	})
}

func BenchmarkLeftistHeapDeletion(b *testing.B) {
	data := make([]*HeapNode[int, int], 0)
	heap := NewLeftistHeap(data, func(a, b int) bool { return a < b })

	for i := 0; i < b.N; i++ {
		heap.Insert(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			heap.Pop()
		}
	})
}

func BenchmarkSimpleLeftistHeapInsertion(b *testing.B) {
	N := 10_000
	data := make([]*HeapNode[int, int], 0)
	heap := NewSimpleLeftistHeap(data, func(a, b int) bool { return a < b })
	b.ReportAllocs()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var num int
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			num = r.Intn(N)
			heap.Insert(num, num)
		}
	})
}

func BenchmarkSimpleLeftistHeapDeletion(b *testing.B) {
	data := make([]*HeapNode[int, int], 0)
	heap := NewSimpleLeftistHeap(data, func(a, b int) bool { return a < b })

	for i := 0; i < b.N; i++ {
		heap.Insert(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			heap.Pop()
		}
	})
}

func TestLeftistHeapInsertReturnsID(t *testing.T) {
	h := NewLeftistHeap([]*HeapNode[int, int]{}, lt)

	// Test sequential ID assignment
	id1 := h.Insert(10, 10)
	id2 := h.Insert(20, 20)
	id3 := h.Insert(30, 30)

	assert.Equal(t, uint(1), id1)
	assert.Equal(t, uint(2), id2)
	assert.Equal(t, uint(3), id3)

	// Verify elements can be retrieved using IDs
	val1, _ := h.GetValue(id1)
	val2, _ := h.GetValue(id2)
	val3, _ := h.GetValue(id3)
	assert.Equal(t, 10, val1)
	assert.Equal(t, 20, val2)
	assert.Equal(t, 30, val3)

	// Test ID continues after operations
	h.Pop()
	id4 := h.Insert(40, 40)
	assert.Equal(t, uint(4), id4)
}

func TestLeftistHeapInsertIDAfterClear(t *testing.T) {
	h := NewLeftistHeap([]*HeapNode[int, int]{}, lt)

	id1 := h.Insert(10, 10)
	h.Clear()
	id2 := h.Insert(20, 20)

	assert.Equal(t, uint(1), id1)
	assert.Equal(t, uint(1), id2) // Should reset to 1
}

func TestSimpleLeftistHeapInsertNoID(t *testing.T) {
	h := NewSimpleLeftistHeap([]*HeapNode[int, int]{}, lt)

	// SimpleLeftistHeap Insert should not return ID
	h.Insert(10, 10)
	h.Insert(20, 20)

	assert.Equal(t, 2, h.Length())
	val1, _ := h.PopValue()
	val2, _ := h.PopValue()
	assert.Equal(t, 10, val1)
	assert.Equal(t, 20, val2)
}
