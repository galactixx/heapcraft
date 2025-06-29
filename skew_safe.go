package heapcraft

import (
	"sync"
)

// SyncSkewHeap is a thread-safe wrapper around SkewHeap.
// All operations are protected by a sync.RWMutex, making it safe for concurrent use.
type SyncFullSkewHeap[V any, P any] struct {
	heap *FullSkewHeap[V, P]
	lock sync.RWMutex
}

// Push inserts a new value with the given priority into the heap.
// It returns the unique ID of the inserted node.
// This method acquires a write lock.
func (s *SyncFullSkewHeap[V, P]) Push(value V, priority P) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.Push(value, priority)
}

// Pop removes and returns the minimum element from the heap.
// It acquires a write lock.
func (s *SyncFullSkewHeap[V, P]) Pop() (V, P, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.Pop()
}

// PopValue removes and returns just the value at the root.
// It acquires a write lock.
func (s *SyncFullSkewHeap[V, P]) PopValue() (V, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.PopValue()
}

// PopPriority removes and returns just the priority at the root.
// It acquires a write lock.
func (s *SyncFullSkewHeap[V, P]) PopPriority() (P, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.PopPriority()
}

// Peek returns the minimum element without removing it.
// It acquires a read lock.
func (s *SyncFullSkewHeap[V, P]) Peek() (V, P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.Peek()
}

// PeekValue returns the value at the root without removing it.
// It acquires a read lock.
func (s *SyncFullSkewHeap[V, P]) PeekValue() (V, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.PeekValue()
}

// PeekPriority returns the priority at the root without removing it.
// It acquires a read lock.
func (s *SyncFullSkewHeap[V, P]) PeekPriority() (P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.PeekPriority()
}

// UpdateValue changes the value of the node with the given ID.
// It acquires a write lock.
func (s *SyncFullSkewHeap[V, P]) UpdateValue(id string, value V) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.UpdateValue(id, value)
}

// UpdatePriority changes the priority of the node with the given ID and restructures the heap.
// It acquires a write lock.
func (s *SyncFullSkewHeap[V, P]) UpdatePriority(id string, priority P) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.UpdatePriority(id, priority)
}

// Get returns the element associated with the given ID.
// It acquires a read lock.
func (s *SyncFullSkewHeap[V, P]) Get(id string) (V, P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.Get(id)
}

// GetValue returns the value associated with the given ID.
// It acquires a read lock.
func (s *SyncFullSkewHeap[V, P]) GetValue(id string) (V, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.GetValue(id)
}

// GetPriority returns the priority associated with the given ID.
// It acquires a read lock.
func (s *SyncFullSkewHeap[V, P]) GetPriority(id string) (P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.GetPriority(id)
}

// Length returns the current number of elements in the heap.
// It acquires a read lock.
func (s *SyncFullSkewHeap[V, P]) Length() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.Length()
}

// IsEmpty returns true if the heap contains no elements.
// It acquires a read lock.
func (s *SyncFullSkewHeap[V, P]) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.IsEmpty()
}

// Clear removes all elements from the heap and resets its state.
// It acquires a write lock.
func (s *SyncFullSkewHeap[V, P]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.heap.Clear()
}

// Clone creates a deep copy of the heap structure and nodes.
// The returned heap is also thread-safe, but shares no data with the original.
// It acquires a read lock.
func (s *SyncFullSkewHeap[V, P]) Clone() *SyncFullSkewHeap[V, P] {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return &SyncFullSkewHeap[V, P]{
		heap: s.heap.Clone(),
	}
}

// SyncSkewHeap is a thread-safe wrapper around SkewHeap.
// All operations are protected by a sync.RWMutex, making it safe for concurrent use.
type SyncSkewHeap[V any, P any] struct {
	heap *SkewHeap[V, P]
	lock sync.RWMutex
}

// Push adds a new element to the simple heap by creating a singleton node
// and merging it with the existing tree.
// It acquires a write lock.
func (s *SyncSkewHeap[V, P]) Push(value V, priority P) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.heap.Push(value, priority)
}

// Pop removes and returns the minimum element from the simple heap.
// The heap property is restored through merging the root's children.
// It acquires a write lock.
func (s *SyncSkewHeap[V, P]) Pop() (V, P, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.Pop()
}

// PopValue removes and returns just the value at the root.
// The heap property is restored through merging the root's children.
// It acquires a write lock.
func (s *SyncSkewHeap[V, P]) PopValue() (V, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.PopValue()
}

// PopPriority removes and returns just the priority at the root.
// The heap property is restored through merging the root's children.
// It acquires a write lock.
func (s *SyncSkewHeap[V, P]) PopPriority() (P, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.heap.PopPriority()
}

// Peek returns the minimum element without removing it.
// It acquires a read lock.
func (s *SyncSkewHeap[V, P]) Peek() (V, P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.Peek()
}

// PeekValue returns the value at the root without removing it.
// It acquires a read lock.
func (s *SyncSkewHeap[V, P]) PeekValue() (V, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.PeekValue()
}

// PeekPriority returns the priority at the root without removing it.
// It acquires a read lock.
func (s *SyncSkewHeap[V, P]) PeekPriority() (P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.PeekPriority()
}

// Length returns the current number of elements in the simple heap.
// It acquires a read lock.
func (s *SyncSkewHeap[V, P]) Length() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.Length()
}

// IsEmpty returns true if the simple heap contains no elements.
// It acquires a read lock.
func (s *SyncSkewHeap[V, P]) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.heap.IsEmpty()
}

// Clear removes all elements from the simple heap.
// The heap is ready for new insertions after clearing.
// It acquires a write lock.
func (s *SyncSkewHeap[V, P]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.heap.Clear()
}

// Clone creates a deep copy of the heap structure and nodes.
// The returned heap is also thread-safe, but shares no data with the original.
// It acquires a read lock.
func (s *SyncSkewHeap[V, P]) Clone() *SyncSkewHeap[V, P] {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return &SyncSkewHeap[V, P]{
		heap: s.heap.Clone(),
	}
}
