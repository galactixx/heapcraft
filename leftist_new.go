package heapcraft

import "github.com/google/uuid"

// NewLeftistHeap constructs a leftist heap from a slice of HeapPairs.
// Uses a queue to iteratively merge singleton nodes until one root remains.
// The comparison function determines the heap order (min or max).
func NewLeftistHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *LeftistHeap[V, P] {
	pool := newPool(usePool, func() *leftistNode[V, P] {
		return &leftistNode[V, P]{}
	})
	heap := LeftistHeap[V, P]{cmp: cmp, size: 0, pool: pool}
	if len(data) == 0 {
		return &heap
	}

	n := len(data)
	queueData := make([]*leftistNode[V, P], 0, n)
	initQueue := leftistQueue[*leftistNode[V, P]]{data: queueData, head: 0, size: 0}

	heap.size = n

	for i := range data {
		node := pool.Get()
		node.value = data[i].value
		node.priority = data[i].priority
		node.s = 1
		initQueue.push(node)
	}

	for initQueue.remainingElements() > 1 {
		merged := heap.merge(initQueue.pop(), initQueue.pop())
		initQueue.push(merged)
	}

	heap.root = initQueue.pop()
	return &heap
}

// NewLeftistHeap constructs a leftist heap with node tracking from a slice of HeapPairs.
// Each node is assigned a unique ID and stored in a map for O(1) access.
// Uses a queue to iteratively merge singleton nodes until one root remains.
// The comparison function determines the heap order (min or max).
func NewFullLeftistHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, config HeapConfig) *FullLeftistHeap[V, P] {
	pool := newPool(config.UsePool, func() *leftistHeapNode[V, P] {
		return &leftistHeapNode[V, P]{}
	})
	elements := make(map[string]*leftistHeapNode[V, P])
	heap := FullLeftistHeap[V, P]{
		cmp:      cmp,
		size:     0,
		elements: elements,
		pool:     pool,
		idGen:    config.GetGenerator(),
	}
	if len(data) == 0 {
		return &heap
	}

	n := len(data)
	queueData := make([]*leftistHeapNode[V, P], 0, n)
	initQueue := leftistQueue[*leftistHeapNode[V, P]]{data: queueData, head: 0, size: 0}

	heap.size = n

	for i := range data {
		node := pool.Get()
		node.id = uuid.New().String()
		node.value = data[i].value
		node.priority = data[i].priority
		node.s = 1
		initQueue.push(node)
		elements[node.id] = node
	}

	for initQueue.remainingElements() > 1 {
		merged := heap.merge(initQueue.pop(), initQueue.pop())
		initQueue.push(merged)
	}

	heap.root = initQueue.pop()
	return &heap
}

// NewSyncFullLeftistHeap constructs a new thread-safe leftist heap from the
// given data and comparison function.
// The resulting heap is safe for concurrent use.
func NewSyncFullLeftistHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, config HeapConfig) *SyncFullLeftistHeap[V, P] {
	return &SyncFullLeftistHeap[V, P]{
		heap: NewFullLeftistHeap(data, cmp, config),
	}
}

// NewSyncLeftistHeap constructs a new thread-safe leftist
// heap from the given data and comparison function.
// The resulting heap is safe for concurrent use.
func NewSyncLeftistHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *SyncLeftistHeap[V, P] {
	return &SyncLeftistHeap[V, P]{
		heap: NewLeftistHeap(data, cmp, usePool),
	}
}
