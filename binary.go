package heapcraft

import (
	"fmt"
)

// HeapifyCopy duplicates the given slice of *HeapPair and constructs a
// new binary heap from that copy.
func HeapifyCopy[V any, P any](data []*HeapPair[V, P], cmp func(a, b P) bool) BinaryHeap[V, P] {
	heap := make([]*HeapPair[V, P], len(data))
	copy(heap, data)
	return Heapify(heap, cmp)
}

// Heapify rearranges the provided slice of *HeapPair into a valid binary heap
// in place and returns it. It starts from the last parent node (index (n-2)/2)
// and calls siftDown on each node up to the root.
func Heapify[V any, P any](data []*HeapPair[V, P], cmp func(a, b P) bool) BinaryHeap[V, P] {
	if len(data) == 0 {
		return BinaryHeap[V, P]{data: make([]*HeapPair[V, P], 0), cmp: cmp}
	}

	h := BinaryHeap[V, P]{data: data, cmp: cmp}
	start := (h.Length() - 2) / 2
	for i := start; i >= 0; i-- {
		h.siftDown(i)
	}
	return h
}

// HeapPair binds a value to its priority for heap operations.
type HeapPair[V any, P any] struct {
	value    V
	priority P
}

func (b HeapPair[V, P]) Value() V    { return b.value }
func (b HeapPair[V, P]) Priority() P { return b.priority }

// CreateHeapPair constructs a new *HeapPair from the given value and priority.
func CreateHeapPair[V any, P any](value V, priority P) *HeapPair[V, P] {
	return &HeapPair[V, P]{value: value, priority: priority}
}

// Callbacks maintains a registry of callback functions (ID → function).
type Callbacks struct {
	callbacks map[int]Callback
	curId     int
}

// run invokes each registered callback function with the provided indices x and y.
func (c *Callbacks) run(x int, y int) {
	for _, callback := range c.callbacks {
		callback.Function(x, y)
	}
}

// Callback stores a unique ID and the function to invoke when swaps occur.
type Callback struct {
	ID       int
	Function func(x int, y int)
}

// BinaryHeap implements a generic binary heap (min-heap or max-heap depending
// on cmp), with support for swap notifications via callbacks.
//   - data: slice of pointers to HeapPair (each containing a value and priority).
//   - cmp: comparison function on priority (e.g., a < b for min-heap).
//   - onSwap: set of callbacks to invoke whenever two elements swap.
type BinaryHeap[V any, P any] struct {
	data   []*HeapPair[V, P]
	cmp    func(a, b P) bool
	onSwap Callbacks
}

// Register adds a callback function to be invoked on each swap. Returns a
// Callback struct containing the assigned ID for deregistration.
func (h *BinaryHeap[V, P]) Register(fn func(x int, y int)) Callback {
	newId := h.onSwap.curId + 1
	newCallback := Callback{ID: newId, Function: fn}
	if h.onSwap.callbacks == nil {
		h.onSwap.callbacks = make(map[int]Callback)
	}
	h.onSwap.callbacks[newId] = newCallback
	h.onSwap.curId = newId
	return newCallback
}

// Deregister removes the callback with the specified ID from the registry.
// Returns an error if no callback with that ID exists.
func (h *BinaryHeap[V, P]) Deregister(id int) error {
	if _, exists := h.onSwap.callbacks[id]; !exists {
		return fmt.Errorf("%d is not an ID of a callback", id)
	}
	delete(h.onSwap.callbacks, id)
	return nil
}

// swap exchanges the elements at indices i and j in the data slice and
// invokes all registered onSwap callbacks with those indices.
func (h *BinaryHeap[V, P]) swap(i int, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
	h.onSwap.run(i, j)
}

// swapWithLast replaces the element at index i with the last element in the slice,
// shortens the slice by one, then calls siftDown at index i to restore heap order.
// Returns the removed *HeapPair.
func (h *BinaryHeap[V, P]) swapWithLast(i int) *HeapPair[V, P] {
	removed := h.data[i]
	h.data[i] = h.data[h.Length()-1]
	h.data = h.data[:h.Length()-1]
	h.siftDown(i)
	return removed
}

// Clear removes all elements from the heap by setting its underlying slice length to zero.
func (h *BinaryHeap[V, P]) Clear() {
	h.data = h.data[:0]
}

// Length returns the number of elements currently stored in the heap.
func (h BinaryHeap[V, P]) Length() int {
	return len(h.data)
}

// IsEmpty returns true if the heap contains no elements.
func (h BinaryHeap[V, P]) IsEmpty() bool {
	return h.Length() == 0
}

// Peek returns a pointer to the root element (min or max per cmp) without removing it.
// Returns nil if the heap is empty.
func (h BinaryHeap[V, P]) Peek() *HeapPair[V, P] {
	if h.IsEmpty() {
		return nil
	}
	return h.data[0]
}

// PopPush inserts a new element (*HeapPair) into the heap and then immediately
// removes and returns the current root.
func (h *BinaryHeap[V, P]) PopPush(value V, priority P) *HeapPair[V, P] {
	element := &HeapPair[V, P]{value: value, priority: priority}
	h.data = append(h.data, element)
	return h.swapWithLast(0)
}

// PushPop compares a new element’s priority with the root: if the new element
// should become the root (per cmp), it returns that new element without
// modifying the heap. Otherwise, it inserts the new element and removes
// and returns the current root.
func (h *BinaryHeap[V, P]) PushPop(value V, priority P) *HeapPair[V, P] {
	element := &HeapPair[V, P]{value: value, priority: priority}
	if !h.IsEmpty() && h.cmp(element.priority, h.Peek().priority) {
		return element
	}
	h.data = append(h.data, element)
	return h.swapWithLast(0)
}

// Clone creates a shallow copy of the heap: the slice of *HeapPair pointers
// is copied, but the individual HeapPair instances are not duplicated.
func (h BinaryHeap[V, P]) Clone() BinaryHeap[V, P] {
	newData := make([]*HeapPair[V, P], h.Length())
	copy(newData, h.data)
	return BinaryHeap[V, P]{data: newData, cmp: h.cmp}
}

// siftUp moves the element at index i up the tree until the heap property
// (as defined by cmp) is satisfied. Stops when reaching the root or when the
// parent already compares appropriately.
func (h *BinaryHeap[V, P]) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if !h.cmp(h.data[i].priority, h.data[parent].priority) {
			break
		}
		h.swap(i, parent)
		i = parent
	}
}

// siftDown moves the element at index i down the tree until it rests in the
// correct position relative to its children per cmp. At each step, selects the
// child that best satisfies the heap order.
func (h *BinaryHeap[V, P]) siftDown(i int) {
	cur := i
	n := h.Length()
	for 2*cur+1 < n {
		left := 2*cur + 1
		right := left + 1
		var swapIdx int

		// Choose the child according to cmp
		if right >= n || h.cmp(h.data[left].priority, h.data[right].priority) {
			swapIdx = left
		} else {
			swapIdx = right
		}

		if !h.cmp(h.data[swapIdx].priority, h.data[cur].priority) {
			break
		}
		h.swap(swapIdx, cur)
		cur = swapIdx
	}
}

// Update replaces the element at index i with a new (value, priority) pair and
// restores heap order by either sifting up or sifting down, depending on how
// the new priority compares to its parent.
func (h *BinaryHeap[V, P]) Update(i int, value V, priority P) (*HeapPair[V, P], error) {
	if i < 0 || i >= h.Length() {
		return nil, fmt.Errorf("index %d is out of bounds", i)
	}
	element := &HeapPair[V, P]{value: value, priority: priority}
	h.data[i] = element
	if i > 0 && h.cmp(element.priority, h.data[(i-1)/2].priority) {
		h.siftUp(i)
	} else {
		h.siftDown(i)
	}
	return element, nil
}

// Remove deletes the element at index i, returns it, and restores heap order
// by sifting down the replacement.
func (h *BinaryHeap[V, P]) Remove(i int) (*HeapPair[V, P], error) {
	if i < 0 || i >= h.Length() {
		return nil, fmt.Errorf("index %d is out of bounds", i)
	}
	removed := h.swapWithLast(i)
	return removed, nil
}

// Pop removes and returns the root *HeapPair (minimum or maximum per cmp),
// or nil if the heap is empty.
func (h *BinaryHeap[V, P]) Pop() *HeapPair[V, P] {
	if h.IsEmpty() {
		return nil
	}
	return h.swapWithLast(0)
}

// Push inserts a new element (*HeapPair) at the end of the heap and sifts it
// up to maintain heap order.
func (h *BinaryHeap[V, P]) Push(value V, priority P) *HeapPair[V, P] {
	element := &HeapPair[V, P]{value: value, priority: priority}
	h.data = append(h.data, element)
	i := h.Length() - 1
	h.siftUp(i)
	return element
}

// nHeap builds a heap of size n from the input slice by first pushing the
// first min(n, len(data)) elements, then using PushPop for the rest to maintain
// heap size n.
func nHeap[V any, P any](n int, data []*HeapPair[V, P], cmp func(a, b P) bool) BinaryHeap[V, P] {
	h := BinaryHeap[V, P]{data: make([]*HeapPair[V, P], 0, n), cmp: cmp}
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

// NLargest returns a min-heap of size n containing the n largest elements from data.
// The comparator lt should return true if a < b (root is the smallest among the n largest).
func NLargest[V any, P any](n int, data []*HeapPair[V, P], lt func(a, b P) bool) BinaryHeap[V, P] {
	return nHeap(n, data, lt)
}

// NSmallest returns a max-heap of size n containing the n smallest elements from data.
// The comparator gt should return true if a > b (root is the largest among the n smallest).
func NSmallest[V any, P any](n int, data []*HeapPair[V, P], gt func(a, b P) bool) BinaryHeap[V, P] {
	return nHeap(n, data, gt)
}
