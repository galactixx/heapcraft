package heapcraft

import (
	"sync"
)

// SafeSkewHeap is a thread-safe wrapper around SkewHeap.
// All operations are protected by a sync.RWMutex, making it safe for concurrent use.
type SafeSkewHeap[V any, P any] struct {
	heap *SkewHeap[V, P]
	lock sync.RWMutex
}

// NewSafeSkewHeap constructs a new thread-safe skew heap from the given data and comparison function.
// The resulting heap is safe for concurrent use.
func NewSafeSkewHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool) *SafeSkewHeap[V, P] {
	return &SafeSkewHeap[V, P]{
		heap: NewSkewHeap(data, cmp),
	}
}

// Push inserts a new value with the given priority into the heap.
// It returns the unique ID of the inserted node.
// This method acquires a write lock.
func (s *SafeSkewHeap[V, P]) Push(value V, priority P) string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.Push(value, priority)
}

// Pop removes and returns the minimum element from the heap.
// It acquires a write lock.
func (s *SafeSkewHeap[V, P]) Pop() (V, P, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.Pop()
}

// PopValue removes and returns just the value at the root.
// It acquires a write lock.
func (s *SafeSkewHeap[V, P]) PopValue() (V, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.PopValue()
}

// PopPriority removes and returns just the priority at the root.
// It acquires a write lock.
func (s *SafeSkewHeap[V, P]) PopPriority() (P, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.PopPriority()
}

// Peek returns the minimum element without removing it.
// It acquires a read lock.
func (s *SafeSkewHeap[V, P]) Peek() (V, P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.Peek()
}

// PeekValue returns the value at the root without removing it.
// It acquires a read lock.
func (s *SafeSkewHeap[V, P]) PeekValue() (V, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.PeekValue()
}

// PeekPriority returns the priority at the root without removing it.
// It acquires a read lock.
func (s *SafeSkewHeap[V, P]) PeekPriority() (P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.PeekPriority()
}

// UpdateValue changes the value of the node with the given ID.
// It acquires a write lock.
func (s *SafeSkewHeap[V, P]) UpdateValue(id string, value V) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.UpdateValue(id, value)
}

// UpdatePriority changes the priority of the node with the given ID and restructures the heap.
// It acquires a write lock.
func (s *SafeSkewHeap[V, P]) UpdatePriority(id string, priority P) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.UpdatePriority(id, priority)
}

// Get returns the element associated with the given ID.
// It acquires a read lock.
func (s *SafeSkewHeap[V, P]) Get(id string) (V, P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.Get(id)
}

// GetValue returns the value associated with the given ID.
// It acquires a read lock.
func (s *SafeSkewHeap[V, P]) GetValue(id string) (V, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.GetValue(id)
}

// GetPriority returns the priority associated with the given ID.
// It acquires a read lock.
func (s *SafeSkewHeap[V, P]) GetPriority(id string) (P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.GetPriority(id)
}

// Length returns the current number of elements in the heap.
// It acquires a read lock.
func (s *SafeSkewHeap[V, P]) Length() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.Length()
}

// IsEmpty returns true if the heap contains no elements.
// It acquires a read lock.
func (s *SafeSkewHeap[V, P]) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.IsEmpty()
}

// Clear removes all elements from the heap and resets its state.
// It acquires a write lock.
func (s *SafeSkewHeap[V, P]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.heap.Clear()
}

// Clone creates a deep copy of the heap structure and nodes.
// The returned heap is also thread-safe, but shares no data with the original.
// It acquires a read lock.
func (s *SafeSkewHeap[V, P]) Clone() *SafeSkewHeap[V, P] {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return &SafeSkewHeap[V, P]{
		heap: s.heap.Clone(),
	}
}

// SafeSimpleSkewHeap is a thread-safe wrapper around SimpleSkewHeap.
// All operations are protected by a sync.RWMutex, making it safe for concurrent use.
type SafeSimpleSkewHeap[V any, P any] struct {
	heap *SimpleSkewHeap[V, P]
	lock sync.RWMutex
}

// NewSafeSimpleSkewHeap constructs a new thread-safe simple skew heap from the given data and comparison function.
// The resulting heap is safe for concurrent use.
func NewSafeSimpleSkewHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool) *SafeSimpleSkewHeap[V, P] {
	return &SafeSimpleSkewHeap[V, P]{
		heap: NewSimpleSkewHeap(data, cmp),
	}
}

// Push adds a new element to the simple heap by creating a singleton node
// and merging it with the existing tree.
// It acquires a write lock.
func (s *SafeSimpleSkewHeap[V, P]) Push(value V, priority P) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.heap.Push(value, priority)
}

// Pop removes and returns the minimum element from the simple heap.
// The heap property is restored through merging the root's children.
// It acquires a write lock.
func (s *SafeSimpleSkewHeap[V, P]) Pop() (V, P, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.Pop()
}

// PopValue removes and returns just the value at the root.
// The heap property is restored through merging the root's children.
// It acquires a write lock.
func (s *SafeSimpleSkewHeap[V, P]) PopValue() (V, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.PopValue()
}

// PopPriority removes and returns just the priority at the root.
// The heap property is restored through merging the root's children.
// It acquires a write lock.
func (s *SafeSimpleSkewHeap[V, P]) PopPriority() (P, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.PopPriority()
}

// Peek returns the minimum element without removing it.
// It acquires a read lock.
func (s *SafeSimpleSkewHeap[V, P]) Peek() (V, P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.Peek()
}

// PeekValue returns the value at the root without removing it.
// It acquires a read lock.
func (s *SafeSimpleSkewHeap[V, P]) PeekValue() (V, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.PeekValue()
}

// PeekPriority returns the priority at the root without removing it.
// It acquires a read lock.
func (s *SafeSimpleSkewHeap[V, P]) PeekPriority() (P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.PeekPriority()
}

// Length returns the current number of elements in the simple heap.
// It acquires a read lock.
func (s *SafeSimpleSkewHeap[V, P]) Length() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.Length()
}

// IsEmpty returns true if the simple heap contains no elements.
// It acquires a read lock.
func (s *SafeSimpleSkewHeap[V, P]) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.IsEmpty()
}

// Clear removes all elements from the simple heap.
// The heap is ready for new insertions after clearing.
// It acquires a write lock.
func (s *SafeSimpleSkewHeap[V, P]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.heap.Clear()
}

// Clone creates a deep copy of the heap structure and nodes.
// The returned heap is also thread-safe, but shares no data with the original.
// It acquires a read lock.
func (s *SafeSimpleSkewHeap[V, P]) Clone() *SafeSimpleSkewHeap[V, P] {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return &SafeSimpleSkewHeap[V, P]{
		heap: s.heap.Clone(),
	}
}
