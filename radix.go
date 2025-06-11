package heapcraft

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"sync"

	"golang.org/x/exp/constraints"
)

// cloneBuckets creates a shallow copy of the buckets array, copying each bucket slice.
func cloneBuckets[V any, P constraints.Unsigned](buckets [][]*HeapNode[V, P]) [][]*HeapNode[V, P] {
	newBuckets := make([][]*HeapNode[V, P], len(buckets))
	for i := range buckets {
		newBuckets[i] = make([]*HeapNode[V, P], 0)
	}
	copy(newBuckets, buckets)
	return newBuckets
}

// NewRadixHeap creates a RadixHeap from a given slice of *HeapNode[V,P].
// It determines the number of buckets from the bit-length of P, initializes
// 'last' to the minimum priority if data is present, and assigns each element
// into its corresponding bucket.
func NewRadixHeap[V any, P constraints.Unsigned](data []*HeapNode[V, P]) *RadixHeap[V, P] {
	var forTypeCheck P
	t := reflect.TypeOf(forTypeCheck)
	bits := t.Bits()
	numBuckets := bits + 1
	buckets := make([][]*HeapNode[V, P], numBuckets)

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
			rPair := &HeapNode[V, P]{value: pair.value, priority: pair.priority}
			bucketInsert(rPair, last, buckets)
		}
	}

	return &RadixHeap[V, P]{buckets: buckets, size: size, last: last}
}

// RadixHeap implements a monotonic priority queue over unsigned priorities.
//   - buckets: array of slices of *HeapNode, each holding items whose priorities
//     fall within a range defined by 'last'.
//   - size: the count of elements in the heap.
//   - last: the most recently extracted minimum priority.
type RadixHeap[V any, P constraints.Unsigned] struct {
	buckets [][]*HeapNode[V, P]
	size    int
	last    P
	lock    sync.RWMutex
}

// Clone produces a shallow copy of the heap's structure (the buckets slice),
// copying the bucket slices but sharing the same elements within each bucket.
func (r *RadixHeap[V, P]) Clone() *RadixHeap[V, P] {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return &RadixHeap[V, P]{
		buckets: cloneBuckets(r.buckets), size: r.size, last: r.last,
	}
}

// Push adds a new value and priority pair into the heap. If priority is less than r.last,
// it returns an error to preserve the monotonic property. Otherwise, it puts the item
// into the appropriate bucket and increments the size.
func (r *RadixHeap[V, P]) Push(value V, priority P) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.push(value, priority)
}

// push is an unexported helper that forms a HeapNode and places it into its bucket.
// It enforces the condition that priority must not be less than r.last.
func (r *RadixHeap[V, P]) push(value V, priority P) error {
	if priority < r.last {
		return fmt.Errorf("insertion of a priority less than last popped")
	}
	newPair := HeapNode[V, P]{value: value, priority: priority}
	bucketInsert(&newPair, r.last, r.buckets)
	r.size++
	return nil
}

// getMin removes and returns the first element from bucket 0.
// It also decreases the total size. The caller must ensure bucket 0 is not empty.
func (r *RadixHeap[V, P]) getMin() HeapNode[V, P] {
	minPair := r.buckets[0][0]
	r.buckets[0] = r.buckets[0][1:]
	r.size--
	return *minPair
}

// pop removes and returns the first element in bucket 0.
// It also decreases the total size. The caller must ensure bucket 0 is not empty.
func (r *RadixHeap[V, P]) pop() (HeapNode[V, P], error) {
	if r.size == 0 {
		var zero HeapNode[V, P]
		return zero, errors.New("the heap is empty and contains no elements")
	}

	// If bucket 0 has entries, pop directly
	if len(r.buckets[0]) > 0 {
		return r.getMin(), nil
	}

	// Otherwise, refill bucket 0 from the next non-empty bucket
	r.rebalance()
	return r.getMin(), nil
}

// Pop extracts and returns the HeapNode with the minimum priority. If bucket
// 0 contains items, it takes from there directly. Otherwise, it calls
// rebalance to refill bucket 0 from the next non-empty bucket, then
// returns the new minimum. Returns an error if the heap is empty.
func (r *RadixHeap[V, P]) Pop() (SimpleNode[V, P], error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.pop()
}

// Clear reinitializes the heap by creating fresh buckets, resetting size to zero,
// and setting 'last' back to its zero value.
func (r *RadixHeap[V, P]) Clear() {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.buckets = make([][]*HeapNode[V, P], len(r.buckets))
	r.size = 0
	r.last = 0
}

// rebalance locates the next bucket with elements (i > 0), updates 'last'
// to the smallest priority found there, and reinserts all items from that bucket
// into new buckets based on the updated 'last'. Afterward, it empties that bucket.
func (r *RadixHeap[V, P]) rebalance() {
	for i := 1; i < len(r.buckets); i++ {
		if len(r.buckets[i]) > 0 {
			r.last = minFromSlice(r.buckets[i]).priority
			for _, pair := range r.buckets[i] {
				bucketInsert(pair, r.last, r.buckets)
			}
			r.buckets[i] = make([]*HeapNode[V, P], 0)
			return
		}
	}
}

// Rebalance fills bucket 0 if it is empty. It returns an error if the heap is
// empty, or if bucket 0 already contains elements (no action was needed).
func (r *RadixHeap[V, P]) Rebalance() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.size == 0 {
		return errors.New("the heap is empty and contains no elements")
	}
	if len(r.buckets[0]) == 0 {
		r.rebalance()
		return nil
	}
	return errors.New("no rebalancing needed")
}

// Length returns the number of items currently stored in the heap.
func (r *RadixHeap[V, P]) Length() int {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.size
}

// IsEmpty returns true if the heap contains no items.
func (r *RadixHeap[V, P]) IsEmpty() bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.size == 0
}

// Peek returns a HeapNode with the minimum priority without removing it.
// If bucket 0 has elements, it returns the first one. Otherwise, it finds
// the minimum element in the next non-empty bucket. Returns an error if the heap is empty.
func (r *RadixHeap[V, P]) Peek() (SimpleNode[V, P], error) {
	var zero HeapNode[V, P]
	r.lock.RLock()
	defer r.lock.RUnlock()
	if r.size == 0 {
		return zero, errors.New("the heap is empty and contains no elements")
	}
	if len(r.buckets[0]) > 0 {
		return *r.buckets[0][0], nil
	}
	for i := 1; i < len(r.buckets); i++ {
		if len(r.buckets[i]) > 0 {
			minPair := minFromSlice(r.buckets[i])
			return *minPair, nil
		}
	}
	return zero, errors.New("the heap is empty and contains no elements")
}

// Merge integrates another RadixHeap into this one. It selects the heap
// with the smaller 'last' as the new baseline, adopts its buckets and 'last',
// then reinserts all items from the other heap to preserve monotonicity.
func (r *RadixHeap[V, P]) Merge(radix *RadixHeap[V, P]) {
	r.lock.Lock()
	defer r.lock.Unlock()
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
// It returns floor(log2(num XOR last)) + 1. If num equals last, callers should
// put it in bucket 0.
func getBucketIndex[T constraints.Unsigned](num T, last T) int {
	bitDiff := num ^ last
	i := math.Floor(math.Log2(float64(bitDiff))) + 1
	return int(i)
}

// bucketInsert puts a *HeapNode into the correct bucket based on its priority
// and 'last'.
// If priority equals last, it goes into bucket 0; otherwise, getBucketIndex
// determines the bucket index.
func bucketInsert[V any, P constraints.Unsigned](pair *HeapNode[V, P], last P, buckets [][]*HeapNode[V, P]) {
	if pair.priority == last {
		buckets[0] = append(buckets[0], pair)
	} else {
		i := getBucketIndex(pair.priority, last)
		buckets[i] = append(buckets[i], pair)
	}
}

// minFromSlice returns the *HeapNode with the minimum priority from a non-empty slice.
func minFromSlice[V any, P constraints.Unsigned](pairs []*HeapNode[V, P]) *HeapNode[V, P] {
	minPair := pairs[0]
	for _, pair := range pairs {
		if pair.priority < minPair.priority {
			minPair = pair
		}
	}
	return minPair
}
