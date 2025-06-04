package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func lt(a, b int) bool { return a < b }
func gt(a, b int) bool { return a > b }

func TestHeapify(t *testing.T) {
	tests := []struct {
		rawData  []*HeapPair[string, int]
		heapData []*HeapPair[string, int]
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

	for _, tt := range tests {
		heap := Heapify(tt.rawData, tt.cmp)
		for i := range heap.data {
			assert.Equal(t, tt.heapData[i].Value(), heap.data[i].Value())
			assert.Equal(t, tt.heapData[i].Priority(), heap.data[i].Priority())
		}
	}
}

func TestPushPopPeekLenIsEmpty(t *testing.T) {
	h := BinaryHeap[string, int]{data: make([]*HeapPair[string, int], 0), cmp: lt}

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

	// heap should be empty again
	assert.True(t, h.IsEmpty())
	assert.Nil(t, h.Peek())
}

func TestClearCloneDeepClone(t *testing.T) {
	h := Heapify([]*HeapPair[string, int]{
		CreateHeapPair("7", 7),
		CreateHeapPair("2", 2),
		CreateHeapPair("9", 9),
		CreateHeapPair("1", 1),
	}, lt)
	assert.Equal(t, 4, h.Length())

	clone := h.Clone()
	for i := range h.data {
		assert.Equal(t, h.data[i].Value(), clone.data[i].Value())
		assert.Equal(t, h.data[i].Priority(), clone.data[i].Priority())
	}

	// modify original to ensure clone is shallow-copy
	h.Push("0", 0)
	assert.NotEqual(t, len(h.data), len(clone.data))

	h2 := Heapify([]*HeapPair[string, int]{
		CreateHeapPair("10", 10),
		CreateHeapPair("5", 5),
		CreateHeapPair("3", 3),
	}, lt)

	h2.Clear()
	assert.True(t, h2.IsEmpty())
}

func TestUpdateRemove(t *testing.T) {
	h := Heapify([]*HeapPair[string, int]{
		CreateHeapPair("4", 4),
		CreateHeapPair("10", 10),
		CreateHeapPair("3", 3),
		CreateHeapPair("5", 5),
		CreateHeapPair("1", 1),
	}, lt)

	initial := make([]*HeapPair[string, int], len(h.data))
	copy(initial, h.data)

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

	removed, err := h.Remove(0)
	assert.NoError(t, err)
	assert.Equal(t, 0, removed.Priority())
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

func TestPopPushPushPop(t *testing.T) {
	h := Heapify([]*HeapPair[string, int]{
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

func TestNLargestNSmallest(t *testing.T) {
	data := []*HeapPair[string, int]{
		CreateHeapPair("7", 7),
		CreateHeapPair("2", 2),
		CreateHeapPair("9", 9),
		CreateHeapPair("1", 1),
		CreateHeapPair("5", 5),
		CreateHeapPair("3", 3),
	}

	hMax := NLargest(3, data, lt)
	assert.Equal(t, 3, hMax.Length())

	res := []int{}
	for !hMax.IsEmpty() {
		res = append(res, hMax.Pop().Priority())
	}
	assert.Equal(t, []int{5, 7, 9}, res)

	hMin := NSmallest(3, data, gt)
	assert.Equal(t, 3, hMin.Length())
	res2 := []int{}
	for !hMin.IsEmpty() {
		res2 = append(res2, hMin.Pop().Priority())
	}
	assert.Equal(t, []int{3, 2, 1}, res2)
}

func TestRegisterDeregisterCallbacks(t *testing.T) {
	h := Heapify([]*HeapPair[string, int]{
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

func TestPeekPopEmpty(t *testing.T) {
	h := BinaryHeap[string, int]{data: []*HeapPair[string, int]{}, cmp: lt}
	assert.Nil(t, h.Peek())
	assert.Nil(t, h.Pop())
}

func TestRemoveOutOfBounds(t *testing.T) {
	h := Heapify([]*HeapPair[string, int]{
		CreateHeapPair("1", 1),
		CreateHeapPair("2", 2),
		CreateHeapPair("3", 3),
	}, lt)
	_, err := h.Remove(5)
	assert.Error(t, err)
}

func TestUpdateOutOfBounds(t *testing.T) {
	h := Heapify([]*HeapPair[string, int]{
		CreateHeapPair("1", 1),
		CreateHeapPair("2", 2),
		CreateHeapPair("3", 3),
	}, lt)
	_, err := h.Update(5, "10", 10)
	assert.Error(t, err)
}
