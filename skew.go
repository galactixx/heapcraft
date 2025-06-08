package heapcraft

import (
	"errors"
	"sync"
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
func NewSkewHeap[V any, P any](data []*HeapPair[V, P], cmp func(a, b P) bool) *SkewHeap[V, P] {
	elements := make(map[uint]*SkewHeapNode[V, P])
	heap := SkewHeap[V, P]{cmp: cmp, size: 0, curID: 1, elements: elements}
	if len(data) == 0 {
		return &heap
	}

	for i := range data {
		heap.Insert(data[i].Value(), data[i].Priority())
	}
	return &heap
}

// SkewHeap implements a skew heap with parent pointers and element tracking.
// It maintains a map of node IDs to nodes for O(1) element access.
type SkewHeap[V any, P any] struct {
	root     *SkewHeapNode[V, P]
	cmp      func(a, b P) bool
	size     int
	curID    uint
	elements map[uint]*SkewHeapNode[V, P]
	lock     sync.RWMutex
}

// Clone creates a shallow copy of the heap.
// The copy shares nodes with the original, so structural modifications
// to either heap will affect both.
func (s *SkewHeap[V, P]) Clone() *SkewHeap[V, P] {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return &SkewHeap[V, P]{root: s.root, cmp: s.cmp, size: s.size, elements: s.elements}
}

// Clear removes all elements from the heap.
func (s *SkewHeap[V, P]) Clear() {
	s.lock.Lock()
	s.root = nil
	s.size = 0
	s.curID = 1
	s.elements = make(map[uint]*SkewHeapNode[V, P])
	s.lock.Unlock()
}

// Length returns the current number of elements in the heap.
func (s *SkewHeap[V, P]) Length() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.size
}

// IsEmpty returns true if the heap contains no elements.
func (s *SkewHeap[V, P]) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.size == 0
}

// peek is an internal method that returns the root node's value and priority without removing it.
// Returns nil if the heap is empty.
func (s *SkewHeap[V, P]) peek() *HeapPair[V, P] {
	if s.size == 0 {
		return nil
	}
	return &HeapPair[V, P]{
		value:    s.root.value,
		priority: s.root.priority,
	}
}

// Peek returns the minimum element without removing it.
// Returns nil if the heap is empty.
func (s *SkewHeap[V, P]) Peek() *HeapPair[V, P] {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.peek()
}

// PeekValue returns the value of the minimum element without removing it.
// Returns nil if the heap is empty.
func (s *SkewHeap[V, P]) PeekValue() *V {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if node := s.peek(); node != nil {
		val := node.Value()
		return &val
	}
	return nil
}

// PeekPriority returns the priority of the minimum element without removing it.
// Returns nil if the heap is empty.
func (s *SkewHeap[V, P]) PeekPriority() *P {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if node := s.peek(); node != nil {
		pri := node.Priority()
		return &pri
	}
	return nil
}

// get is an internal method that retrieves a HeapPair for the node with the given ID.
// Returns an error if the ID doesn't exist in the heap.
func (s *SkewHeap[V, P]) get(id uint) (*HeapPair[V, P], error) {
	if node, exists := s.elements[id]; exists {
		return &HeapPair[V, P]{
			value:    node.value,
			priority: node.priority,
		}, nil
	}
	return nil, errors.New("id does not link to existing node")
}

// Get returns the element with the given ID.
// Returns ErrIDNotFound if the ID does not exist.
func (s *SkewHeap[V, P]) Get(id uint) (*HeapPair[V, P], error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.get(id)
}

// GetValue returns the value of the element with the given ID.
// Returns ErrIDNotFound if the ID does not exist.
func (s *SkewHeap[V, P]) GetValue(id uint) (*V, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	pair, err := s.get(id)
	if err != nil {
		return nil, err
	}
	val := pair.Value()
	return &val, nil
}

// GetPriority returns the priority of the element with the given ID.
// Returns ErrIDNotFound if the ID does not exist.
func (s *SkewHeap[V, P]) GetPriority(id uint) (*P, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	pair, err := s.get(id)
	if err != nil {
		return nil, err
	}
	pri := pair.Priority()
	return &pri, nil
}

// pop is an internal method that removes and returns the minimum element from the heap.
// Returns nil if the heap is empty.
func (s *SkewHeap[V, P]) pop() *SkewHeapNode[V, P] {
	if s.size == 0 {
		return nil
	}

	rootNode := s.root
	s.root = s.merge(s.root.left, s.root.right)
	if s.root != nil {
		s.root.parent = nil
	}
	s.size--
	delete(s.elements, rootNode.id)
	return rootNode
}

// Pop removes and returns the minimum element from the heap.
// Returns nil if the heap is empty.
func (s *SkewHeap[V, P]) Pop() *HeapPair[V, P] {
	s.lock.Lock()
	defer s.lock.Unlock()
	if rootNode := s.pop(); rootNode != nil {
		return &HeapPair[V, P]{
			value:    rootNode.value,
			priority: rootNode.priority,
		}
	}
	return nil
}

// PopValue removes and returns the value of the minimum element.
// Returns nil if the heap is empty.
func (s *SkewHeap[V, P]) PopValue() *V {
	s.lock.Lock()
	defer s.lock.Unlock()
	if rootNode := s.pop(); rootNode != nil {
		val := rootNode.value
		return &val
	}
	return nil
}

// PopPriority removes and returns the priority of the minimum element.
// Returns nil if the heap is empty.
func (s *SkewHeap[V, P]) PopPriority() *P {
	s.lock.Lock()
	defer s.lock.Unlock()
	if rootNode := s.pop(); rootNode != nil {
		pri := rootNode.priority
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
	s.lock.Lock()
	defer s.lock.Unlock()
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
func NewSimpleSkewHeap[V any, P any](data []*HeapPair[V, P], cmp func(a, b P) bool) *SimpleSkewHeap[V, P] {
	heap := SimpleSkewHeap[V, P]{cmp: cmp, size: 0}
	if len(data) == 0 {
		return &heap
	}

	for i := range data {
		heap.Insert(data[i].Value(), data[i].Priority())
	}
	return &heap
}

// SimpleSkewHeap implements a basic skew heap without parent pointers.
// It provides the same core functionality as SkewHeap but without element tracking.
type SimpleSkewHeap[V any, P any] struct {
	root *SkewNode[V, P]
	cmp  func(a, b P) bool
	size int
	lock sync.RWMutex
}

// Clone creates a shallow copy of the heap.
// The copy shares nodes with the original, so structural modifications
// to either heap will affect both.
func (s *SimpleSkewHeap[V, P]) Clone() *SimpleSkewHeap[V, P] {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return &SimpleSkewHeap[V, P]{root: s.root, cmp: s.cmp, size: s.size}
}

// Clear removes all elements from the heap.
func (s *SimpleSkewHeap[V, P]) Clear() {
	s.lock.Lock()
	s.root = nil
	s.size = 0
	s.lock.Unlock()
}

// Length returns the current number of elements in the heap.
func (s *SimpleSkewHeap[V, P]) Length() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.size
}

// IsEmpty returns true if the heap contains no elements.
func (s *SimpleSkewHeap[V, P]) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.size == 0
}

// peek is an internal method that returns the root node's value and priority without removing it.
// Returns nil if the heap is empty.
func (s *SimpleSkewHeap[V, P]) peek() *HeapPair[V, P] {
	if s.size == 0 {
		return nil
	}
	return &HeapPair[V, P]{
		value:    s.root.value,
		priority: s.root.priority,
	}
}

// Peek returns the minimum element without removing it.
// Returns nil if the heap is empty.
func (s *SimpleSkewHeap[V, P]) Peek() *HeapPair[V, P] {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.peek()
}

// PeekValue returns the value of the minimum element without removing it.
// Returns nil if the heap is empty.
func (s *SimpleSkewHeap[V, P]) PeekValue() *V {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if node := s.peek(); node != nil {
		val := node.Value()
		return &val
	}
	return nil
}

// PeekPriority returns the priority of the minimum element without removing it.
// Returns nil if the heap is empty.
func (s *SimpleSkewHeap[V, P]) PeekPriority() *P {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if node := s.peek(); node != nil {
		pri := node.Priority()
		return &pri
	}
	return nil
}

// pop is an internal method that removes and returns the minimum element from the heap.
// Returns nil if the heap is empty.
func (s *SimpleSkewHeap[V, P]) pop() *HeapPair[V, P] {
	if s.size == 0 {
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

// Pop removes and returns the minimum element from the heap.
// Returns nil if the heap is empty.
func (s *SimpleSkewHeap[V, P]) Pop() *HeapPair[V, P] {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.pop()
}

// PopValue removes and returns the value of the minimum element.
// Returns nil if the heap is empty.
func (s *SimpleSkewHeap[V, P]) PopValue() *V {
	s.lock.Lock()
	defer s.lock.Unlock()
	if rootNode := s.pop(); rootNode != nil {
		val := rootNode.Value()
		return &val
	}
	return nil
}

// PopPriority removes and returns the priority of the minimum element.
// Returns nil if the heap is empty.
func (s *SimpleSkewHeap[V, P]) PopPriority() *P {
	s.lock.Lock()
	defer s.lock.Unlock()
	if rootNode := s.pop(); rootNode != nil {
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
	s.lock.Lock()
	defer s.lock.Unlock()
	newNode := SkewNode[V, P]{
		value:    value,
		priority: priority,
	}
	s.root = s.merge(&newNode, s.root)
	s.size++
}

// UpdateValue updates the value of the element with the given ID.
// Returns an error if the ID does not exist.
func (s *SkewHeap[V, P]) UpdateValue(id uint, value V) error {
	s.lock.Lock()
	defer s.lock.Unlock()
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
	s.lock.Lock()
	defer s.lock.Unlock()
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
