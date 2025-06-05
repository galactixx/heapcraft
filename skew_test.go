package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ltInt(a, b int) bool { return a < b }

// helper to pop all elements into a slice
func collectSkew[T any](h *SkewHeap[T]) []T {
	var result []T
	for !h.IsEmpty() {
		valPtr := h.Pop()
		result = append(result, *valPtr)
	}
	return result
}

func TestNewSkewHeapPopOrder(t *testing.T) {
	data := []int{9, 4, 6, 1, 7, 3}
	h := NewSkewHeap(data, ltInt)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(data), h.Length())

	expected := []int{1, 3, 4, 6, 7, 9}
	actual := collectSkew(&h)
	assert.Equal(t, expected, actual)
	assert.True(t, h.IsEmpty())

	assert.Nil(t, h.Pop())
}

func TestInsertPopPeekLenIsEmptySkew(t *testing.T) {
	h := NewSkewHeap([]int{}, ltInt)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	assert.Nil(t, h.Peek())

	input := []int{5, 2, 8, 3, 6}
	expectedOrder := []int{2, 3, 5, 6, 8}

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

func TestPeekPopEmptySkew(t *testing.T) {
	h := NewSkewHeap([]int{}, ltInt)
	assert.Nil(t, h.Peek())
	assert.Nil(t, h.Pop())
}

func TestLengthIsEmptySkew(t *testing.T) {
	h := NewSkewHeap([]int{}, ltInt)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())

	h.Insert(10)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, 1, h.Length())
}
