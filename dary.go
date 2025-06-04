package heapcraft

import (
	"fmt"
)

// NewDaryHeapCopy duplicates the provided slice of *HeapPair and builds a
// new d-ary heap from it.
func NewDaryHeapCopy[V any, P any](d int, data []*HeapPair[V, P], cmp func(a, b P) bool) DaryHeap[V, P] {
	heap := make([]*HeapPair[V, P], len(data))
	copy(heap, data)
	return NewDaryHeap(d, heap, cmp)
}

// NewDaryHeap transforms the given slice of *HeapPair into a valid d-ary heap in-place
// and returns it. Uses siftDown starting from the last parent toward the root.
func NewDaryHeap[V any, P any](d int, data []*HeapPair[V, P], cmp func(a, b P) bool) DaryHeap[V, P] {
	if len(data) == 0 {
		emptyHeap := make([]*HeapPair[V, P], 0)
		return DaryHeap[V, P]{data: emptyHeap, cmp: cmp, d: d}
	}
	h := DaryHeap[V, P]{data: data, cmp: cmp, d: d}
	// Start sifting down from the last parent node toward the root.
	start := (h.Length() - 2) / d
	for i := start; i >= 0; i-- {
		h.siftDown(i)
	}
	return h
}

// DaryHeap represents a generic min-heap or max-heap (depending on cmp), with
// support for swap callbacks.
//   - data: slice of pointers to HeapPair (value, priority).
//   - cmp: comparison function on priority to enforce heap order.
//   - onSwap: callbacks invoked whenever two elements swap.
//   - d: the arity of the heap (e.g., 2 for binary, 3 for ternary).
type DaryHeap[V any, P any] struct {
	data   []*HeapPair[V, P]
	cmp    func(a, b P) bool
	onSwap Callbacks
	d      int
}

// Register adds a callback function to be called on each swap and returns its
// callback ID.
func (h *DaryHeap[V, P]) Register(fn func(x int, y int)) Callback {
	newId := h.onSwap.curId + 1
	newCallback := Callback{ID: newId, Function: fn}
	if h.onSwap.callbacks == nil {
		h.onSwap.callbacks = make(map[int]Callback)
	}

	h.onSwap.callbacks[newId] = newCallback
	h.onSwap.curId = newId
	return newCallback
}

// Deregister removes the callback with the specified ID, returning an error
// if it does not exist.
func (h *DaryHeap[V, P]) Deregister(id int) error {
	if _, exists := h.onSwap.callbacks[id]; !exists {
		return fmt.Errorf("%d is not an ID of a callback", id)
	}
	delete(h.onSwap.callbacks, id)
	return nil
}

// swap exchanges the elements at indices i and j (both are pointers),
// and invokes all registered onSwap callbacks.
func (h *DaryHeap[V, P]) swap(i int, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
	h.onSwap.run(i, j)
}

// swapWithLast swaps the element at index i with the last element (both pointers),
// removes the last element, sifts down the item now at index i to restore heap order,
// and returns the removed *HeapPair.
func (h *DaryHeap[V, P]) swapWithLast(i int) *HeapPair[V, P] {
	removed := h.data[i]
	h.data[i] = h.data[h.Length()-1]
	h.data = h.data[:h.Length()-1]
	h.siftDown(i)
	return removed
}

// Clear removes all elements from the heap by resetting its underlying slice
// to length zero.
func (h *DaryHeap[V, P]) Clear() {
	h.data = h.data[:0]
}

// Length returns the current number of elements (pointers to HeapPair) in the heap.
func (h DaryHeap[V, P]) Length() int {
	return len(h.data)
}

// IsEmpty returns true if there are no elements in the heap.
func (h DaryHeap[V, P]) IsEmpty() bool {
	return h.Length() == 0
}

// Peek returns a pointer to the root *HeapPair without removing it, or nil if
// the heap is empty.
func (h DaryHeap[V, P]) Peek() *HeapPair[V, P] {
	if h.IsEmpty() {
		return nil
	}
	return h.data[0]
}

// PopPush inserts a new element (*HeapPair) into the heap, then removes and returns the
// root, in one operation.
func (h *DaryHeap[V, P]) PopPush(value V, priority P) *HeapPair[V, P] {
	element := CreateHeapPair(value, priority)
	h.data = append(h.data, element)
	return h.swapWithLast(0)
}

// PushPop compares the new element's priority with the current root: if the new element
// belongs at the root (per cmp), it returns the new element directly; otherwise, it
// inserts the new element and removes the old root, returning that old root.
func (h *DaryHeap[V, P]) PushPop(value V, priority P) *HeapPair[V, P] {
	element := CreateHeapPair(value, priority)
	if !h.IsEmpty() && h.cmp(element.Priority(), h.Peek().Priority()) {
		return element
	}
	h.data = append(h.data, element)
	return h.swapWithLast(0)
}

// Clone creates a shallow copy of the heap: the slice of *HeapPair is copied,
// but the individual *HeapPair elements themselves are not duplicated.
func (h DaryHeap[V, P]) Clone() DaryHeap[V, P] {
	newData := make([]*HeapPair[V, P], h.Length())
	copy(newData, h.data)
	return DaryHeap[V, P]{data: newData, cmp: h.cmp, d: h.d}
}

// siftUp moves the element at index i (pointer to HeapPair) up the tree until the heap
// property is restored (parent priority compares appropriately per cmp).
func (h *DaryHeap[V, P]) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / h.d
		if !h.cmp(h.data[i].Priority(), h.data[parent].Priority()) {
			break
		}
		h.swap(i, parent)
		i = parent
	}
}

// siftDown moves the element at index i down the tree until all children satisfy the
// heap order (per cmp on priorities).
func (h *DaryHeap[V, P]) siftDown(i int) {
	cur := i
	n := h.Length()
	for h.d*cur+1 < n {
		left := h.d*cur + 1
		right := min(left+h.d, n)

		swapIdx := left
		for k := left + 1; k < right; k++ {
			if h.cmp(h.data[k].Priority(), h.data[swapIdx].Priority()) {
				swapIdx = k
			}
		}

		if !h.cmp(h.data[swapIdx].Priority(), h.data[cur].Priority()) {
			break
		}
		h.swap(swapIdx, cur)
		cur = swapIdx
	}
}

// Update replaces the element at index i with a new value/priority (creates a new *HeapPair),
// and then restores heap order by sifting up or down as needed.
func (h *DaryHeap[V, P]) Update(i int, value V, priority P) (*HeapPair[V, P], error) {
	if i < 0 || i >= h.Length() {
		return nil, fmt.Errorf("index %d is out of bounds", i)
	}
	element := CreateHeapPair(value, priority)
	h.data[i] = element
	// Decide whether to sift up or down depending on the new priority.
	if i > 0 && h.cmp(element.Priority(), h.data[(i-1)/h.d].Priority()) {
		h.siftUp(i)
	} else {
		h.siftDown(i)
	}
	return element, nil
}

// Remove deletes the element at index i, returns it, and restores heap order by sifting
// down the replacement that was swapped in.
func (h *DaryHeap[V, P]) Remove(i int) (*HeapPair[V, P], error) {
	if i < 0 || i >= h.Length() {
		return nil, fmt.Errorf("index %d is out of bounds", i)
	}
	removed := h.swapWithLast(i)
	return removed, nil
}

// Pop removes and returns the root *HeapPair of the heap (minimum or maximum per cmp),
// or nil if the heap is empty.
func (h *DaryHeap[V, P]) Pop() *HeapPair[V, P] {
	if h.IsEmpty() {
		return nil
	}
	removed := h.swapWithLast(0)
	return removed
}

// Push inserts a new element (*HeapPair) at the end of the heap and sifts it up to
// maintain heap order.
func (h *DaryHeap[V, P]) Push(value V, priority P) *HeapPair[V, P] {
	element := CreateHeapPair(value, priority)
	h.data = append(h.data, element)
	i := h.Length() - 1
	h.siftUp(i)
	return element
}

// nDary builds a heap of size n from the data slice using Push for the first n elements,
// and PushPop for the remainder, returning a DaryHeap of size n.
func nDary[V any, P any](n int, d int, data []*HeapPair[V, P], cmp func(a, b P) bool) DaryHeap[V, P] {
	h := DaryHeap[V, P]{data: make([]*HeapPair[V, P], 0, n), cmp: cmp, d: d}
	i := 0
	m := len(data)
	minNum := min(n, m)

	// Build initial heap with the first min(n, m) elements.
	for ; i < minNum; i++ {
		element := data[i]
		h.Push(element.Value(), element.Priority())
	}

	// For remaining elements, use PushPop to keep the heap size at n.
	for ; i < m; i++ {
		element := data[i]
		h.PushPop(element.Value(), element.Priority())
	}
	return h
}

// NLargestDary returns a min-heap of size n containing the n largest
// elements from data. lt should compare priorities by returning true if a < b.
func NLargestDary[V any, P any](n int, d int, data []*HeapPair[V, P], lt func(a, b P) bool) DaryHeap[V, P] {
	return nDary(n, d, data, lt)
}

// NSmallestDary returns a max-heap of size n containing the n smallest
// elements from data. gt should compare priorities by returning true if a > b.
func NSmallestDary[V any, P any](n int, d int, data []*HeapPair[V, P], gt func(a, b P) bool) DaryHeap[V, P] {
	return nDary(n, d, data, gt)
}
