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
)
