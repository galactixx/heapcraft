package heapcraft

import (
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/mohae/deepcopy"
	"golang.org/x/exp/constraints"
)

// getBucketIndex computes which bucket an element belongs to by
// taking the bitwise difference between num and last, then returning
// ⌊log₂(bitDiff)⌋ + 1. This determines the bucket index for insertion.
func getBucketIndex[T constraints.Unsigned](num T, last T) int {
	bitDiff := num ^ last
	i := math.Floor(math.Log2(float64(bitDiff))) + 1
	return int(i)
}

// bucketInsert places num into the appropriate bucket based on last.
// If num equals last, it goes into bucket 0; otherwise, it uses getBucketIndex.
func bucketInsert[T constraints.Unsigned](num T, last T, buckets [][]T) {
	if num == last {
		buckets[0] = append(buckets[0], num)
	} else {
		i := getBucketIndex(num, last)
		buckets[i] = append(buckets[i], num)
	}
}

// minFromSlice returns the smallest element in nums by scanning sequentially.
func minFromSlice[T constraints.Unsigned](nums []T) T {
	minNum := nums[0]
	for _, num := range nums {
		minNum = min(num, minNum)
	}
	return minNum
}

// NewRadixHeap initializes a RadixHeap from an existing slice.
// It determines the number of buckets based on the bit-size of T,
// finds the initial last (minimum) if data is nonempty, and distributes
// all elements into their correct buckets.
func NewRadixHeap[T constraints.Unsigned](data []T) RadixHeap[T] {
	var forTypeCheck T
	t := reflect.TypeOf(forTypeCheck)
	bits := t.Bits()
	numBuckets := bits + 1
	buckets := make([][]T, numBuckets)

	var last T
	var size int

	if len(data) == 0 {
		last = 0
		size = 0
	} else {
		last = minFromSlice(data)
		size = len(data)

		// insert all elements into their buckets relative to last
		for _, number := range data {
			bucketInsert(number, last, buckets)
		}
	}

	return RadixHeap[T]{buckets: buckets, size: size, last: last}
}

// RadixHeap represents a monotone priority queue over unsigned integers.
// - buckets: a slice of slices, each bucket holds keys in a certain range
// - size: total number of elements across all buckets
// - last: the last-extracted minimum key
type RadixHeap[T constraints.Unsigned] struct {
	buckets [][]T
	size    int
	last    T
}

// DeepClone performs a deep copy of the RadixHeap, including all buckets.
// It uses a third-party deepcopy to duplicate each bucket slice.
func (r RadixHeap[T]) DeepClone() RadixHeap[T] {
	newBuckets := make([][]T, len(r.buckets))
	for i, element := range r.buckets {
		bucketCopy := deepcopy.Copy(element)
		newBuckets[i] = bucketCopy.([]T)
	}
	return RadixHeap[T]{buckets: newBuckets, size: r.size, last: r.last}
}

// Clone makes a shallow copy of the RadixHeap’s buckets slice structure.
// Note: this does not deep-copy individual bucket contents (unlike DeepClone).
func (r RadixHeap[T]) Clone() RadixHeap[T] {
	newBuckets := make([][]T, len(r.buckets))
	for i := range r.buckets {
		newBuckets[i] = make([]T, 0)
	}

	copy(newBuckets, r.buckets)
	return RadixHeap[T]{buckets: newBuckets, size: r.size, last: r.last}
}

// Push inserts a new number into the RadixHeap.
// It returns an error if number < last to enforce monotonicity.
// Otherwise, it places number into the correct bucket and increments size.
func (r *RadixHeap[T]) Push(number T) error {
	if number < r.last {
		return fmt.Errorf("insertion of a number less than last popped")
	}

	bucketInsert(number, r.last, r.buckets)
	r.size++
	return nil
}

// Pop removes and returns the current minimum key.
// If bucket 0 is nonempty, it pops there directly. Otherwise, it finds
// the first nonempty bucket, rebuilds bucket 0 by re-bucketing its contents
// around a new last, and then returns that new minimum. Errors if empty.
func (r *RadixHeap[T]) Pop() (*T, error) {
	// if bucket 0 has elements equal to last, pop one immediately
	if len(r.buckets[0]) > 0 {
		minNum := r.buckets[0][0]
		r.buckets[0] = r.buckets[0][:len(r.buckets[0])-1]
		r.size--
		return &minNum, nil
	}

	// otherwise, find the next nonempty bucket to rebuild bucket 0
	for i := 1; i < len(r.buckets); i++ {
		if len(r.buckets[i]) > 0 {
			r.last = minFromSlice(r.buckets[i])
			for _, number := range r.buckets[i] {
				bucketInsert(number, r.last, r.buckets)
			}
			r.buckets[i] = make([]T, 0)
			r.buckets[0] = r.buckets[0][:len(r.buckets[0])-1]
			r.size--
			return &r.last, nil
		}
	}

	// no buckets had any elements; the heap is empty
	return nil, errors.New("heap has no element and is empty")
}

// Clear resets the RadixHeap to an empty state by reinitializing all buckets,
// setting size to 0, and resetting last to zero value.
func (r *RadixHeap[T]) Clear() {
	r.buckets = make([][]T, len(r.buckets))
	r.size = 0
	r.last = 0
}

// Length returns the number of elements currently in the heap.
func (r RadixHeap[T]) Length() int { return r.size }

// IsEmpty reports whether the heap contains zero elements.
func (r RadixHeap[T]) IsEmpty() bool { return r.Length() == 0 }

// Peek returns the next minimum without removing it.
// If bucket 0 is nonempty, that minimum equals last.
// Otherwise, it scans the first nonempty bucket for its minimum.
// Returns nil if the heap is empty.
func (r *RadixHeap[T]) Peek() *T {
	if r.IsEmpty() {
		return nil
	}

	if len(r.buckets[0]) > 0 {
		return &r.last
	}

	for i := 1; i < len(r.buckets); i++ {
		if len(r.buckets[i]) > 0 {
			minNum := minFromSlice(r.buckets[i])
			return &minNum
		}
	}

	return nil
}

// Merge takes another RadixHeap and merges its contents into r.
// It chooses the smaller baseline (last) between the two, rebuckets all
// elements from the larger-last heap relative to the new last, and pushes them.
func (r *RadixHeap[T]) Merge(radix RadixHeap[T]) {
	var newRadix RadixHeap[T]

	// identify which heap has the smaller baseline
	if r.last > radix.last {
		newRadix = r.Clone()
		r.buckets = radix.buckets
		r.last = radix.last
		r.size = radix.size
	} else {
		newRadix = radix
	}

	// push all elements from newRadix into r, enforcing monotonicity
	for i := 0; i < len(newRadix.buckets); i++ {
		for _, number := range newRadix.buckets[i] {
			r.Push(number)
		}
	}
}
