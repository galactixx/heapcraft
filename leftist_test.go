package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// helper to pop all elements into a slice
func collectPop[T any](h *LeftistHeap[T]) []T {
	var result []T
	for !h.IsEmpty() {
		valPtr := h.Pop()
		result = append(result, *valPtr)
	}
	return result
}

func TestNewLeftistHeap_PopOrder(t *testing.T) {
	data := []int{8, 3, 5, 1, 7, 2}
	h := NewLeftistHeap(data, lt)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(data), h.Length())

	expected := []int{1, 2, 3, 5, 7, 8}
	actual := collectPop(&h)
	assert.Equal(t, expected, actual)
	assert.True(t, h.IsEmpty())

	assert.Nil(t, h.Pop())
}

func TestInsertPopPeekLenIsEmptyLeftist(t *testing.T) {
	h := NewLeftistHeap([]int{}, lt)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	assert.Nil(t, h.Peek())

	input := []int{6, 4, 9, 2, 5}
	expectedOrder := []int{2, 4, 5, 6, 9}

	for _, v := range input {
		h.Insert(v)
	}

	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(input), h.Length())
	assert.Equal(t, 2, *h.Peek())

	for i, expected := range expectedOrder {
		popped := h.Pop()
		assert.NotNil(t, popped)
		assert.Equal(t, expected, *popped)
		assert.Equal(t, len(input)-(i+1), h.Length())
	}

	assert.True(t, h.IsEmpty())
	assert.Nil(t, h.Peek())
}

func TestClearCloneDeepCloneLeftist(t *testing.T) {
	data := []int{4, 1, 3, 2}
	h := NewLeftistHeap(data, lt)
	assert.Equal(t, 4, h.Length())

	clone := h.Clone()
	assert.Equal(t, h.Length(), clone.Length())
	assert.Equal(t, *h.Peek(), *clone.Peek())

	h.Insert(0)
	assert.Equal(t, 0, *h.Peek())
	assert.Equal(t, 1, *clone.Peek())

	h2 := NewLeftistHeap([]int{7, 5, 9}, lt)
	deep := h2.DeepClone()
	assert.Equal(t, h2.Length(), deep.Length())
	assert.Equal(t, *h2.Peek(), *deep.Peek())

	h2.Insert(3)
	assert.Equal(t, 3, *h2.Peek())
	assert.Equal(t, 5, *deep.Peek())

	h2.Clear()
	assert.True(t, h2.IsEmpty())
}

func TestPeekPopEmptyLeftist(t *testing.T) {
	h := NewLeftistHeap([]int{}, lt)
	assert.Nil(t, h.Peek())
	assert.Nil(t, h.Pop())
}

func TestLengthIsEmptyLeftist(t *testing.T) {
	h := NewLeftistHeap([]int{}, lt)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())

	h.Insert(10)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, 1, h.Length())
}
