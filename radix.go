package heapcraft

import (
	"errors"
	"fmt"
	"math"
	"reflect"

	"golang.org/x/exp/constraints"
)

// NewRadixHeap creates a RadixHeap from a given slice of *RadixPair[V,P].
// It determines the number of buckets from the bit-length of P, initializes
// 'last' to the smallest priority if data is present, and assigns each element
// into its corresponding bucket.
func NewRadixHeap[V any, P constraints.Unsigned](data []*RadixPair[V, P]) RadixHeap[V, P] {
	var forTypeCheck P
	t := reflect.TypeOf(forTypeCheck)
	bits := t.Bits()
	numBuckets := bits + 1
	buckets := make([][]*RadixPair[V, P], numBuckets)

	var last P
	var size int

	if len(data) == 0 {
		last = 0
		size = 0
	} else {
		// Determine the minimum priority among the input items
		last = minFromSlice(data).priority
		size = len(data)

		// Insert each item into the appropriate bucket relative to 'last'
		for _, pair := range data {
			rPair := &RadixPair[V, P]{value: pair.value, priority: pair.priority}
			bucketInsert(rPair, last, buckets)
		}
	}

	return RadixHeap[V, P]{buckets: buckets, size: size, last: last}
}

// CreateRadixPair constructs a new *RadixPair with the given value and priority.
func CreateRadixPair[V any, P constraints.Unsigned](value V, priority P) *RadixPair[V, P] {
	return &RadixPair[V, P]{value: value, priority: priority}
}

// RadixPair binds a generic value to an unsigned priority.
type RadixPair[V any, P constraints.Unsigned] struct {
	value    V
	priority P
}

func (r RadixPair[V, P]) Value() V    { return r.value }
func (r RadixPair[V, P]) Priority() P { return r.priority }

// RadixHeap implements a monotonic priority queue over unsigned priorities.
//   - buckets: array of slices of *RadixPair, each holding items whose priorities
//     fall within a range defined by 'last'.
//   - size: the count of elements in the heap.
//   - last: the most recently extracted minimum priority.
type RadixHeap[V any, P constraints.Unsigned] struct {
	buckets [][]*RadixPair[V, P]
	size    int
	last    P
}

// Clone produces a shallow copy of the heap's structure (the buckets slice),
// but does not duplicate the actual elements within each bucket.
func (r RadixHeap[V, P]) Clone() RadixHeap[V, P] {
	newBuckets := make([][]*RadixPair[V, P], len(r.buckets))
	for i := range r.buckets {
		newBuckets[i] = make([]*RadixPair[V, P], 0)
	}
	copy(newBuckets, r.buckets)
	return RadixHeap[V, P]{buckets: newBuckets, size: r.size, last: r.last}
}

// Push adds a new value and priority pair into the heap. If priority is less than r.last,
// it returns an error to preserve the monotonic property. Otherwise, it puts the item
// into the appropriate bucket and increments the size.
func (r *RadixHeap[V, P]) Push(value V, priority P) (*RadixPair[V, P], error) {
	return r.internalPush(value, priority)
}

// internalPush is an unexported helper that forms a RadixPair and places it into its bucket.
// It enforces the condition that priority must not be less than r.last.
func (r *RadixHeap[V, P]) internalPush(value V, priority P) (*RadixPair[V, P], error) {
	if priority < r.last {
		return nil, fmt.Errorf("insertion of a priority less than last popped")
	}
	newPair := &RadixPair[V, P]{value: value, priority: priority}
	bucketInsert(newPair, r.last, r.buckets)
	r.size++
	return newPair, nil
}

// popMinElement removes and returns the first element in bucket 0.
// It also decreases the total size. The caller must ensure bucket 0 is not empty.
func (r *RadixHeap[V, P]) popMinElement() *RadixPair[V, P] {
	minPair := r.buckets[0][0]
	r.buckets[0] = r.buckets[0][1:]
	r.size--
	return minPair
}

// Pop extracts and returns the RadixPair with the smallest priority. If bucket
// 0 contains items, it takes from there directly. Otherwise, it calls
// rebalanceBuckets to refill bucket 0 from the next non-empty bucket, then
// returns the new minimum. Returns an error if the heap is empty.
func (r *RadixHeap[V, P]) Pop() (*RadixPair[V, P], error) {
	if r.IsEmpty() {
		return nil, errors.New("heap has no elements and is empty")
	}

	// If bucket 0 has entries, pop directly
	if len(r.buckets[0]) > 0 {
		return r.popMinElement(), nil
	}

	// Otherwise, refill bucket 0 from the next non-empty bucket
	r.rebalanceBuckets()
	return r.popMinElement(), nil
}

// Clear reinitializes the heap by creating fresh buckets, resetting size to zero,
// and setting 'last' back to its zero value.
func (r *RadixHeap[V, P]) Clear() {
	r.buckets = make([][]*RadixPair[V, P], len(r.buckets))
	r.size = 0
	r.last = 0
}

// rebalanceBuckets locates the next bucket with elements (i > 0), updates 'last'
// to the smallest priority found there, and reinserts all items from that bucket
// into new buckets based on the updated 'last'. Afterward, it empties that bucket.
func (r *RadixHeap[V, P]) rebalanceBuckets() {
	for i := 1; i < len(r.buckets); i++ {
		if len(r.buckets[i]) > 0 {
			r.last = minFromSlice(r.buckets[i]).priority
			for _, pair := range r.buckets[i] {
				bucketInsert(pair, r.last, r.buckets)
			}
			r.buckets[i] = make([]*RadixPair[V, P], 0)
			return
		}
	}
}

// Rebalance fills bucket 0 if it is empty. It returns an error if the heap is
// empty, or if bucket 0 already contains elements (no action was needed).
func (r *RadixHeap[V, P]) Rebalance() error {
	if r.IsEmpty() {
		return errors.New("heap has no elements and is empty")
	}
	if len(r.buckets[0]) == 0 {
		r.rebalanceBuckets()
		return nil
	}
	return errors.New("no rebalancing needed")
}

// Length returns the number of items currently stored in the heap.
func (r RadixHeap[V, P]) Length() int { return r.size }

// IsEmpty returns true if the heap contains no items.
func (r RadixHeap[V, P]) IsEmpty() bool { return r.Length() == 0 }

// Peek returns a *RadixPair with the smallest priority without removing it.
// If bucket 0 has elements, it returns the first one. Otherwise, it finds
// the minimal element in the next non-empty bucket. Returns nil if the heap is empty.
func (r *RadixHeap[V, P]) Peek() *RadixPair[V, P] {
	if r.IsEmpty() {
		return nil
	}
	if len(r.buckets[0]) > 0 {
		return r.buckets[0][0]
	}
	for i := 1; i < len(r.buckets); i++ {
		if len(r.buckets[i]) > 0 {
			minPair := minFromSlice(r.buckets[i])
			return minPair
		}
	}
	return nil
}

// Merge integrates another RadixHeap into this one. It selects the heap
// with the smaller 'last' as the new baseline, adopts its buckets and 'last',
// then reinserts all items from the other heap to preserve monotonicity.
func (r *RadixHeap[V, P]) Merge(radix RadixHeap[V, P]) {
	var newRadix RadixHeap[V, P]
	if r.last > radix.last {
		newRadix = r.Clone()
		r.buckets = radix.buckets
		r.last = radix.last
		r.size = radix.size
	} else {
		newRadix = radix
	}
	for i := range newRadix.buckets {
		for _, pair := range newRadix.buckets[i] {
			r.Push(pair.value, pair.priority)
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

// bucketInsert puts a *RadixPair into the correct bucket based on its priority
// and 'last'.
// If priority equals last, it goes into bucket 0; otherwise, getBucketIndex
// determines the bucket index.
func bucketInsert[V any, P constraints.Unsigned](pair *RadixPair[V, P], last P, buckets [][]*RadixPair[V, P]) {
	if pair.priority == last {
		buckets[0] = append(buckets[0], pair)
	} else {
		i := getBucketIndex(pair.priority, last)
		buckets[i] = append(buckets[i], pair)
	}
}

// minFromSlice returns the *RadixPair with the lowest priority from a non-empty slice.
func minFromSlice[V any, P constraints.Unsigned](pairs []*RadixPair[V, P]) *RadixPair[V, P] {
	minPair := pairs[0]
	for _, pair := range pairs {
		if pair.Priority() < minPair.Priority() {
			minPair = pair
		}
	}
	return minPair
}
