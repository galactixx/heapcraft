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
	data := []*HeapNode[int, int]{
		CreateHeapPairPtr(9, 9),
		CreateHeapPairPtr(4, 4),
		CreateHeapPairPtr(6, 6),
		CreateHeapPairPtr(1, 1),
		CreateHeapPairPtr(7, 7),
		CreateHeapPairPtr(3, 3),
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
	h := NewSimpleSkewHeap([]*HeapNode[int, int]{}, lt)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	_, err := h.Peek()
	assert.NotNil(t, err)

	input := []*HeapNode[int, int]{
		CreateHeapPairPtr(5, 5),
		CreateHeapPairPtr(2, 2),
		CreateHeapPairPtr(8, 8),
		CreateHeapPairPtr(3, 3),
		CreateHeapPairPtr(6, 6),
	}
	expectedOrder := []int{2, 3, 5, 6, 8}

	for _, pair := range input {
		h.Insert(pair.Value(), pair.Priority())
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
	data := []*HeapNode[int, int]{
		CreateHeapPairPtr(4, 4),
		CreateHeapPairPtr(1, 1),
		CreateHeapPairPtr(3, 3),
		CreateHeapPairPtr(2, 2),
	}
	h := NewSimpleSkewHeap(data, lt)
	assert.Equal(t, 4, h.Length())

	clone := h.Clone()
	assert.Equal(t, h.Length(), clone.Length())
	hPeekNode, _ := h.Peek()
	clonePeekNode, _ := clone.Peek()
	assert.Equal(t, hPeekNode.Value(), clonePeekNode.Value())

	h.Insert(0, 0)
	hPeekNodeAfterInsert, _ := h.Peek()
	assert.Equal(t, 0, hPeekNodeAfterInsert.Value())
	clonePeekNodeAfterInsert, _ := clone.Peek()
	assert.Equal(t, 1, clonePeekNodeAfterInsert.Value())

	h.Clear()
	assert.True(t, h.IsEmpty())
}

func TestPeekPopEmptySkew(t *testing.T) {
	h := NewSimpleSkewHeap([]*HeapNode[int, int]{}, lt)
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
	h := NewSimpleSkewHeap([]*HeapNode[int, int]{}, lt)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())

	h.Insert(10, 10)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, 1, h.Length())
}

func TestPeekValueAndPrioritySkew(t *testing.T) {
	h := NewSimpleSkewHeap([]*HeapNode[int, int]{}, lt)
	peekValueEmpty, _ := h.PeekValue()
	assert.Equal(t, 0, peekValueEmpty)
	peekPriorityEmpty, _ := h.PeekPriority()
	assert.Equal(t, 0, peekPriorityEmpty)

	h.Insert(42, 10)
	peekValue42, _ := h.PeekValue()
	assert.Equal(t, 42, peekValue42)
	peekPriority10, _ := h.PeekPriority()
	assert.Equal(t, 10, peekPriority10)

	h.Insert(15, 5)
	peekValue15, _ := h.PeekValue()
	assert.Equal(t, 15, peekValue15)
	peekPriority5, _ := h.PeekPriority()
	assert.Equal(t, 5, peekPriority5)

	h.Insert(100, 1)
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
	h := NewSimpleSkewHeap([]*HeapNode[int, int]{
		CreateHeapPairPtr(42, 10),
		CreateHeapPairPtr(15, 5),
		CreateHeapPairPtr(100, 1),
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
	h := NewSkewHeap([]*HeapNode[int, int]{
		CreateHeapPairPtr(42, 10),
		CreateHeapPairPtr(15, 5),
		CreateHeapPairPtr(100, 1),
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
	h := NewSkewHeap([]*HeapNode[int, int]{
		CreateHeapPairPtr(42, 10),
		CreateHeapPairPtr(15, 5),
		CreateHeapPairPtr(100, 1),
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
	h := NewSkewHeap([]*HeapNode[int, int]{
		CreateHeapPairPtr(1, 1),
		CreateHeapPairPtr(2, 2),
		CreateHeapPairPtr(3, 3),
		CreateHeapPairPtr(4, 4),
		CreateHeapPairPtr(5, 5),
		CreateHeapPairPtr(6, 6),
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
	h := NewSkewHeap([]*HeapNode[int, int]{
		CreateHeapPairPtr(1, 1),
		CreateHeapPairPtr(2, 2),
		CreateHeapPairPtr(3, 3),
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
	data := make([]*HeapNode[int, int], 0)
	heap := NewSkewHeap(data, func(a, b int) bool { return a < b })
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

func BenchmarkSkewHeapDeletion(b *testing.B) {
	data := make([]*HeapNode[int, int], 0)
	heap := NewSkewHeap(data, func(a, b int) bool { return a < b })

	for i := 0; i < b.N; i++ {
		heap.Insert(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Pop()
	}
}

func BenchmarkSimpleSkewHeapInsertion(b *testing.B) {
	N := 10_000
	data := make([]*HeapNode[int, int], 0)
	heap := NewSimpleSkewHeap(data, func(a, b int) bool { return a < b })
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

func BenchmarkSimpleSkewHeapDeletion(b *testing.B) {
	data := make([]*HeapNode[int, int], 0)
	heap := NewSimpleSkewHeap(data, func(a, b int) bool { return a < b })

	for i := 0; i < b.N; i++ {
		heap.Insert(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Pop()
	}
}
