package heapcraft

import (
	"github.com/google/uuid"
)

// skewNode represents a node in a simple skew heap without parent pointers.
// Each node contains a value, priority, and links to its left and right children.
type skewNode[V any, P any] struct {
	value    V
	priority P
	left     *skewNode[V, P]
	right    *skewNode[V, P]
}

// Value returns the value stored in the node.
func (n *skewNode[V, P]) Value() V { return n.value }

// Priority returns the priority of the node.
func (n *skewNode[V, P]) Priority() P { return n.priority }

// skewHeapNode represents a node in a skew heap with parent pointers and ID tracking.
// Each node contains a value, priority, unique ID, and links to its parent and children.
type skewHeapNode[V any, P any] struct {
	id       string
	value    V
	priority P
	parent   *skewHeapNode[V, P]
	left     *skewHeapNode[V, P]
	right    *skewHeapNode[V, P]
}

// Value returns the value stored in the node.
func (n *skewHeapNode[V, P]) Value() V { return n.value }

// Priority returns the priority of the node.
func (n *skewHeapNode[V, P]) Priority() P { return n.priority }

// NewSkewHeap creates a new skew heap from the given data slice.
// Each element is inserted individually using the provided comparison function
// to determine heap order (min or max). Returns an empty heap if the input
// slice is empty.
func NewSkewHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *SkewHeap[V, P] {
	pool := newPool(usePool, func() *skewHeapNode[V, P] {
		return &skewHeapNode[V, P]{}
	})
	elements := make(map[string]*skewHeapNode[V, P], len(data))
	heap := SkewHeap[V, P]{cmp: cmp, size: 0, elements: elements, pool: pool}
	if len(data) == 0 {
		return &heap
	}

	for i := range data {
		heap.Push(data[i].value, data[i].priority)
	}
	return &heap
}

// SkewHeap implements a skew heap with parent pointers and element tracking.
// It maintains a map of node IDs to nodes for O(1) element access and updates.
// The heap can be either a min-heap or max-heap depending on the comparison function.
type SkewHeap[V any, P any] struct {
	root     *skewHeapNode[V, P]
	cmp      func(a, b P) bool
	size     int
	elements map[string]*skewHeapNode[V, P]
	pool     pool[*skewHeapNode[V, P]]
}

// Clone creates a deep copy of the heap structure and nodes. If values or
// priorities are reference types, those reference values are shared between the
// original and cloned heaps.
func (s *SkewHeap[V, P]) Clone() *SkewHeap[V, P] {
	elements := make(map[string]*skewHeapNode[V, P], len(s.elements))
	for _, node := range s.elements {
		cloned := s.pool.Get()
		cloned.id = node.id
		cloned.value = node.value
		cloned.priority = node.priority
		cloned.parent = node.parent
		cloned.left = node.left
		cloned.right = node.right
		elements[node.id] = cloned
	}

	// Restore parent pointers and children links after cloning
	for _, node := range elements {
		// Restore parent pointer if it exists
		if node.parent != nil {
			node.parent = elements[node.parent.id]
		}
		// Restore left child pointer if it exists
		if node.left != nil {
			node.left = elements[node.left.id]
		}
		// Restore right child pointer if it exists
		if node.right != nil {
			node.right = elements[node.right.id]
		}
	}

	return &SkewHeap[V, P]{
		root:     elements[s.root.id],
		cmp:      s.cmp,
		size:     s.size,
		elements: elements,
		pool:     s.pool,
	}
}

// Clear removes all elements from the heap.
// Resets the root to nil, size to zero, and initializes a new empty element map.
// The next node ID is reset to 1.
func (s *SkewHeap[V, P]) Clear() {
	s.root = nil
	s.size = 0
	s.elements = make(map[string]*skewHeapNode[V, P])
}

// Length returns the current number of elements in the heap.
func (s *SkewHeap[V, P]) Length() int { return s.size }

// IsEmpty returns true if the heap contains no elements.
func (s *SkewHeap[V, P]) IsEmpty() bool { return s.size == 0 }

// peek is an internal method that returns the root node's value and priority without removing it.
// Returns nil and an error if the heap is empty.
func (s *SkewHeap[V, P]) peek() (V, P, error) {
	if s.size == 0 {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrHeapEmpty
	}
	return s.root.value, s.root.priority, nil
}

// Peek returns the minimum element without removing it.
// Returns nil and an error if the heap is empty.
func (s *SkewHeap[V, P]) Peek() (V, P, error) { return s.peek() }

// PeekValue returns the value of the minimum element without removing it.
// Returns zero value and an error if the heap is empty.
func (s *SkewHeap[V, P]) PeekValue() (V, error) {
	return valueFromNode(s.peek())
}

// PeekPriority returns the priority of the minimum element without removing it.
// Returns zero value and an error if the heap is empty.
func (s *SkewHeap[V, P]) PeekPriority() (P, error) {
	return priorityFromNode(s.peek())
}

// get is an internal method that retrieves a HeapNode for the node with the given ID.
// Returns nil and an error if the ID doesn't exist in the heap.
func (s *SkewHeap[V, P]) get(id string) (V, P, error) {
	if node, exists := s.elements[id]; exists {
		return node.value, node.priority, nil
	}
	v, p := zeroValuePair[V, P]()
	return v, p, ErrNodeNotFound
}

// Get returns the element with the given ID.
// Returns nil and an error if the ID does not exist.
func (s *SkewHeap[V, P]) Get(id string) (V, P, error) { return s.get(id) }

// GetValue returns the value of the element with the given ID.
// Returns zero value and an error if the ID does not exist.
func (s *SkewHeap[V, P]) GetValue(id string) (V, error) {
	return valueFromNode(s.get(id))
}

// GetPriority returns the priority of the element with the given ID.
// Returns zero value and an error if the ID does not exist.
func (s *SkewHeap[V, P]) GetPriority(id string) (P, error) {
	return priorityFromNode(s.get(id))
}

// pop is an internal method that removes and returns the minimum element from the heap.
// Returns nil and an error if the heap is empty.
func (s *SkewHeap[V, P]) pop() (V, P, error) {
	if s.size == 0 {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrHeapEmpty
	}

	removed := s.root
	s.root = s.merge(s.root.left, s.root.right)
	if s.root != nil {
		s.root.parent = nil
	}
	s.size--
	delete(s.elements, removed.id)
	removed.left, removed.right, removed.parent = nil, nil, nil
	v, p := removed.value, removed.priority
	s.pool.Put(removed)
	return v, p, nil
}

// Pop removes and returns the minimum element from the heap.
// Returns nil and an error if the heap is empty.
func (s *SkewHeap[V, P]) Pop() (V, P, error) { return s.pop() }

// PopValue removes and returns the value of the minimum element.
// Returns zero value and an error if the heap is empty.
func (s *SkewHeap[V, P]) PopValue() (V, error) {
	return valueFromNode(s.pop())
}

// PopPriority removes and returns the priority of the minimum element.
// Returns zero value and an error if the heap is empty.
func (s *SkewHeap[V, P]) PopPriority() (P, error) {
	return priorityFromNode(s.pop())
}

// merge combines two skew heap subtrees into a single heap.
// The root with the higher priority (according to cmp) becomes the new root.
// Children are swapped to maintain the skew heap property.
// Returns the new root of the merged tree.
func (s *SkewHeap[V, P]) merge(new *skewHeapNode[V, P], root *skewHeapNode[V, P]) *skewHeapNode[V, P] {
	if new == nil {
		return root
	}

	if root == nil {
		return new
	}

	first := new
	second := root

	if s.cmp(first.priority, second.priority) {
		tempNode := first.right
		first.right = first.left
		first.left = s.merge(second, tempNode)

		if first.right != nil {
			first.right.parent = first
		}

		if first.left != nil {
			first.left.parent = first
		}
		return first
	} else {
		// When priorities are equal or second has higher priority,
		// merge second with first's children
		tempNode := second.right
		second.right = second.left
		second.left = s.merge(first, tempNode)

		if second.right != nil {
			second.right.parent = second
		}

		if second.left != nil {
			second.left.parent = second
		}
		return second
	}
}

// Push adds a new element to the heap.
// The element is assigned a unique ID and stored in the elements map.
// Returns the ID of the inserted node.
func (s *SkewHeap[V, P]) Push(value V, priority P) string {
	newNode := s.pool.Get()
	newNode.id = uuid.New().String()
	newNode.value = value
	newNode.priority = priority
	s.elements[newNode.id] = newNode
	s.root = s.merge(newNode, s.root)
	s.size++
	return newNode.id
}

// UpdateValue updates the value of the element with the given ID.
// Returns an error if the ID does not exist.
// The heap structure remains unchanged as this operation only modifies the value.
func (s *SkewHeap[V, P]) UpdateValue(id string, value V) error {
	if _, exists := s.elements[id]; !exists {
		return ErrNodeNotFound
	}

	s.elements[id].value = value
	return nil
}

// UpdatePriority updates the priority of the element with the given ID.
// The heap is restructured to maintain the heap property.
// Returns an error if the ID does not exist.
func (s *SkewHeap[V, P]) UpdatePriority(id string, priority P) error {
	if _, exists := s.elements[id]; !exists {
		return ErrNodeNotFound
	}

	updated := s.elements[id]
	updated.priority = priority

	if updated.id == s.root.id {
		s.root = s.merge(updated.left, updated.right)
		s.root.parent = nil
	} else {
		var new *skewHeapNode[V, P]
		parent := updated.parent
		if updated.left == nil && updated.right == nil {
			new = nil
		} else {
			new = s.merge(updated.left, updated.right)
			new.parent = parent
		}

		if parent.left == updated {
			parent.left = new
		} else {
			parent.right = new
		}
	}

	updated.parent, updated.left, updated.right = nil, nil, nil
	s.root = s.merge(updated, s.root)
	return nil
}

// NewSimpleSkewHeap creates a new simple skew heap from the given data slice.
// Each element is inserted individually using the provided comparison function
// to determine heap order (min or max). Returns an empty heap if the input
// slice is empty.
func NewSimpleSkewHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool, usePool bool) *SimpleSkewHeap[V, P] {
	pool := newPool(usePool, func() *skewNode[V, P] {
		return &skewNode[V, P]{}
	})
	heap := SimpleSkewHeap[V, P]{cmp: cmp, size: 0, pool: pool}
	if len(data) == 0 {
		return &heap
	}

	for i := range data {
		heap.Push(data[i].value, data[i].priority)
	}
	return &heap
}

// SimpleSkewHeap implements a basic skew heap without parent pointers.
// It provides the same core functionality as SkewHeap but without element tracking.
// The heap can be either a min-heap or max-heap depending on the comparison function.
type SimpleSkewHeap[V any, P any] struct {
	root *skewNode[V, P]
	cmp  func(a, b P) bool
	size int
	pool pool[*skewNode[V, P]]
}

// Clone creates a deep copy of the heap structure and nodes. If values or
// priorities are reference types, those reference values are shared between the
// original and cloned heaps.
func (s *SimpleSkewHeap[V, P]) Clone() *SimpleSkewHeap[V, P] {
	return &SimpleSkewHeap[V, P]{
		root: s.cloneNode(s.root),
		cmp:  s.cmp,
		size: s.size,
		pool: s.pool,
	}
}

// cloneNode creates a deep copy of a skew node.
// It recursively clones the left and right children.
func (s *SimpleSkewHeap[V, P]) cloneNode(node *skewNode[V, P]) *skewNode[V, P] {
	if node == nil {
		return nil
	}

	cloned := s.pool.Get()
	cloned.value = node.value
	cloned.priority = node.priority
	cloned.left = s.cloneNode(node.left)
	cloned.right = s.cloneNode(node.right)
	return cloned
}

// Clear removes all elements from the heap.
// Resets the root to nil and size to zero.
func (s *SimpleSkewHeap[V, P]) Clear() {
	s.root = nil
	s.size = 0
}

// Length returns the current number of elements in the heap.
func (s *SimpleSkewHeap[V, P]) Length() int { return s.size }

// IsEmpty returns true if the heap contains no elements.
func (s *SimpleSkewHeap[V, P]) IsEmpty() bool { return s.size == 0 }

// peek is an internal method that returns the root node's value and priority without removing it.
// Returns nil and an error if the heap is empty.
func (s *SimpleSkewHeap[V, P]) peek() (V, P, error) {
	if s.size == 0 {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrHeapEmpty
	}
	return s.root.value, s.root.priority, nil
}

// Peek returns the minimum element without removing it.
// Returns nil and an error if the heap is empty.
func (s *SimpleSkewHeap[V, P]) Peek() (V, P, error) { return s.peek() }

// PeekValue returns the value of the minimum element without removing it.
// Returns zero value and an error if the heap is empty.
func (s *SimpleSkewHeap[V, P]) PeekValue() (V, error) {
	return valueFromNode(s.peek())
}

// PeekPriority returns the priority of the minimum element without removing it.
// Returns zero value and an error if the heap is empty.
func (s *SimpleSkewHeap[V, P]) PeekPriority() (P, error) {
	return priorityFromNode(s.peek())
}

// pop is an internal method that removes and returns the minimum element from the heap.
// Returns nil and an error if the heap is empty.
func (s *SimpleSkewHeap[V, P]) pop() (V, P, error) {
	if s.size == 0 {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrHeapEmpty
	}

	rootNode := s.root
	s.root = s.merge(s.root.left, s.root.right)
	rootNode.left, rootNode.right = nil, nil
	s.size--
	v, p := rootNode.value, rootNode.priority
	s.pool.Put(rootNode)
	return v, p, nil
}

// Pop removes and returns the minimum element from the heap.
// Returns nil and an error if the heap is empty.
func (s *SimpleSkewHeap[V, P]) Pop() (V, P, error) { return s.pop() }

// PopValue removes and returns the value of the minimum element.
// Returns zero value and an error if the heap is empty.
func (s *SimpleSkewHeap[V, P]) PopValue() (V, error) {
	return valueFromNode(s.pop())
}

// PopPriority removes and returns the priority of the minimum element.
// Returns zero value and an error if the heap is empty.
func (s *SimpleSkewHeap[V, P]) PopPriority() (P, error) {
	return priorityFromNode(s.pop())
}

// merge combines two skew heap subtrees into a single heap.
// The root with the higher priority (according to cmp) becomes the new root.
// Children are swapped to maintain the skew heap property.
// Returns the new root of the merged tree.
func (s *SimpleSkewHeap[V, P]) merge(new *skewNode[V, P], root *skewNode[V, P]) *skewNode[V, P] {
	if new == nil {
		return root
	}

	if root == nil {
		return new
	}

	first := new
	second := root

	if s.cmp(first.priority, second.priority) {
		tempNode := first.right
		first.right = first.left
		first.left = s.merge(second, tempNode)
		return first
	} else {
		// When priorities are equal or second has higher priority,
		// merge second with first's children
		tempNode := second.right
		second.right = second.left
		second.left = s.merge(first, tempNode)
		return second
	}
}

// Push adds a new element to the heap.
// The element is merged with the existing root to maintain the heap property.
func (s *SimpleSkewHeap[V, P]) Push(value V, priority P) {
	newNode := s.pool.Get()
	newNode.value = value
	newNode.priority = priority
	s.root = s.merge(newNode, s.root)
	s.size++
}
