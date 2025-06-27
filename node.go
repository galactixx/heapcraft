package heapcraft

// SimpleNode represents a node with a value and priority.
type SimpleNode[V any, P any] interface {
	Value() V
	Priority() P
}

// Node represents a node with an ID, value, and priority.
type Node[V any, P any] interface {
	ID() string
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

// CreateHeapNode constructs a new HeapNode from the given value and priority.
func CreateHeapNode[V any, P any](value V, priority P) HeapNode[V, P] {
	return HeapNode[V, P]{value: value, priority: priority}
}
