package heapcraft

import (
	"errors"
	"sync"
)

// leftistQueue is a generic FIFO queue used for building heaps via pairwise merging.
// It efficiently manages a slice of elements with a head pointer to avoid unnecessary
// allocations when elements are removed.
type leftistQueue[N any] struct {
	data []N
	head int
	size int
}

// push adds an element to the end of the queue, growing the underlying slice if needed.
func (l *leftistQueue[N]) push(element N) {
	l.data = append(l.data, element)
	l.size++
}

// remainingElements returns the count of elements that have not been popped from the queue.
func (l leftistQueue[N]) remainingElements() int {
	return l.size
}

// length returns the total capacity of the underlying slice, including popped elements.
func (l leftistQueue[N]) length() int { return len(l.data) }

// pop removes and returns the element at the head of the queue.
// If the queue is empty, returns the zero value of type N.
// Periodically compacts the underlying slice when the head pointer
// reaches the midpoint to maintain memory efficiency.
func (l *leftistQueue[N]) pop() N {
	if l.remainingElements() == 0 {
		var zero N
		return zero
	}

	popNode := l.data[l.head]
	l.head++

	if l.head >= l.length()/2 {
		l.data = l.data[l.head:]
		l.head = 0
	}
	l.size--
	return popNode
}

// NewSimpleLeftistHeap constructs a leftist heap from a slice of HeapPairs.
// Uses a queue to iteratively merge singleton nodes until one root remains.
// The comparison function determines the heap order (min or max).
func NewSimpleLeftistHeap[V any, P any](data []*HeapNode[V, P], cmp func(a, b P) bool) *SimpleLeftistHeap[V, P] {
	heap := SimpleLeftistHeap[V, P]{cmp: cmp, size: 0}
	if len(data) == 0 {
		return &heap
	}

	n := len(data)
	queueData := make([]*leftistNode[V, P], 0, n)
	initQueue := leftistQueue[*leftistNode[V, P]]{data: queueData, head: 0, size: 0}

	heap.size = n

	for i := range data {
		initQueue.push(&leftistNode[V, P]{
			value:    data[i].value,
			priority: data[i].priority,
			s:        1,
		})
	}

	for initQueue.remainingElements() > 1 {
		merged := heap.merge(initQueue.pop(), initQueue.pop())
		initQueue.push(merged)
	}

	heap.root = initQueue.pop()
	return &heap
}

// NewLeftistHeap constructs a leftist heap with node tracking from a slice of HeapPairs.
// Each node is assigned a unique ID and stored in a map for O(1) access.
// Uses a queue to iteratively merge singleton nodes until one root remains.
// The comparison function determines the heap order (min or max).
func NewLeftistHeap[V any, P any](data []*HeapNode[V, P], cmp func(a, b P) bool) *LeftistHeap[V, P] {
	elements := make(map[uint]*leftistHeapNode[V, P])
	heap := LeftistHeap[V, P]{cmp: cmp, size: 0, curID: 1, elements: elements}
	if len(data) == 0 {
		return &heap
	}

	n := len(data)
	queueData := make([]*leftistHeapNode[V, P], 0, n)
	initQueue := leftistQueue[*leftistHeapNode[V, P]]{data: queueData, head: 0, size: 0}

	var curID uint = 1
	heap.size = n

	for i := range data {
		node := &leftistHeapNode[V, P]{
			id:       curID,
			value:    data[i].value,
			priority: data[i].priority,
			s:        1,
		}
		initQueue.push(node)
		elements[node.id] = node
		curID++
	}

	for initQueue.remainingElements() > 1 {
		merged := heap.merge(initQueue.pop(), initQueue.pop())
		initQueue.push(merged)
	}

	heap.root = initQueue.pop()
	heap.curID = curID
	return &heap
}

// LeftistNode represents a node in a simple leftist heap.
// Each node stores a value, priority, and maintains the leftist property
// through its s-value (null-path length) and child pointers.
type leftistNode[V any, P any] struct {
	value    V
	priority P
	left     *leftistNode[V, P]
	right    *leftistNode[V, P]
	s        int
}

// Value returns the value stored in the node.
func (n *leftistNode[V, P]) Value() V { return n.value }

// Priority returns the priority of the node.
func (n *leftistNode[V, P]) Priority() P { return n.priority }

// leftistHeapNode represents a node in a tracked leftist heap.
// Extends leftistNode with an ID and parent pointer for node tracking
// and efficient updates.
type leftistHeapNode[V any, P any] struct {
	id       uint
	value    V
	priority P
	parent   *leftistHeapNode[V, P]
	left     *leftistHeapNode[V, P]
	right    *leftistHeapNode[V, P]
	s        int
}

// ID returns the unique identifier of the node.
func (n *leftistHeapNode[V, P]) ID() uint { return n.id }

// Value returns the value stored in the node.
func (n *leftistHeapNode[V, P]) Value() V { return n.value }

// Priority returns the priority of the node.
func (n *leftistHeapNode[V, P]) Priority() P { return n.priority }

// LeftistHeap implements a leftist heap with node tracking capabilities.
// Maintains a map of node IDs to nodes for O(1) access and updates.
// The heap property is maintained through the comparison function.
type LeftistHeap[V any, P any] struct {
	root     *leftistHeapNode[V, P]
	cmp      func(a, b P) bool
	size     int
	curID    uint
	elements map[uint]*leftistHeapNode[V, P]
	lock     sync.RWMutex
}

// UpdateValue changes the value of the node with the given ID.
// Returns an error if the ID doesn't exist in the heap.
func (l *LeftistHeap[V, P]) UpdateValue(id uint, value V) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if _, exists := l.elements[id]; !exists {
		return errors.New("id does not link to existing node")
	}

	l.elements[id].value = value
	return nil
}

// UpdatePriority changes the priority of the node with the given ID and
// restructures the heap to maintain the heap property.
// Returns an error if the ID doesn't exist in the heap.
func (l *LeftistHeap[V, P]) UpdatePriority(id uint, priority P) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if _, exists := l.elements[id]; !exists {
		return errors.New("id does not link to existing node")
	}

	updated := l.elements[id]
	updated.priority = priority

	if updated.id == l.root.id {
		l.root = l.merge(l.root.left, l.root.right)
		l.root.parent = nil
	} else {
		var new *leftistHeapNode[V, P]
		parent := updated.parent
		if updated.left == nil && updated.right == nil {
			new = nil
		} else {
			new = l.merge(updated.left, updated.right)
			new.parent = parent
		}

		if parent.left == updated {
			parent.left = new
		} else {
			parent.right = new
		}
	}

	updated.parent, updated.left, updated.right = nil, nil, nil
	l.root = l.merge(updated, l.root)
	return nil
}

// Clone creates a shallow copy of the heap.
// The new heap shares the same nodes and underlying structure.
func (l *LeftistHeap[V, P]) Clone() *LeftistHeap[V, P] {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return &LeftistHeap[V, P]{root: l.root, cmp: l.cmp, size: l.size, elements: l.elements}
}

// Clear removes all elements from the heap and resets its state.
// The heap is ready for new insertions after clearing.
func (l *LeftistHeap[V, P]) Clear() {
	l.lock.Lock()
	l.root = nil
	l.size = 0
	l.curID = 1
	l.elements = make(map[uint]*leftistHeapNode[V, P])
	l.lock.Unlock()
}

// Length returns the current number of elements in the heap.
func (l *LeftistHeap[V, P]) Length() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.size
}

// IsEmpty returns true if the heap contains no elements.
func (l *LeftistHeap[V, P]) IsEmpty() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.size == 0
}

// peek is an internal method that returns the root node without removing it.
// Returns nil and error if the heap is empty.
func (l *LeftistHeap[V, P]) peek() (*leftistHeapNode[V, P], error) {
	if l.size == 0 {
		return nil, errors.New("the heap is empty and contains no elements")
	}
	return l.root, nil
}

// Peek returns the minimum element without removing it.
// Returns nil and error if the heap is empty.
func (l *LeftistHeap[V, P]) Peek() (Node[V, P], error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.peek()
}

// PeekValue returns the value at the root without removing it.
// Returns zero value and error if the heap is empty.
func (l *LeftistHeap[V, P]) PeekValue() (V, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	pair, err := l.peek()
	if err != nil {
		var zero V
		return zero, err
	}
	return pair.value, nil
}

// PeekPriority returns the priority at the root without removing it.
// Returns zero value and error if the heap is empty.
func (l *LeftistHeap[V, P]) PeekPriority() (P, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	pair, err := l.peek()
	if err != nil {
		var zero P
		return zero, err
	}
	return pair.priority, nil
}

// get is an internal method that retrieves a node with the given ID.
// Returns an error if the ID doesn't exist in the heap.
func (l *LeftistHeap[V, P]) get(id uint) (*leftistHeapNode[V, P], error) {
	if node, exists := l.elements[id]; exists {
		return node, nil
	}
	return nil, errors.New("id does not link to existing node")
}

// Get returns the element associated with the given ID.
// Returns an error if the ID doesn't exist in the heap.
func (l *LeftistHeap[V, P]) Get(id uint) (Node[V, P], error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.get(id)
}

// GetValue returns the value associated with the given ID.
// Returns zero value and error if the ID doesn't exist in the heap.
func (l *LeftistHeap[V, P]) GetValue(id uint) (V, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	pair, err := l.get(id)
	if err != nil {
		var zero V
		return zero, err
	}
	return pair.value, nil
}

// GetPriority returns the priority associated with the given ID.
// Returns zero value and error if the ID doesn't exist in the heap.
func (l *LeftistHeap[V, P]) GetPriority(id uint) (P, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	pair, err := l.get(id)
	if err != nil {
		var zero P
		return zero, err
	}
	return pair.priority, nil
}

// Pop removes and returns the minimum element from the heap.
// The heap property is restored through merging the root's children.
// Returns nil and error if the heap is empty.
func (l *LeftistHeap[V, P]) Pop() (Node[V, P], error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.pop()
}

// PopValue removes and returns just the value at the root.
// The heap property is restored through merging the root's children.
// Returns zero value and error if the heap is empty.
func (l *LeftistHeap[V, P]) PopValue() (V, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	pair, err := l.pop()
	if err != nil {
		var zero V
		return zero, err
	}
	return pair.value, nil
}

// PopPriority removes and returns just the priority at the root.
// The heap property is restored through merging the root's children.
// Returns zero value and error if the heap is empty.
func (l *LeftistHeap[V, P]) PopPriority() (P, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	pair, err := l.pop()
	if err != nil {
		var zero P
		return zero, err
	}
	return pair.priority, nil
}

// pop is an internal method that removes the root node and returns it.
// Handles the common logic of removing the root and merging its children.
// Returns nil and error if the heap is empty.
func (l *LeftistHeap[V, P]) pop() (*leftistHeapNode[V, P], error) {
	if l.size == 0 {
		return nil, errors.New("the heap is empty and contains no elements")
	}

	rootNode := l.root
	l.root = l.merge(l.root.right, l.root.left)
	if l.root != nil {
		l.root.parent = nil
	}
	delete(l.elements, rootNode.id)
	l.size--
	return rootNode, nil
}

// merge combines two leftist subheaps while maintaining the heap property
// and leftist structure. The root of the resulting heap is the node with
// the minimum priority according to the comparison function.
func (l *LeftistHeap[V, P]) merge(a, b *leftistHeapNode[V, P]) *leftistHeapNode[V, P] {
	if a == nil {
		return b
	}

	if b == nil {
		return a
	}

	if l.cmp(a.priority, b.priority) {
		return l.merge(b, a)
	}

	b.right = l.merge(b.right, a)
	b.right.parent = b
	if b.left == nil {
		b.left = b.right
		b.right = nil
	} else {
		if b.left.s < b.right.s {
			b.left, b.right = b.right, b.left
		}
		b.s = b.right.s + 1
	}
	b.left.parent = b
	return b
}

// Push adds a new element to the heap by creating a singleton node
// and merging it with the existing tree. The new node is assigned
// a unique ID and stored in the elements map. Returns the ID of the inserted node.
func (l *LeftistHeap[V, P]) Push(value V, priority P) uint {
	l.lock.Lock()
	defer l.lock.Unlock()
	newNode := &leftistHeapNode[V, P]{
		id:       l.curID,
		value:    value,
		priority: priority,
		s:        1,
	}
	l.root = l.merge(newNode, l.root)
	l.elements[newNode.id] = newNode
	l.size++
	l.curID++
	return newNode.id
}

// SimpleLeftistHeap implements a basic leftist heap without node tracking.
// Maintains the heap property through the comparison function and
// the leftist property through s-values.
type SimpleLeftistHeap[V any, P any] struct {
	root *leftistNode[V, P]
	cmp  func(a, b P) bool
	size int
	lock sync.RWMutex
}

// Clone creates a shallow copy of the simple heap.
// The new heap shares the same nodes and underlying structure.
func (l *SimpleLeftistHeap[V, P]) Clone() *SimpleLeftistHeap[V, P] {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return &SimpleLeftistHeap[V, P]{root: l.root, cmp: l.cmp, size: l.size}
}

// Clear removes all elements from the simple heap.
// The heap is ready for new insertions after clearing.
func (l *SimpleLeftistHeap[V, P]) Clear() {
	l.lock.Lock()
	l.root = nil
	l.size = 0
	l.lock.Unlock()
}

// Length returns the current number of elements in the simple heap.
func (l *SimpleLeftistHeap[V, P]) Length() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.size
}

// IsEmpty returns true if the simple heap contains no elements.
func (l *SimpleLeftistHeap[V, P]) IsEmpty() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.size == 0
}

// peek is an internal method that returns the root node without removing it.
// Returns nil and error if the heap is empty.
func (l *SimpleLeftistHeap[V, P]) peek() (*leftistNode[V, P], error) {
	if l.size == 0 {
		return nil, errors.New("the heap is empty and contains no elements")
	}
	return l.root, nil
}

// Peek returns the minimum element without removing it.
// Returns nil and error if the heap is empty.
func (l *SimpleLeftistHeap[V, P]) Peek() (SimpleNode[V, P], error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.peek()
}

// PeekValue returns the value at the root without removing it.
// Returns zero value and error if the heap is empty.
func (l *SimpleLeftistHeap[V, P]) PeekValue() (V, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	pair, err := l.peek()
	if err != nil {
		var zero V
		return zero, err
	}
	return pair.value, nil
}

// PeekPriority returns the priority at the root without removing it.
// Returns zero value and error if the heap is empty.
func (l *SimpleLeftistHeap[V, P]) PeekPriority() (P, error) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	pair, err := l.peek()
	if err != nil {
		var zero P
		return zero, err
	}
	return pair.priority, nil
}

// pop is an internal method that removes the root node and returns it.
// Handles the common logic of removing the root and merging its children.
// Returns nil and error if the heap is empty.
func (l *SimpleLeftistHeap[V, P]) pop() (*leftistNode[V, P], error) {
	if l.size == 0 {
		return nil, errors.New("the heap is empty and contains no elements")
	}

	rootNode := l.root
	l.root = l.merge(l.root.right, l.root.left)
	l.size--
	return rootNode, nil
}

// Pop removes and returns the minimum element from the simple heap.
// The heap property is restored through merging the root's children.
// Returns nil and error if the heap is empty.
func (l *SimpleLeftistHeap[V, P]) Pop() (SimpleNode[V, P], error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.pop()
}

// PopValue removes and returns just the value at the root.
// The heap property is restored through merging the root's children.
// Returns zero value and error if the heap is empty.
func (l *SimpleLeftistHeap[V, P]) PopValue() (V, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	pair, err := l.pop()
	if err != nil {
		var zero V
		return zero, err
	}
	return pair.value, nil
}

// PopPriority removes and returns just the priority at the root.
// The heap property is restored through merging the root's children.
// Returns zero value and error if the heap is empty.
func (l *SimpleLeftistHeap[V, P]) PopPriority() (P, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	pair, err := l.pop()
	if err != nil {
		var zero P
		return zero, err
	}
	return pair.priority, nil
}

// merge combines two leftist subheaps while maintaining the heap property
// and leftist structure. The root of the resulting heap is the node with
// the minimum priority according to the comparison function.
func (l *SimpleLeftistHeap[V, P]) merge(a, b *leftistNode[V, P]) *leftistNode[V, P] {
	if a == nil {
		return b
	}

	if b == nil {
		return a
	}

	if l.cmp(a.priority, b.priority) {
		return l.merge(b, a)
	}

	b.right = l.merge(b.right, a)
	if b.left == nil {
		b.left = b.right
		b.right = nil
	} else {
		if b.left.s < b.right.s {
			b.left, b.right = b.right, b.left
		}
		b.s = b.right.s + 1
	}
	return b
}

// Push adds a new element to the simple heap by creating a singleton node
// and merging it with the existing tree.
func (l *SimpleLeftistHeap[V, P]) Push(value V, priority P) {
	l.lock.Lock()
	defer l.lock.Unlock()
	newNode := &leftistNode[V, P]{
		value:    value,
		priority: priority,
		s:        1,
	}
	l.root = l.merge(newNode, l.root)
	l.size++
}
