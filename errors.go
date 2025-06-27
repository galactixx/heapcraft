package heapcraft

import "errors"

var (
	// ErrCallbackNotFound is returned when attempting to deregister a callback that
	// doesn't exist.
	ErrCallbackNotFound = errors.New("callback not found")

	// ErrHeapEmpty is returned when attempting to access elements from an empty heap.
	ErrHeapEmpty = errors.New("the heap is empty and contains no elements")

	// ErrIndexOutOfBounds is returned when attempting to access an index that is outside
	// the valid range of the heap.
	ErrIndexOutOfBounds = errors.New("index out of bounds")

	// ErrPriorityLessThanLast is returned when attempting to insert a priority that is
	// less than the last extracted priority, which would violate the monotonic property
	// of the radix heap.
	ErrPriorityLessThanLast = errors.New("insertion of a priority less than last popped")

	// ErrNoRebalancingNeeded is returned when attempting to rebalance a radix heap
	// that doesn't need rebalancing (bucket 0 already contains elements).
	ErrNoRebalancingNeeded = errors.New("no rebalancing needed")

	// ErrNodeNotFound is returned when attempting to access a node with an ID that
	// does not exist in the pairing heap.
	ErrNodeNotFound = errors.New("id does not link to existing node")
)
