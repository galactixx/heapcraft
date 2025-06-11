package heapcraft

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRadixHeapPopOrder(t *testing.T) {
	raw := []*RadixPair[string, uint]{
		CreateRadixPair("value10", uint(10)),
		CreateRadixPair("value3", uint(3)),
		CreateRadixPair("value7", uint(7)),
		CreateRadixPair("value1", uint(1)),
		CreateRadixPair("value5", uint(5)),
		CreateRadixPair("value2", uint(2)),
	}
	rh := NewRadixHeap(raw)
	assert.False(t, rh.IsEmpty())
	assert.Equal(t, len(raw), rh.Length())

	expected := []*RadixPair[string, uint]{
		CreateRadixPair("value1", uint(1)),
		CreateRadixPair("value2", uint(2)),
		CreateRadixPair("value3", uint(3)),
		CreateRadixPair("value5", uint(5)),
		CreateRadixPair("value7", uint(7)),
		CreateRadixPair("value10", uint(10)),
	}
	actualValues := []string{}
	actualPriorities := []uint{}
	for !rh.IsEmpty() {
		v, err := rh.Pop()
		assert.NoError(t, err)
		actualValues = append(actualValues, v.Value())
		actualPriorities = append(actualPriorities, v.Priority())
	}
	for i := range expected {
		assert.Equal(t, expected[i].Value(), actualValues[i])
		assert.Equal(t, expected[i].Priority(), actualPriorities[i])
	}
	assert.True(t, rh.IsEmpty())

	_, err := rh.Pop()
	assert.Error(t, err)
}

func TestRadixHeapPushMonotonicity(t *testing.T) {
	rh := NewRadixHeap([]*RadixPair[string, uint]{
		CreateRadixPair("value2", uint(2)),
		CreateRadixPair("value4", uint(4)),
		CreateRadixPair("value6", uint(6)),
	})

	minPtr, err := rh.Pop()
	assert.NoError(t, err)
	assert.Equal(t, uint(2), minPtr.Priority())

	err = rh.Push("value3", uint(3))
	assert.NoError(t, err)
	peeked, err := rh.Peek()
	assert.NoError(t, err)
	assert.Equal(t, uint(3), peeked.Priority())

	err = rh.Push("value1", uint(1))
	assert.Error(t, err)
}

func TestRadixHeapPeek(t *testing.T) {
	rh := NewRadixHeap([]*RadixPair[string, uint]{
		CreateRadixPair("value8", uint(8)),
		CreateRadixPair("value2", uint(2)),
		CreateRadixPair("value5", uint(5)),
	})
	peeked, err := rh.Peek()
	assert.NoError(t, err)
	assert.Equal(t, uint(2), peeked.Priority())

	// removes 2
	_, _ = rh.Pop()
	peeked, err = rh.Peek()
	assert.NoError(t, err)
	assert.Equal(t, uint(5), peeked.Priority())

	// clearing then Peek should return error
	rh.Clear()
	_, err = rh.Peek()
	assert.Error(t, err)
}

func TestRadixHeapClearCloneDeepClone(t *testing.T) {
	original := []*RadixPair[string, uint]{
		CreateRadixPair("value4", uint(4)),
		CreateRadixPair("value1", uint(1)),
		CreateRadixPair("value3", uint(3)),
	}
	rh := NewRadixHeap(original)
	assert.Equal(t, 3, rh.Length())

	clone := rh.Clone()
	assert.Equal(t, rh.Length(), clone.Length())

	// now last = 1, size = 2
	_, _ = rh.Pop()

	// valid since 2 >= last
	err := rh.Push("value2", uint(2))
	assert.NoError(t, err)

	cloneVals := []uint{}
	for !clone.IsEmpty() {
		vPtr, _ := clone.Pop()
		cloneVals = append(cloneVals, vPtr.Priority())
	}
	assert.Equal(t, []uint{1, 3, 4}, cloneVals)
}

func TestRadixHeapMerge(t *testing.T) {
	rh1 := NewRadixHeap([]*RadixPair[string, uint]{
		CreateRadixPair("value1", uint(1)),
		CreateRadixPair("value4", uint(4)),
		CreateRadixPair("value6", uint(6)),
	})
	rh2 := NewRadixHeap([]*RadixPair[string, uint]{
		CreateRadixPair("value2", uint(2)),
		CreateRadixPair("value3", uint(3)),
		CreateRadixPair("value5", uint(5)),
	})
	rh1.Merge(rh2)

	result := []uint{}
	for !rh1.IsEmpty() {
		vPtr, err := rh1.Pop()
		assert.NoError(t, err)
		result = append(result, vPtr.Priority())
	}
	assert.Equal(t, []uint{1, 2, 3, 4, 5, 6}, result)
}

func TestRadixHeapRemoveAndErrors(t *testing.T) {
	rh := NewRadixHeap([]*RadixPair[string, uint]{})
	assert.True(t, rh.IsEmpty())
	_, err := rh.Pop()
	assert.Error(t, err)

	rh.Clear()
	err = rh.Push("value0", uint(0))
	assert.NoError(t, err)
	peeked, err := rh.Peek()
	assert.NoError(t, err)
	assert.Equal(t, uint(0), peeked.Priority())
}

func TestRadixHeapLengthIsEmpty(t *testing.T) {
	rh := NewRadixHeap([]*RadixPair[string, uint]{})
	assert.True(t, rh.IsEmpty())
	assert.Equal(t, 0, rh.Length())

	_ = rh.Push("value7", uint(7))
	assert.False(t, rh.IsEmpty())
	assert.Equal(t, 1, rh.Length())
}

// Radix Heap Benchmarks
func BenchmarkRadixHeapInsertion(b *testing.B) {
	N := 10_000
	data := make([]*RadixPair[int, uint], 0)
	heap := NewRadixHeap(data)
	b.ReportAllocs()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var num uint
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			num = uint(r.Intn(N))
			heap.Push(int(num), num)
		}
	})
}

func BenchmarkRadixHeapDeletion(b *testing.B) {
	data := make([]*RadixPair[int, uint], 0)
	heap := NewRadixHeap(data)

	for i := 0; i < b.N; i++ {
		err := heap.Push(i, uint(i))
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			heap.Pop()
		}
	})
}
