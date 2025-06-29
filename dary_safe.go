package heapcraft

import (
	"sync"
)

// SyncDaryHeap represents a thread-safe wrapper around DaryHeap.
// It provides the same interface as DaryHeap but with mutex-protected operations.
type SyncDaryHeap[V any, P any] struct {
	heap *DaryHeap[V, P]
	lock sync.RWMutex
}

// Deregister removes the callback with the specified ID from the heap's swap
// callbacks. Returns an error if no callback exists with the given ID.
func (h *SyncDaryHeap[V, P]) Deregister(id string) error {
	return h.heap.Deregister(id)
}

// Register adds a callback function to be called whenever elements in the heap
// swap positions. Returns a callback that can be used to deregister the
// function later.
func (h *SyncDaryHeap[V, P]) Register(fn func(x, y int)) callback {
	return h.heap.Register(fn)
}

// Clear removes all elements from the heap by resetting its underlying slice to
// length zero.
func (h *SyncDaryHeap[V, P]) Clear() {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.heap.Clear()
}

// Length returns the current number of elements in the heap.
func (h *SyncDaryHeap[V, P]) Length() int {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.heap.Length()
}

// IsEmpty returns true if the heap contains no elements.
func (h *SyncDaryHeap[V, P]) IsEmpty() bool {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.heap.IsEmpty()
}

// Pop removes and returns the root element of the heap (minimum or maximum per
// cmp). If the heap is empty, returns a zero value and priority with an error.
func (h *SyncDaryHeap[V, P]) Pop() (V, P, error) {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.heap.Pop()
}

// Peek returns the root HeapNode without removing it.
// If the heap is empty, returns a zero value and priority with an error.
func (h *SyncDaryHeap[V, P]) Peek() (V, P, error) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.heap.Peek()
}

// PopValue removes and returns just the value of the root element.
// If the heap is empty, returns a zero value with an error.
func (h *SyncDaryHeap[V, P]) PopValue() (V, error) {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.heap.PopValue()
}

// PopPriority removes and returns just the priority of the root element.
// If the heap is empty, returns a zero value with an error.
func (h *SyncDaryHeap[V, P]) PopPriority() (P, error) {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.heap.PopPriority()
}

// PeekValue returns just the value of the root element without removing it.
// If the heap is empty, returns a zero value with an error.
func (h *SyncDaryHeap[V, P]) PeekValue() (V, error) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.heap.PeekValue()
}

// PeekPriority returns just the priority of the root element without removing it.
// If the heap is empty, returns a zero value with an error.
func (h *SyncDaryHeap[V, P]) PeekPriority() (P, error) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.heap.PeekPriority()
}

// Push inserts a new element with the given value and priority into the heap.
// The element is added at the end and then sifted up to maintain the heap property.
func (h *SyncDaryHeap[V, P]) Push(value V, priority P) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.heap.Push(value, priority)
}

// Update replaces the element at index i with a new value and priority.
// It then restores the heap property by either sifting up (if the new priority
// is more appropriate than its parent) or sifting down (if the new priority is
// less appropriate than its children).
// Returns an error if the index is out of bounds.
func (h *SyncDaryHeap[V, P]) Update(i int, value V, priority P) error {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.heap.Update(i, value, priority)
}

// Remove deletes the element at index i from the heap and returns it.
// The heap property is restored by replacing the removed element with the last
// element and sifting it down to its appropriate position.
// Returns the removed element and an error if the index is out of bounds.
func (h *SyncDaryHeap[V, P]) Remove(i int) (V, P, error) {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.heap.Remove(i)
}

// PopPush atomically removes the root element and inserts a new element into the heap.
// Returns the removed root element.
func (h *SyncDaryHeap[V, P]) PopPush(value V, priority P) (V, P) {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.heap.PopPush(value, priority)
}

// PushPop atomically inserts a new element and removes the root element if the
// new element doesn't belong at the root. If the new element belongs at the
// root, it is returned directly.
// Returns either the new element or the old root element.
func (h *SyncDaryHeap[V, P]) PushPop(value V, priority P) (V, P) {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.heap.PushPop(value, priority)
}

// Clone creates a deep copy of the heap structure. The new heap preserves the
// original size. If values or priorities are reference types, those reference
// values are shared between the original and cloned heaps.
func (h *SyncDaryHeap[V, P]) Clone() *SyncDaryHeap[V, P] {
	h.lock.RLock()
	defer h.lock.RUnlock()
	clonedHeap := h.heap.Clone()
	clonedHeap.onSwap = &syncCallbacks{callbacks: clonedHeap.onSwap.(baseCallbacks)}
	return &SyncDaryHeap[V, P]{heap: clonedHeap}
}
