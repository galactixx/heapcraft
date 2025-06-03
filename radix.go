package heapcraft

import (
	"errors"
	"fmt"
	"math"
	"reflect"

	"golang.org/x/exp/constraints"
)

// NewSimpleRadixHeap constructs a SimpleRadixHeap from an initial slice of RadixPair[V,P].
// It allocates buckets based on the bit-size of P, sets 'last' to the minimum priority
// if data is non-empty, and distributes all items into their appropriate buckets.
func NewSimpleRadixHeap[V any, P constraints.Unsigned](data []Pair[V, P]) SimpleRadixHeap[V, P] {
	var forTypeCheck P
	t := reflect.TypeOf(forTypeCheck)
	bits := t.Bits()
	numBuckets := bits + 1
	buckets := make([][]RadixPair[V, P], numBuckets)

	var last P
	var size int
	var curId uint = 1

	if len(data) == 0 {
		last = 0
		size = 0
	} else {
		// find the smallest priority among input items
		last = minFromSlice(data).priority
		size = len(data)

		// place each item into the correct bucket relative to 'last'
		for _, pair := range data {
			rPair := RadixPair[V, P]{ID: curId, value: pair.value, priority: pair.priority}
			bucketInsert(rPair, last, buckets)
			curId++
		}
	}

	return SimpleRadixHeap[V, P]{buckets: buckets, size: size, last: last, curId: curId}
}

// NewRadixHeap constructs a RadixHeap (with dictionary) from an initial slice of RadixPair[V,P].
// It initializes the underlying SimpleRadixHeap and populates the elements map for fast lookups.
func NewRadixHeap[V any, P constraints.Unsigned](data []Pair[V, P]) RadixHeap[V, P] {
	simpleRadixHeap := NewSimpleRadixHeap(data)
	elements := make(map[uint]*RadixPair[V, P], simpleRadixHeap.Length())
	for i := range simpleRadixHeap.buckets {
		for _, element := range simpleRadixHeap.buckets[i] {
			elements[element.ID] = &element
		}
	}
	return RadixHeap[V, P]{heap: simpleRadixHeap, elements: elements}
}

// RadixHeap wraps a SimpleRadixHeap and maintains a map of element IDs to RadixPair
// pointers.
type RadixHeap[V any, P constraints.Unsigned] struct {
	heap     SimpleRadixHeap[V, P]
	elements map[uint]*RadixPair[V, P]
}

// Contains reports whether an element with the given ID exists in the heap.
func (r *RadixHeap[V, P]) Contains(id uint) bool {
	_, exists := r.elements[id]
	return exists
}

// Push inserts a new value with the given priority into the heap, returning its RadixPair.
// It also adds the new pair to the elements map for O(1) access.
func (r *RadixHeap[V, P]) Push(value V, priority P) (*RadixPair[V, P], error) {
	newPair, err := r.heap.internalPush(value, priority)
	if err == nil {
		r.elements[newPair.ID] = newPair
	}
	return newPair, err
}

// Pop removes and returns the RadixPair with the minimum priority from the heap.
// It also deletes that pair from the elements map.
func (r *RadixHeap[V, P]) Pop() (*RadixPair[V, P], error) {
	popped, err := r.heap.Pop()
	if err == nil {
		delete(r.elements, popped.ID)
	}
	return popped, err
}

// GetElement retrieves the RadixPair pointer for the given ID, or returns an error
// if not found.
func (r *RadixHeap[V, P]) GetElement(id uint) (*RadixPair[V, P], error) {
	if _, exists := r.elements[id]; !exists {
		return nil, fmt.Errorf("id %d does not link to a valid element", id)
	}

	element := r.elements[id]
	return element, nil
}

// GetValue returns the value associated with the given ID, or an error if the ID is invalid.
func (r *RadixHeap[V, P]) GetValue(id uint) (*V, error) {
	element, err := r.GetElement(id)
	if err == nil {
		return &element.value, err
	}
	return nil, err
}

// GetPriority returns the priority associated with the given ID, or an error if the ID is invalid.
func (r *RadixHeap[V, P]) GetPriority(id uint) (*P, error) {
	element, err := r.GetElement(id)
	if err == nil {
		return &element.priority, err
	}
	return nil, err
}

// UpdatePriority changes the priority of the element with the given ID by removing it and re-pushing.
// Returns the new RadixPair or an error if the ID does not exist.
func (r *RadixHeap[V, P]) UpdatePriority(id uint, priority P) (*RadixPair[V, P], error) {
	var err error
	var removed, newPair *RadixPair[V, P]

	removed, err = r.GetElement(id)

	if err == nil {
		newPair, err = r.Push(removed.value, priority)

		if err == nil {
			r.Remove(id)
			return newPair, err
		}
	}

	return nil, err
}

// Remove deletes the RadixPair with the given ID from the heap and map, and returns it.
// Does not rebalance buckets; updates heap size directly.
func (r *RadixHeap[V, P]) Remove(id uint) (*RadixPair[V, P], error) {
	if _, exists := r.elements[id]; !exists {
		return nil, fmt.Errorf("id %d does not link to a valid element", id)
	}

	removed := r.elements[id]
	delete(r.elements, id)
	r.heap.size--
	return removed, nil
}

// Clear resets the heap (and its underlying SimpleRadixHeap) to an empty state.
func (r *RadixHeap[V, P]) Clear() { r.heap.Clear() }

// Peek returns one RadixPair with the minimum priority without removing it
// (calls underlying heap).
func (r *RadixHeap[V, P]) Peek() *RadixPair[V, P] { return r.heap.Peek() }

// Rebalance triggers a rebalance of the underlying SimpleRadixHeap if bucket
// 0 is empty.
func (r *RadixHeap[V, P]) Rebalance() error { return r.heap.Rebalance() }

// Length returns the total number of elements currently in the heap.
func (r *RadixHeap[V, P]) Length() int { return r.heap.Length() }

// IsEmpty reports whether the heap contains zero elements.
func (r *RadixHeap[V, P]) IsEmpty() bool { return r.heap.IsEmpty() }

type ElementPair[V any, P constraints.Unsigned] interface {
	Value() V
	Priority() P
}

type Pair[V any, P constraints.Unsigned] struct {
	value    V
	priority P
}

func (p Pair[V, P]) Value() V    { return p.value }
func (p Pair[V, P]) Priority() P { return p.priority }

// RadixPair associates a generic Value with an unsigned Priority.
type RadixPair[V any, P constraints.Unsigned] struct {
	ID       uint
	value    V
	priority P
}

func (r RadixPair[V, P]) Value() V    { return r.value }
func (r RadixPair[V, P]) Priority() P { return r.priority }

// SimpleRadixHeap is a monotone priority queue over unsigned priorities.
//   - buckets: slices of RadixPair, each bucket holds items whose priorities
//     fall into a specific range relative to 'last'.
//   - size: total number of items stored.
//   - last: the last-extracted minimum priority.
type SimpleRadixHeap[V any, P constraints.Unsigned] struct {
	buckets [][]RadixPair[V, P]
	size    int
	last    P
	curId   uint
}

// Clone returns a shallow copy of the heap structure (buckets slice),
// but does not deep-copy the contents of each bucket.
func (r SimpleRadixHeap[V, P]) Clone() SimpleRadixHeap[V, P] {
	newBuckets := make([][]RadixPair[V, P], len(r.buckets))
	for i := range r.buckets {
		newBuckets[i] = make([]RadixPair[V, P], 0)
	}
	copy(newBuckets, r.buckets)
	return SimpleRadixHeap[V, P]{buckets: newBuckets, size: r.size, last: r.last}
}

// Push inserts a new (value, priority) into the heap. It returns an error
// if priority < r.last to enforce the monotonicity invariant.
// Otherwise, it places the item into the correct bucket and increments size.
func (r *SimpleRadixHeap[V, P]) Push(value V, priority P) (*RadixPair[V, P], error) {
	return r.internalPush(value, priority)
}

// internalPush is an unexported helper that creates a RadixPair and buckets it.
// Returns an error if priority < r.last.
func (r *SimpleRadixHeap[V, P]) internalPush(value V, priority P) (*RadixPair[V, P], error) {
	if priority < r.last {
		return nil, fmt.Errorf("insertion of a priority less than last popped")
	}
	newPair := RadixPair[V, P]{value: value, priority: priority}
	bucketInsert(newPair, r.last, r.buckets)
	r.size++
	return &newPair, nil
}

// Pop removes and returns the RadixPair with the minimum priority. If bucket 0
// is non-empty, it pops directly. Otherwise, it calls rebalanceBuckets to refill
// bucket 0 from the next non-empty bucket, then returns the new minimum.
// Returns an error if the heap is empty.
func (r *SimpleRadixHeap[V, P]) Pop() (*RadixPair[V, P], error) {
	if r.IsEmpty() {
		return nil, errors.New("heap has no elements and is empty")
	}

	// if bucket 0 has items, pop the first one
	if len(r.buckets[0]) > 0 {
		minPair := r.buckets[0][0]
		r.buckets[0] = r.buckets[0][:len(r.buckets[0])-1]
		r.size--
		return &minPair, nil
	}

	// rebalance bucket 0 from the next non-empty bucket
	r.rebalanceBuckets()
	removed := r.buckets[0][0]
	r.buckets[0] = r.buckets[0][:len(r.buckets[0])-1]
	r.size--
	return &removed, nil
}

// Clear resets the heap to an empty state by reinitializing buckets, setting size to 0,
// and resetting 'last' to zero value.
func (r *SimpleRadixHeap[V, P]) Clear() {
	r.buckets = make([][]RadixPair[V, P], len(r.buckets))
	r.size = 0
	r.last = 0
}

// rebalanceBuckets finds the next non-empty bucket i>0, sets 'last' to the minimum
// priority in that bucket, and re-inserts all items from bucket i into their new buckets
// based on the updated 'last'. Afterward, bucket i is cleared.
func (r *SimpleRadixHeap[V, P]) rebalanceBuckets() {
	for i := 1; i < len(r.buckets); i++ {
		if len(r.buckets[i]) > 0 {
			r.last = minFromSlice(r.buckets[i]).priority
			for _, pair := range r.buckets[i] {
				bucketInsert(pair, r.last, r.buckets)
			}
			r.buckets[i] = make([]RadixPair[V, P], 0)
			return
		}
	}
}

// Rebalance ensures bucket 0 is filled if needed. It returns an error if the heap
// is empty or if no rebalancing was required (i.e., bucket 0 was already non-empty).
func (r *SimpleRadixHeap[V, P]) Rebalance() error {
	if r.IsEmpty() {
		return errors.New("heap has no elements and is empty")
	}
	if len(r.buckets[0]) == 0 {
		r.rebalanceBuckets()
		return nil
	}
	return errors.New("no rebalancing needed")
}

// Length returns the total number of items in the heap.
func (r SimpleRadixHeap[V, P]) Length() int { return r.size }

// IsEmpty reports whether the heap currently contains zero items.
func (r SimpleRadixHeap[V, P]) IsEmpty() bool { return r.Length() == 0 }

// Peek returns one RadixPair with the minimum priority without removing it.
// If bucket 0 is non-empty, it returns that item. Otherwise, it scans the first
// non-empty bucket for its minimum. Returns nil if the heap is empty.
func (r *SimpleRadixHeap[V, P]) Peek() *RadixPair[V, P] {
	if r.IsEmpty() {
		return nil
	}
	if len(r.buckets[0]) > 0 {
		return &r.buckets[0][0]
	}
	for i := 1; i < len(r.buckets); i++ {
		if len(r.buckets[i]) > 0 {
			minPair := minFromSlice(r.buckets[i])
			return &minPair
		}
	}
	return nil
}

// Merge merges another SimpleRadixHeap into this one. It chooses the smaller 'last'
// baseline, adopts that heap’s buckets and 'last', then pushes all items from the other
// heap to maintain monotonicity.
func (r *SimpleRadixHeap[V, P]) Merge(radix SimpleRadixHeap[V, P]) {
	var newRadix SimpleRadixHeap[V, P]
	if r.last > radix.last {
		newRadix = r.Clone()
		r.buckets = radix.buckets
		r.last = radix.last
		r.size = radix.size
	} else {
		newRadix = radix
	}
	for i := 0; i < len(newRadix.buckets); i++ {
		for _, pair := range newRadix.buckets[i] {
			r.Push(pair.value, pair.priority)
		}
	}
}

// getBucketIndex computes the appropriate bucket index for a given priority 'num'
// relative to 'last'. It returns floor(log2(num XOR last)) + 1. If num == last,
// callers should use bucket 0 instead.
func getBucketIndex[T constraints.Unsigned](num T, last T) int {
	bitDiff := num ^ last
	i := math.Floor(math.Log2(float64(bitDiff))) + 1
	return int(i)
}

// bucketInsert places a RadixPair into the correct bucket based on its priority
// and the current 'last'. If priority == last, it goes into bucket 0; otherwise,
// getBucketIndex is used.
func bucketInsert[V any, P constraints.Unsigned](pair RadixPair[V, P], last P, buckets [][]RadixPair[V, P]) {
	if pair.priority == last {
		buckets[0] = append(buckets[0], pair)
	} else {
		i := getBucketIndex(pair.priority, last)
		buckets[i] = append(buckets[i], pair)
	}
}

// minFromSlice returns the RadixPair with the smallest Priority from a non-empty slice.
func minFromSlice[E ElementPair[V, P], V any, P constraints.Unsigned](pairs []E) E {
	minPair := pairs[0]
	for _, pair := range pairs {
		if pair.Priority() < minPair.Priority() {
			minPair = pair
		}
	}
	return minPair
}
