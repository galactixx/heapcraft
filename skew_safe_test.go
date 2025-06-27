package heapcraft

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSafeSkewHeap_BasicOperations(t *testing.T) {
	heap := NewSafeSkewHeap[int, int](nil, func(a, b int) bool { return a < b })

	assert.True(t, heap.IsEmpty())
	assert.Equal(t, 0, heap.Length())

	id1 := heap.Push(10, 1)
	id2 := heap.Push(20, 2)
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

func TestSafeSkewHeap_ConcurrentAccess(t *testing.T) {
	heap := NewSafeSkewHeap[int, int](nil, func(a, b int) bool { return a < b })
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

func TestSafeSkewHeap_Clone(t *testing.T) {
	heap := NewSafeSkewHeap[int, int](nil, func(a, b int) bool { return a < b })
	heap.Push(10, 1)
	heap.Push(20, 2)

	clone := heap.Clone()

	assert.Equal(t, heap.Length(), clone.Length())

	heap.Push(30, 3)

	assert.Equal(t, 2, clone.Length())
}

func TestSafeSkewHeap_EmptyOperations(t *testing.T) {
	heap := NewSafeSkewHeap[int, int](nil, func(a, b int) bool { return a < b })

	_, _, err := heap.Pop()
	assert.Equal(t, ErrHeapEmpty, err)

	_, _, err = heap.Peek()
	assert.Equal(t, ErrHeapEmpty, err)

	_, _, err = heap.Get("nonexistent")
	assert.Equal(t, ErrNodeNotFound, err)
}

func TestSafeSimpleSkewHeap_BasicOperations(t *testing.T) {
	heap := NewSafeSimpleSkewHeap[int, int](nil, func(a, b int) bool { return a < b })

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

func TestSafeSimpleSkewHeap_ConcurrentAccess(t *testing.T) {
	heap := NewSafeSimpleSkewHeap[int, int](nil, func(a, b int) bool { return a < b })
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

func TestSafeSimpleSkewHeap_Clone(t *testing.T) {
	heap := NewSafeSimpleSkewHeap[int, int](nil, func(a, b int) bool { return a < b })
	heap.Push(10, 1)
	heap.Push(20, 2)

	clone := heap.Clone()

	assert.Equal(t, heap.Length(), clone.Length())

	heap.Push(30, 3)

	assert.Equal(t, 2, clone.Length())
}

func TestSafeSimpleSkewHeap_EmptyOperations(t *testing.T) {
	heap := NewSafeSimpleSkewHeap[int, int](nil, func(a, b int) bool { return a < b })

	_, _, err := heap.Pop()
	assert.Equal(t, ErrHeapEmpty, err)

	_, _, err = heap.Peek()
	assert.Equal(t, ErrHeapEmpty, err)
}
