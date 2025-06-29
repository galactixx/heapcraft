package heapcraft

// NewFullPairingHeap creates a new pairing heap from a slice of HeapPairs.
// The heap is initialized with the provided elements and uses the given comparison
// function to determine heap order. The comparison function determines the heap order (min or max).
// Returns an empty heap if the input slice is empty.
func NewFullPairingHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, config HeapConfig) *FullPairingHeap[V, P] {
	pool := newPool(config.UsePool, func() *pairingHeapNode[V, P] {
		return &pairingHeapNode[V, P]{}
	})
	elements := make(map[string]*pairingHeapNode[V, P])
	heap := FullPairingHeap[V, P]{
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

// NewPairingHeap creates a new simple pairing heap from a slice of HeapPairs.
// Unlike PairingHeap, this implementation does not track node IDs or support
// node updates. It uses the provided comparison function to determine heap order (min or max).
// Returns an empty heap if the input slice is empty.
func NewPairingHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *PairingHeap[V, P] {
	pool := newPool(usePool, func() *pairingNode[V, P] {
		return &pairingNode[V, P]{}
	})
	heap := PairingHeap[V, P]{cmp: cmp, size: 0, pool: pool}
	if len(data) == 0 {
		return &heap
	}

	for i := range data {
		heap.Push(data[i].value, data[i].priority)
	}
	return &heap
}

// NewSyncPairingHeap creates a new thread-safe pairing heap from a slice of HeapPairs.
// The heap is initialized with the provided elements and uses the given comparison
// function to determine heap order. The comparison function determines the heap order (min or max).
// Returns an empty heap if the input slice is empty.
func NewSyncFullPairingHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, config HeapConfig) *SyncFullPairingHeap[V, P] {
	return &SyncFullPairingHeap[V, P]{heap: NewFullPairingHeap(data, cmp, config)}
}

// NewSyncPairingHeap creates a new thread-safe simple pairing heap from a slice of HeapPairs.
// Unlike SyncPairingHeap, this implementation does not track node IDs or support
// node updates. It uses the provided comparison function to determine heap order (min or max).
// Returns an empty heap if the input slice is empty.
func NewSyncPairingHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *SyncPairingHeap[V, P] {
	return &SyncPairingHeap[V, P]{heap: NewPairingHeap(data, cmp, usePool)}
}
