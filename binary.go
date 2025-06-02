package heapcraft

import (
	"fmt"

	"github.com/mohae/deepcopy"
)

// Callbacks holds a registry of callback functions keyed
// by an integer ID.
type Callbacks struct {
	callbacks map[int]Callback
	curId     int
}

// run invokes all registered callbacks, passing x and y
// as arguments.
func (c *Callbacks) run(x int, y int) {
	for _, callback := range c.callbacks {
		callback.Function(x, y)
	}
}

// Callback represents a single callback entry with a unique
// ID and the function to call.
type Callback struct {
	ID       int
	Function func(x int, y int)
}

// HeapifyCopy takes a slice of T and returns a new heap
// built from a copy of that slice.
func HeapifyCopy[T any](data []T, cmp func(a, b T) bool) Heap[T] {
	heap := make([]T, len(data))
	copy(heap, data)
	return Heapify(heap, cmp)
}

// Heapify transforms an existing slice into a valid heap
// in-place and returns it.
func Heapify[T any](data []T, cmp func(a, b T) bool) Heap[T] {
	if len(data) == 0 {
		emptyHeap := make([]T, 0)
		return Heap[T]{data: emptyHeap, cmp: cmp}
	}
	h := Heap[T]{data: data, cmp: cmp}
	start := (h.Length() - 2) / 2
	for i := start; i >= 0; i-- {
		h.siftDown(i)
	}
	return h
}

// Heap is a generic min-heap (or max-heap if cmp defines
// reverse order) with optional swap callbacks.
type Heap[T any] struct {
	data   []T
	cmp    func(a, b T) bool
	onSwap Callbacks
}

// Register adds a new swap callback function and returns
// its Callback entry (including ID).
func (h *Heap[T]) Register(fn func(x int, y int)) Callback {
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
func (h *Heap[T]) Deregister(id int) error {
	if _, exists := h.onSwap.callbacks[id]; !exists {
		return fmt.Errorf("%d is not an ID of a callback", id)
	}
	delete(h.onSwap.callbacks, id)
	return nil
}

// swap exchanges elements at indices cmpIdx and cur, then
// runs any registered callbacks.
func (h *Heap[T]) swap(cmpIdx int, cur int) {
	h.data[cmpIdx], h.data[cur] = h.data[cur], h.data[cmpIdx]
	h.onSwap.run(cmpIdx, cur)
}

// swapWithLast swaps the element at index i with the last
// element, removes the last, then sifts down the element now
// at index i. Returns the removed element.
func (h *Heap[T]) swapWithLast(i int) T {
	removed := h.data[i]
	h.data[i] = h.data[h.Length()-1]
	h.data = h.data[:h.Length()-1]
	h.siftDown(i)
	return removed
}

// Clear empties the heap by resetting the slice to length zero.
func (h *Heap[T]) Clear() {
	h.data = h.data[:0]
}

// Length returns the number of elements currently in the heap.
func (h Heap[T]) Length() int {
	return len(h.data)
}

// IsEmpty returns true if the heap has no elements.
func (h Heap[T]) IsEmpty() bool {
	return h.Length() == 0
}

// Peek returns a pointer to the root element without removing
// it; returns nil if empty.
func (h Heap[T]) Peek() *T {
	if h.IsEmpty() {
		return nil
	}
	return &h.data[0]
}

// PopPush pushes element onto the heap, then removes and
// returns the root in one step.
func (h *Heap[T]) PopPush(element T) T {
	h.data = append(h.data, element)
	return h.swapWithLast(0)
}

// PushPop compares element with the root: if element should
// be root (cmp returns true), it returns element and does nothing;
// otherwise, it pushes element and pops the old root.
func (h *Heap[T]) PushPop(element T) T {
	if !h.IsEmpty() && h.cmp(element, *h.Peek()) {
		return element
	}
	h.data = append(h.data, element)
	return h.swapWithLast(0)
}

// DeepClone returns a copy of the heap where each element is
// deep-copied via deepcopy.Copy.
func (h Heap[T]) DeepClone() Heap[T] {
	newData := make([]T, h.Length())
	for i, element := range h.data {
		elementCopy := deepcopy.Copy(element)
		newData[i] = elementCopy.(T)
	}
	return Heap[T]{data: newData, cmp: h.cmp}
}

// Clone returns a shallow copy of the heap (copies the slice
// but not the elements themselves).
func (h Heap[T]) Clone() Heap[T] {
	newData := make([]T, h.Length())
	copy(newData, h.data)
	return Heap[T]{data: newData, cmp: h.cmp}
}

// siftUp restores heap property by moving the element at index
// i upward until its parent is smaller.
func (h *Heap[T]) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if !h.cmp(h.data[i], h.data[parent]) {
			break
		}
		h.swap(i, parent)
		i = parent
	}
}

// siftDown restores heap property by moving the element at
// index i downward until both children are larger.
func (h *Heap[T]) siftDown(i int) {
	cur := i
	n := h.Length()
	for 2*cur+1 < n {
		left := 2*cur + 1
		right := left + 1
		var swapIdx int

		// pick the child according to cmp
		if right >= n || h.cmp(h.data[left], h.data[right]) {
			swapIdx = left
		} else {
			swapIdx = right
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
func (h *Heap[T]) Update(i int, element T) error {
	if i < 0 || i >= h.Length() {
		return fmt.Errorf("index %d is out of bounds", i)
	}
	h.data[i] = element
	if i > 0 && h.cmp(element, h.data[(i-1)/2]) {
		h.siftUp(i)
	} else {
		h.siftDown(i)
	}
	return nil
}

// Remove deletes the element at index i, returns its value
// via pointer, and restores heap property.
func (h *Heap[T]) Remove(i int) (*T, error) {
	if i < 0 || i >= h.Length() {
		return nil, fmt.Errorf("index %d is out of bounds", i)
	}
	removed := h.swapWithLast(i)
	return &removed, nil
}

// Pop removes and returns the root element; returns nil
// if heap is empty.
func (h *Heap[T]) Pop() *T {
	if h.IsEmpty() {
		return nil
	}
	removed := h.swapWithLast(0)
	return &removed
}

// Push inserts a new element into the heap and restores heap
// property by sifting up.
func (h *Heap[T]) Push(element T) {
	h.data = append(h.data, element)
	i := h.Length() - 1
	h.siftUp(i)
}

// nHeap is a helper for NLargest and NSmallest: it builds
// a size-n heap from data.
func nHeap[T any](n int, data []T, cmp func(a, b T) bool) Heap[T] {
	h := Heap[T]{data: make([]T, 0, n), cmp: cmp}
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
func NLargest[T any](n int, data []T, lt func(a, b T) bool) Heap[T] {
	return nHeap(n, data, lt)
}

// NSmallest returns a heap of the n smallest elements from
// data (max-heap of size n).
// gt should return true if x > y.
func NSmallest[T any](n int, data []T, gt func(a, b T) bool) Heap[T] {
	return nHeap(n, data, gt)
}
