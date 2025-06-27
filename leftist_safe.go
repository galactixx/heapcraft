package heapcraft

import (
	"sync"
)

// SafeLeftistHeap is a thread-safe wrapper around LeftistHeap.
// All operations are protected by a sync.RWMutex, making it safe for concurrent use.
type SafeLeftistHeap[V any, P any] struct {
	heap *LeftistHeap[V, P]
	lock sync.RWMutex
}

// NewSafeLeftistHeap constructs a new thread-safe leftist heap from the given data and comparison function.
// The resulting heap is safe for concurrent use.
func NewSafeLeftistHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *SafeLeftistHeap[V, P] {
	return &SafeLeftistHeap[V, P]{
		heap: NewLeftistHeap(data, cmp, usePool),
	}
}

// Push inserts a new value with the given priority into the heap.
// It returns the unique ID of the inserted node.
// This method acquires a write lock.
func (s *SafeLeftistHeap[V, P]) Push(value V, priority P) string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.Push(value, priority)
}

// Pop removes and returns the minimum element from the heap.
// It acquires a write lock.
func (s *SafeLeftistHeap[V, P]) Pop() (V, P, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.Pop()
}

// PopValue removes and returns just the value at the root.
// It acquires a write lock.
func (s *SafeLeftistHeap[V, P]) PopValue() (V, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.PopValue()
}

// PopPriority removes and returns just the priority at the root.
// It acquires a write lock.
func (s *SafeLeftistHeap[V, P]) PopPriority() (P, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.PopPriority()
}

// Peek returns the minimum element without removing it.
// It acquires a read lock.
func (s *SafeLeftistHeap[V, P]) Peek() (V, P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.Peek()
}

// PeekValue returns the value at the root without removing it.
// It acquires a read lock.
func (s *SafeLeftistHeap[V, P]) PeekValue() (V, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.PeekValue()
}

// PeekPriority returns the priority at the root without removing it.
// It acquires a read lock.
func (s *SafeLeftistHeap[V, P]) PeekPriority() (P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.PeekPriority()
}

// UpdateValue changes the value of the node with the given ID.
// It acquires a write lock.
func (s *SafeLeftistHeap[V, P]) UpdateValue(id string, value V) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.UpdateValue(id, value)
}

// UpdatePriority changes the priority of the node with the given ID and restructures the heap.
// It acquires a write lock.
func (s *SafeLeftistHeap[V, P]) UpdatePriority(id string, priority P) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.UpdatePriority(id, priority)
}

// Get returns the element associated with the given ID.
// It acquires a read lock.
func (s *SafeLeftistHeap[V, P]) Get(id string) (V, P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.Get(id)
}

// GetValue returns the value associated with the given ID.
// It acquires a read lock.
func (s *SafeLeftistHeap[V, P]) GetValue(id string) (V, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.GetValue(id)
}

// GetPriority returns the priority associated with the given ID.
// It acquires a read lock.
func (s *SafeLeftistHeap[V, P]) GetPriority(id string) (P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.GetPriority(id)
}

// Length returns the current number of elements in the heap.
// It acquires a read lock.
func (s *SafeLeftistHeap[V, P]) Length() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.Length()
}

// IsEmpty returns true if the heap contains no elements.
// It acquires a read lock.
func (s *SafeLeftistHeap[V, P]) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.IsEmpty()
}

// Clear removes all elements from the heap and resets its state.
// It acquires a write lock.
func (s *SafeLeftistHeap[V, P]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.heap.Clear()
}

// Clone creates a deep copy of the heap structure and nodes.
// The returned heap is also thread-safe, but shares no data with the original.
// It acquires a read lock.
func (s *SafeLeftistHeap[V, P]) Clone() *SafeLeftistHeap[V, P] {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return &SafeLeftistHeap[V, P]{
		heap: s.heap.Clone(),
	}
}
