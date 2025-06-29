package heapcraft

import (
	"sync"
)

// SyncPairingHeap provides a thread-safe wrapper around PairingHeap.
// It uses a read-write mutex to allow concurrent reads and exclusive writes.
type SyncPairingHeap[V any, P any] struct {
	heap *PairingHeap[V, P]
	mu   sync.RWMutex
}

// UpdateValue updates the value of a node with the given ID.
// Returns an error if the ID does not exist in the heap.
// The heap structure remains unchanged as this operation only modifies the value.
func (s *SyncPairingHeap[V, P]) UpdateValue(id string, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.heap.UpdateValue(id, value)
}

// UpdatePriority updates the priority of a node with the given ID.
// Returns an error if the ID does not exist in the heap.
// The node is removed from its current position and reinserted into the heap
// to maintain the heap property. This operation may change the heap structure.
func (s *SyncPairingHeap[V, P]) UpdatePriority(id string, priority P) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.heap.UpdatePriority(id, priority)
}

// Clone creates a deep copy of the heap structure and nodes. If values or
// priorities are reference types, those reference values are shared between the
// original and cloned heaps.
func (s *SyncPairingHeap[V, P]) Clone() *SyncPairingHeap[V, P] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &SyncPairingHeap[V, P]{heap: s.heap.Clone()}
}

// Clear removes all elements from the heap.
// Resets the root to nil, size to zero, and initializes a new empty element map.
// The next node ID is reset to 1.
func (s *SyncPairingHeap[V, P]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.heap.Clear()
}

// Length returns the current number of elements in the heap.
func (s *SyncPairingHeap[V, P]) Length() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.Length()
}

// IsEmpty returns true if the heap contains no elements.
func (s *SyncPairingHeap[V, P]) IsEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.IsEmpty()
}

// Peek returns a HeapNode containing the value and priority
// of the root node without removing it. Returns nil and an error if the heap is empty.
func (s *SyncPairingHeap[V, P]) Peek() (V, P, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.Peek()
}

// PeekValue returns the value at the root without removing it.
// Returns zero value and an error if the heap is empty.
func (s *SyncPairingHeap[V, P]) PeekValue() (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.PeekValue()
}

// PeekPriority returns the priority at the root without removing it.
// Returns zero value and an error if the heap is empty.
func (s *SyncPairingHeap[V, P]) PeekPriority() (P, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.PeekPriority()
}

// Get retrieves a HeapNode for the node with the given ID.
// Returns an error if the ID does not exist in the heap.
func (s *SyncPairingHeap[V, P]) Get(id string) (V, P, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.Get(id)
}

// GetValue retrieves the value of the node with the given ID.
// Returns zero value and an error if the ID does not exist in the heap.
func (s *SyncPairingHeap[V, P]) GetValue(id string) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.GetValue(id)
}

// GetPriority retrieves the priority of the node with the given ID.
// Returns zero value and an error if the ID does not exist in the heap.
func (s *SyncPairingHeap[V, P]) GetPriority(id string) (P, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.GetPriority(id)
}

// Pop removes and returns a HeapNode containing the value and priority
// of the root node. The root's children are merged to form the new heap.
// Returns nil and an error if the heap is empty.
func (s *SyncPairingHeap[V, P]) Pop() (V, P, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.heap.Pop()
}

// PopValue removes and returns just the value at the root.
// The root's children are merged to form the new heap.
// Returns zero value and an error if the heap is empty.
func (s *SyncPairingHeap[V, P]) PopValue() (V, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.heap.PopValue()
}

// PopPriority removes and returns just the priority at the root.
// The root's children are merged to form the new heap.
// Returns zero value and an error if the heap is empty.
func (s *SyncPairingHeap[V, P]) PopPriority() (P, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.heap.PopPriority()
}

// Push adds a new element with the given value and priority to the heap.
// A new node is created with a unique ID and melded with the existing root.
// The new node becomes the root if its priority is higher than the current root's.
// Returns the ID of the inserted node.
func (s *SyncPairingHeap[V, P]) Push(value V, priority P) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.heap.Push(value, priority)
}

// SyncSimplePairingHeap provides a thread-safe wrapper around SimplePairingHeap.
// It uses a read-write mutex to allow concurrent reads and exclusive writes.
type SyncSimplePairingHeap[V any, P any] struct {
	heap *SimplePairingHeap[V, P]
	mu   sync.RWMutex
}

// Clone creates a deep copy of the simple heap structure and nodes. If values or
// priorities are reference types, those reference values are shared between the
// original and cloned heaps.
func (s *SyncSimplePairingHeap[V, P]) Clone() *SyncSimplePairingHeap[V, P] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &SyncSimplePairingHeap[V, P]{
		heap: s.heap.Clone(),
	}
}

// Clear removes all elements from the simple heap.
// The heap is ready for new insertions after clearing.
func (s *SyncSimplePairingHeap[V, P]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.heap.Clear()
}

// Length returns the current number of elements in the simple heap.
func (s *SyncSimplePairingHeap[V, P]) Length() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.Length()
}

// IsEmpty returns true if the simple heap contains no elements.
func (s *SyncSimplePairingHeap[V, P]) IsEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.IsEmpty()
}

// Peek returns a HeapNode containing the value and priority
// of the root node without removing it. Returns nil and an error if the heap is empty.
func (s *SyncSimplePairingHeap[V, P]) Peek() (V, P, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.Peek()
}

// PeekValue returns the value at the root without removing it.
// Returns zero value and an error if the heap is empty.
func (s *SyncSimplePairingHeap[V, P]) PeekValue() (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.PeekValue()
}

// PeekPriority returns the priority at the root without removing it.
// Returns zero value and an error if the heap is empty.
func (s *SyncSimplePairingHeap[V, P]) PeekPriority() (P, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.PeekPriority()
}

// Pop removes and returns a HeapNode containing the value and priority
// of the root node. The root's children are merged to form the new heap.
// Returns nil and an error if the heap is empty.
func (s *SyncSimplePairingHeap[V, P]) Pop() (V, P, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.heap.Pop()
}

// PopValue removes and returns just the value at the root.
// The root's children are merged to form the new heap.
// Returns zero value and an error if the heap is empty.
func (s *SyncSimplePairingHeap[V, P]) PopValue() (V, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.heap.PopValue()
}

// PopPriority removes and returns just the priority at the root.
// The root's children are merged to form the new heap.
// Returns zero value and an error if the heap is empty.
func (s *SyncSimplePairingHeap[V, P]) PopPriority() (P, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.heap.PopPriority()
}

// Push adds a new element with its priority by creating a single-node heap
// and melding it with the existing root. The new node becomes the root if
// its priority is higher than the current root's priority.
func (s *SyncSimplePairingHeap[V, P]) Push(value V, priority P) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.heap.Push(value, priority)
}
