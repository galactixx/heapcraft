package heapcraft

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// lt returns true if a is less than b
func lt(a, b int) bool { return a < b }

// gt returns true if a is greater than b
func gt(a, b int) bool { return a > b }

func TestNewDaryHeap(t *testing.T) {
	tests := []struct {
		rawData  []*HeapPair[string, int]
		heapData []*HeapPair[string, int]
		d        int
		cmp      func(a, b int) bool
	}{
		{
			rawData: []*HeapPair[string, int]{
				CreateHeapPair("1", 1),
				CreateHeapPair("2", 2),
				CreateHeapPair("3", 3),
				CreateHeapPair("4", 4),
				CreateHeapPair("5", 5),
			},
			d:   3,
			cmp: lt,
			heapData: []*HeapPair[string, int]{
				CreateHeapPair("1", 1),
				CreateHeapPair("2", 2),
				CreateHeapPair("3", 3),
				CreateHeapPair("4", 4),
				CreateHeapPair("5", 5),
			},
		},
		{
			rawData: []*HeapPair[string, int]{
				CreateHeapPair("10", 10),
				CreateHeapPair("-1", -1),
				CreateHeapPair("0", 0),
				CreateHeapPair("42", 42),
				CreateHeapPair("7", 7),
				CreateHeapPair("-5", -5),
				CreateHeapPair("3", 3),
			},
			d:   4,
			cmp: lt,
			heapData: []*HeapPair[string, int]{
				CreateHeapPair("-5", -5),
				CreateHeapPair("-1", -1),
				CreateHeapPair("0", 0),
				CreateHeapPair("42", 42),
				CreateHeapPair("7", 7),
				CreateHeapPair("10", 10),
				CreateHeapPair("3", 3),
			},
		},
		{
			rawData: []*HeapPair[string, int]{
				CreateHeapPair("5", 5),
				CreateHeapPair("4", 4),
				CreateHeapPair("3", 3),
				CreateHeapPair("2", 2),
				CreateHeapPair("1", 1),
			},
			d:   2,
			cmp: lt,
			heapData: []*HeapPair[string, int]{
				CreateHeapPair("1", 1),
				CreateHeapPair("2", 2),
				CreateHeapPair("3", 3),
				CreateHeapPair("5", 5),
				CreateHeapPair("4", 4),
			},
		},
	}

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("New Dary Heap Test %d", idx+1), func(t *testing.T) {
			h := NewDaryHeap(tt.d, tt.rawData, tt.cmp)
			for i := range h.data {
				assert.Equal(t, tt.heapData[i].Value(), h.data[i].Value())
				assert.Equal(t, tt.heapData[i].Priority(), h.data[i].Priority())
			}
		})
	}
}

func TestPushPopPeekLenIsEmptyDary(t *testing.T) {
	h := DaryHeap[string, int]{data: make([]*HeapPair[string, int], 0), cmp: lt, d: 3}

	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	assert.Nil(t, h.Peek())

	input := []*HeapPair[string, int]{
		CreateHeapPair("5", 5),
		CreateHeapPair("3", 3),
		CreateHeapPair("8", 8),
		CreateHeapPair("1", 1),
		CreateHeapPair("4", 4),
	}
	expectedOrder := []*HeapPair[string, int]{
		CreateHeapPair("1", 1),
		CreateHeapPair("3", 3),
		CreateHeapPair("4", 4),
		CreateHeapPair("5", 5),
		CreateHeapPair("8", 8),
	}

	for _, v := range input {
		h.Push(v.Value(), v.Priority())
	}

	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(input), h.Length())
	assert.Equal(t, 1, h.Peek().Priority())

	for i, expected := range expectedOrder {
		popped := h.Pop()
		assert.NotNil(t, popped)
		assert.Equal(t, expected.Value(), popped.Value())
		assert.Equal(t, expected.Priority(), popped.Priority())
		assert.Equal(t, len(input)-(i+1), h.Length())
	}

	assert.True(t, h.IsEmpty())
	assert.Nil(t, h.Peek())
}

func TestClearDary(t *testing.T) {
	h := NewDaryHeap(3, []*HeapPair[string, int]{
		CreateHeapPair("7", 7),
		CreateHeapPair("2", 2),
		CreateHeapPair("9", 9),
		CreateHeapPair("1", 1),
	}, lt)
	assert.Equal(t, 4, h.Length())

	h.Clear()
	assert.True(t, h.IsEmpty())
}

func TestUpdateRemoveDary(t *testing.T) {
	h := NewDaryHeap(3, []*HeapPair[string, int]{
		CreateHeapPair("4", 4),
		CreateHeapPair("10", 10),
		CreateHeapPair("3", 3),
		CreateHeapPair("5", 5),
		CreateHeapPair("1", 1),
	}, lt)

	var idx4 int
	for i, v := range h.data {
		if v.Priority() == 4 {
			idx4 = i
			break
		}
	}
	_, err := h.Update(idx4, "0", 0)
	assert.NoError(t, err)
	assert.Equal(t, 0, h.Peek().Priority())

	removedRoot, err := h.Remove(0)
	assert.NoError(t, err)
	assert.Equal(t, 0, removedRoot.Priority())
	assert.Equal(t, 1, h.Peek().Priority())

	var idx5 int
	for i, v := range h.data {
		if v.Priority() == 5 {
			idx5 = i
			break
		}
	}
	removed5, err := h.Remove(idx5)
	assert.NoError(t, err)
	assert.Equal(t, 5, removed5.Priority())

	result := []int{}
	for !h.IsEmpty() {
		val := h.Pop()
		result = append(result, val.Priority())
	}
	assert.Equal(t, []int{1, 3, 10}, result)
}

func TestPopPushPushPopDary(t *testing.T) {
	h := NewDaryHeap(3, []*HeapPair[string, int]{
		CreateHeapPair("2", 2),
		CreateHeapPair("6", 6),
		CreateHeapPair("4", 4),
	}, lt)

	popped := h.PopPush("1", 1)
	assert.Equal(t, 2, popped.Priority())

	returned := h.PushPop("5", 5)
	assert.Equal(t, 1, returned.Priority())

	expected := []*HeapPair[string, int]{
		CreateHeapPair("4", 4),
		CreateHeapPair("6", 6),
		CreateHeapPair("5", 5),
	}
	for i := range h.data {
		assert.Equal(t, expected[i].Value(), h.data[i].Value())
		assert.Equal(t, expected[i].Priority(), h.data[i].Priority())
	}
}

func TestNLargestNSmallestDary(t *testing.T) {
	data := []*HeapPair[string, int]{
		CreateHeapPair("7", 7),
		CreateHeapPair("2", 2),
		CreateHeapPair("9", 9),
		CreateHeapPair("1", 1),
		CreateHeapPair("5", 5),
		CreateHeapPair("3", 3),
	}

	hMax := NLargestDary(3, 3, data, lt)
	assert.Equal(t, 3, hMax.Length())

	res := []int{}
	for !hMax.IsEmpty() {
		res = append(res, hMax.Pop().Priority())
	}
	assert.Equal(t, []int{5, 7, 9}, res)

	hMin := NSmallestDary(3, 3, data, gt)
	assert.Equal(t, 3, hMin.Length())
	res2 := []int{}
	for !hMin.IsEmpty() {
		res2 = append(res2, hMin.Pop().Priority())
	}
	assert.Equal(t, []int{3, 2, 1}, res2)
}

func TestRegisterDeregisterCallbacksDary(t *testing.T) {
	h := NewDaryHeap(3, []*HeapPair[string, int]{
		CreateHeapPair("3", 3),
		CreateHeapPair("1", 1),
		CreateHeapPair("4", 4),
		CreateHeapPair("2", 2),
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
	h := DaryHeap[string, int]{data: []*HeapPair[string, int]{}, cmp: lt, d: 2}
	assert.Nil(t, h.Peek())
	assert.Nil(t, h.Pop())
}

func TestRemoveOutOfBoundsDary(t *testing.T) {
	h := NewDaryHeap(3, []*HeapPair[string, int]{
		CreateHeapPair("1", 1),
		CreateHeapPair("2", 2),
		CreateHeapPair("3", 3),
	}, lt)
	_, err := h.Remove(5)
	assert.Error(t, err)
}

func TestUpdateOutOfBoundsDary(t *testing.T) {
	h := NewDaryHeap(3, []*HeapPair[string, int]{
		CreateHeapPair("1", 1),
		CreateHeapPair("2", 2),
		CreateHeapPair("3", 3),
	}, lt)
	_, err := h.Update(5, "10", 10)
	assert.Error(t, err)
}

func TestNLargestNSmallestBinary(t *testing.T) {
	data := []*HeapPair[string, int]{
		CreateHeapPair("7", 7),
		CreateHeapPair("2", 2),
		CreateHeapPair("9", 9),
		CreateHeapPair("1", 1),
		CreateHeapPair("5", 5),
		CreateHeapPair("3", 3),
	}

	// Test NLargestBinary - should get the 3 largest numbers
	hMax := NLargestBinary(3, data, lt)
	assert.Equal(t, 3, hMax.Length())

	res := []int{}
	for !hMax.IsEmpty() {
		res = append(res, hMax.Pop().Priority())
	}
	assert.Equal(t, []int{5, 7, 9}, res)

	// Test NSmallestBinary - should get the 3 smallest numbers
	hMin := NSmallestBinary(3, data, gt)
	assert.Equal(t, 3, hMin.Length())
	res2 := []int{}
	for !hMin.IsEmpty() {
		res2 = append(res2, hMin.Pop().Priority())
	}
	assert.Equal(t, []int{3, 2, 1}, res2)
}
