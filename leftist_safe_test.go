package heapcraft

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSafeLeftistHeap_BasicOperations(t *testing.T) {
	heap := NewSafeLeftistHeap[int, int](nil, func(a, b int) bool { return a < b })

	// Test empty heap
	assert.True(t, heap.IsEmpty())
	assert.Equal(t, 0, heap.Length())

	// Test Push
	id1 := heap.Push(10, 1)
	id2 := heap.Push(20, 2)
	heap.Push(5, 0)

	assert.False(t, heap.IsEmpty())
	assert.Equal(t, 3, heap.Length())

	// Test Peek
	value, err := heap.PeekValue()
	require.NoError(t, err)
	assert.Equal(t, 5, value)

	// Test Pop
	value, err = heap.PopValue()
	require.NoError(t, err)
	assert.Equal(t, 5, value)
	assert.Equal(t, 2, heap.Length())

	// Test Get
	value, priority, err := heap.Get(id1)
	require.NoError(t, err)
	assert.Equal(t, 10, value)
	assert.Equal(t, 1, priority)

	// Test UpdateValue
	err = heap.UpdateValue(id2, 25)
	require.NoError(t, err)
	value, err = heap.GetValue(id2)
	require.NoError(t, err)
	assert.Equal(t, 25, value)

	// Test UpdatePriority
	err = heap.UpdatePriority(id1, 0)
	require.NoError(t, err)
	value, err = heap.PeekValue()
	require.NoError(t, err)
	assert.Equal(t, 10, value)

	// Test Clear
	heap.Clear()
	assert.True(t, heap.IsEmpty())
}

func TestSafeLeftistHeap_ConcurrentAccess(t *testing.T) {
	heap := NewSafeLeftistHeap[int, int](nil, func(a, b int) bool { return a < b })
	var wg sync.WaitGroup

	// Concurrent pushes
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			heap.Push(val, val)
		}(i)
	}

	// Concurrent peeks
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

func TestSafeLeftistHeap_Clone(t *testing.T) {
	heap := NewSafeLeftistHeap[int, int](nil, func(a, b int) bool { return a < b })
	heap.Push(10, 1)
	heap.Push(20, 2)

	clone := heap.Clone()

	assert.Equal(t, heap.Length(), clone.Length())

	// Modify original
	heap.Push(30, 3)

	// Clone should be unaffected
	assert.Equal(t, 2, clone.Length())
}

func TestSafeLeftistHeap_EmptyOperations(t *testing.T) {
	heap := NewSafeLeftistHeap[int, int](nil, func(a, b int) bool { return a < b })

	// Test Pop on empty heap
	_, _, err := heap.Pop()
	assert.Equal(t, ErrHeapEmpty, err)

	// Test Peek on empty heap
	_, _, err = heap.Peek()
	assert.Equal(t, ErrHeapEmpty, err)

	// Test Get on non-existent ID
	_, _, err = heap.Get("nonexistent")
	assert.Equal(t, ErrNodeNotFound, err)
}
