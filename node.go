package heapcraft

// HeapNode binds a value to its priority for heap operations.
type HeapNode[V any, P any] struct {
	value    V
	priority P
}

// CreateHeapNode constructs a new HeapNode from the given value and priority.
func CreateHeapNode[V any, P any](value V, priority P) HeapNode[V, P] {
	return HeapNode[V, P]{value: value, priority: priority}
}
