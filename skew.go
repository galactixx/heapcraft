package heapcraft

import (
	"errors"
)

var ErrIDNotFound = errors.New("element with given ID not found in heap")

// SkewNode represents a node in a simple skew heap without parent pointers
type SkewNode[V any, P any] struct {
	value    V
	priority P
	left     *SkewNode[V, P]
	right    *SkewNode[V, P]
}

// SkewHeapNode represents a node in a skew heap with parent pointers and ID tracking
type SkewHeapNode[V any, P any] struct {
	id       uint
	value    V
	priority P
	parent   *SkewHeapNode[V, P]
	left     *SkewHeapNode[V, P]
	right    *SkewHeapNode[V, P]
}

// NewSkewHeap creates a new skew heap from the given data slice.
// Each element is inserted individually using the provided comparison function.
// Returns an empty heap if the input slice is empty.
func NewSkewHeap[V any, P any](data []*HeapPair[V, P], cmp func(a, b P) bool) SkewHeap[V, P] {
	elements := make(map[uint]*SkewHeapNode[V, P])
	heap := SkewHeap[V, P]{cmp: cmp, size: 0, curID: 1, elements: elements}
	if len(data) == 0 {
		return heap
	}

	for i := range data {
		heap.Insert(data[i].Value(), data[i].Priority())
	}
	return heap
}

// SkewHeap implements a skew heap with parent pointers and element tracking.
// It maintains a map of node IDs to nodes for O(1) element access.
type SkewHeap[V any, P any] struct {
	root     *SkewHeapNode[V, P]
	cmp      func(a, b P) bool
	size     int
	curID    uint
	elements map[uint]*SkewHeapNode[V, P]
}

// Clone creates a shallow copy of the heap.
// The copy shares nodes with the original, so structural modifications
// to either heap will affect both.
func (s SkewHeap[V, P]) Clone() SkewHeap[V, P] {
	return SkewHeap[V, P]{root: s.root, cmp: s.cmp, size: s.size}
}

// Clear removes all elements from the heap.
func (s *SkewHeap[V, P]) Clear() { s.root = nil; s.size = 0 }

// Length returns the current number of elements in the heap.
func (s SkewHeap[V, P]) Length() int { return s.size }

// IsEmpty returns true if the heap contains no elements.
func (s SkewHeap[V, P]) IsEmpty() bool { return s.Length() == 0 }

// Peek returns the minimum element without removing it.
// Returns nil if the heap is empty.
func (s *SkewHeap[V, P]) Peek() *HeapPair[V, P] {
	if s.IsEmpty() {
		return nil
	}

	return &HeapPair[V, P]{
		value:    s.root.value,
		priority: s.root.priority,
	}
}

// PeekValue returns the value of the minimum element without removing it.
// Returns nil if the heap is empty.
func (s *SkewHeap[V, P]) PeekValue() *V {
	if node := s.Peek(); node != nil {
		val := node.Value()
		return &val
	}
	return nil
}

// PeekPriority returns the priority of the minimum element without removing it.
// Returns nil if the heap is empty.
func (s *SkewHeap[V, P]) PeekPriority() *P {
	if node := s.Peek(); node != nil {
		pri := node.Priority()
		return &pri
	}
	return nil
}

// UpdateValue updates the value of the element with the given ID.
// Returns an error if the ID does not exist.
func (s *SkewHeap[V, P]) UpdateValue(id uint, value V) error {
	if _, exists := s.elements[id]; !exists {
		return errors.New("id does not link to existing node")
	}

	s.elements[id].value = value
	return nil
}

// UpdatePriority updates the priority of the element with the given ID.
// The heap is restructured to maintain the heap property.
// Returns an error if the ID does not exist.
func (s *SkewHeap[V, P]) UpdatePriority(id uint, priority P) error {
	if _, exists := s.elements[id]; !exists {
		return errors.New("id does not link to existing node")
	}

	updated := s.elements[id]
	updated.priority = priority

	if updated.id == s.root.id {
		s.root = s.merge(updated.left, updated.right)
		s.root.parent = nil
	} else {
		var new *SkewHeapNode[V, P]
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

// Get returns the element with the given ID.
// Returns ErrIDNotFound if the ID does not exist.
func (s *SkewHeap[V, P]) Get(id uint) (*HeapPair[V, P], error) {
	node, exists := s.elements[id]
	if !exists {
		return nil, ErrIDNotFound
	}
	return &HeapPair[V, P]{
		value:    node.value,
		priority: node.priority,
	}, nil
}

// GetValue returns the value of the element with the given ID.
// Returns ErrIDNotFound if the ID does not exist.
func (s *SkewHeap[V, P]) GetValue(id uint) (*V, error) {
	pair, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	val := pair.Value()
	return &val, nil
}

// GetPriority returns the priority of the element with the given ID.
// Returns ErrIDNotFound if the ID does not exist.
func (s *SkewHeap[V, P]) GetPriority(id uint) (*P, error) {
	pair, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	pri := pair.Priority()
	return &pri, nil
}

// Pop removes and returns the minimum element from the heap.
// Returns nil if the heap is empty.
func (s *SkewHeap[V, P]) Pop() *HeapPair[V, P] {
	if s.IsEmpty() {
		return nil
	}

	rootNode := s.root
	s.root = s.merge(s.root.left, s.root.right)
	if s.root != nil {
		s.root.parent = nil
	}
	s.size--
	delete(s.elements, rootNode.id)
	return &HeapPair[V, P]{
		value:    rootNode.value,
		priority: rootNode.priority,
	}
}

// PopValue removes and returns the value of the minimum element.
// Returns nil if the heap is empty.
func (s *SkewHeap[V, P]) PopValue() *V {
	if rootNode := s.Pop(); rootNode != nil {
		val := rootNode.Value()
		return &val
	}
	return nil
}

// PopPriority removes and returns the priority of the minimum element.
// Returns nil if the heap is empty.
func (s *SkewHeap[V, P]) PopPriority() *P {
	if rootNode := s.Pop(); rootNode != nil {
		pri := rootNode.Priority()
		return &pri
	}
	return nil
}

// merge combines two skew heap subtrees into a single heap.
// The root with the smaller priority becomes the new root.
// Children are swapped to maintain the skew heap property.
func (s *SkewHeap[V, P]) merge(new *SkewHeapNode[V, P], root *SkewHeapNode[V, P]) *SkewHeapNode[V, P] {
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
		return s.merge(second, first)
	}
}

// Insert adds a new element to the heap.
// The element is assigned a unique ID and stored in the elements map.
func (s *SkewHeap[V, P]) Insert(value V, priority P) {
	newNode := SkewHeapNode[V, P]{
		id:       s.curID,
		value:    value,
		priority: priority,
	}
	s.elements[newNode.id] = &newNode
	s.root = s.merge(&newNode, s.root)
	s.size++
	s.curID++
}

// NewSimpleSkewHeap creates a new simple skew heap from the given data slice.
// Each element is inserted individually using the provided comparison function.
// Returns an empty heap if the input slice is empty.
func NewSimpleSkewHeap[V any, P any](data []*HeapPair[V, P], cmp func(a, b P) bool) SimpleSkewHeap[V, P] {
	heap := SimpleSkewHeap[V, P]{cmp: cmp, size: 0}
	if len(data) == 0 {
		return heap
	}

	for i := range data {
		heap.Insert(data[i].Value(), data[i].Priority())
	}
	return heap
}

// SimpleSkewHeap implements a basic skew heap without parent pointers.
// It provides the same core functionality as SkewHeap but without element tracking.
type SimpleSkewHeap[V any, P any] struct {
	root *SkewNode[V, P]
	cmp  func(a, b P) bool
	size int
}

// Clone creates a shallow copy of the heap.
// The copy shares nodes with the original, so structural modifications
// to either heap will affect both.
func (s SimpleSkewHeap[V, P]) Clone() SimpleSkewHeap[V, P] {
	return SimpleSkewHeap[V, P]{root: s.root, cmp: s.cmp, size: s.size}
}

// Clear removes all elements from the heap.
func (s *SimpleSkewHeap[V, P]) Clear() { s.root = nil; s.size = 0 }

// Length returns the current number of elements in the heap.
func (s SimpleSkewHeap[V, P]) Length() int { return s.size }

// IsEmpty returns true if the heap contains no elements.
func (s SimpleSkewHeap[V, P]) IsEmpty() bool { return s.Length() == 0 }

// Peek returns the minimum element without removing it.
// Returns nil if the heap is empty.
func (s *SimpleSkewHeap[V, P]) Peek() *HeapPair[V, P] {
	if s.IsEmpty() {
		return nil
	}

	return &HeapPair[V, P]{
		value:    s.root.value,
		priority: s.root.priority,
	}
}

// PeekValue returns the value of the minimum element without removing it.
// Returns nil if the heap is empty.
func (s *SimpleSkewHeap[V, P]) PeekValue() *V {
	if node := s.Peek(); node != nil {
		val := node.Value()
		return &val
	}
	return nil
}

// PeekPriority returns the priority of the minimum element without removing it.
// Returns nil if the heap is empty.
func (s *SimpleSkewHeap[V, P]) PeekPriority() *P {
	if node := s.Peek(); node != nil {
		pri := node.Priority()
		return &pri
	}
	return nil
}

// Pop removes and returns the minimum element from the heap.
// Returns nil if the heap is empty.
func (s *SimpleSkewHeap[V, P]) Pop() *HeapPair[V, P] {
	if s.IsEmpty() {
		return nil
	}

	rootNode := s.root
	s.root = s.merge(s.root.left, s.root.right)
	s.size--
	return &HeapPair[V, P]{
		value:    rootNode.value,
		priority: rootNode.priority,
	}
}

// PopValue removes and returns the value of the minimum element.
// Returns nil if the heap is empty.
func (s *SimpleSkewHeap[V, P]) PopValue() *V {
	if rootNode := s.Pop(); rootNode != nil {
		val := rootNode.Value()
		return &val
	}
	return nil
}

// PopPriority removes and returns the priority of the minimum element.
// Returns nil if the heap is empty.
func (s *SimpleSkewHeap[V, P]) PopPriority() *P {
	if rootNode := s.Pop(); rootNode != nil {
		pri := rootNode.Priority()
		return &pri
	}
	return nil
}

// merge combines two skew heap subtrees into a single heap.
// The root with the smaller priority becomes the new root.
// Children are swapped to maintain the skew heap property.
func (s *SimpleSkewHeap[V, P]) merge(new *SkewNode[V, P], root *SkewNode[V, P]) *SkewNode[V, P] {
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
		return s.merge(second, first)
	}
}

// Insert adds a new element to the heap.
func (s *SimpleSkewHeap[V, P]) Insert(value V, priority P) {
	newNode := SkewNode[V, P]{
		value:    value,
		priority: priority,
	}
	s.root = s.merge(&newNode, s.root)
	s.size++
}
