package heapcraft

// NewSkewHeap creates a new skew heap from the given data slice.
// Each element is inserted individually using the provided comparison function
// to determine heap order (min or max). Returns an empty heap if the input
// slice is empty.
func NewSkewHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, config HeapConfig) *SkewHeap[V, P] {
	pool := newPool(config.UsePool, func() *skewHeapNode[V, P] {
		return &skewHeapNode[V, P]{}
	})
	elements := make(map[string]*skewHeapNode[V, P], len(data))
	heap := SkewHeap[V, P]{
		cmp:      cmp,
		size:     0,
		elements: elements,
		pool:     pool,
		idGen:    config.GetGenerator(),
	}
	if len(data) == 0 {
		return &heap
	}

	for i := range data {
		heap.Push(data[i].value, data[i].priority)
	}
	return &heap
}

// NewSimpleSkewHeap creates a new simple skew heap from the given data slice.
// Each element is inserted individually using the provided comparison function
// to determine heap order (min or max). Returns an empty heap if the input
// slice is empty.
func NewSimpleSkewHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *SimpleSkewHeap[V, P] {
	pool := newPool(usePool, func() *skewNode[V, P] {
		return &skewNode[V, P]{}
	})
	heap := SimpleSkewHeap[V, P]{cmp: cmp, size: 0, pool: pool}
	if len(data) == 0 {
		return &heap
	}

	for i := range data {
		heap.Push(data[i].value, data[i].priority)
	}
	return &heap
}

// NewSyncSkewHeap constructs a new thread-safe skew heap from the given data and comparison function.
// The resulting heap is safe for concurrent use.
func NewSyncSkewHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, config HeapConfig) *SyncSkewHeap[V, P] {
	return &SyncSkewHeap[V, P]{
		heap: NewSkewHeap(data, cmp, config),
	}
}

// NewSyncSimpleSkewHeap constructs a new thread-safe simple skew heap from the given data and comparison function.
// The resulting heap is safe for concurrent use.
func NewSyncSimpleSkewHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *SyncSimpleSkewHeap[V, P] {
	return &SyncSimpleSkewHeap[V, P]{
		heap: NewSimpleSkewHeap(data, cmp, usePool),
	}
}
