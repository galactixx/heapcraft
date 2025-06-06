package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ltSkew(a, b int) bool { return a < b }

func collectSimpleSkew(h *SimpleSkewHeap[int, int]) []int {
	result := make([]int, 0)
	for !h.IsEmpty() {
		result = append(result, *h.PopValue())
	}
	return result
}

func collectSkewHeap(h *SkewHeap[int, int]) []int {
	result := make([]int, 0)
	for !h.IsEmpty() {
		result = append(result, *h.PopValue())
	}
	return result
}

func TestNewSkewHeapPopOrder(t *testing.T) {
	data := []*HeapPair[int, int]{
		CreateHeapPair(9, 9),
		CreateHeapPair(4, 4),
		CreateHeapPair(6, 6),
		CreateHeapPair(1, 1),
		CreateHeapPair(7, 7),
		CreateHeapPair(3, 3),
	}
	h := NewSimpleSkewHeap(data, ltSkew)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(data), h.Length())

	expected := []int{1, 3, 4, 6, 7, 9}
	actual := collectSimpleSkew(&h)
	assert.Equal(t, expected, actual)
	assert.True(t, h.IsEmpty())

	assert.Nil(t, h.Pop())
}

func TestInsertPopPeekLenIsEmptySkew(t *testing.T) {
	h := NewSimpleSkewHeap([]*HeapPair[int, int]{}, ltSkew)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	assert.Nil(t, h.Peek())

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
	assert.Equal(t, 2, h.Peek().Value())

	for i, expected := range expectedOrder {
		popped := h.Pop()
		assert.NotNil(t, popped)
		assert.Equal(t, expected, popped.Value())
		assert.Equal(t, expected, popped.Priority())
		assert.Equal(t, len(input)-(i+1), h.Length())
	}

	assert.True(t, h.IsEmpty())
	assert.Nil(t, h.Peek())
}

func TestClearCloneSkew(t *testing.T) {
	data := []*HeapPair[int, int]{
		CreateHeapPair(4, 4),
		CreateHeapPair(1, 1),
		CreateHeapPair(3, 3),
		CreateHeapPair(2, 2),
	}
	h := NewSimpleSkewHeap(data, ltSkew)
	assert.Equal(t, 4, h.Length())

	clone := h.Clone()
	assert.Equal(t, h.Length(), clone.Length())
	assert.Equal(t, h.Peek().Value(), clone.Peek().Value())

	h.Insert(0, 0)
	assert.Equal(t, 0, h.Peek().Value())
	assert.Equal(t, 1, clone.Peek().Value())

	h.Clear()
	assert.True(t, h.IsEmpty())
}

func TestPeekPopEmptySkew(t *testing.T) {
	h := NewSimpleSkewHeap([]*HeapPair[int, int]{}, ltSkew)
	assert.Nil(t, h.Peek())
	assert.Nil(t, h.Pop())
	assert.Nil(t, h.PopValue())
	assert.Nil(t, h.PopPriority())
}

func TestLengthIsEmptySkew(t *testing.T) {
	h := NewSimpleSkewHeap([]*HeapPair[int, int]{}, ltSkew)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())

	h.Insert(10, 10)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, 1, h.Length())
}

func TestPeekValueAndPrioritySkew(t *testing.T) {
	h := NewSimpleSkewHeap([]*HeapPair[int, int]{}, ltSkew)
	assert.Nil(t, h.PeekValue())
	assert.Nil(t, h.PeekPriority())

	h.Insert(42, 10)
	assert.Equal(t, 42, *h.PeekValue())
	assert.Equal(t, 10, *h.PeekPriority())

	h.Insert(15, 5)
	assert.Equal(t, 15, *h.PeekValue())
	assert.Equal(t, 5, *h.PeekPriority())

	h.Insert(100, 1)
	assert.Equal(t, 100, *h.PeekValue())
	assert.Equal(t, 1, *h.PeekPriority())

	h.Pop()
	assert.Equal(t, 15, *h.PeekValue())
	assert.Equal(t, 5, *h.PeekPriority())

	h.Clear()
	assert.Nil(t, h.PeekValue())
	assert.Nil(t, h.PeekPriority())
}

func TestPopValueAndPrioritySkew(t *testing.T) {
	h := NewSimpleSkewHeap([]*HeapPair[int, int]{
		CreateHeapPair(42, 10),
		CreateHeapPair(15, 5),
		CreateHeapPair(100, 1),
	}, ltSkew)

	val := h.PopValue()
	assert.Equal(t, 100, *val)
	assert.Equal(t, 15, *h.PeekValue())

	pri := h.PopPriority()
	assert.Equal(t, 5, *pri)
	assert.Equal(t, 42, *h.PeekValue())

	h.Clear()
	assert.Nil(t, h.PopValue())
	assert.Nil(t, h.PopPriority())
}

func TestSkewHeapGetOperations(t *testing.T) {
	h := NewSkewHeap([]*HeapPair[int, int]{
		CreateHeapPair(42, 10),
		CreateHeapPair(15, 5),
		CreateHeapPair(100, 1),
	}, ltSkew)

	val, err := h.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 42, *val)

	pri, err := h.GetPriority(2)
	assert.Nil(t, err)
	assert.Equal(t, 5, *pri)

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
	h := NewSkewHeap([]*HeapPair[int, int]{
		CreateHeapPair(42, 10),
		CreateHeapPair(15, 5),
		CreateHeapPair(100, 1),
	}, ltSkew)

	err := h.UpdateValue(2, 25)
	assert.Nil(t, err)
	val, _ := h.GetValue(2)
	assert.Equal(t, 25, *val)

	err = h.UpdatePriority(1, 2)
	assert.Nil(t, err)
	pri, _ := h.GetPriority(1)
	assert.Equal(t, 2, *pri)

	err = h.UpdateValue(999, 0)
	assert.Error(t, err)
	err = h.UpdatePriority(0, 0)
	assert.Error(t, err)
}

func TestSkewHeapUpdatePriorityPositions(t *testing.T) {
	h := NewSkewHeap([]*HeapPair[int, int]{
		CreateHeapPair(1, 1),
		CreateHeapPair(2, 2),
		CreateHeapPair(3, 3),
		CreateHeapPair(4, 4),
		CreateHeapPair(5, 5),
		CreateHeapPair(6, 6),
	}, ltSkew)

	err := h.UpdatePriority(1, 7)
	assert.Nil(t, err)
	assert.Equal(t, 2, h.Peek().Value())

	err = h.UpdatePriority(4, 0)
	assert.Nil(t, err)
	assert.Equal(t, 4, h.Peek().Value())

	err = h.UpdatePriority(2, 8)
	assert.Nil(t, err)
	val, _ := h.GetValue(2)
	assert.Equal(t, 2, *val)

	err = h.UpdatePriority(3, 9)
	assert.Nil(t, err)
	val, _ = h.GetValue(3)
	assert.Equal(t, 3, *val)

	expected := []int{4, 5, 6, 1, 2, 3}
	actual := collectSkewHeap(&h)
	assert.Equal(t, expected, actual)
}

func TestSkewHeapParentPointers(t *testing.T) {
	h := NewSkewHeap([]*HeapPair[int, int]{
		CreateHeapPair(1, 1),
		CreateHeapPair(2, 2),
		CreateHeapPair(3, 3),
	}, ltSkew)

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
