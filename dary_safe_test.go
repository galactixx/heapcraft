package heapcraft

import (
	"sync"
	"testing"
)

// TestNewSyncDaryHeap tests the creation of thread-safe d-ary heaps.
func TestNewSyncDaryHeap(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 3, priority: 3},
		{value: 1, priority: 1},
		{value: 2, priority: 2},
	}

	// Test binary heap creation
	heap := NewSyncBinaryHeap(data, func(a, b int) bool { return a < b })
	if heap == nil {
		t.Fatal("NewSyncBinaryHeap returned nil")
	}
	if heap.Length() != 3 {
		t.Errorf("Expected length 3, got %d", heap.Length())
	}

	// Test d-ary heap creation
	heap3 := NewSyncDaryHeap(3, data, func(a, b int) bool { return a < b })
	if heap3 == nil {
		t.Fatal("NewSyncDaryHeap returned nil")
	}
	if heap3.Length() != 3 {
		t.Errorf("Expected length 3, got %d", heap3.Length())
	}
}

// TestNewSyncDaryHeapCopy tests the creation of thread-safe d-ary heaps with data copying.
func TestNewSyncDaryHeapCopy(t *testing.T) {
	original := []HeapNode[int, int]{
		{value: 3, priority: 3},
		{value: 1, priority: 1},
		{value: 2, priority: 2},
	}

	// Test binary heap copy creation
	heap := NewSyncBinaryHeapCopy(original, func(a, b int) bool { return a < b })
	if heap == nil {
		t.Fatal("NewSyncBinaryHeapCopy returned nil")
	}

	// Verify original data is unchanged
	if len(original) != 3 {
		t.Errorf("Original data length changed, expected 3, got %d", len(original))
	}

	// Test d-ary heap copy creation
	heap3 := NewSyncDaryHeapCopy(3, original, func(a, b int) bool { return a < b })
	if heap3 == nil {
		t.Fatal("NewSyncDaryHeapCopy returned nil")
	}
}

// TestSyncDaryHeapBasicOperations tests basic heap operations in a thread-safe manner.
func TestSyncDaryHeapBasicOperations(t *testing.T) {
	heap := NewSyncBinaryHeap([]HeapNode[int, int]{}, func(a, b int) bool { return a < b })

	// Test empty heap
	if !heap.IsEmpty() {
		t.Error("New heap should be empty")
	}
	if heap.Length() != 0 {
		t.Errorf("Expected length 0, got %d", heap.Length())
	}

	// Test Push
	heap.Push(3, 3)
	heap.Push(1, 1)
	heap.Push(2, 2)

	if heap.Length() != 3 {
		t.Errorf("Expected length 3, got %d", heap.Length())
	}

	// Test Peek
	peeked, err := heap.Peek()
	if err != nil {
		t.Errorf("Peek failed: %v", err)
	}
	if peeked.Priority() != 1 {
		t.Errorf("Expected priority 1, got %d", peeked.Priority())
	}

	// Test Pop
	popped, err := heap.Pop()
	if err != nil {
		t.Errorf("Pop failed: %v", err)
	}
	if popped.Priority() != 1 {
		t.Errorf("Expected priority 1, got %d", popped.Priority())
	}
	if heap.Length() != 2 {
		t.Errorf("Expected length 2, got %d", heap.Length())
	}
}

// TestSyncDaryHeapConcurrentAccess tests concurrent access to the heap.
func TestSyncDaryHeapConcurrentAccess(t *testing.T) {
	heap := NewSyncBinaryHeap([]HeapNode[int, int]{}, func(a, b int) bool { return a < b })
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
	if heap.Length() < 0 {
		t.Error("Heap length should not be negative")
	}

	// Pop all remaining elements and verify they're in order
	lastPriority := -1
	for !heap.IsEmpty() {
		popped, err := heap.Pop()
		if err != nil {
			t.Errorf("Pop failed: %v", err)
		}
		if lastPriority != -1 && popped.Priority() < lastPriority {
			t.Errorf("Heap property violated: %d came after %d", popped.Priority(), lastPriority)
		}
		lastPriority = popped.Priority()
	}
}

// TestSyncDaryHeapUpdateAndRemove tests Update and Remove operations.
func TestSyncDaryHeapUpdateAndRemove(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 3, priority: 3},
		{value: 1, priority: 1},
		{value: 2, priority: 2},
	}
	heap := NewSyncBinaryHeap(data, func(a, b int) bool { return a < b })

	// Test Update
	err := heap.Update(1, 5, 5)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	// Test Remove
	removed, err := heap.Remove(0)
	if err != nil {
		t.Errorf("Remove failed: %v", err)
	}
	if removed.Priority() != 1 {
		t.Errorf("Expected priority 1, got %d", removed.Priority())
	}

	// Test error cases
	err = heap.Update(10, 1, 1)
	if err != ErrIndexOutOfBounds {
		t.Errorf("Expected ErrIndexOutOfBounds, got %v", err)
	}

	_, err = heap.Remove(10)
	if err != ErrIndexOutOfBounds {
		t.Errorf("Expected ErrIndexOutOfBounds, got %v", err)
	}
}

// TestSyncDaryHeapPopPushAndPushPop tests PopPush and PushPop operations.
func TestSyncDaryHeapPopPushAndPushPop(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 3, priority: 3},
		{value: 1, priority: 1},
		{value: 2, priority: 2},
	}
	heap := NewSyncBinaryHeap(data, func(a, b int) bool { return a < b })

	// Test PopPush
	removed := heap.PopPush(4, 4)
	if removed.Priority() != 1 {
		t.Errorf("Expected priority 1, got %d", removed.Priority())
	}

	// Test PushPop with higher priority
	result := heap.PushPop(0, 0)
	if result.Priority() != 0 {
		t.Errorf("Expected priority 0, got %d", result.Priority())
	}

	// Test PushPop with lower priority
	result = heap.PushPop(5, 5)
	if result.Priority() != 2 {
		t.Errorf("Expected priority 2, got %d", result.Priority())
	}
}

// TestSyncDaryHeapCallbacks tests callback registration and deregistration.
func TestSyncDaryHeapCallbacks(t *testing.T) {
	heap := NewSyncBinaryHeap([]HeapNode[int, int]{}, func(a, b int) bool { return a < b })

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

	if callbackCount == 0 {
		t.Error("No callbacks were triggered")
	}

	// Deregister callback
	err := heap.Deregister(callback.ID)
	if err != nil {
		t.Errorf("Deregister failed: %v", err)
	}

	// Test deregistering non-existent callback
	err = heap.Deregister("non-existent")
	if err != ErrCallbackNotFound {
		t.Errorf("Expected ErrCallbackNotFound, got %v", err)
	}
}

// TestSyncDaryHeapClone tests heap cloning.
func TestSyncDaryHeapClone(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 3, priority: 3},
		{value: 1, priority: 1},
		{value: 2, priority: 2},
	}
	original := NewSyncBinaryHeap(data, func(a, b int) bool { return a < b })

	// Register a callback
	callback := original.Register(func(x, y int) {})

	// Clone the heap
	cloned := original.Clone()
	if cloned == nil {
		t.Fatal("Clone returned nil")
	}

	// Verify they have the same length
	if original.Length() != cloned.Length() {
		t.Errorf("Original length %d != cloned length %d", original.Length(), cloned.Length())
	}

	// Modify original and verify clone is unaffected
	original.Push(4, 4)
	if original.Length() == cloned.Length() {
		t.Error("Clone should be independent of original")
	}

	// Verify callback is preserved in clone
	err := cloned.Deregister(callback.ID)
	if err != nil {
		t.Errorf("Failed to deregister callback in clone: %v", err)
	}
}

// TestSyncDaryHeapStress tests stress conditions with many concurrent operations.
func TestSyncDaryHeapStress(t *testing.T) {
	heap := NewSyncBinaryHeap([]HeapNode[int, int]{}, func(a, b int) bool { return a < b })
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
	if heap.Length() < 0 {
		t.Error("Heap length should not be negative")
	}
}

// TestSyncDaryHeapEmptyOperations tests operations on empty heaps.
func TestSyncDaryHeapEmptyOperations(t *testing.T) {
	heap := NewSyncBinaryHeap([]HeapNode[int, int]{}, func(a, b int) bool { return a < b })

	// Test Pop on empty heap
	_, err := heap.Pop()
	if err != ErrHeapEmpty {
		t.Errorf("Expected ErrHeapEmpty, got %v", err)
	}

	// Test Peek on empty heap
	_, err = heap.Peek()
	if err != ErrHeapEmpty {
		t.Errorf("Expected ErrHeapEmpty, got %v", err)
	}

	// Test PopValue on empty heap
	_, err = heap.PopValue()
	if err != ErrHeapEmpty {
		t.Errorf("Expected ErrHeapEmpty, got %v", err)
	}

	// Test PopPriority on empty heap
	_, err = heap.PopPriority()
	if err != ErrHeapEmpty {
		t.Errorf("Expected ErrHeapEmpty, got %v", err)
	}

	// Test PeekValue on empty heap
	_, err = heap.PeekValue()
	if err != ErrHeapEmpty {
		t.Errorf("Expected ErrHeapEmpty, got %v", err)
	}

	// Test PeekPriority on empty heap
	_, err = heap.PeekPriority()
	if err != ErrHeapEmpty {
		t.Errorf("Expected ErrHeapEmpty, got %v", err)
	}
}

// BenchmarkSyncBinaryHeapPush benchmarks concurrent push operations.
func BenchmarkSyncBinaryHeapPush(b *testing.B) {
	heap := NewSyncBinaryHeap([]HeapNode[int, int]{}, func(a, b int) bool { return a < b })
	b.ReportAllocs()

	insertions := generateRandomNumbers(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Push(insertions[i], insertions[i])
	}
}

// BenchmarkSyncBinaryHeapPushPop benchmarks concurrent push/pop operations.
func BenchmarkSyncBinaryHeapPushPop(b *testing.B) {
	heap := NewSyncBinaryHeap([]HeapNode[int, int]{}, func(a, b int) bool { return a < b })

	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Pop()
	}
}
