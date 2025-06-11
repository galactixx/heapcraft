package heapcraft

import (
	"errors"
	"fmt"
	"sync"
)

// NewBinaryHeap creates a new binary heap (d=2) from the given data slice and comparison function.
// It is a convenience wrapper around NewDaryHeap with d=2.
func NewBinaryHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool) *DaryHeap[V, P] {
	return NewDaryHeap(2, data, cmp)
}

// NewBinaryHeapCopy creates a new binary heap (d=2) from a copy of the given data slice.
// Unlike NewBinaryHeap, this function creates a new slice and copies the data before
// heapifying it, leaving the original data unchanged. It is a convenience wrapper
// around NewDaryHeapCopy with d=2.
func NewBinaryHeapCopy[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool) *DaryHeap[V, P] {
	return NewDaryHeapCopy(2, data, cmp)
}

// NewDaryHeapCopy duplicates the provided slice of HeapNode and builds a
// new d-ary heap from it.
func NewDaryHeapCopy[V any, P any](d int, data []HeapNode[V, P], cmp func(a, b P) bool) *DaryHeap[V, P] {
	heap := make([]HeapNode[V, P], len(data))
	copy(heap, data)
	return NewDaryHeap(d, heap, cmp)
}

// NewDaryHeap transforms the given slice of HeapNode into a valid d-ary heap in-place
// and returns it. Uses siftDown starting from the last parent toward the root.
func NewDaryHeap[V any, P any](d int, data []HeapNode[V, P], cmp func(a, b P) bool) *DaryHeap[V, P] {
	callbacks := &Callbacks{callbacks: make(map[int]Callback, 0), curId: 1}
	if len(data) == 0 {
		emptyHeap := make([]HeapNode[V, P], 0)
		return &DaryHeap[V, P]{data: emptyHeap, cmp: cmp, onSwap: callbacks, d: d}
	}
	h := DaryHeap[V, P]{data: data, cmp: cmp, onSwap: callbacks, d: d}

	// Start sifting down from the last parent node toward the root.
	start := (h.Length() - 2) / d
	for i := start; i >= 0; i-- {
		h.siftDown(i)
	}
	return &h
}

// DaryHeap represents a generic min-heap or max-heap (depending on cmp), with
// support for swap callbacks.
//   - data: slice of HeapNode (value, priority).
//   - cmp: comparison function on priority to enforce heap order.
//   - onSwap: callbacks invoked whenever two elements swap.
//   - d: the arity of the heap (e.g., 2 for binary, 3 for ternary).
type DaryHeap[V any, P any] struct {
	data   []HeapNode[V, P]
	cmp    func(a, b P) bool
	onSwap *Callbacks
	d      int
	lock   sync.RWMutex
}

// Deregister removes the callback with the specified ID, returning an error
// if it does not exist.
func (h *DaryHeap[V, P]) Deregister(id int) error { return h.onSwap.deregister(id) }

// Register adds a callback function to be called on each swap and returns its
// callback ID.
func (h *DaryHeap[V, P]) Register(fn func(x, y int)) Callback { return h.onSwap.register(fn) }

// swap exchanges the elements at indices i and j,
// and invokes all registered onSwap callbacks.
func (h *DaryHeap[V, P]) swap(i int, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
	h.onSwap.run(i, j)
}

// swapWithLast swaps the element at index i with the last element,
// removes the last element, sifts down the item now at index i to restore heap order,
// and returns the removed HeapNode.
func (h *DaryHeap[V, P]) swapWithLast(i int) HeapNode[V, P] {
	n := len(h.data)
	removed := h.data[i]
	h.swap(i, n-1)
	h.data = h.data[:n-1]
	h.siftDown(i)
	return removed
}

// Clear removes all elements from the heap by resetting its underlying slice
// to length zero.
func (h *DaryHeap[V, P]) Clear() {
	h.lock.Lock()
	h.data = h.data[:0]
	h.lock.Unlock()
}

// Length returns the current number of elements (HeapNode) in the heap.
func (h *DaryHeap[V, P]) Length() int {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return len(h.data)
}

// IsEmpty returns true if there are no elements in the heap.
func (h *DaryHeap[V, P]) IsEmpty() bool {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return len(h.data) == 0
}

// Peek returns the root HeapNode without removing it, or a zero value if
// the heap is empty.
func (h *DaryHeap[V, P]) Peek() (SimpleNode[V, P], error) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	if len(h.data) == 0 {
		var zero HeapNode[V, P]
		return zero, errors.New("the heap is empty and contains no elements")
	}
	return h.data[0], nil
}

// PopPush inserts a new element (HeapNode) into the heap, then removes and returns the
// root, in one operation.
func (h *DaryHeap[V, P]) PopPush(value V, priority P) SimpleNode[V, P] {
	h.lock.Lock()
	defer h.lock.Unlock()
	element := CreateHeapPair(value, priority)
	h.data = append(h.data, element)
	return h.swapWithLast(0)
}

// PushPop compares the new element's priority with the current root: if the new element
// belongs at the root (per cmp), it returns the new element directly; otherwise, it
// inserts the new element and removes the old root, returning that old root.
func (h *DaryHeap[V, P]) PushPop(value V, priority P) SimpleNode[V, P] {
	h.lock.Lock()
	defer h.lock.Unlock()
	element := CreateHeapPair(value, priority)
	if len(h.data) != 0 && h.cmp(element.priority, h.data[0].priority) {
		return element
	}
	h.data = append(h.data, element)
	return h.swapWithLast(0)
}

// Clone creates a shallow copy of the heap: the slice of HeapNode is copied,
// but the individual HeapNode elements themselves are not duplicated.
// Returns a new DaryHeap with the same comparison function and arity.
func (h *DaryHeap[V, P]) Clone() *DaryHeap[V, P] {
	h.lock.RLock()
	defer h.lock.RUnlock()
	newData := make([]HeapNode[V, P], h.Length())
	copy(newData, h.data)
	callbacks := &Callbacks{callbacks: make(map[int]Callback, 0), curId: 1}
	return &DaryHeap[V, P]{data: newData, cmp: h.cmp, onSwap: callbacks, d: h.d}
}

// siftUp moves the element at index i up the tree until the heap property is restored.
// The heap property is determined by the comparison function cmp, where a parent's priority
// should compare appropriately with its children's priorities.
func (h *DaryHeap[V, P]) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / h.d
		if !h.cmp(h.data[i].priority, h.data[parent].priority) {
			break
		}
		h.swap(i, parent)
		i = parent
	}
}

// siftDown moves the element at index i down the tree until all children satisfy the heap order.
// For each node, it finds the child with the most appropriate priority (per cmp) and swaps
// if necessary to maintain the heap property.
func (h *DaryHeap[V, P]) siftDown(i int) {
	cur := i
	n := len(h.data)
	for h.d*cur+1 < n {
		left := h.d*cur + 1
		right := min(left+h.d, n)

		swapIdx := left
		for k := left + 1; k < right; k++ {
			if h.cmp(h.data[k].priority, h.data[swapIdx].priority) {
				swapIdx = k
			}
		}

		if !h.cmp(h.data[swapIdx].priority, h.data[cur].priority) {
			break
		}
		h.swap(swapIdx, cur)
		cur = swapIdx
	}
}

// restoreHeapProperty restores the heap property after an element at index i has been updated.
// It decides whether to sift up or down based on the element's priority relative to its parent.
func (h *DaryHeap[V, P]) restoreHeap(i int, element HeapNode[V, P]) {
	h.data[i] = element
	if i > 0 && h.cmp(element.priority, h.data[(i-1)/h.d].priority) {
		h.siftUp(i)
	} else {
		h.siftDown(i)
	}
}

// Update replaces the element at index i with a new value and priority.
// It then restores the heap property by either sifting up (if the new priority is more
// appropriate than its parent) or sifting down (if the new priority is less appropriate
// than its children).
// Returns the updated element and an error if the index is out of bounds.
func (h *DaryHeap[V, P]) Update(i int, value V, priority P) error {
	h.lock.Lock()
	defer h.lock.Unlock()
	if i < 0 || i >= len(h.data) {
		return fmt.Errorf("index %d is out of bounds", i)
	}
	element := CreateHeapPair(value, priority)
	h.restoreHeap(i, element)
	return nil
}

// Remove deletes the element at index i from the heap and returns it.
// The heap property is restored by replacing the removed element with the last element
// and sifting it down to its appropriate position.
// Returns the removed element and an error if the index is out of bounds.
func (h *DaryHeap[V, P]) Remove(i int) (SimpleNode[V, P], error) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if i < 0 || i >= len(h.data) {
		var zero HeapNode[V, P]
		return zero, fmt.Errorf("index %d is out of bounds", i)
	}
	removed := h.swapWithLast(i)
	return removed, nil
}

// Pop removes and returns the root element of the heap (minimum or maximum per cmp).
// If the heap is empty, returns a zero value with error. The heap property is restored by replacing
// the root with the last element and sifting it down.
func (h *DaryHeap[V, P]) Pop() (SimpleNode[V, P], error) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if len(h.data) == 0 {
		var zero HeapNode[V, P]
		return zero, errors.New("the heap is empty and contains no elements")
	}
	removed := h.swapWithLast(0)
	return removed, nil
}

// Push inserts a new element with the given value and priority into the heap.
// The element is added at the end and then sifted up to maintain the heap property.
// Returns the newly created SimpleNode.
func (h *DaryHeap[V, P]) Push(value V, priority P) {
	h.lock.Lock()
	defer h.lock.Unlock()
	element := CreateHeapPair(value, priority)
	h.data = append(h.data, element)
	i := len(h.data) - 1
	h.siftUp(i)
}

// nDary builds a heap of size n from the data slice.
// It uses Push for the first n elements and PushPop for the remainder to maintain
// a heap of exactly size n. This is used as the underlying implementation for
// both NLargestDary and NSmallestDary.
func nDary[V any, P any](n int, d int, data []HeapNode[V, P], cmp func(a, b P) bool) *DaryHeap[V, P] {
	callbacks := &Callbacks{callbacks: make(map[int]Callback, 0), curId: 1}
	heap := DaryHeap[V, P]{data: make([]HeapNode[V, P], 0, n), cmp: cmp, onSwap: callbacks, d: d}
	i := 0
	m := len(data)
	minNum := min(n, m)

	// Build initial heap with the first min(n, m) elements.
	for ; i < minNum; i++ {
		element := data[i]
		heap.Push(element.value, element.priority)
	}

	// For remaining elements, use PushPop to keep the heap size at n.
	for ; i < m; i++ {
		element := data[i]
		heap.PushPop(element.value, element.priority)
	}
	return &heap
}

// NLargestDary returns a min-heap of size n containing the n largest
// elements from data. lt should compare priorities by returning true if a < b.
func NLargestDary[V any, P any](n int, d int, data []HeapNode[V, P], lt func(a, b P) bool) *DaryHeap[V, P] {
	return nDary(n, d, data, lt)
}

// NLargestBinary returns a min-heap of size n containing the n largest
// elements from data, using a binary heap (d=2). lt should compare priorities
// by returning true if a < b. This is a convenience wrapper around NLargestDary.
func NLargestBinary[V any, P any](n int, data []HeapNode[V, P], lt func(a, b P) bool) *DaryHeap[V, P] {
	return NLargestDary(n, 2, data, lt)
}

// NSmallestDary returns a max-heap of size n containing the n smallest
// elements from data. gt should compare priorities by returning true if a > b.
func NSmallestDary[V any, P any](n int, d int, data []HeapNode[V, P], gt func(a, b P) bool) *DaryHeap[V, P] {
	return nDary(n, d, data, gt)
}

// NSmallestBinary returns a max-heap of size n containing the n smallest
// elements from data, using a binary heap (d=2). gt should compare priorities
// by returning true if a > b. This is a convenience wrapper around NSmallestDary.
func NSmallestBinary[V any, P any](n int, data []HeapNode[V, P], gt func(a, b P) bool) *DaryHeap[V, P] {
	return NSmallestDary(n, 2, data, gt)
}
