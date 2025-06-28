package heapcraft

import (
	"reflect"

	"golang.org/x/exp/constraints"
)

// NewRadixHeap creates a RadixHeap from a given slice of HeapNode[V,P].
// It determines the number of buckets from the bit-length of P, initializess
// 'last' to the minimum priority if data is present, and assigns each element
// into its corresponding bucket. The heap maintains a monotonic property where
// priorities must be non-decreasing.
func NewRadixHeap[V any, P constraints.Unsigned](data []HeapNode[V, P], usePool bool) *RadixHeap[V, P] {
	pool := newPool(usePool, func() HeapNode[V, P] {
		return HeapNode[V, P]{}
	})
	var pType P
	t := reflect.TypeOf(pType)
	bits := t.Bits()
	numBuckets := bits + 1
	buckets := make([][]HeapNode[V, P], numBuckets)

	var last P
	var size int

	if len(data) == 0 {
		last = 0
		size = 0
	} else {
		// Determine the minimum priority among the input items
		last = minFromSlice(data).priority
		size = len(data)

		// Push each item into the appropriate bucket relative to 'last'
		for _, pair := range data {
			rPair := pool.Get()
			rPair.value = pair.value
			rPair.priority = pair.priority
			bucketInsert(rPair, last, buckets)
		}
	}

	return &RadixHeap[V, P]{
		buckets: buckets, size: size, last: last, pool: pool,
	}
}

// NewSyncRadixHeap creates a new thread-safe RadixHeap from a given slice of HeapNode[V,P].
func NewSyncRadixHeap[V any, P constraints.Unsigned](data []HeapNode[V, P], usePool bool) *SyncRadixHeap[V, P] {
	return &SyncRadixHeap[V, P]{heap: NewRadixHeap(data, usePool)}
}
