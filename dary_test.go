package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDaryHeap(t *testing.T) {
	tests := []struct {
		rawData  []int
		heapData []int
		d        int
		cmp      func(a, b int) bool
	}{
		{
			rawData:  []int{1, 2, 3, 4, 5},
			d:        3,
			cmp:      lt,
			heapData: []int{1, 2, 3, 4, 5},
		},
		{
			rawData:  []int{10, -1, 0, 42, 7, -5, 3},
			d:        4,
			cmp:      lt,
			heapData: []int{-5, -1, 0, 42, 7, 10, 3},
		},
		{
			rawData:  []int{5, 4, 3, 2, 1},
			d:        2,
			cmp:      lt,
			heapData: []int{1, 2, 3, 5, 4},
		},
	}

	for _, tt := range tests {
		h := NewDaryHeap(tt.d, tt.rawData, tt.cmp)
		assert.Equal(t, tt.heapData, h.data)
	}
}

func TestPushPopPeekLenIsEmptyDary(t *testing.T) {
	h := DaryHeap[int]{data: make([]int, 0), cmp: lt, d: 3}

	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	assert.Nil(t, h.Peek())

	input := []int{5, 3, 8, 1, 4}
	expectedOrder := []int{1, 3, 4, 5, 8}

	for _, v := range input {
		h.Push(v)
	}

	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(input), h.Length())
	assert.Equal(t, 1, *h.Peek())

	for i, expected := range expectedOrder {
		popped := h.Pop()
		assert.NotNil(t, popped)
		assert.Equal(t, expected, *popped)
		assert.Equal(t, len(input)-(i+1), h.Length())
	}

	assert.True(t, h.IsEmpty())
	assert.Nil(t, h.Peek())
}

func TestClearDary(t *testing.T) {
	h := NewDaryHeap(3, []int{7, 2, 9, 1}, lt)
	assert.Equal(t, 4, h.Length())

	h.Clear()
	assert.True(t, h.IsEmpty())
}

func TestUpdateRemoveDary(t *testing.T) {
	h := NewDaryHeap(3, []int{4, 10, 3, 5, 1}, lt)
	var idx4 int
	for i, v := range h.data {
		if v == 4 {
			idx4 = i
			break
		}
	}
	err := h.Update(idx4, 0)
	assert.NoError(t, err)
	assert.Equal(t, 0, *h.Peek())

	removedRoot, err := h.Remove(0)
	assert.NoError(t, err)
	assert.Equal(t, 0, *removedRoot)
	assert.Equal(t, 1, *h.Peek())

	// remove a middle element, e.g., value 5
	var idx5 int
	for i, v := range h.data {
		if v == 5 {
			idx5 = i
			break
		}
	}
	removed5, err := h.Remove(idx5)
	assert.NoError(t, err)
	assert.Equal(t, 5, *removed5)

	result := []int{}
	for !h.IsEmpty() {
		val := h.Pop()
		result = append(result, *val)
	}
	assert.Equal(t, []int{1, 3, 10}, result)
}

func TestPopPushPushPopDary(t *testing.T) {
	h := NewDaryHeap(3, []int{2, 6, 4}, lt)

	popped := h.PopPush(1)
	assert.Equal(t, 2, popped)

	returned := h.PushPop(5)
	assert.Equal(t, 1, returned)
	assert.Equal(t, []int{4, 6, 5}, h.data)
}

func TestNLargestNSmallestDary(t *testing.T) {
	data := []int{7, 2, 9, 1, 5, 3}

	hMax := NLargestDary(3, 3, data, lt)
	assert.Equal(t, 3, hMax.Length())
	res := []int{}
	for !hMax.IsEmpty() {
		res = append(res, *hMax.Pop())
	}
	assert.Equal(t, []int{5, 7, 9}, res)

	hMin := NSmallestDary(3, 3, data, gt)
	assert.Equal(t, 3, hMin.Length())
	res2 := []int{}
	for !hMin.IsEmpty() {
		res2 = append(res2, *hMin.Pop())
	}
	assert.Equal(t, []int{3, 2, 1}, res2)
}

func TestRegisterDeregisterCallbacksDary(t *testing.T) {
	h := NewDaryHeap(3, []int{3, 1, 4, 2}, lt)
	events := [][2]int{}
	cb := h.Register(func(x, y int) {
		events = append(events, [2]int{x, y})
	})

	h.Push(0)
	assert.NotEmpty(t, events)

	events = nil
	err := h.Deregister(cb.ID)
	assert.NoError(t, err)
	h.Push(-1)
	assert.Empty(t, events)

	err = h.Deregister(999)
	assert.Error(t, err)
}

func TestPeekPopEmptyDary(t *testing.T) {
	h := DaryHeap[int]{data: []int{}, cmp: lt, d: 2}
	assert.Nil(t, h.Peek())
	assert.Nil(t, h.Pop())
}

func TestRemoveOutOfBoundsDary(t *testing.T) {
	h := NewDaryHeap(3, []int{1, 2, 3}, lt)
	_, err := h.Remove(5)
	assert.Error(t, err)
}

func TestUpdateOutOfBoundsDary(t *testing.T) {
	h := NewDaryHeap(3, []int{1, 2, 3}, lt)
	err := h.Update(5, 10)
	assert.Error(t, err)
}
