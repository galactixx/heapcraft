package heapcraft

import (
	"fmt"

	"github.com/mohae/deepcopy"
)

// HeapifyCopy takes a slice of T and returns a new heap
// built from a copy of that slice.
func NewDaryHeapCopy[T any](d int, data []T, cmp func(a, b T) bool) DaryHeap[T] {
	heap := make([]T, len(data))
	copy(heap, data)
	return NewDaryHeap(d, heap, cmp)
}

// Heapify transforms an existing slice into a valid heap
// in-place and returns it.
func NewDaryHeap[T any](d int, data []T, cmp func(a, b T) bool) DaryHeap[T] {
	if len(data) == 0 {
		emptyHeap := make([]T, 0)
		return DaryHeap[T]{data: emptyHeap, cmp: cmp, d: d}
	}
	h := DaryHeap[T]{data: data, cmp: cmp, d: d}
	start := (h.Length() - 2) / d
	for i := start; i >= 0; i-- {
		h.siftDown(i)
	}
	return h
}

// Heap is a generic min-heap (or max-heap if cmp defines
// reverse order) with optional swap callbacks.
type DaryHeap[T any] struct {
	data   []T
	cmp    func(a, b T) bool
	onSwap Callbacks
	d      int
}

// Register adds a new swap callback function and returns
// its Callback entry (including ID).
func (h *DaryHeap[T]) Register(fn func(x int, y int)) Callback {
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
func (h *DaryHeap[T]) Deregister(id int) error {
	if _, exists := h.onSwap.callbacks[id]; !exists {
		return fmt.Errorf("%d is not an ID of a callback", id)
	}
	delete(h.onSwap.callbacks, id)
	return nil
}

// swap exchanges elements at indices cmpIdx and cur, then
// runs any registered callbacks.
func (h *DaryHeap[T]) swap(cmpIdx int, cur int) {
	h.data[cmpIdx], h.data[cur] = h.data[cur], h.data[cmpIdx]
	h.onSwap.run(cmpIdx, cur)
}

// swapWithLast swaps the element at index i with the last
// element, removes the last, then sifts down the element now
// at index i. Returns the removed element.
func (h *DaryHeap[T]) swapWithLast(i int) T {
	removed := h.data[i]
	h.data[i] = h.data[h.Length()-1]
	h.data = h.data[:h.Length()-1]
	h.siftDown(i)
	return removed
}

// Clear empties the heap by resetting the slice to length zero.
func (h *DaryHeap[T]) Clear() {
	h.data = h.data[:0]
}

// Length returns the number of elements currently in the heap.
func (h DaryHeap[T]) Length() int {
	return len(h.data)
}

// IsEmpty returns true if the heap has no elements.
func (h DaryHeap[T]) IsEmpty() bool {
	return h.Length() == 0
}

// Peek returns a pointer to the root element without removing
// it; returns nil if empty.
func (h DaryHeap[T]) Peek() *T {
	if h.IsEmpty() {
		return nil
	}
	return &h.data[0]
}

// PopPush pushes element onto the heap, then removes and
// returns the root in one step.
func (h *DaryHeap[T]) PopPush(element T) T {
	h.data = append(h.data, element)
	return h.swapWithLast(0)
}

// PushPop compares element with the root: if element should
// be root (cmp returns true), it returns element and does nothing;
// otherwise, it pushes element and pops the old root.
func (h *DaryHeap[T]) PushPop(element T) T {
	if !h.IsEmpty() && h.cmp(element, *h.Peek()) {
		return element
	}
	h.data = append(h.data, element)
	return h.swapWithLast(0)
}

// DeepClone returns a copy of the heap where each element is
// deep-copied via deepcopy.Copy.
func (h DaryHeap[T]) DeepClone() Heap[T] {
	newData := make([]T, h.Length())
	for i, element := range h.data {
		elementCopy := deepcopy.Copy(element)
		newData[i] = elementCopy.(T)
	}
	return Heap[T]{data: newData, cmp: h.cmp}
}

// Clone returns a shallow copy of the heap (copies the slice
// but not the elements themselves).
func (h DaryHeap[T]) Clone() Heap[T] {
	newData := make([]T, h.Length())
	copy(newData, h.data)
	return Heap[T]{data: newData, cmp: h.cmp}
}

// siftUp restores heap property by moving the element at index
// i upward until its parent is smaller.
func (h *DaryHeap[T]) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / h.d
		if !h.cmp(h.data[i], h.data[parent]) {
			break
		}
		h.swap(i, parent)
		i = parent
	}
}

// siftDown restores heap property by moving the element at
// index i downward until both children are larger.
func (h *DaryHeap[T]) siftDown(i int) {
	cur := i
	n := h.Length()
	for h.d*cur+1 < n {
		left := h.d*cur + 1

		swapIdx := left
		for k := left; k < left+h.d; k++ {
			if k < n && h.cmp(h.data[k], h.data[swapIdx]) {
				swapIdx = k
			}
		}

		if !h.cmp(h.data[swapIdx], h.data[cur]) {
			break
		}
		h.swap(swapIdx, cur)
		cur = swapIdx
	}
}

// Update changes the value at index i to element, then restores
// heap property by sifting up or down.
func (h *DaryHeap[T]) Update(i int, element T) error {
	if i < 0 || i >= h.Length() {
		return fmt.Errorf("index %d is out of bounds", i)
	}
	h.data[i] = element
	if i > 0 && h.cmp(element, h.data[(i-1)/h.d]) {
		h.siftUp(i)
	} else {
		h.siftDown(i)
	}
	return nil
}

// Remove deletes the element at index i, returns its value
// via pointer, and restores heap property.
func (h *DaryHeap[T]) Remove(i int) (*T, error) {
	if i < 0 || i >= h.Length() {
		return nil, fmt.Errorf("index %d is out of bounds", i)
	}
	removed := h.swapWithLast(i)
	return &removed, nil
}

// Pop removes and returns the root element; returns nil
// if heap is empty.
func (h *DaryHeap[T]) Pop() *T {
	if h.IsEmpty() {
		return nil
	}
	removed := h.swapWithLast(0)
	return &removed
}

// Push inserts a new element into the heap and restores heap
// property by sifting up.
func (h *DaryHeap[T]) Push(element T) {
	h.data = append(h.data, element)
	i := h.Length() - 1
	h.siftUp(i)
}

// nDary is a helper for NLargest and NSmallest: it builds
// a size-n heap from data.
func nDary[T any](n int, d int, data []T, cmp func(a, b T) bool) DaryHeap[T] {
	h := DaryHeap[T]{data: make([]T, 0, n), cmp: cmp, d: d}
	i := 0
	m := len(data)
	minNum := min(n, m)

	for ; i < minNum; i++ {
		h.Push(data[i])
	}

	for ; i < m; i++ {
		h.PushPop(data[i])
	}
	return h
}

// NLargest returns a heap of the n largest elements from
// data (min-heap of size n).
// lt should return true if x < y.
func NLargestDary[T any](n int, d int, data []T, lt func(a, b T) bool) DaryHeap[T] {
	return nDary(n, d, data, lt)
}

// NSmallest returns a heap of the n smallest elements from
// data (max-heap of size n).
// gt should return true if x > y.
func NSmallestDary[T any](n int, d int, data []T, gt func(a, b T) bool) DaryHeap[T] {
	return nDary(n, d, data, gt)
}
