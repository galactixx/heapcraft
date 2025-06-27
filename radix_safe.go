package heapcraft

import (
	"sync"
	"unsafe"

	"golang.org/x/exp/constraints"
)

// getHeapAddr returns the address of the heap.
func getHeapAddr[V any, P constraints.Unsigned](h *SyncRadixHeap[V, P]) uintptr {
	return uintptr(unsafe.Pointer(h))
}

// SyncRadixHeap provides a thread-safe wrapper around RadixHeap.
// It uses a read-write mutex to allow concurrent reads and exclusive writes.
type SyncRadixHeap[V any, P constraints.Unsigned] struct {
	heap *RadixHeap[V, P]
	mu   sync.RWMutex
}

// NewSyncRadixHeap creates a new thread-safe RadixHeap from a given slice of HeapNode[V,P].
func NewSyncRadixHeap[V any, P constraints.Unsigned](data []HeapNode[V, P], usePool bool) *SyncRadixHeap[V, P] {
	return &SyncRadixHeap[V, P]{heap: NewRadixHeap(data, usePool)}
}

// Clone creates a deep copy of the heap structure. The new heap preserves the
// original size and last value. If values or priorities are reference types, those
// reference values are shared between the original and cloned heaps.
func (s *SyncRadixHeap[V, P]) Clone() *SyncRadixHeap[V, P] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &SyncRadixHeap[V, P]{
		heap: s.heap.Clone(),
	}
}

// Push adds a new value and priority pair into the heap.
// Returns an error if the priority is less than the last extracted priority, as this would violate
// the monotonic property. Otherwise, puts the item into the appropriate bucket
// and increments the size.
func (s *SyncRadixHeap[V, P]) Push(value V, priority P) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.heap.Push(value, priority)
}

// Pop extracts and returns the HeapNode with the minimum priority.
// Returns nil and an error if the heap is empty.
func (s *SyncRadixHeap[V, P]) Pop() (V, P, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.heap.Pop()
}

// Peek returns a HeapNode with the minimum priority without removing it.
// Returns nil and an error if the heap is empty.
func (s *SyncRadixHeap[V, P]) Peek() (V, P, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.Peek()
}

// PopValue removes and returns just the value of the root element.
// Returns zero value and an error if the heap is empty.
func (s *SyncRadixHeap[V, P]) PopValue() (V, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.heap.PopValue()
}

// PopPriority removes and returns just the priority of the root element.
// Returns zero value and an error if the heap is empty.
func (s *SyncRadixHeap[V, P]) PopPriority() (P, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.heap.PopPriority()
}

// PeekValue returns just the value of the root element without removing it.
// Returns zero value and an error if the heap is empty.
func (s *SyncRadixHeap[V, P]) PeekValue() (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.PeekValue()
}

// PeekPriority returns just the priority of the root element without removing it.
// Returns zero value and an error if the heap is empty.
func (s *SyncRadixHeap[V, P]) PeekPriority() (P, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.PeekPriority()
}

// Clear reinitializes the heap by creating fresh buckets, resetting size to zero,
// and setting 'last' back to its zero value.
func (s *SyncRadixHeap[V, P]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.heap.Clear()
}

// Rebalance fills bucket 0 if it is empty.
// Returns an error if the heap is empty, or if bucket 0 already contains elements
// (no action was needed).
func (s *SyncRadixHeap[V, P]) Rebalance() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.heap.Rebalance()
}

// Length returns the number of items currently stored in the heap.
func (s *SyncRadixHeap[V, P]) Length() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.Length()
}

// IsEmpty returns true if the heap contains no items.
func (s *SyncRadixHeap[V, P]) IsEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.heap.IsEmpty()
}

// Merge integrates another SafeRadixHeap into this one.
// It selects the heap with the smaller 'last' as the new baseline, adopts its
// buckets and 'last', then reinserts all items from the other heap to preserve
// the monotonic property.
func (s *SyncRadixHeap[V, P]) Merge(other *SyncRadixHeap[V, P]) {
	if getHeapAddr(s) > getHeapAddr(other) {
		s.mu.Lock()
		defer s.mu.Unlock()
		other.mu.Lock()
		defer other.mu.Unlock()
	} else {
		other.mu.Lock()
		defer other.mu.Unlock()
		s.mu.Lock()
		defer s.mu.Unlock()
	}
	s.heap.Merge(other.heap)
}
