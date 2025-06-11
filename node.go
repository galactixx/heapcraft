package heapcraft

// SimpleNode represents a node with a value and priority.
type SimpleNode[V any, P any] interface {
	Value() V
	Priority() P
}

// Node represents a node with an ID, value, and priority.
type Node[V any, P any] interface {
	ID() uint
	Value() V
	Priority() P
}

// HeapNode binds a value to its priority for heap operations.
type HeapNode[V any, P any] struct {
	value    V
	priority P
}

// Value returns the value stored in the HeapNode.
func (b HeapNode[V, P]) Value() V { return b.value }

// Priority returns the priority of the HeapNode.
func (b HeapNode[V, P]) Priority() P { return b.priority }

// CreateHeapPair constructs a new HeapNode from the given value and priority.
func CreateHeapPair[V any, P any](value V, priority P) HeapNode[V, P] {
	return HeapNode[V, P]{value: value, priority: priority}
}

// CreateHeapPairPtr constructs a new *HeapNode from the given value and priority.
// This function is specifically for use with Leftist, Skew, and Pairing heaps
// that expect pointer slices.
func CreateHeapPairPtr[V any, P any](value V, priority P) *HeapNode[V, P] {
	return &HeapNode[V, P]{value: value, priority: priority}
}

// RadixPair binds a generic value to an unsigned priority.
type RadixPair[V any, P any] struct {
	value    V
	priority P
}

// Value returns the value stored in the RadixPair.
func (r RadixPair[V, P]) Value() V { return r.value }

// Priority returns the priority of the RadixPair.
func (r RadixPair[V, P]) Priority() P { return r.priority }

// CreateRadixPair constructs a new RadixPair from the given value and priority.
func CreateRadixPair[V any, P any](value V, priority P) *RadixPair[V, P] {
	return &RadixPair[V, P]{value: value, priority: priority}
}
