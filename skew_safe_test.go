package heapcraft

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncFullSkewHeap_BasicOperations(t *testing.T) {
	heap := NewSyncFullSkewHeap[int](nil, lt, HeapConfig{UsePool: false})

	assert.True(t, heap.IsEmpty())
	assert.Equal(t, 0, heap.Length())

	id1, _ := heap.Push(10, 1)
	id2, _ := heap.Push(20, 2)
	heap.Push(5, 0)

	assert.False(t, heap.IsEmpty())
	assert.Equal(t, 3, heap.Length())

	value, err := heap.PeekValue()
	require.NoError(t, err)
	assert.Equal(t, 5, value)

	value, err = heap.PopValue()
	require.NoError(t, err)
	assert.Equal(t, 5, value)
	assert.Equal(t, 2, heap.Length())

	value, priority, err := heap.Get(id1)
	require.NoError(t, err)
	assert.Equal(t, 10, value)
	assert.Equal(t, 1, priority)

	err = heap.UpdateValue(id2, 25)
	require.NoError(t, err)
	value, err = heap.GetValue(id2)
	require.NoError(t, err)
	assert.Equal(t, 25, value)

	err = heap.UpdatePriority(id1, 0)
	require.NoError(t, err)
	value, err = heap.PeekValue()
	require.NoError(t, err)
	assert.Equal(t, 10, value)

	heap.Clear()
	assert.True(t, heap.IsEmpty())
}

func TestSyncFullSkewHeap_ConcurrentAccess(t *testing.T) {
	heap := NewSyncFullSkewHeap[int](nil, lt, HeapConfig{UsePool: false})
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			heap.Push(val, val)
		}(i)
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			heap.PeekValue()
		}()
	}

	wg.Wait()

	assert.Equal(t, 10, heap.Length())
}

func TestSyncFullSkewHeap_Clone(t *testing.T) {
	heap := NewSyncFullSkewHeap[int](nil, lt, HeapConfig{UsePool: false})
	heap.Push(10, 1)
	heap.Push(20, 2)

	clone := heap.Clone()

	assert.Equal(t, heap.Length(), clone.Length())

	heap.Push(30, 3)

	assert.Equal(t, 2, clone.Length())
}

func TestSyncFullSkewHeap_EmptyOperations(t *testing.T) {
	heap := NewSyncFullSkewHeap[int](nil, lt, HeapConfig{UsePool: false})

	_, _, err := heap.Pop()
	assert.Equal(t, ErrHeapEmpty, err)

	_, _, err = heap.Peek()
	assert.Equal(t, ErrHeapEmpty, err)

	_, _, err = heap.Get("nonexistent")
	assert.Equal(t, ErrNodeNotFound, err)
}

func TestSyncSkewHeap_BasicOperations(t *testing.T) {
	heap := NewSyncSkewHeap[int](nil, lt, false)

	assert.True(t, heap.IsEmpty())
	assert.Equal(t, 0, heap.Length())

	heap.Push(10, 1)
	heap.Push(20, 2)
	heap.Push(5, 0)

	assert.False(t, heap.IsEmpty())
	assert.Equal(t, 3, heap.Length())

	value, err := heap.PeekValue()
	require.NoError(t, err)
	assert.Equal(t, 5, value)

	value, err = heap.PopValue()
	require.NoError(t, err)
	assert.Equal(t, 5, value)
	assert.Equal(t, 2, heap.Length())

	heap.Clear()
	assert.True(t, heap.IsEmpty())
}

func TestSyncSkewHeap_ConcurrentAccess(t *testing.T) {
	heap := NewSyncSkewHeap[int](nil, lt, false)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			heap.Push(val, val)
		}(i)
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			heap.PeekValue()
		}()
	}

	wg.Wait()

	assert.Equal(t, 10, heap.Length())
}

func TestSyncSkewHeap_Clone(t *testing.T) {
	heap := NewSyncSkewHeap[int](nil, lt, false)
	heap.Push(10, 1)
	heap.Push(20, 2)

	clone := heap.Clone()

	assert.Equal(t, heap.Length(), clone.Length())

	heap.Push(30, 3)

	assert.Equal(t, 2, clone.Length())
}

func TestSyncSkewHeap_EmptyOperations(t *testing.T) {
	heap := NewSyncSkewHeap[int](nil, lt, false)

	_, _, err := heap.Pop()
	assert.Equal(t, ErrHeapEmpty, err)

	_, _, err = heap.Peek()
	assert.Equal(t, ErrHeapEmpty, err)
}
