package heapcraft

import (
	"math"

	"golang.org/x/exp/constraints"
)

// cloneBuckets creates a shallow copy of the buckets array, copying each bucket slice.
// The elements within each bucket are shared between the original and the copy.
func cloneBuckets[V any, P constraints.Unsigned](buckets [][]HeapNode[V, P]) [][]HeapNode[V, P] {
	newBuckets := make([][]HeapNode[V, P], len(buckets))
	for i, bucket := range buckets {
		newBuckets[i] = make([]HeapNode[V, P], len(bucket))
		copy(newBuckets[i], bucket)
	}
	return newBuckets
}

// RadixHeap implements a monotonic priority queue over unsigned priorities.
// The heap maintains the invariant that priorities must be non-decreasing.
//   - buckets: array of slices of HeapNode, each holding items whose priorities
//     fall within a range defined by 'last'.
//   - size: the count of elements in the heap.
//   - last: the most recently extracted minimum priority.
type RadixHeap[V any, P constraints.Unsigned] struct {
	buckets [][]HeapNode[V, P]
	size    int
	last    P
	pool    pool[HeapNode[V, P]]
}

// Clone creates a deep copy of the heap structure. The new heap preserves the
// original size and last value. If values or priorities are reference types, those
// reference values are shared between the original and cloned heaps.
func (r *RadixHeap[V, P]) Clone() *RadixHeap[V, P] {
	return &RadixHeap[V, P]{
		buckets: cloneBuckets(r.buckets),
		size:    r.size,
		last:    r.last,
		pool:    r.pool,
	}
}

// Push adds a new value and priority pair into the heap.
// Returns an error if the priority is less than r.last, as this would violate
// the monotonic property. Otherwise, puts the item into the appropriate bucket
// and increments the size.
func (r *RadixHeap[V, P]) Push(value V, priority P) error {
	return r.push(value, priority)
}

// push is an unexported helper that forms a HeapNode and places it into its bucket.
// It enforces the condition that priority must not be less than r.last to maintain
// the monotonic property of the heap.
func (r *RadixHeap[V, P]) push(value V, priority P) error {
	if r.size == 0 {
		r.last = priority
	}

	if priority < r.last {
		return ErrPriorityLessThanLast
	}
	newPair := r.pool.Get()
	newPair.value = value
	newPair.priority = priority
	bucketInsert(newPair, r.last, r.buckets)
	r.size++
	return nil
}

// getMin removes and returns the first element from bucket 0.
// It also decreases the total size. The caller must ensure bucket 0 is not empty.
func (r *RadixHeap[V, P]) getMin() HeapNode[V, P] {
	minPair := r.buckets[0][0]
	r.buckets[0] = r.buckets[0][1:]
	r.size--
	return minPair
}

// pop removes and returns the first element in bucket 0.
// If bucket 0 is empty, it rebalances the heap before returning the minimum.
// Returns nil and an error if the heap is empty.
func (r *RadixHeap[V, P]) pop() (V, P, error) {
	if r.size == 0 {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrHeapEmpty
	}

	// If bucket 0 has entries, pop directly
	if len(r.buckets[0]) > 0 {
		removed := r.getMin()
		v, p := removed.value, removed.priority
		r.pool.Put(removed)
		return v, p, nil
	}

	// Otherwise, refill bucket 0 from the next non-empty bucket
	r.rebalance()
	removed := r.getMin()
	v, p := removed.value, removed.priority
	r.pool.Put(removed)
	return v, p, nil
}

// peek returns the HeapNode with the minimum priority without removing it.
// If bucket 0 has elements, it returns the first one. Otherwise, it finds
// the minimum element in the next non-empty bucket.
// Returns nil and an error if the heap is empty.
func (r *RadixHeap[V, P]) peek() (V, P, error) {
	if r.size == 0 {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrHeapEmpty
	}
	if len(r.buckets[0]) > 0 {
		v, p := pairFromNode(r.buckets[0][0])
		return v, p, nil
	}
	var bucket []HeapNode[V, P]
	for i := 1; i < len(r.buckets); i++ {
		if len(r.buckets[i]) > 0 {
			bucket = r.buckets[i]
			break
		}
	}
	minPair := minFromSlice(bucket)
	v, p := pairFromNode(minPair)
	return v, p, nil
}

// Pop extracts and returns the HeapNode with the minimum priority.
// Returns nil and an error if the heap is empty.
func (r *RadixHeap[V, P]) Pop() (V, P, error) { return r.pop() }

// Peek returns a HeapNode with the minimum priority without removing it.
// Returns nil and an error if the heap is empty.
func (r *RadixHeap[V, P]) Peek() (V, P, error) { return r.peek() }

// PopValue removes and returns just the value of the root element.
// Returns zero value and an error if the heap is empty.
func (r *RadixHeap[V, P]) PopValue() (V, error) {
	return valueFromNode(r.pop())
}

// PopPriority removes and returns just the priority of the root element.
// Returns zero value and an error if the heap is empty.
func (r *RadixHeap[V, P]) PopPriority() (P, error) {
	return priorityFromNode(r.pop())
}

// PeekValue returns just the value of the root element without removing it.
// Returns zero value and an error if the heap is empty.
func (r *RadixHeap[V, P]) PeekValue() (V, error) {
	return valueFromNode(r.peek())
}

// PeekPriority returns just the priority of the root element without removing it.
// Returns zero value and an error if the heap is empty.
func (r *RadixHeap[V, P]) PeekPriority() (P, error) {
	return priorityFromNode(r.peek())
}

// Clear reinitializes the heap by creating fresh buckets, resetting size to zero,
// and setting 'last' back to its zero value.
func (r *RadixHeap[V, P]) Clear() {
	r.buckets = make([][]HeapNode[V, P], len(r.buckets))
	r.size = 0
	r.last = 0
}

// rebalance locates the next bucket with elements (i > 0), updates 'last'
// to the smallest priority found there, and reinserts all items from that bucket
// into new buckets based on the updated 'last'. Afterward, it empties that bucket.
// This operation maintains the monotonic property of the heap.
func (r *RadixHeap[V, P]) rebalance() {
	for i := 1; i < len(r.buckets); i++ {
		if len(r.buckets[i]) > 0 {
			r.last = minFromSlice(r.buckets[i]).priority
			for _, pair := range r.buckets[i] {
				bucketInsert(pair, r.last, r.buckets)
			}
			r.buckets[i] = make([]HeapNode[V, P], 0)
			return
		}
	}
}

// Rebalance fills bucket 0 if it is empty.
// Returns an error if the heap is empty, or if bucket 0 already contains elements
// (no action was needed).
func (r *RadixHeap[V, P]) Rebalance() error {
	if r.size == 0 {
		return ErrHeapEmpty
	}
	if len(r.buckets[0]) == 0 {
		r.rebalance()
		return nil
	}
	return ErrNoRebalancingNeeded
}

// Length returns the number of items currently stored in the heap.
func (r *RadixHeap[V, P]) Length() int { return r.size }

// IsEmpty returns true if the heap contains no items.
func (r *RadixHeap[V, P]) IsEmpty() bool { return r.size == 0 }

// Merge integrates another RadixHeap into this one.
// It selects the heap with the smaller 'last' as the new baseline, adopts its
// buckets and 'last', then reinserts all items from the other heap to preserve
// the monotonic property.
func (r *RadixHeap[V, P]) Merge(radix *RadixHeap[V, P]) {
	var newRadix *RadixHeap[V, P]
	if r.last > radix.last {
		newRadix = &RadixHeap[V, P]{
			buckets: cloneBuckets(r.buckets),
			size:    r.size,
			last:    r.last,
		}
		r.buckets = radix.buckets
		r.last = radix.last
		r.size = radix.size
	} else {
		newRadix = radix
	}
	for i := range newRadix.buckets {
		for _, pair := range newRadix.buckets[i] {
			r.push(pair.value, pair.priority)
		}
	}
}

// getBucketIndex calculates which bucket index a priority 'num' belongs to,
// relative to 'last'.
// Returns floor(log2(num XOR last)) + 1. If num equals last, callers should
// put it in bucket 0.
func getBucketIndex[T constraints.Unsigned](num T, last T) int {
	bitDiff := num ^ last
	i := math.Floor(math.Log2(float64(bitDiff))) + 1
	return int(i)
}

// bucketInsert puts a HeapNode into the correct bucket based on its priority
// and 'last'.
// If priority equals last, it goes into bucket 0; otherwise, getBucketIndex
// determines the bucket index.
func bucketInsert[V any, P constraints.Unsigned](pair HeapNode[V, P], last P, buckets [][]HeapNode[V, P]) {
	if pair.priority == last {
		buckets[0] = append(buckets[0], pair)
	} else {
		i := getBucketIndex(pair.priority, last)
		buckets[i] = append(buckets[i], pair)
	}
}

// minFromSlice returns the HeapNode with the minimum priority from a non-empty slice.
// The caller must ensure the slice is not empty.
func minFromSlice[V any, P constraints.Unsigned, T Node[V, P]](pairs []T) T {
	minPair := pairs[0]
	for _, pair := range pairs {
		if pair.Priority() < minPair.Priority() {
			minPair = pair
		}
	}
	return minPair
}
