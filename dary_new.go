package heapcraft

// NewBinaryHeap creates a new binary heap (d=2) from the given data slice and
// comparison function. The comparison function determines the heap order (min or
// max). It is a convenience wrapper around NewDaryHeap with d=2.
func NewBinaryHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *DaryHeap[V, P] {
	return NewDaryHeap(2, data, cmp, usePool)
}

// NewBinaryHeapCopy creates a new binary heap (d=2) from a copy of the given data
// slice. Unlike NewBinaryHeap, this function creates a new slice and copies the
// data before heapifying it, leaving the original data unchanged. The comparison
// function determines the heap order (min or max). It is a convenience wrapper
// around NewDaryHeapCopy with d=2.
func NewBinaryHeapCopy[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *DaryHeap[V, P] {
	return NewDaryHeapCopy(2, data, cmp, usePool)
}

// NewDaryHeapCopy creates a new d-ary heap from a copy of the provided data
// slice. The comparison function determines the heap order (min or max). The
// original data slice remains unchanged.
func NewDaryHeapCopy[V any, P any](d int, data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *DaryHeap[V, P] {
	heap := make([]HeapNode[V, P], len(data))
	copy(heap, data)
	return NewDaryHeap(d, heap, cmp, usePool)
}

// NewDaryHeap transforms the given slice of HeapNode into a valid d-ary heap
// in-place. The comparison function determines the heap order (min or max).
// Uses siftDown starting from the last parent toward the root to build the heap.
func NewDaryHeap[V any, P any](d int, data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *DaryHeap[V, P] {
	pool := newPool(usePool, func() HeapNode[V, P] {
		return HeapNode[V, P]{}
	})

	callbacks := make(baseCallbacks, 0)
	h := DaryHeap[V, P]{
		data:   data,
		cmp:    cmp,
		onSwap: callbacks,
		d:      d,
		pool:   pool,
	}

	// Start sifting down from the last parent node toward the root.
	start := (h.Length() - 2) / d
	for i := start; i >= 0; i-- {
		h.siftDown(i)
	}
	return &h
}

// nDary builds a heap of size n from the data slice.
// It uses Push for the first n elements and PushPop for the remainder to
// maintain a heap of exactly size n. This is used as the underlying
// implementation for both NLargestDary and NSmallestDary.
func nDary[V any, P any](n int, d int, data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *DaryHeap[V, P] {
	pool := newPool(usePool, func() HeapNode[V, P] {
		return HeapNode[V, P]{}
	})

	callbacks := make(baseCallbacks, 0)
	heap := DaryHeap[V, P]{
		data:   make([]HeapNode[V, P], 0, n),
		cmp:    cmp,
		onSwap: callbacks,
		d:      d,
		pool:   pool,
	}
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
// elements from data. The comparison function lt should return true if a < b.
func NLargestDary[V any, P any](n int, d int, data []HeapNode[V, P], lt func(a, b P) bool, usePool bool) *DaryHeap[V, P] {
	return nDary(n, d, data, lt, usePool)
}

// NLargestBinary returns a min-heap of size n containing the n largest
// elements from data, using a binary heap (d=2). The comparison function lt
// should return true if a < b. This is a convenience wrapper around
// NLargestDary.
func NLargestBinary[V any, P any](n int, data []HeapNode[V, P], lt func(a, b P) bool, usePool bool) *DaryHeap[V, P] {
	return NLargestDary(n, 2, data, lt, usePool)
}

// NSmallestDary returns a max-heap of size n containing the n smallest
// elements from data. The comparison function gt should return true if a > b.
func NSmallestDary[V any, P any](n int, d int, data []HeapNode[V, P], gt func(a, b P) bool, usePool bool) *DaryHeap[V, P] {
	return nDary(n, d, data, gt, usePool)
}

// NSmallestBinary returns a max-heap of size n containing the n smallest
// elements from data, using a binary heap (d=2). The comparison function gt
// should return true if a > b. This is a convenience wrapper around
// NSmallestDary.
func NSmallestBinary[V any, P any](n int, data []HeapNode[V, P], gt func(a, b P) bool, usePool bool) *DaryHeap[V, P] {
	return NSmallestDary(n, 2, data, gt, usePool)
}

// NewSyncBinaryHeap creates a new thread-safe binary heap (d=2) from the given
// data slice and comparison function. The comparison function determines the
// heap order (min or max).
func NewSyncBinaryHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *SyncDaryHeap[V, P] {
	return NewSyncDaryHeap(2, data, cmp, usePool)
}

// NewSyncBinaryHeapCopy creates a new thread-safe binary heap (d=2) from a copy
// of the given data slice. Unlike NewSyncBinaryHeap, this function creates a
// new slice and copies the data before heapifying it, leaving the original data
// unchanged.
func NewSyncBinaryHeapCopy[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *SyncDaryHeap[V, P] {
	return NewSyncDaryHeapCopy(2, data, cmp, usePool)
}

// NewSyncDaryHeapCopy creates a new thread-safe d-ary heap from a copy of the
// provided data slice. The comparison function determines the heap order (min or
// max). The original data slice remains unchanged.
func NewSyncDaryHeapCopy[V any, P any](d int, data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *SyncDaryHeap[V, P] {
	heap := NewDaryHeapCopy(d, data, cmp, usePool)
	heap.onSwap = NewSyncCallbacks()
	return &SyncDaryHeap[V, P]{heap: heap}
}

// NewSyncDaryHeap creates a new thread-safe d-ary heap from the given data
// slice and comparison function. The comparison function determines the heap
// order (min or max).
func NewSyncDaryHeap[V any, P any](d int, data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *SyncDaryHeap[V, P] {
	heap := NewDaryHeap(d, data, cmp, usePool)
	heap.onSwap = NewSyncCallbacks()
	return &SyncDaryHeap[V, P]{heap: heap}
}
