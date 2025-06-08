package heapcraft

import (
	"errors"
	"sync"
)

// clearNodeLinks resets all the linking pointers of a node to nil.
// This is used when removing a node from its current position in the heap
// before reinserting it elsewhere.
func clearNodeLinks[V any, P any](node *PairingHeapNode[V, P]) {
	node.nextSibling = nil
	node.parent = nil
	node.prevSibling = nil
}

// NewPairingHeap creates a new pairing heap from a slice of HeapPairs.
// The heap is initialized with the provided elements and uses the given comparison
// function to determine heap order. The comparison function should return true
// when the first priority is considered higher priority than the second.
// Returns an empty heap if the input slice is empty.
func NewPairingHeap[V any, P any](data []*HeapPair[V, P], cmp func(a, b P) bool) *PairingHeap[V, P] {
	elements := make(map[uint]*PairingHeapNode[V, P])
	heap := PairingHeap[V, P]{cmp: cmp, size: 0, curID: 1, elements: elements}
	if len(data) == 0 {
		return &heap
	}

	for i := range data {
		heap.Insert(data[i].Value(), data[i].Priority())
	}
	return &heap
}

// NewSimplePairingHeap creates a new simple pairing heap from a slice of HeapPairs.
// Unlike PairingHeap, this implementation does not track node IDs or support
// node updates. It uses the provided comparison function to determine heap order.
// Returns an empty heap if the input slice is empty.
func NewSimplePairingHeap[V any, P any](data []*HeapPair[V, P], cmp func(a, b P) bool) *SimplePairingHeap[V, P] {
	heap := SimplePairingHeap[V, P]{cmp: cmp, size: 0}
	if len(data) == 0 {
		return &heap
	}

	for i := range data {
		heap.Insert(data[i].Value(), data[i].Priority())
	}
	return &heap
}

// PairingHeapNode represents a node in the pairing heap data structure.
// Each node contains a value, priority, and maintains links to its parent,
// children, and siblings. The node also has a unique identifier for tracking.
// The doubly-linked sibling list allows for efficient node removal and updates.
type PairingHeapNode[V any, P any] struct {
	id          uint
	value       V
	priority    P
	parent      *PairingHeapNode[V, P]
	firstChild  *PairingHeapNode[V, P]
	nextSibling *PairingHeapNode[V, P]
	prevSibling *PairingHeapNode[V, P]
}

// ID returns the unique identifier of the node.
// This identifier is used for tracking and updating nodes in the heap.
func (n *PairingHeapNode[V, P]) ID() uint { return n.id }

// Value returns the value stored in the node.
func (n *PairingHeapNode[V, P]) Value() V { return n.value }

// Priority returns the priority of the node.
func (n *PairingHeapNode[V, P]) Priority() P { return n.priority }

// PairingHeap implements a pairing heap data structure with node tracking.
// It maintains a multi-way tree structure where each node can have multiple children.
// The heap supports efficient insertion, deletion, and priority updates of nodes.
// Nodes are tracked by unique IDs, allowing for O(1) access and updates.
type PairingHeap[V any, P any] struct {
	root     *PairingHeapNode[V, P]
	cmp      func(a, b P) bool
	size     int
	curID    uint
	elements map[uint]*PairingHeapNode[V, P]
	lock     sync.RWMutex
}

// UpdateValue updates the value of a node with the given ID.
// Returns an error if the ID does not exist in the heap.
// The heap structure remains unchanged as this operation only modifies the value.
func (p *PairingHeap[V, P]) UpdateValue(id uint, value V) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, exists := p.elements[id]; !exists {
		return errors.New("id does not link to existing node")
	}

	p.elements[id].value = value
	return nil
}

// UpdatePriority updates the priority of a node with the given ID.
// Returns an error if the ID does not exist in the heap.
// The node is removed from its current position and reinserted into the heap
// to maintain the heap property. This operation may change the heap structure.
func (p *PairingHeap[V, P]) UpdatePriority(id uint, priority P) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, exists := p.elements[id]; !exists {
		return errors.New("id does not link to existing node")
	}

	updated := p.elements[id]
	updated.priority = priority

	if updated.id == p.root.id {
		newRoot := updated.firstChild
		if newRoot != nil {
			newRoot.prevSibling, newRoot.parent = nil, nil
		}
		updated.firstChild = nil
		p.root = p.merge(newRoot)
	} else if updated.prevSibling != nil {
		prev, next := updated.prevSibling, updated.nextSibling
		if next != nil {
			next.prevSibling = prev
		}

		prev.nextSibling = next
	} else {
		next := updated.nextSibling
		if next != nil {
			next.prevSibling, next.parent = nil, updated.parent
		}
		updated.parent.firstChild = next
	}

	clearNodeLinks(updated)
	p.root = p.meld(updated, p.root)
	return nil
}

// Clone creates a shallow copy of the heap.
// The new heap shares the same nodes as the original but has its own
// root pointer and element map. Modifications to the clone will affect
// the original heap's nodes but not its structure.
func (p *PairingHeap[V, P]) Clone() *PairingHeap[V, P] {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return &PairingHeap[V, P]{
		root: p.root, cmp: p.cmp, size: p.size, curID: p.curID, elements: p.elements,
	}
}

// Clear removes all elements from the heap.
// Resets the root to nil, size to zero, and initializes a new empty element map.
// The next node ID is reset to 1.
func (p *PairingHeap[V, P]) Clear() {
	p.lock.Lock()
	p.root = nil
	p.size = 0
	p.curID = 1
	p.elements = make(map[uint]*PairingHeapNode[V, P], 0)
	p.lock.Unlock()
}

// Length returns the current number of elements in the heap.
func (p *PairingHeap[V, P]) Length() int {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.size
}

// IsEmpty returns true if the heap contains no elements.
func (p *PairingHeap[V, P]) IsEmpty() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.size == 0
}

// peek is an internal method that returns the root node's value and priority without removing it.
// Returns nil if the heap is empty.
func (p *PairingHeap[V, P]) peek() *HeapPair[V, P] {
	if p.size == 0 {
		return nil
	}
	return &HeapPair[V, P]{value: p.root.value, priority: p.root.priority}
}

// Peek returns a pointer to a HeapPair containing the value and priority
// of the root node without removing it. Returns nil if the heap is empty.
func (p *PairingHeap[V, P]) Peek() *HeapPair[V, P] {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.peek()
}

// PeekValue returns a pointer to the value at the root without removing it.
// Returns nil if the heap is empty.
func (p *PairingHeap[V, P]) PeekValue() *V {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if node := p.peek(); node != nil {
		val := node.Value()
		return &val
	}
	return nil
}

// PeekPriority returns a pointer to the priority at the root without removing it.
// Returns nil if the heap is empty.
func (p *PairingHeap[V, P]) PeekPriority() *P {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if node := p.peek(); node != nil {
		pri := node.Priority()
		return &pri
	}
	return nil
}

// get is an internal method that retrieves a HeapPair for the node with the given ID.
// Returns an error if the ID does not exist in the heap.
func (p *PairingHeap[V, P]) get(id uint) (*HeapPair[V, P], error) {
	node, exists := p.elements[id]
	if !exists {
		return nil, errors.New("node with id does not exist")
	}
	return &HeapPair[V, P]{value: node.value, priority: node.priority}, nil
}

// Get retrieves a HeapPair for the node with the given ID.
// Returns an error if the ID does not exist in the heap.
func (p *PairingHeap[V, P]) Get(id uint) (*HeapPair[V, P], error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.get(id)
}

// GetValue retrieves the value of the node with the given ID.
// Returns an error if the ID does not exist in the heap.
func (p *PairingHeap[V, P]) GetValue(id uint) (*V, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	pair, err := p.get(id)
	if err != nil {
		return nil, err
	}
	val := pair.Value()
	return &val, nil
}

// GetPriority retrieves the priority of the node with the given ID.
// Returns an error if the ID does not exist in the heap.
func (p *PairingHeap[V, P]) GetPriority(id uint) (*P, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	pair, err := p.get(id)
	if err != nil {
		return nil, err
	}
	pri := pair.Priority()
	return &pri, nil
}

// meld combines two pairing heap trees into a single tree.
// The tree with the higher priority (according to cmp) becomes the root,
// and the other tree becomes its first child. The operation maintains
// the doubly-linked sibling list structure.
// Returns the new root of the combined tree.
func (p *PairingHeap[V, P]) meld(new *PairingHeapNode[V, P], root *PairingHeapNode[V, P]) *PairingHeapNode[V, P] {
	if root == nil {
		return new
	}

	if new == nil {
		return root
	}

	var prior, noPrior *PairingHeapNode[V, P]

	if p.cmp(new.priority, root.priority) {
		prior, noPrior = new, root
	} else {
		prior, noPrior = root, new
	}

	if prior.firstChild != nil {
		prior.firstChild.prevSibling = noPrior
		prior.firstChild.parent = prior
	}

	noPrior.nextSibling = prior.firstChild
	noPrior.parent = prior
	noPrior.prevSibling = nil
	prior.firstChild = noPrior
	return prior
}

// merge performs the two-pass pairing process on a list of siblings.
// It pairs adjacent siblings, melds them, and recursively merges the
// remaining siblings. This operation is used during Pop to combine
// the root's children into a new heap structure.
// Returns the new root of the merged tree.
func (p *PairingHeap[V, P]) merge(node *PairingHeapNode[V, P]) *PairingHeapNode[V, P] {
	if node == nil {
		return node
	}

	if node.nextSibling == nil {
		clearNodeLinks(node)
		return node
	}

	firstNode := node
	secondNode := node.nextSibling
	remaining := node.nextSibling.nextSibling

	clearNodeLinks(firstNode)
	clearNodeLinks(secondNode)
	return p.meld(p.meld(firstNode, secondNode), p.merge(remaining))
}

// pop is an internal method that removes and returns the root node.
// It handles the common logic of removing the root, merging its children,
// updating the size, and removing the node from the element map.
// Returns nil if the heap is empty.
func (p *PairingHeap[V, P]) pop() *PairingHeapNode[V, P] {
	if p.size == 0 {
		return nil
	}

	rootNode := p.root
	p.root = p.merge(p.root.firstChild)
	p.size--
	delete(p.elements, rootNode.id)
	return rootNode
}

// Pop removes and returns a HeapPair containing the value and priority
// of the root node. The root's children are merged to form the new heap.
// Returns nil if the heap is empty.
func (p *PairingHeap[V, P]) Pop() *HeapPair[V, P] {
	p.lock.Lock()
	defer p.lock.Unlock()
	if rootNode := p.pop(); rootNode != nil {
		return &HeapPair[V, P]{value: rootNode.value, priority: rootNode.priority}
	}
	return nil
}

// PopValue removes and returns a pointer to just the value at the root.
// The root's children are merged to form the new heap.
// Returns nil if the heap is empty.
func (p *PairingHeap[V, P]) PopValue() *V {
	p.lock.Lock()
	defer p.lock.Unlock()
	if rootNode := p.pop(); rootNode != nil {
		val := rootNode.value
		return &val
	}
	return nil
}

// PopPriority removes and returns a pointer to just the priority at the root.
// The root's children are merged to form the new heap.
// Returns nil if the heap is empty.
func (p *PairingHeap[V, P]) PopPriority() *P {
	p.lock.Lock()
	defer p.lock.Unlock()
	if rootNode := p.pop(); rootNode != nil {
		pri := rootNode.priority
		return &pri
	}
	return nil
}

// Insert adds a new element with the given value and priority to the heap.
// A new node is created with a unique ID and melded with the existing root.
// The new node becomes the root if its priority is higher than the current root's.
func (p *PairingHeap[V, P]) Insert(value V, priority P) {
	p.lock.Lock()
	defer p.lock.Unlock()
	newNode := &PairingHeapNode[V, P]{
		id:       p.curID,
		value:    value,
		priority: priority,
	}
	p.elements[newNode.id] = newNode
	p.root = p.meld(newNode, p.root)
	p.curID++
	p.size++
}

// PairingNode represents a node in the simple pairing heap.
// Unlike PairingHeapNode, this node does not have an ID or parent/prevSibling
// pointers, making it simpler but less feature-rich.
type PairingNode[V any, P any] struct {
	value       V
	priority    P
	firstChild  *PairingNode[V, P]
	nextSibling *PairingNode[V, P]
}

// Value returns the value stored in the node.
func (n *PairingNode[V, P]) Value() V { return n.value }

// Priority returns the priority of the node.
func (n *PairingNode[V, P]) Priority() P { return n.priority }

// SimplePairingHeap implements a basic pairing heap without node tracking.
// It maintains a multi-way tree structure but does not support node updates
// or removal of arbitrary nodes. This implementation is simpler but less
// feature-rich than PairingHeap.
type SimplePairingHeap[V any, P any] struct {
	root *PairingNode[V, P]
	cmp  func(a, b P) bool
	size int
	lock sync.RWMutex
}

// Clone creates a shallow copy of the heap, sharing the same nodes (no
// duplication). The new heap will have the same root, comparison function,
// and size as the original.
func (p *SimplePairingHeap[V, P]) Clone() *SimplePairingHeap[V, P] {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return &SimplePairingHeap[V, P]{root: p.root, cmp: p.cmp, size: p.size}
}

// Clear removes all elements by resetting root to nil and size to zero.
func (p *SimplePairingHeap[V, P]) Clear() {
	p.lock.Lock()
	p.root = nil
	p.size = 0
	p.lock.Unlock()
}

// Length returns the current number of elements in the heap.
func (p *SimplePairingHeap[V, P]) Length() int {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.size
}

// IsEmpty returns true if the heap contains no elements.
func (p *SimplePairingHeap[V, P]) IsEmpty() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.size == 0
}

// peek is an internal method that returns the root node's value and priority without removing it.
// Returns nil if the heap is empty.
func (p *SimplePairingHeap[V, P]) peek() *HeapPair[V, P] {
	if p.size == 0 {
		return nil
	}
	return &HeapPair[V, P]{value: p.root.value, priority: p.root.priority}
}

// Peek returns a pointer to the root node without removing it.
// Returns nil if the heap is empty.
func (p *SimplePairingHeap[V, P]) Peek() *HeapPair[V, P] {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.peek()
}

// PeekValue returns a pointer to the value at the root without removing it.
// Returns nil if the heap is empty.
func (p *SimplePairingHeap[V, P]) PeekValue() *V {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if node := p.peek(); node != nil {
		val := node.Value()
		return &val
	}
	return nil
}

// PeekPriority returns a pointer to the priority at the root without removing it.
// Returns nil if the heap is empty.
func (p *SimplePairingHeap[V, P]) PeekPriority() *P {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if node := p.peek(); node != nil {
		pri := node.Priority()
		return &pri
	}
	return nil
}

// meld links two pairing-heap trees and returns the new root.
// The tree with the smaller priority (according to cmp) becomes the new root,
// and the other tree becomes its first child. The nextSibling pointer of the
// new child is set to the original first child of the new root.
// The prevChild pointer is updated to maintain the doubly-linked child list.
func (p *SimplePairingHeap[V, P]) meld(new *PairingNode[V, P], root *PairingNode[V, P]) *PairingNode[V, P] {
	if root == nil {
		return new
	}

	if new == nil {
		return root
	}

	newRoot := root

	if p.cmp(new.priority, newRoot.priority) {
		newRoot.nextSibling = new.firstChild
		new.firstChild = newRoot
		newRoot = new
	} else {
		new.nextSibling = newRoot.firstChild
		newRoot.firstChild = new
	}
	return newRoot
}

// merge performs the two-pass pairing process on the sibling list.
// It pairs adjacent siblings, melds them, and recursively merges the
// remaining siblings. This is used during Pop to combine the root's
// children into a new heap.
func (p *SimplePairingHeap[V, P]) merge(node *PairingNode[V, P]) *PairingNode[V, P] {
	if node == nil || node.nextSibling == nil {
		return node
	}

	firstNode := node
	secondNode := node.nextSibling
	remaining := node.nextSibling.nextSibling

	firstNode.nextSibling = nil
	secondNode.nextSibling = nil

	return p.meld(p.meld(firstNode, secondNode), p.merge(remaining))
}

// Pop removes and returns a pointer to the value at the root.
// It then merges the root's children to form the new heap.
// Returns nil if the heap is empty.
func (p *SimplePairingHeap[V, P]) Pop() *HeapPair[V, P] {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.size == 0 {
		return nil
	}

	rootNode := p.pop()
	return &HeapPair[V, P]{value: rootNode.value, priority: rootNode.priority}
}

// pop is an internal method that removes the root node and returns it.
// It handles the common logic of removing the root and merging children.
// Returns nil if the heap is empty.
func (p *SimplePairingHeap[V, P]) pop() *PairingNode[V, P] {
	if p.size == 0 {
		return nil
	}

	rootNode := p.root
	p.root = p.merge(p.root.firstChild)
	p.size--
	return rootNode
}

// PopValue removes and returns a pointer to just the value at the root.
// It then merges the root's children to form the new heap.
// Returns nil if the heap is empty.
func (p *SimplePairingHeap[V, P]) PopValue() *V {
	p.lock.Lock()
	defer p.lock.Unlock()
	if rootNode := p.pop(); rootNode != nil {
		val := rootNode.value
		return &val
	}
	return nil
}

// PopPriority removes and returns a pointer to just the priority at the root.
// It then merges the root's children to form the new heap.
// Returns nil if the heap is empty.
func (p *SimplePairingHeap[V, P]) PopPriority() *P {
	p.lock.Lock()
	defer p.lock.Unlock()
	if rootNode := p.pop(); rootNode != nil {
		pri := rootNode.priority
		return &pri
	}
	return nil
}

// Insert adds a new element with its priority by creating a single-node heap
// and melding it with the existing root. The new node becomes the root if
// its priority is smaller than the current root's priority.
func (p *SimplePairingHeap[V, P]) Insert(value V, priority P) {
	p.lock.Lock()
	defer p.lock.Unlock()
	newNode := &PairingNode[V, P]{
		value:    value,
		priority: priority,
	}
	p.root = p.meld(newNode, p.root)
	p.size++
}

// MergeWith combines another pairing heap into this one by melding their
// roots and updating the size. The comparison function of this heap is used
// for determining the new root.
func (p *SimplePairingHeap[V, P]) MergeWith(heap *SimplePairingHeap[V, P]) {
	heap.lock.Lock()
	defer heap.lock.Unlock()
	p.lock.Lock()
	defer p.lock.Unlock()
	p.root = p.meld(heap.root, p.root)
	p.size = p.size + heap.size
}
