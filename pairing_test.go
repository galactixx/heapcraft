package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func collectPairing[T any](h *PairingHeap[T]) []T {
	var result []T
	for !h.IsEmpty() {
		valPtr := h.Pop()
		result = append(result, *valPtr)
	}
	return result
}

func TestNewPairingHeapPopOrder(t *testing.T) {
	data := []int{9, 4, 6, 1, 7, 3}
	h := NewPairingHeap(data, lt)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(data), h.Length())

	expected := []int{1, 3, 4, 6, 7, 9}
	actual := collectPairing(&h)
	assert.Equal(t, expected, actual)
	assert.True(t, h.IsEmpty())

	assert.Nil(t, h.Pop())
}

func TestInsertPopPeekLenIsEmptyPairing(t *testing.T) {
	h := NewPairingHeap([]int{}, lt)
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

func TestClearCloneDeepClonePairing(t *testing.T) {
	data := []int{4, 1, 3, 2}
	h := NewPairingHeap(data, lt)
	assert.Equal(t, 4, h.Length())

	clone := h.Clone()
	assert.Equal(t, h.Length(), clone.Length())
	assert.Equal(t, *h.Peek(), *clone.Peek())

	h.Insert(0)
	assert.Equal(t, 0, *h.Peek())
	assert.Equal(t, 1, *clone.Peek())

	h2 := NewPairingHeap([]int{7, 5, 9}, lt)
	deep := h2.DeepClone()

	assert.Equal(t, h2.Length(), deep.Length())
	assert.Equal(t, *h2.Peek(), *deep.Peek())

	h2.Insert(3)
	assert.Equal(t, 3, *h2.Peek())
	assert.Equal(t, 5, *deep.Peek())

	h2.Clear()
	assert.True(t, h2.IsEmpty())
}

func TestPeekPopEmptyPairing(t *testing.T) {
	h := NewPairingHeap([]int{}, lt)
	assert.Nil(t, h.Peek())
	assert.Nil(t, h.Pop())
}

func TestLengthIsEmptyPairing(t *testing.T) {
	h := NewPairingHeap([]int{}, lt)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())

	h.Insert(10)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, 1, h.Length())
}

func TestMergeWithPairing(t *testing.T) {
	h1 := NewPairingHeap([]int{1, 4, 7}, lt)
	h2 := NewPairingHeap([]int{2, 3, 5, 6}, lt)

	h1.MergeWith(h2)
	result := collectPairing(&h1)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, result)
}
