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
		rawData  []HeapNode[string, int]
		heapData []HeapNode[string, int]
		d        int
		cmp      func(a, b int) bool
	}{
		{
			rawData: []HeapNode[string, int]{
				CreateHeapNode("1", 1),
				CreateHeapNode("2", 2),
				CreateHeapNode("3", 3),
				CreateHeapNode("4", 4),
				CreateHeapNode("5", 5),
			},
			d:   3,
			cmp: lt,
			heapData: []HeapNode[string, int]{
				CreateHeapNode("1", 1),
				CreateHeapNode("2", 2),
				CreateHeapNode("3", 3),
				CreateHeapNode("4", 4),
				CreateHeapNode("5", 5),
			},
		},
		{
			rawData: []HeapNode[string, int]{
				CreateHeapNode("10", 10),
				CreateHeapNode("-1", -1),
				CreateHeapNode("0", 0),
				CreateHeapNode("42", 42),
				CreateHeapNode("7", 7),
				CreateHeapNode("-5", -5),
				CreateHeapNode("3", 3),
			},
			d:   4,
			cmp: lt,
			heapData: []HeapNode[string, int]{
				CreateHeapNode("-5", -5),
				CreateHeapNode("-1", -1),
				CreateHeapNode("0", 0),
				CreateHeapNode("42", 42),
				CreateHeapNode("7", 7),
				CreateHeapNode("10", 10),
				CreateHeapNode("3", 3),
			},
		},
		{
			rawData: []HeapNode[string, int]{
				CreateHeapNode("5", 5),
				CreateHeapNode("4", 4),
				CreateHeapNode("3", 3),
				CreateHeapNode("2", 2),
				CreateHeapNode("1", 1),
			},
			d:   2,
			cmp: lt,
			heapData: []HeapNode[string, int]{
				CreateHeapNode("1", 1),
				CreateHeapNode("2", 2),
				CreateHeapNode("3", 3),
				CreateHeapNode("5", 5),
				CreateHeapNode("4", 4),
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
	h := NewDaryHeap(3, []HeapNode[string, int]{}, lt)

	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	_, err := h.Peek()
	assert.Error(t, err)

	input := []HeapNode[string, int]{
		CreateHeapNode("5", 5),
		CreateHeapNode("3", 3),
		CreateHeapNode("8", 8),
		CreateHeapNode("1", 1),
		CreateHeapNode("4", 4),
	}
	expectedOrder := []HeapNode[string, int]{
		CreateHeapNode("1", 1),
		CreateHeapNode("3", 3),
		CreateHeapNode("4", 4),
		CreateHeapNode("5", 5),
		CreateHeapNode("8", 8),
	}

	for _, v := range input {
		h.Push(v.Value(), v.Priority())
	}

	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(input), h.Length())
	peeked, err := h.Peek()
	assert.NoError(t, err)
	assert.Equal(t, 1, peeked.Priority())

	for i, expected := range expectedOrder {
		popped, err := h.Pop()
		assert.NoError(t, err)
		assert.Equal(t, expected.Value(), popped.Value())
		assert.Equal(t, expected.Priority(), popped.Priority())
		assert.Equal(t, len(input)-(i+1), h.Length())
	}

	assert.True(t, h.IsEmpty())
	_, err = h.Peek()
	assert.Error(t, err)
}

func TestClearDary(t *testing.T) {
	h := NewDaryHeap(3, []HeapNode[string, int]{
		CreateHeapNode("7", 7),
		CreateHeapNode("2", 2),
		CreateHeapNode("9", 9),
		CreateHeapNode("1", 1),
	}, lt)
	assert.Equal(t, 4, h.Length())

	h.Clear()
	assert.True(t, h.IsEmpty())
}

func TestUpdateRemoveDary(t *testing.T) {
	h := NewDaryHeap(3, []HeapNode[string, int]{
		CreateHeapNode("4", 4),
		CreateHeapNode("10", 10),
		CreateHeapNode("3", 3),
		CreateHeapNode("5", 5),
		CreateHeapNode("1", 1),
	}, lt)

	var idx4 int
	for i, v := range h.data {
		if v.Priority() == 4 {
			idx4 = i
			break
		}
	}
	err := h.Update(idx4, "0", 0)
	assert.NoError(t, err)
	peeked, err := h.Peek()
	assert.NoError(t, err)
	assert.Equal(t, 0, peeked.Priority())

	removedRoot, err := h.Remove(0)
	assert.NoError(t, err)
	assert.Equal(t, 0, removedRoot.Priority())
	peeked, err = h.Peek()
	assert.NoError(t, err)
	assert.Equal(t, 1, peeked.Priority())

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
		val, err := h.Pop()
		assert.NoError(t, err)
		result = append(result, val.Priority())
	}
	assert.Equal(t, []int{1, 3, 10}, result)
}

func TestPopPushPushPopDary(t *testing.T) {
	h := NewDaryHeap(3, []HeapNode[string, int]{
		CreateHeapNode("2", 2),
		CreateHeapNode("6", 6),
		CreateHeapNode("4", 4),
	}, lt)

	popped := h.PopPush("1", 1)
	assert.Equal(t, 2, popped.Priority())

	returned := h.PushPop("5", 5)
	assert.Equal(t, 1, returned.Priority())

	expected := []HeapNode[string, int]{
		CreateHeapNode("4", 4),
		CreateHeapNode("6", 6),
		CreateHeapNode("5", 5),
	}
	for i := range h.data {
		assert.Equal(t, expected[i].Value(), h.data[i].Value())
		assert.Equal(t, expected[i].Priority(), h.data[i].Priority())
	}
}

func TestNLargestNSmallestDary(t *testing.T) {
	data := []HeapNode[string, int]{
		CreateHeapNode("7", 7),
		CreateHeapNode("2", 2),
		CreateHeapNode("9", 9),
		CreateHeapNode("1", 1),
		CreateHeapNode("5", 5),
		CreateHeapNode("3", 3),
	}

	hMax := NLargestDary(3, 3, data, lt)
	assert.Equal(t, 3, hMax.Length())

	res := []int{}
	for !hMax.IsEmpty() {
		popped, err := hMax.Pop()
		assert.NoError(t, err)
		res = append(res, popped.Priority())
	}
	assert.Equal(t, []int{5, 7, 9}, res)

	hMin := NSmallestDary(3, 3, data, gt)
	assert.Equal(t, 3, hMin.Length())
	res2 := []int{}
	for !hMin.IsEmpty() {
		popped, err := hMin.Pop()
		assert.NoError(t, err)
		res2 = append(res2, popped.Priority())
	}
	assert.Equal(t, []int{3, 2, 1}, res2)
}

func TestRegisterDeregisterCallbacksDary(t *testing.T) {
	h := NewDaryHeap(3, []HeapNode[string, int]{
		CreateHeapNode("3", 3),
		CreateHeapNode("1", 1),
		CreateHeapNode("4", 4),
		CreateHeapNode("2", 2),
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
	h := DaryHeap[string, int]{data: []HeapNode[string, int]{}, cmp: lt, d: 2}
	_, err := h.Peek()
	assert.Error(t, err)
	_, err = h.Pop()
	assert.Error(t, err)
}

func TestRemoveOutOfBoundsDary(t *testing.T) {
	h := NewDaryHeap(3, []HeapNode[string, int]{
		CreateHeapNode("1", 1),
		CreateHeapNode("2", 2),
		CreateHeapNode("3", 3),
	}, lt)
	_, err := h.Remove(5)
	assert.Error(t, err)
}

func TestUpdateOutOfBoundsDary(t *testing.T) {
	h := NewDaryHeap(3, []HeapNode[string, int]{
		CreateHeapNode("1", 1),
		CreateHeapNode("2", 2),
		CreateHeapNode("3", 3),
	}, lt)
	err := h.Update(5, "10", 10)
	assert.Error(t, err)
}

func TestNLargestNSmallestBinary(t *testing.T) {
	data := []HeapNode[string, int]{
		CreateHeapNode("7", 7),
		CreateHeapNode("2", 2),
		CreateHeapNode("9", 9),
		CreateHeapNode("1", 1),
		CreateHeapNode("5", 5),
		CreateHeapNode("3", 3),
	}

	// Test NLargestBinary - should get the 3 largest numbers
	hMax := NLargestBinary(3, data, lt)
	assert.Equal(t, 3, hMax.Length())

	res := []int{}
	for !hMax.IsEmpty() {
		popped, err := hMax.Pop()
		assert.NoError(t, err)
		res = append(res, popped.Priority())
	}
	assert.Equal(t, []int{5, 7, 9}, res)

	// Test NSmallestBinary - should get the 3 smallest numbers
	hMin := NSmallestBinary(3, data, gt)
	assert.Equal(t, 3, hMin.Length())
	res2 := []int{}
	for !hMin.IsEmpty() {
		popped, err := hMin.Pop()
		assert.NoError(t, err)
		res2 = append(res2, popped.Priority())
	}
	assert.Equal(t, []int{3, 2, 1}, res2)
}

func BenchmarkBinaryHeapInsertion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewBinaryHeap(data, func(a, b int) bool { return a < b })
	b.ReportAllocs()

	insertions := generateRandomNumbers(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Push(insertions[i], insertions[i])
	}
}

func BenchmarkBinaryHeapDeletion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewBinaryHeap(data, func(a, b int) bool { return a < b })

	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Pop()
	}
}

// D-ary Heap Benchmarks
func BenchmarkDaryHeap3Insertion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewDaryHeap(3, data, func(a, b int) bool { return a < b })
	b.ReportAllocs()

	insertions := generateRandomNumbers(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Push(insertions[i], insertions[i])
	}
}

func BenchmarkDaryHeap3Deletion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewDaryHeap(3, data, func(a, b int) bool { return a < b })

	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Pop()
	}
}

func BenchmarkDaryHeap4Insertion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewDaryHeap(4, data, func(a, b int) bool { return a < b })
	b.ReportAllocs()

	insertions := generateRandomNumbers(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Push(insertions[i], insertions[i])
	}
}

func BenchmarkDaryHeap4Deletion(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewDaryHeap(4, data, func(a, b int) bool { return a < b })

	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Pop()
	}
}

func BenchmarkBinaryPushPop(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewDaryHeap(3, data, func(a, b int) bool { return a < b })

	insertions := generateRandomNumbers(b)
	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.PushPop(insertions[i], insertions[i])
	}
}

func BenchmarkBinaryPopPush(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewDaryHeap(3, data, func(a, b int) bool { return a < b })

	insertions := generateRandomNumbers(b)
	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.PopPush(insertions[i], insertions[i])
	}
}

func BenchmarkDaryHeap3PushPop(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewDaryHeap(3, data, func(a, b int) bool { return a < b })

	insertions := generateRandomNumbers(b)
	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.PushPop(insertions[i], insertions[i])
	}
}

func BenchmarkDaryHeap3PopPush(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewDaryHeap(3, data, func(a, b int) bool { return a < b })

	insertions := generateRandomNumbers(b)
	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.PopPush(insertions[i], insertions[i])
	}
}

func BenchmarkDaryHeap4PushPop(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewDaryHeap(4, data, func(a, b int) bool { return a < b })

	insertions := generateRandomNumbers(b)
	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.PushPop(insertions[i], insertions[i])
	}
}

func BenchmarkDaryHeap4PopPush(b *testing.B) {
	data := make([]HeapNode[int, int], 0)
	heap := NewDaryHeap(4, data, func(a, b int) bool { return a < b })

	insertions := generateRandomNumbers(b)
	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.PopPush(insertions[i], insertions[i])
	}
}
