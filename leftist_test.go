package heapcraft

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// helper to pop all elements into a slice of values
func collectPop(h *SimpleLeftistHeap[int, int]) []int {
	result := make([]int, 0)
	for !h.IsEmpty() {
		result = append(result, *h.PopValue())
	}
	return result
}

func TestNewLeftistHeapPopOrder(t *testing.T) {
	data := []*HeapPair[int, int]{
		CreateHeapPair(8, 8),
		CreateHeapPair(3, 3),
		CreateHeapPair(5, 5),
		CreateHeapPair(1, 1),
		CreateHeapPair(7, 7),
		CreateHeapPair(2, 2),
	}
	h := NewSimpleLeftistHeap(data, lt)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(data), h.Length())

	expected := []int{1, 2, 3, 5, 7, 8}
	actual := collectPop(h)
	assert.Equal(t, expected, actual)
	assert.True(t, h.IsEmpty())

	assert.Nil(t, h.Pop())
}

func TestInsertPopPeekLenIsEmptyLeftist(t *testing.T) {
	h := NewSimpleLeftistHeap([]*HeapPair[int, int]{}, lt)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	assert.Nil(t, h.Peek())

	input := []*HeapPair[int, int]{
		CreateHeapPair(6, 6),
		CreateHeapPair(4, 4),
		CreateHeapPair(9, 9),
		CreateHeapPair(2, 2),
		CreateHeapPair(5, 5),
	}
	expectedOrder := []int{2, 4, 5, 6, 9}

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

func TestClearCloneLeftist(t *testing.T) {
	data := []*HeapPair[int, int]{
		CreateHeapPair(4, 4),
		CreateHeapPair(1, 1),
		CreateHeapPair(3, 3),
		CreateHeapPair(2, 2),
	}
	h := NewSimpleLeftistHeap(data, lt)
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

func TestPeekPopEmptyLeftist(t *testing.T) {
	h := NewSimpleLeftistHeap([]*HeapPair[int, int]{}, lt)
	assert.Nil(t, h.Peek())
	assert.Nil(t, h.Pop())
	assert.Nil(t, h.PopValue())
	assert.Nil(t, h.PopPriority())
}

func TestLengthIsEmptyLeftist(t *testing.T) {
	h := NewSimpleLeftistHeap([]*HeapPair[int, int]{}, lt)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())

	h.Insert(10, 10)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, 1, h.Length())
}

func TestPeekValueAndPriorityLeftist(t *testing.T) {
	h := NewSimpleLeftistHeap([]*HeapPair[int, int]{}, lt)
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

func TestPopValueAndPriorityLeftist(t *testing.T) {
	h := NewSimpleLeftistHeap([]*HeapPair[int, int]{
		CreateHeapPair(42, 10),
		CreateHeapPair(15, 5),
		CreateHeapPair(100, 1),
	}, lt)

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

func TestNewLeftistHeapConstruction(t *testing.T) {
	data := []*HeapPair[int, int]{
		CreateHeapPair(8, 8),
		CreateHeapPair(3, 3),
		CreateHeapPair(5, 5),
	}
	h := NewLeftistHeap(data, lt)
	assert.Equal(t, 3, h.Length())
	assert.Equal(t, 3, *h.PeekValue())
}

func TestLeftistHeapGetOperations(t *testing.T) {
	data := []*HeapPair[int, int]{
		CreateHeapPair(8, 8),
		CreateHeapPair(3, 3),
		CreateHeapPair(5, 5),
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
	data := []*HeapPair[int, int]{
		CreateHeapPair(8, 8),
		CreateHeapPair(3, 3),
		CreateHeapPair(5, 5),
	}
	h := NewLeftistHeap(data, lt)

	err := h.UpdateValue(1, 10)
	assert.NoError(t, err)
	value, _ := h.GetValue(1)
	assert.Equal(t, 10, value)

	err = h.UpdatePriority(2, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, *h.PeekPriority())

	err = h.UpdateValue(999, 10)
	assert.Error(t, err)
	err = h.UpdatePriority(999, 10)
	assert.Error(t, err)
}

func TestLeftistHeapInsertAndPop(t *testing.T) {
	h := NewLeftistHeap([]*HeapPair[int, int]{}, lt)

	h.Insert(5, 5)
	h.Insert(3, 3)
	h.Insert(7, 7)
	assert.Equal(t, 3, *h.PeekValue())

	popped := h.Pop()
	assert.Equal(t, 3, popped.Value())
	assert.Equal(t, 3, popped.Priority())
	assert.Equal(t, 5, *h.PeekValue())
}

func TestLeftistHeapClearAndClone(t *testing.T) {
	data := []*HeapPair[int, int]{
		CreateHeapPair(8, 8),
		CreateHeapPair(3, 3),
	}
	h := NewLeftistHeap(data, lt)

	clone := h.Clone()
	assert.Equal(t, h.Length(), clone.Length())
	assert.Equal(t, h.Peek().Value(), clone.Peek().Value())

	h.Clear()
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	assert.Nil(t, h.Peek())
}

func TestLeftistHeapComplexUpdate(t *testing.T) {
	data := []*HeapPair[int, int]{
		CreateHeapPair(8, 8),
		CreateHeapPair(3, 3),
		CreateHeapPair(5, 5),
		CreateHeapPair(1, 1),
	}
	h := NewLeftistHeap(data, lt)

	err := h.UpdatePriority(2, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, *h.PeekValue())

	err = h.UpdatePriority(4, 0)
	assert.NoError(t, err)
	assert.Equal(t, 0, *h.PeekPriority())

	values := make([]int, 0)
	for !h.IsEmpty() {
		values = append(values, *h.PopPriority())
	}
	assert.True(t, sort.IntsAreSorted(values))
}

func TestLeftistHeapUpdatePriorityPositions(t *testing.T) {
	data := []*HeapPair[int, int]{
		CreateHeapPair(1, 1),
		CreateHeapPair(2, 2),
		CreateHeapPair(3, 3),
		CreateHeapPair(4, 4),
		CreateHeapPair(5, 5),
		CreateHeapPair(6, 6),
		CreateHeapPair(7, 7),
	}
	h := NewLeftistHeap(data, lt)

	assert.Equal(t, 1, *h.PeekValue())
	rootID := uint(1)
	leftChildID := uint(2)
	rightChildID := uint(3)
	leafID := uint(4)

	err := h.UpdatePriority(rootID, 10)
	assert.NoError(t, err)
	assert.Equal(t, 2, *h.PeekValue())
	value, _ := h.GetValue(rootID)
	assert.Equal(t, 1, value)

	err = h.UpdatePriority(leafID, 0)
	assert.NoError(t, err)
	assert.Equal(t, 0, *h.PeekPriority())
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
		values = append(values, *h.PopPriority())
	}
	assert.True(t, sort.IntsAreSorted(values))
}
