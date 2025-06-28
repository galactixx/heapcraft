package heapcraft

// DaryHeap represents a generic d-ary heap with support for swap callbacks. The
// heap can be either a min-heap or max-heap depending on the comparison
// function.   - data: slice of HeapNode containing value-priority pairs   - cmp:
// comparison function that determines the heap order (min or max)   - onSwap:
// callbacks invoked whenever two elements swap positions   - d: the arity of the
// heap (e
type DaryHeap[V any, P any] struct {
	data   []HeapNode[V, P]
	cmp    func(a, b P) bool
	onSwap callbacks
	d      int
	pool   pool[HeapNode[V, P]]
}

// getNewNode creates a new HeapNode with the given value and priority.
// It is used to create new nodes when inserting elements into the heap.
func (h *DaryHeap[V, P]) getNewNode(value V, priority P) HeapNode[V, P] {
	node := h.pool.Get()
	node.value = value
	node.priority = priority
	return node
}

// Deregister removes the callback with the specified ID from the heap's swap
// callbacks. Returns an error if no callback exists with the given ID.
func (h *DaryHeap[V, P]) Deregister(id string) error { return h.onSwap.deregister(id) }

// Register adds a callback function to be called whenever elements in the heap
// swap positions. Returns a callback that can be used to deregister the
// function later.
func (h *DaryHeap[V, P]) Register(fn func(x, y int)) callback { return h.onSwap.register(fn) }

// swap exchanges the elements at indices i and j in the heap, and invokes all
// registered swap callbacks with the indices.
func (h *DaryHeap[V, P]) swap(i int, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
	h.onSwap.run(i, j)
}

// swapWithLast swaps the element at index i with the last element in the heap,
// removes the last element, and sifts down the element now at index i to restore
// heap order. Returns the removed HeapNode.
func (h *DaryHeap[V, P]) swapWithLast(i int) HeapNode[V, P] {
	n := len(h.data)
	removed := h.data[i]
	h.swap(i, n-1)
	h.data = h.data[:n-1]
	h.siftDown(i)
	return removed
}

// Clear removes all elements from the heap by resetting its underlying slice to
// length zero.
func (h *DaryHeap[V, P]) Clear() { h.data = nil }

// Length returns the current number of elements in the heap.
func (h *DaryHeap[V, P]) Length() int { return len(h.data) }

// IsEmpty returns true if the heap contains no elements.
func (h *DaryHeap[V, P]) IsEmpty() bool { return len(h.data) == 0 }

// pop removes and returns the root element of the heap.
// If the heap is empty, returns a zero value SimpleNode with an error.
func (h *DaryHeap[V, P]) pop() (V, P, error) {
	if len(h.data) == 0 {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrHeapEmpty
	}
	removed := h.swapWithLast(0)
	v, p := removed.value, removed.priority
	h.pool.Put(removed)
	return v, p, nil
}

// peek returns the root HeapNode without removing it.
// If the heap is empty, returns a zero value SimpleNode with an error.
func (h *DaryHeap[V, P]) peek() (V, P, error) {
	if len(h.data) == 0 {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrHeapEmpty
	}
	v, p := pairFromNode(h.data[0])
	return v, p, nil
}

// Pop removes and returns the root element of the heap (minimum or maximum per
// cmp). If the heap is empty, returns a zero value SimpleNode with an error.
func (h *DaryHeap[V, P]) Pop() (V, P, error) { return h.pop() }

// Peek returns the root HeapNode without removing it.
// If the heap is empty, returns a zero value SimpleNode with an error.
func (h *DaryHeap[V, P]) Peek() (V, P, error) { return h.peek() }

// PopValue removes and returns just the value of the root element.
// If the heap is empty, returns a zero value with an error.
func (h *DaryHeap[V, P]) PopValue() (V, error) {
	return valueFromNode(h.pop())
}

// PopPriority removes and returns just the priority of the root element.
// If the heap is empty, returns a zero value with an error.
func (h *DaryHeap[V, P]) PopPriority() (P, error) {
	return priorityFromNode(h.pop())
}

// PeekValue returns just the value of the root element without removing it.
// If the heap is empty, returns a zero value with an error.
func (h *DaryHeap[V, P]) PeekValue() (V, error) {
	return valueFromNode(h.peek())
}

// PeekPriority returns just the priority of the root element without removing
// it. If the heap is empty, returns a zero value with an error.
func (h *DaryHeap[V, P]) PeekPriority() (P, error) {
	return priorityFromNode(h.peek())
}

// Push inserts a new element with the given value and priority into the heap.
// The element is added at the end and then sifted up to maintain the heap
// property.
func (h *DaryHeap[V, P]) Push(value V, priority P) {
	element := h.getNewNode(value, priority)
	h.data = append(h.data, element)
	i := len(h.data) - 1
	h.siftUp(i)
}

// siftUp moves the element at index i up the tree until the heap property is
// restored. The heap property is determined by the comparison function cmp,
// where a parent's priority should compare appropriately with its children's
// priorities.
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

// siftDown moves the element at index i down the tree until all children satisfy
// the heap order. For each node, it finds the child with the most appropriate
// priority (per cmp) and swaps if necessary to maintain the heap property.
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

// restoreHeap restores the heap property after an element at index i has been
// updated. It decides whether to sift up or down based on the element's priority
// relative to its parent.
func (h *DaryHeap[V, P]) restoreHeap(i int) {
	if i > 0 && h.cmp(h.data[i].priority, h.data[(i-1)/h.d].priority) {
		h.siftUp(i)
	} else {
		h.siftDown(i)
	}
}

// Update replaces the element at index i with a new value and priority.
// It then restores the heap property by either sifting up (if the new priority
// is more appropriate than its parent) or sifting down (if the new priority is
// less appropriate than its children).
// Returns an error if the index is out of bounds.
func (h *DaryHeap[V, P]) Update(i int, value V, priority P) error {
	if i < 0 || i >= len(h.data) {
		return ErrIndexOutOfBounds
	}
	element := h.getNewNode(value, priority)
	h.data[i] = element
	h.restoreHeap(i)
	return nil
}

// Remove deletes the element at index i from the heap and returns it.
// The heap property is restored by replacing the removed element with the last
// element and sifting it down to its appropriate position.
// Returns the removed element and an error if the index is out of bounds.
func (h *DaryHeap[V, P]) Remove(i int) (V, P, error) {
	if i < 0 || i >= len(h.data) {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrIndexOutOfBounds
	}

	removed := h.data[i]
	h.data[i] = h.data[len(h.data)-1]
	h.data = h.data[:len(h.data)-1]

	idx := i
	if i > 0 {
		idx = i - 1
	}

	v, p := removed.value, removed.priority
	h.restoreHeap(idx)
	h.pool.Put(removed)
	return v, p, nil
}

// PopPush atomically removes the root element and inserts a new element into
// the heap. Returns the removed root element.
func (h *DaryHeap[V, P]) PopPush(value V, priority P) (V, P) {
	element := h.getNewNode(value, priority)
	h.data = append(h.data, element)
	removed := h.swapWithLast(0)
	v, p := removed.value, removed.priority
	h.pool.Put(removed)
	return v, p
}

// PushPop atomically inserts a new element and removes the root element if the
// new element doesn't belong at the root. If the new element belongs at the
// root, it is returned directly. Returns either the new element or the old root
// element.
func (h *DaryHeap[V, P]) PushPop(value V, priority P) (V, P) {
	element := h.getNewNode(value, priority)
	if len(h.data) != 0 && h.cmp(element.priority, h.data[0].priority) {
		return element.value, element.priority
	}
	h.data = append(h.data, element)
	removed := h.swapWithLast(0)
	v, p := removed.value, removed.priority
	h.pool.Put(removed)
	return v, p
}

// Clone creates a deep copy of the heap structure. The new heap preserves the
// original size. If values or priorities are reference types, those reference
// values are shared between the original and cloned heaps.
func (h *DaryHeap[V, P]) Clone() *DaryHeap[V, P] {
	newData := make([]HeapNode[V, P], h.Length())
	copy(newData, h.data)
	callbacks := h.onSwap.getCallbacks()
	return &DaryHeap[V, P]{
		data:   newData,
		cmp:    h.cmp,
		onSwap: callbacks,
		d:      h.d,
		pool:   h.pool,
	}
}
