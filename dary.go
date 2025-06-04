package heapcraft

import (
	"fmt"
)

// HeapifyCopy takes a slice of HeapPair and returns a new heap
// built from a copy of that slice.
func NewDaryHeapCopy[V any, P any](d int, data []HeapPair[V, P], cmp func(a, b P) bool) DaryHeap[V, P] {
	heap := make([]HeapPair[V, P], len(data))
	copy(heap, data)
	return NewDaryHeap(d, heap, cmp)
}

// Heapify transforms an existing slice into a valid heap
// in-place and returns it.
func NewDaryHeap[V any, P any](d int, data []HeapPair[V, P], cmp func(a, b P) bool) DaryHeap[V, P] {
	if len(data) == 0 {
		emptyHeap := make([]HeapPair[V, P], 0)
		return DaryHeap[V, P]{data: emptyHeap, cmp: cmp, d: d}
	}
	h := DaryHeap[V, P]{data: data, cmp: cmp, d: d}
	start := (h.Length() - 2) / d
	for i := start; i >= 0; i-- {
		h.siftDown(i)
	}
	return h
}

// Heap is a generic min-heap (or max-heap if cmp defines
// reverse order) with optional swap callbacks.
type DaryHeap[V any, P any] struct {
	data   []HeapPair[V, P]
	cmp    func(a, b P) bool
	onSwap Callbacks
	d      int
}

// Register adds a new swap callback function and returns
// its Callback entry (including ID).
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

// Deregister removes a previously registered callback by its ID.
func (h *DaryHeap[V, P]) Deregister(id int) error {
	if _, exists := h.onSwap.callbacks[id]; !exists {
		return fmt.Errorf("%d is not an ID of a callback", id)
	}
	delete(h.onSwap.callbacks, id)
	return nil
}

// swap exchanges elements at indices cmpIdx and cur, then
// runs any registered callbacks.
func (h *DaryHeap[V, P]) swap(cmpIdx int, cur int) {
	h.data[cmpIdx], h.data[cur] = h.data[cur], h.data[cmpIdx]
	h.onSwap.run(cmpIdx, cur)
}

// swapWithLast swaps the element at index i with the last
// element, removes the last, then sifts down the element now
// at index i. Returns the removed element.
func (h *DaryHeap[V, P]) swapWithLast(i int) *HeapPair[V, P] {
	removed := h.data[i]
	h.data[i] = h.data[h.Length()-1]
	h.data = h.data[:h.Length()-1]
	h.siftDown(i)
	return &removed
}

// Clear empties the heap by resetting the slice to length zero.
func (h *DaryHeap[V, P]) Clear() {
	h.data = h.data[:0]
}

// Length returns the number of elements currently in the heap.
func (h DaryHeap[V, P]) Length() int {
	return len(h.data)
}

// IsEmpty returns true if the heap has no elements.
func (h DaryHeap[V, P]) IsEmpty() bool {
	return h.Length() == 0
}

// Peek returns a pointer to the root element without removing
// it; returns nil if empty.
func (h DaryHeap[V, P]) Peek() *HeapPair[V, P] {
	if h.IsEmpty() {
		return nil
	}
	return &h.data[0]
}

// PopPush pushes element onto the heap, then removes and
// returns the root in one step.
func (h *DaryHeap[V, P]) PopPush(value V, priority P) *HeapPair[V, P] {
	element := HeapPair[V, P]{value: value, priority: priority}
	h.data = append(h.data, element)
	return h.swapWithLast(0)
}

// PushPop compares element with the root: if element should
// be root (cmp returns true), it returns element and does nothing;
// otherwise, it pushes element and pops the old root.
func (h *DaryHeap[V, P]) PushPop(value V, priority P) *HeapPair[V, P] {
	element := HeapPair[V, P]{value: value, priority: priority}
	if !h.IsEmpty() && h.cmp(element.priority, h.Peek().priority) {
		return &element
	}
	h.data = append(h.data, element)
	return h.swapWithLast(0)
}

// Clone returns a shallow copy of the heap (copies the slice
// but not the elements themselves).
func (h DaryHeap[V, P]) Clone() DaryHeap[V, P] {
	newData := make([]HeapPair[V, P], h.Length())
	copy(newData, h.data)
	return DaryHeap[V, P]{data: newData, cmp: h.cmp, d: h.d}
}

// siftUp restores heap property by moving the element at index
// i upward until its parent is smaller.
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

// siftDown restores heap property by moving the element at
// index i downward until both children are larger.
func (h *DaryHeap[V, P]) siftDown(i int) {
	cur := i
	n := h.Length()
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

// Update changes the value at index i to element, then restores
// heap property by sifting up or down.
func (h *DaryHeap[V, P]) Update(i int, value V, priority P) (*HeapPair[V, P], error) {
	if i < 0 || i >= h.Length() {
		return nil, fmt.Errorf("index %d is out of bounds", i)
	}
	element := HeapPair[V, P]{value: value, priority: priority}
	h.data[i] = element
	if i > 0 && h.cmp(element.priority, h.data[(i-1)/h.d].priority) {
		h.siftUp(i)
	} else {
		h.siftDown(i)
	}
	return &element, nil
}

// Remove deletes the element at index i, returns its value
// via pointer, and restores heap property.
func (h *DaryHeap[V, P]) Remove(i int) (*HeapPair[V, P], error) {
	if i < 0 || i >= h.Length() {
		return nil, fmt.Errorf("index %d is out of bounds", i)
	}
	removed := h.swapWithLast(i)
	return removed, nil
}

// Pop removes and returns the root element; returns nil
// if heap is empty.
func (h *DaryHeap[V, P]) Pop() *HeapPair[V, P] {
	if h.IsEmpty() {
		return nil
	}
	removed := h.swapWithLast(0)
	return removed
}

// Push inserts a new element into the heap and restores heap
// property by sifting up.
func (h *DaryHeap[V, P]) Push(value V, priority P) *HeapPair[V, P] {
	element := HeapPair[V, P]{value: value, priority: priority}
	h.data = append(h.data, element)
	i := h.Length() - 1
	h.siftUp(i)
	return &element
}

// nDary is a helper for NLargest and NSmallest: it builds
// a size-n heap from data.
func nDary[V any, P any](n int, d int, data []HeapPair[V, P], cmp func(a, b P) bool) DaryHeap[V, P] {
	h := DaryHeap[V, P]{data: make([]HeapPair[V, P], 0, n), cmp: cmp, d: d}
	i := 0
	m := len(data)
	minNum := min(n, m)

	for ; i < minNum; i++ {
		element := data[i]
		h.Push(element.value, element.priority)
	}

	for ; i < m; i++ {
		element := data[i]
		h.PushPop(element.value, element.priority)
	}
	return h
}

// NLargest returns a heap of the n largest elements from
// data (min-heap of size n).
// lt should return true if x < y.
func NLargestDary[V any, P any](n int, d int, data []HeapPair[V, P], lt func(a, b P) bool) DaryHeap[V, P] {
	return nDary(n, d, data, lt)
}

// NSmallest returns a heap of the n smallest elements from
// data (max-heap of size n).
// gt should return true if x > y.
func NSmallestDary[V any, P any](n int, d int, data []HeapPair[V, P], gt func(a, b P) bool) DaryHeap[V, P] {
	return nDary(n, d, data, gt)
}
