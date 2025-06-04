package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDaryHeap(t *testing.T) {
	tests := []struct {
		rawData  []HeapPair[string, int]
		heapData []HeapPair[string, int]
		d        int
		cmp      func(a, b int) bool
	}{
		{
			rawData: []HeapPair[string, int]{
				{value: "1", priority: 1},
				{value: "2", priority: 2},
				{value: "3", priority: 3},
				{value: "4", priority: 4},
				{value: "5", priority: 5},
			},
			d:   3,
			cmp: lt,
			heapData: []HeapPair[string, int]{
				{value: "1", priority: 1},
				{value: "2", priority: 2},
				{value: "3", priority: 3},
				{value: "4", priority: 4},
				{value: "5", priority: 5},
			},
		},
		{
			rawData: []HeapPair[string, int]{
				{value: "10", priority: 10},
				{value: "-1", priority: -1},
				{value: "0", priority: 0},
				{value: "42", priority: 42},
				{value: "7", priority: 7},
				{value: "-5", priority: -5},
				{value: "3", priority: 3},
			},
			d:   4,
			cmp: lt,
			heapData: []HeapPair[string, int]{
				{value: "-5", priority: -5},
				{value: "-1", priority: -1},
				{value: "0", priority: 0},
				{value: "42", priority: 42},
				{value: "7", priority: 7},
				{value: "10", priority: 10},
				{value: "3", priority: 3},
			},
		},
		{
			rawData: []HeapPair[string, int]{
				{value: "5", priority: 5},
				{value: "4", priority: 4},
				{value: "3", priority: 3},
				{value: "2", priority: 2},
				{value: "1", priority: 1},
			},
			d:   2,
			cmp: lt,
			heapData: []HeapPair[string, int]{
				{value: "1", priority: 1},
				{value: "2", priority: 2},
				{value: "3", priority: 3},
				{value: "5", priority: 5},
				{value: "4", priority: 4},
			},
		},
	}

	for _, tt := range tests {
		h := NewDaryHeap(tt.d, tt.rawData, tt.cmp)
		assert.Equal(t, tt.heapData, h.data)
	}
}

func TestPushPopPeekLenIsEmptyDary(t *testing.T) {
	h := DaryHeap[string, int]{data: make([]HeapPair[string, int], 0), cmp: lt, d: 3}

	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	assert.Nil(t, h.Peek())

	input := []HeapPair[string, int]{
		{value: "5", priority: 5},
		{value: "3", priority: 3},
		{value: "8", priority: 8},
		{value: "1", priority: 1},
		{value: "4", priority: 4},
	}
	expectedOrder := []HeapPair[string, int]{
		{value: "1", priority: 1},
		{value: "3", priority: 3},
		{value: "4", priority: 4},
		{value: "5", priority: 5},
		{value: "8", priority: 8},
	}

	for _, v := range input {
		h.Push(v.value, v.priority)
	}

	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(input), h.Length())
	assert.Equal(t, 1, h.Peek().priority)

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
	h := NewDaryHeap(3, []HeapPair[string, int]{
		{value: "7", priority: 7},
		{value: "2", priority: 2},
		{value: "9", priority: 9},
		{value: "1", priority: 1},
	}, lt)
	assert.Equal(t, 4, h.Length())

	h.Clear()
	assert.True(t, h.IsEmpty())
}

func TestUpdateRemoveDary(t *testing.T) {
	h := NewDaryHeap(3, []HeapPair[string, int]{
		{value: "4", priority: 4},
		{value: "10", priority: 10},
		{value: "3", priority: 3},
		{value: "5", priority: 5},
		{value: "1", priority: 1},
	}, lt)

	var idx4 int
	for i, v := range h.data {
		if v.priority == 4 {
			idx4 = i
			break
		}
	}
	_, err := h.Update(idx4, "0", 0)
	assert.NoError(t, err)
	assert.Equal(t, 0, h.Peek().priority)

	removedRoot, err := h.Remove(0)
	assert.NoError(t, err)
	assert.Equal(t, 0, removedRoot.priority)
	assert.Equal(t, 1, h.Peek().priority)

	var idx5 int
	for i, v := range h.data {
		if v.priority == 5 {
			idx5 = i
			break
		}
	}
	removed5, err := h.Remove(idx5)
	assert.NoError(t, err)
	assert.Equal(t, 5, removed5.priority)

	result := []int{}
	for !h.IsEmpty() {
		val := h.Pop()
		result = append(result, val.priority)
	}
	assert.Equal(t, []int{1, 3, 10}, result)
}

func TestPopPushPushPopDary(t *testing.T) {
	h := NewDaryHeap(3, []HeapPair[string, int]{
		{value: "2", priority: 2},
		{value: "6", priority: 6},
		{value: "4", priority: 4},
	}, lt)

	popped := h.PopPush("1", 1)
	assert.Equal(t, 2, popped.priority)

	returned := h.PushPop("5", 5)
	assert.Equal(t, 1, returned.priority)
	assert.Equal(t, []HeapPair[string, int]{
		{value: "4", priority: 4},
		{value: "6", priority: 6},
		{value: "5", priority: 5},
	}, h.data)
}

func TestNLargestNSmallestDary(t *testing.T) {
	data := []HeapPair[string, int]{
		{value: "7", priority: 7},
		{value: "2", priority: 2},
		{value: "9", priority: 9},
		{value: "1", priority: 1},
		{value: "5", priority: 5},
		{value: "3", priority: 3},
	}

	hMax := NLargestDary(3, 3, data, lt)
	assert.Equal(t, 3, hMax.Length())

	res := []int{}
	for !hMax.IsEmpty() {
		res = append(res, hMax.Pop().priority)
	}
	assert.Equal(t, []int{5, 7, 9}, res)

	hMin := NSmallestDary(3, 3, data, gt)
	assert.Equal(t, 3, hMin.Length())
	res2 := []int{}
	for !hMin.IsEmpty() {
		res2 = append(res2, hMin.Pop().priority)
	}
	assert.Equal(t, []int{3, 2, 1}, res2)
}

func TestRegisterDeregisterCallbacksDary(t *testing.T) {
	h := NewDaryHeap(3, []HeapPair[string, int]{
		{value: "3", priority: 3},
		{value: "1", priority: 1},
		{value: "4", priority: 4},
		{value: "2", priority: 2},
	}, lt)
	events := [][2]int{}
	cb := h.Register(func(x, y int) {
		events = append(events, [2]int{x, y})
	})

	h.Push("0", 0)
	assert.NotEmpty(t, events)

	events = nil
	err := h.Deregister(cb.ID)
	assert.NoError(t, err)
	h.Push("-1", -1)

	assert.Empty(t, events)

	err = h.Deregister(999)
	assert.Error(t, err)
}

func TestPeekPopEmptyDary(t *testing.T) {
	h := DaryHeap[string, int]{data: []HeapPair[string, int]{}, cmp: lt, d: 2}
	assert.Nil(t, h.Peek())
	assert.Nil(t, h.Pop())
}

func TestRemoveOutOfBoundsDary(t *testing.T) {
	h := NewDaryHeap(3, []HeapPair[string, int]{
		{value: "1", priority: 1},
		{value: "2", priority: 2},
		{value: "3", priority: 3},
	}, lt)
	_, err := h.Remove(5)
	assert.Error(t, err)
}

func TestUpdateOutOfBoundsDary(t *testing.T) {
	h := NewDaryHeap(3, []HeapPair[string, int]{
		{value: "1", priority: 1},
		{value: "2", priority: 2},
		{value: "3", priority: 3},
	}, lt)
	_, err := h.Update(5, "10", 10)
	assert.Error(t, err)
}
