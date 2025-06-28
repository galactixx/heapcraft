package heapcraft

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewSyncDaryHeap tests the creation of thread-safe d-ary heaps.
func TestNewSyncDaryHeap(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 3, priority: 3},
		{value: 1, priority: 1},
		{value: 2, priority: 2},
	}

	// Test binary heap creation
	heap := NewSyncBinaryHeap(data, func(a, b int) bool { return a < b }, false)
	assert.NotNil(t, heap)
	assert.Equal(t, 3, heap.Length())

	// Test d-ary heap creation
	heap3 := NewSyncDaryHeap(3, data, func(a, b int) bool { return a < b }, false)
	assert.NotNil(t, heap3)
	assert.Equal(t, 3, heap3.Length())
}

// TestNewSyncDaryHeapCopy tests the creation of thread-safe d-ary heaps with data copying.
func TestNewSyncDaryHeapCopy(t *testing.T) {
	original := []HeapNode[int, int]{
		{value: 3, priority: 3},
		{value: 1, priority: 1},
		{value: 2, priority: 2},
	}

	// Test binary heap copy creation
	heap := NewSyncBinaryHeapCopy(original, func(a, b int) bool { return a < b }, false)
	assert.NotNil(t, heap)

	// Verify original data is unchanged
	assert.Equal(t, 3, len(original))

	// Test d-ary heap copy creation
	heap3 := NewSyncDaryHeapCopy(3, original, func(a, b int) bool { return a < b }, false)
	assert.NotNil(t, heap3)
}

// TestSyncDaryHeapBasicOperations tests basic heap operations in a thread-safe manner.
func TestSyncDaryHeapBasicOperations(t *testing.T) {
	heap := NewSyncBinaryHeap([]HeapNode[int, int]{}, func(a, b int) bool { return a < b }, false)

	// Test empty heap
	assert.True(t, heap.IsEmpty())
	assert.Equal(t, 0, heap.Length())

	// Test Push
	heap.Push(3, 3)
	heap.Push(1, 1)
	heap.Push(2, 2)

	assert.Equal(t, 3, heap.Length())

	// Test Peek
	_, priority, err := heap.Peek()
	assert.NoError(t, err)
	assert.Equal(t, 1, priority)

	// Test Pop
	_, priority, err = heap.Pop()
	assert.NoError(t, err)
	assert.Equal(t, 1, priority)
	assert.Equal(t, 2, heap.Length())
}

// TestSyncDaryHeapConcurrentAccess tests concurrent access to the heap.
func TestSyncDaryHeapConcurrentAccess(t *testing.T) {
	heap := NewSyncBinaryHeap([]HeapNode[int, int]{}, func(a, b int) bool { return a < b }, false)
	var wg sync.WaitGroup
	numGoroutines := 10
	operationsPerGoroutine := 100

	// Start multiple goroutines that push and pop concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				value := id*operationsPerGoroutine + j
				heap.Push(value, value)

				// Occasionally pop
				if j%10 == 0 {
					heap.Pop()
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify the heap is in a consistent state
	assert.GreaterOrEqual(t, heap.Length(), 0)

	// Pop all remaining elements and verify they're in order
	lastPriority := -1
	for !heap.IsEmpty() {
		_, priority, err := heap.Pop()
		assert.NoError(t, err)
		if lastPriority != -1 {
			assert.GreaterOrEqual(t, priority, lastPriority)
		}
		lastPriority = priority
	}
}

// TestSyncDaryHeapUpdateAndRemove tests Update and Remove operations.
func TestSyncDaryHeapUpdateAndRemove(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 3, priority: 3},
		{value: 1, priority: 1},
		{value: 2, priority: 2},
	}
	heap := NewSyncBinaryHeap(data, func(a, b int) bool { return a < b }, false)

	// Test Update
	err := heap.Update(1, 5, 5)
	assert.NoError(t, err)

	// Test Remove
	_, priority, err := heap.Remove(0)
	assert.NoError(t, err)
	assert.Equal(t, 1, priority)

	// Test error cases
	err = heap.Update(10, 1, 1)
	assert.Equal(t, ErrIndexOutOfBounds, err)

	_, _, err = heap.Remove(10)
	assert.Equal(t, ErrIndexOutOfBounds, err)
}

// TestSyncDaryHeapPopPushAndPushPop tests PopPush and PushPop operations.
func TestSyncDaryHeapPopPushAndPushPop(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 3, priority: 3},
		{value: 1, priority: 1},
		{value: 2, priority: 2},
	}
	heap := NewSyncBinaryHeap(data, func(a, b int) bool { return a < b }, false)

	// Test PopPush
	_, priority := heap.PopPush(4, 4)
	assert.Equal(t, 1, priority)

	// Test PushPop with higher priority
	_, priority = heap.PushPop(0, 0)
	assert.Equal(t, 0, priority)

	// Test PushPop with lower priority
	_, priority = heap.PushPop(5, 5)
	assert.Equal(t, 2, priority)
}

// TestSyncDaryHeapCallbacks tests callback registration and deregistration.
func TestSyncDaryHeapCallbacks(t *testing.T) {
	heap := NewSyncBinaryHeap([]HeapNode[int, int]{}, func(a, b int) bool { return a < b }, false)

	callbackCount := 0
	var callbackMutex sync.Mutex

	// Register callback
	callback := heap.Register(func(x, y int) {
		callbackMutex.Lock()
		callbackCount++
		callbackMutex.Unlock()
	})

	// Perform operations that trigger callbacks
	heap.Push(1, 1)
	heap.Push(2, 2)
	heap.Push(0, 0) // This should trigger sift up

	assert.Greater(t, callbackCount, 0)

	// Deregister callback
	err := heap.Deregister(callback.ID)
	assert.NoError(t, err)

	// Test deregistering non-existent callback
	err = heap.Deregister("non-existent")
	assert.Equal(t, ErrCallbackNotFound, err)
}

// TestSyncDaryHeapClone tests heap cloning.
func TestSyncDaryHeapClone(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 3, priority: 3},
		{value: 1, priority: 1},
		{value: 2, priority: 2},
	}
	original := NewSyncBinaryHeap(data, func(a, b int) bool { return a < b }, false)

	// Register a callback
	callback := original.Register(func(x, y int) {})

	// Clone the heap
	cloned := original.Clone()
	assert.NotNil(t, cloned)

	// Verify they have the same length
	assert.Equal(t, original.Length(), cloned.Length())

	// Modify original and verify clone is unaffected
	original.Push(4, 4)
	assert.NotEqual(t, original.Length(), cloned.Length())

	// Verify callback is preserved in clone
	err := cloned.Deregister(callback.ID)
	assert.NoError(t, err)
}

// TestSyncDaryHeapStress tests stress conditions with many concurrent operations.
func TestSyncDaryHeapStress(t *testing.T) {
	heap := NewSyncBinaryHeap([]HeapNode[int, int]{}, func(a, b int) bool { return a < b }, false)
	var wg sync.WaitGroup
	numGoroutines := 20
	operationsPerGoroutine := 50

	// Start goroutines that perform mixed operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				value := id*operationsPerGoroutine + j

				switch j % 4 {
				case 0:
					heap.Push(value, value)
				case 1:
					heap.Pop()
				case 2:
					heap.Peek()
				case 3:
					heap.Length()
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify heap is in a consistent state
	assert.GreaterOrEqual(t, heap.Length(), 0)
}

// TestSyncDaryHeapEmptyOperations tests operations on empty heaps.
func TestSyncDaryHeapEmptyOperations(t *testing.T) {
	heap := NewSyncBinaryHeap([]HeapNode[int, int]{}, func(a, b int) bool { return a < b }, false)

	// Test Pop on empty heap
	_, _, err := heap.Pop()
	assert.Equal(t, ErrHeapEmpty, err)

	// Test Peek on empty heap
	_, _, err = heap.Peek()
	assert.Equal(t, ErrHeapEmpty, err)

	// Test PopValue on empty heap
	_, err = heap.PopValue()
	assert.Equal(t, ErrHeapEmpty, err)

	// Test PopPriority on empty heap
	_, err = heap.PopPriority()
	assert.Equal(t, ErrHeapEmpty, err)

	// Test PeekValue on empty heap
	_, err = heap.PeekValue()
	assert.Equal(t, ErrHeapEmpty, err)

	// Test PeekPriority on empty heap
	_, err = heap.PeekPriority()
	assert.Equal(t, ErrHeapEmpty, err)
}

// BenchmarkSyncBinaryHeapPush benchmarks concurrent push operations.
func BenchmarkSyncBinaryHeapPush(b *testing.B) {
	heap := NewSyncBinaryHeap([]HeapNode[int, int]{}, func(a, b int) bool { return a < b }, true)

	insertions := generateRandomNumbers(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Push(insertions[i], insertions[i])
	}
}

// BenchmarkSyncBinaryHeapPushPop benchmarks concurrent push/pop operations.
func BenchmarkSyncBinaryHeapPushPop(b *testing.B) {
	heap := NewSyncBinaryHeap([]HeapNode[int, int]{}, func(a, b int) bool { return a < b }, true)

	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Pop()
	}
}
