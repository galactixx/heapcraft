package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRadixHeapPopOrder(t *testing.T) {
	raw := []HeapNode[string, uint]{
		CreateHeapNode("value10", uint(10)),
		CreateHeapNode("value3", uint(3)),
		CreateHeapNode("value7", uint(7)),
		CreateHeapNode("value1", uint(1)),
		CreateHeapNode("value5", uint(5)),
		CreateHeapNode("value2", uint(2)),
	}
	rh := NewRadixHeap(raw, false)
	assert.False(t, rh.IsEmpty())
	assert.Equal(t, len(raw), rh.Length())

	expected := []HeapNode[string, uint]{
		CreateHeapNode("value1", uint(1)),
		CreateHeapNode("value2", uint(2)),
		CreateHeapNode("value3", uint(3)),
		CreateHeapNode("value5", uint(5)),
		CreateHeapNode("value7", uint(7)),
		CreateHeapNode("value10", uint(10)),
	}
	actualValues := []string{}
	actualPriorities := []uint{}
	for !rh.IsEmpty() {
		v, p, err := rh.Pop()
		assert.NoError(t, err)
		actualValues = append(actualValues, v)
		actualPriorities = append(actualPriorities, p)
	}
	for i := range expected {
		assert.Equal(t, expected[i].value, actualValues[i])
		assert.Equal(t, expected[i].priority, actualPriorities[i])
	}
	assert.True(t, rh.IsEmpty())

	_, _, err := rh.Pop()
	assert.Error(t, err)
}

func TestRadixHeapPushMonotonicity(t *testing.T) {
	rh := NewRadixHeap([]HeapNode[string, uint]{
		CreateHeapNode("value2", uint(2)),
		CreateHeapNode("value4", uint(4)),
		CreateHeapNode("value6", uint(6)),
	}, false)

	_, priority, err := rh.Pop()
	assert.NoError(t, err)
	assert.Equal(t, uint(2), priority)

	err = rh.Push("value3", uint(3))
	assert.NoError(t, err)
	_, priority, err = rh.Peek()
	assert.NoError(t, err)
	assert.Equal(t, uint(3), priority)

	err = rh.Push("value1", uint(1))
	assert.Error(t, err)
}

func TestRadixHeapPeek(t *testing.T) {
	rh := NewRadixHeap([]HeapNode[string, uint]{
		CreateHeapNode("value8", uint(8)),
		CreateHeapNode("value2", uint(2)),
		CreateHeapNode("value5", uint(5)),
	}, false)
	_, priority, err := rh.Peek()
	assert.NoError(t, err)
	assert.Equal(t, uint(2), priority)

	// removes 2
	_, _, _ = rh.Pop()
	_, priority, err = rh.Peek()
	assert.NoError(t, err)
	assert.Equal(t, uint(5), priority)

	// clearing then Peek should return error
	rh.Clear()
	_, _, err = rh.Peek()
	assert.Error(t, err)
}

func TestRadixHeapClearCloneDeepClone(t *testing.T) {
	original := []HeapNode[string, uint]{
		CreateHeapNode("value4", uint(4)),
		CreateHeapNode("value1", uint(1)),
		CreateHeapNode("value3", uint(3)),
	}
	rh := NewRadixHeap(original, false)
	assert.Equal(t, 3, rh.Length())

	clone := rh.Clone()
	assert.Equal(t, rh.Length(), clone.Length())

	// now last = 1, size = 2
	_, _, _ = rh.Pop()

	// valid since 2 >= last
	err := rh.Push("value2", uint(2))
	assert.NoError(t, err)

	cloneVals := []uint{}
	for !clone.IsEmpty() {
		_, priority, _ := clone.Pop()
		cloneVals = append(cloneVals, priority)
	}
	assert.Equal(t, []uint{1, 3, 4}, cloneVals)
}

func TestRadixHeapMerge(t *testing.T) {
	rh1 := NewRadixHeap([]HeapNode[string, uint]{
		CreateHeapNode("value1", uint(1)),
		CreateHeapNode("value4", uint(4)),
		CreateHeapNode("value6", uint(6)),
	}, false)
	rh2 := NewRadixHeap([]HeapNode[string, uint]{
		CreateHeapNode("value2", uint(2)),
		CreateHeapNode("value3", uint(3)),
		CreateHeapNode("value5", uint(5)),
	}, false)
	rh1.Merge(rh2)

	result := []uint{}
	for !rh1.IsEmpty() {
		_, priority, err := rh1.Pop()
		assert.NoError(t, err)
		result = append(result, priority)
	}
	assert.Equal(t, []uint{1, 2, 3, 4, 5, 6}, result)
}

func TestRadixHeapRemoveAndErrors(t *testing.T) {
	rh := NewRadixHeap([]HeapNode[string, uint]{}, false)
	assert.True(t, rh.IsEmpty())
	_, _, err := rh.Pop()
	assert.Error(t, err)

	rh.Clear()
	err = rh.Push("value0", uint(0))
	assert.NoError(t, err)
	_, priority, err := rh.Peek()
	assert.NoError(t, err)
	assert.Equal(t, uint(0), priority)
}

func TestRadixHeapLengthIsEmpty(t *testing.T) {
	rh := NewRadixHeap([]HeapNode[string, uint]{}, false)
	assert.True(t, rh.IsEmpty())
	assert.Equal(t, 0, rh.Length())

	_ = rh.Push("value7", uint(7))
	assert.False(t, rh.IsEmpty())
	assert.Equal(t, 1, rh.Length())
}

// -------------------------------- Radix Heap Benchmarks --------------------------------

func BenchmarkRadixHeapInsertion(b *testing.B) {
	data := make([]HeapNode[int, uint], 0)
	heap := NewRadixHeap(data, false)

	insertions := generateRandomNumbersv1(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Push(insertions[i], uint(insertions[i]))
	}
}

func BenchmarkRadixHeapDeletion(b *testing.B) {
	data := make([]HeapNode[int, uint], 0)
	heap := NewRadixHeap(data, false)

	for i := 0; i < b.N; i++ {
		heap.Push(i, uint(i))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Pop()
	}
}
