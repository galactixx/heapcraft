package heapcraft

import "github.com/google/uuid"

// clearNodeLinks resets all the linking pointers of a node to nil.
// This is used when removing a node from its current position in the heap
// before reinserting it elsewhere.
func clearNodeLinks[V any, P any](node *pairingHeapNode[V, P]) {
	node.nextSibling = nil
	node.parent = nil
	node.prevSibling = nil
}

// NewPairingHeap creates a new pairing heap from a slice of HeapPairs.
// The heap is initialized with the provided elements and uses the given comparison
// function to determine heap order. The comparison function determines the heap order (min or max).
// Returns an empty heap if the input slice is empty.
func NewPairingHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool) *PairingHeap[V, P] {
	elements := make(map[string]*pairingHeapNode[V, P])
	heap := PairingHeap[V, P]{cmp: cmp, size: 0, elements: elements}
	if len(data) == 0 {
		return &heap
	}

	for i := range data {
		heap.Push(data[i].value, data[i].priority)
	}
	return &heap
}

// NewSimplePairingHeap creates a new simple pairing heap from a slice of HeapPairs.
// Unlike PairingHeap, this implementation does not track node IDs or support
// node updates. It uses the provided comparison function to determine heap order (min or max).
// Returns an empty heap if the input slice is empty.
func NewSimplePairingHeap[V any, P any](data []HeapNode[V, P], cmp func(a, b P) bool) *SimplePairingHeap[V, P] {
	heap := SimplePairingHeap[V, P]{cmp: cmp, size: 0}
	if len(data) == 0 {
		return &heap
	}

	for i := range data {
		heap.Push(data[i].value, data[i].priority)
	}
	return &heap
}

// pairingHeapNode represents a node in the pairing heap data structure.
// Each node contains a value, priority, and maintains links to its parent,
// children, and siblings. The node also has a unique identifier for tracking.
// The doubly-linked sibling list allows for efficient node removal and updates.
type pairingHeapNode[V any, P any] struct {
	id          string
	value       V
	priority    P
	parent      *pairingHeapNode[V, P]
	firstChild  *pairingHeapNode[V, P]
	nextSibling *pairingHeapNode[V, P]
	prevSibling *pairingHeapNode[V, P]
}

// ID returns the unique identifier of the node.
// This identifier is used for tracking and updating nodes in the heap.
func (n *pairingHeapNode[V, P]) ID() string { return n.id }

// Value returns the value stored in the node.
func (n *pairingHeapNode[V, P]) Value() V { return n.value }

// Priority returns the priority of the node.
func (n *pairingHeapNode[V, P]) Priority() P { return n.priority }

// PairingHeap implements a pairing heap data structure with node tracking.
// It maintains a multi-way tree structure where each node can have multiple children.
// The heap supports efficient insertion, deletion, and priority updates of nodes.
// Nodes are tracked by unique IDs, allowing for O(1) access and updates.
type PairingHeap[V any, P any] struct {
	root     *pairingHeapNode[V, P]
	cmp      func(a, b P) bool
	size     int
	elements map[string]*pairingHeapNode[V, P]
}

// UpdateValue updates the value of a node with the given ID.
// Returns an error if the ID does not exist in the heap.
// The heap structure remains unchanged as this operation only modifies the value.
func (p *PairingHeap[V, P]) UpdateValue(id string, value V) error {
	if _, exists := p.elements[id]; !exists {
		return ErrNodeNotFound
	}

	p.elements[id].value = value
	return nil
}

// UpdatePriority updates the priority of a node with the given ID.
// Returns an error if the ID does not exist in the heap.
// The node is removed from its current position and reinserted into the heap
// to maintain the heap property. This operation may change the heap structure.
func (p *PairingHeap[V, P]) UpdatePriority(id string, priority P) error {
	if _, exists := p.elements[id]; !exists {
		return ErrNodeNotFound
	}

	updated := p.elements[id]
	updated.priority = priority

	switch {
	case updated.id == p.root.id:
		newRoot := updated.firstChild
		if newRoot != nil {
			newRoot.prevSibling, newRoot.parent = nil, nil
		}
		updated.firstChild = nil
		p.root = p.merge(newRoot)

	case updated.prevSibling != nil:
		prev, next := updated.prevSibling, updated.nextSibling
		if next != nil {
			next.prevSibling = prev
		}

		prev.nextSibling = next
	default:
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

// Clone creates a deep copy of the heap structure and nodes. If values or
// priorities are reference types, those reference values are shared between the
// original and cloned heaps.
func (p *PairingHeap[V, P]) Clone() *PairingHeap[V, P] {
	elements := make(map[string]*pairingHeapNode[V, P], len(p.elements))
	for _, node := range p.elements {
		elements[node.id] = &pairingHeapNode[V, P]{
			id:          node.id,
			value:       node.value,
			priority:    node.priority,
			parent:      node.parent,
			firstChild:  node.firstChild,
			nextSibling: node.nextSibling,
			prevSibling: node.prevSibling,
		}
	}

	for _, node := range elements {
		if node.parent != nil {
			node.parent = elements[node.parent.id]
		}
		if node.firstChild != nil {
			node.firstChild = elements[node.firstChild.id]
		}
		if node.nextSibling != nil {
			node.nextSibling = elements[node.nextSibling.id]
		}
		if node.prevSibling != nil {
			node.prevSibling = elements[node.prevSibling.id]
		}
	}

	return &PairingHeap[V, P]{
		root:     elements[p.root.id],
		cmp:      p.cmp,
		size:     p.size,
		elements: elements,
	}
}

// Clear removes all elements from the heap.
// Resets the root to nil, size to zero, and initializes a new empty element map.
// The next node ID is reset to 1.
func (p *PairingHeap[V, P]) Clear() {
	p.root = nil
	p.size = 0
	p.elements = make(map[string]*pairingHeapNode[V, P], 0)
}

// Length returns the current number of elements in the heap.
func (p *PairingHeap[V, P]) Length() int { return p.size }

// IsEmpty returns true if the heap contains no elements.
func (p *PairingHeap[V, P]) IsEmpty() bool { return p.size == 0 }

// peek is an internal method that returns the root node's value and priority without removing it.
// Returns nil and an error if the heap is empty.
func (p *PairingHeap[V, P]) peek() (Node[V, P], error) {
	if p.size == 0 {
		return nil, ErrHeapEmpty
	}
	return p.root, nil
}

// Peek returns a HeapNode containing the value and priority
// of the root node without removing it. Returns nil and an error if the heap is empty.
func (p *PairingHeap[V, P]) Peek() (Node[V, P], error) { return p.peek() }

// PeekValue returns the value at the root without removing it.
// Returns zero value and an error if the heap is empty.
func (p *PairingHeap[V, P]) PeekValue() (V, error) {
	return valueFromNode(p.peek())
}

// PeekPriority returns the priority at the root without removing it.
// Returns zero value and an error if the heap is empty.
func (p *PairingHeap[V, P]) PeekPriority() (P, error) {
	return priorityFromNode(p.peek())
}

// get is an internal method that retrieves a HeapNode for the node with the given ID.
// Returns an error if the ID does not exist in the heap.
func (p *PairingHeap[V, P]) get(id string) (Node[V, P], error) {
	node, exists := p.elements[id]
	if !exists {
		return nil, ErrNodeNotFound
	}
	return node, nil
}

// Get retrieves a HeapNode for the node with the given ID.
// Returns an error if the ID does not exist in the heap.
func (p *PairingHeap[V, P]) Get(id string) (Node[V, P], error) { return p.get(id) }

// GetValue retrieves the value of the node with the given ID.
// Returns zero value and an error if the ID does not exist in the heap.
func (p *PairingHeap[V, P]) GetValue(id string) (V, error) {
	return valueFromNode(p.get(id))
}

// GetPriority retrieves the priority of the node with the given ID.
// Returns zero value and an error if the ID does not exist in the heap.
func (p *PairingHeap[V, P]) GetPriority(id string) (P, error) {
	return priorityFromNode(p.get(id))
}

// meld combines two pairing heap trees into a single tree.
// The tree with the higher priority (according to cmp) becomes the root,
// and the other tree becomes its first child. The operation maintains
// the doubly-linked sibling list structure.
// Returns the new root of the combined tree.
func (p *PairingHeap[V, P]) meld(new *pairingHeapNode[V, P], root *pairingHeapNode[V, P]) *pairingHeapNode[V, P] {
	if root == nil {
		return new
	}

	if new == nil {
		return root
	}

	var prior, noPrior *pairingHeapNode[V, P]

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
func (p *PairingHeap[V, P]) merge(node *pairingHeapNode[V, P]) *pairingHeapNode[V, P] {
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
// Returns nil and an error if the heap is empty.
func (p *PairingHeap[V, P]) pop() (Node[V, P], error) {
	if p.size == 0 {
		return nil, ErrHeapEmpty
	}

	rootNode := p.root
	p.root = p.merge(p.root.firstChild)
	p.size--
	delete(p.elements, rootNode.id)
	return rootNode, nil
}

// Pop removes and returns a HeapNode containing the value and priority
// of the root node. The root's children are merged to form the new heap.
// Returns nil and an error if the heap is empty.
func (p *PairingHeap[V, P]) Pop() (Node[V, P], error) { return p.pop() }

// PopValue removes and returns just the value at the root.
// The root's children are merged to form the new heap.
// Returns zero value and an error if the heap is empty.
func (p *PairingHeap[V, P]) PopValue() (V, error) {
	return valueFromNode(p.pop())
}

// PopPriority removes and returns just the priority at the root.
// The root's children are merged to form the new heap.
// Returns zero value and an error if the heap is empty.
func (p *PairingHeap[V, P]) PopPriority() (P, error) {
	return priorityFromNode(p.pop())
}

// Push adds a new element with the given value and priority to the heap.
// A new node is created with a unique ID and melded with the existing root.
// The new node becomes the root if its priority is higher than the current root's.
// Returns the ID of the inserted node.
func (p *PairingHeap[V, P]) Push(value V, priority P) string {
	newNode := &pairingHeapNode[V, P]{
		id:       uuid.New().String(),
		value:    value,
		priority: priority,
	}
	p.elements[newNode.id] = newNode
	p.root = p.meld(newNode, p.root)
	p.size++
	return newNode.id
}

// pairingNode represents a node in the simple pairing heap.
// Unlike pairingHeapNode, this node does not have an ID or parent/prevSibling
// pointers, making it simpler but less feature-rich.
type pairingNode[V any, P any] struct {
	value       V
	priority    P
	firstChild  *pairingNode[V, P]
	nextSibling *pairingNode[V, P]
}

// Value returns the value stored in the node.
func (n *pairingNode[V, P]) Value() V { return n.value }

// Priority returns the priority of the node.
func (n *pairingNode[V, P]) Priority() P { return n.priority }

// SimplePairingHeap implements a basic pairing heap without node tracking.
// It maintains a multi-way tree structure but does not support node updates
// or removal of arbitrary nodes. This implementation is simpler but less
// feature-rich than PairingHeap.
type SimplePairingHeap[V any, P any] struct {
	root *pairingNode[V, P]
	cmp  func(a, b P) bool
	size int
}

// cloneNode creates a deep copy of a pairing node.
// It recursively clones the first child and next sibling.
func (p *SimplePairingHeap[V, P]) cloneNode(node *pairingNode[V, P]) *pairingNode[V, P] {
	if node == nil {
		return nil
	}

	return &pairingNode[V, P]{
		value:       node.value,
		priority:    node.priority,
		firstChild:  p.cloneNode(node.firstChild),
		nextSibling: p.cloneNode(node.nextSibling),
	}
}

// Clone creates a deep copy of the heap structure and nodes. If values or
// priorities are reference types, those reference values are shared between the
// original and cloned heaps.
func (p *SimplePairingHeap[V, P]) Clone() *SimplePairingHeap[V, P] {
	return &SimplePairingHeap[V, P]{
		root: p.cloneNode(p.root),
		cmp:  p.cmp,
		size: p.size,
	}
}

// Clear removes all elements from the simple heap.
// The heap is ready for new insertions after clearing.
func (p *SimplePairingHeap[V, P]) Clear() {
	p.root = nil
	p.size = 0
}

// Length returns the current number of elements in the heap.
func (p *SimplePairingHeap[V, P]) Length() int { return p.size }

// IsEmpty returns true if the heap contains no elements.
func (p *SimplePairingHeap[V, P]) IsEmpty() bool { return p.size == 0 }

// peek is an internal method that returns the root node's value and priority without removing it.
// Returns nil and an error if the heap is empty.
func (p *SimplePairingHeap[V, P]) peek() (SimpleNode[V, P], error) {
	if p.size == 0 {
		return nil, ErrHeapEmpty
	}
	return p.root, nil
}

// Peek returns a HeapNode containing the value and priority
// of the root node without removing it. Returns nil and an error if the heap is empty.
func (p *SimplePairingHeap[V, P]) Peek() (SimpleNode[V, P], error) {
	return p.peek()
}

// PeekValue returns the value at the root without removing it.
// Returns zero value and an error if the heap is empty.
func (p *SimplePairingHeap[V, P]) PeekValue() (V, error) {
	return valueFromNode(p.peek())
}

// PeekPriority returns the priority at the root without removing it.
// Returns zero value and an error if the heap is empty.
func (p *SimplePairingHeap[V, P]) PeekPriority() (P, error) {
	return priorityFromNode(p.peek())
}

// meld links two pairing-heap trees and returns the new root.
// The tree with the higher priority (according to cmp) becomes the new root,
// and the other tree becomes its first child. The nextSibling pointer of the
// new child is set to the original first child of the new root.
func (p *SimplePairingHeap[V, P]) meld(new *pairingNode[V, P], root *pairingNode[V, P]) *pairingNode[V, P] {
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
func (p *SimplePairingHeap[V, P]) merge(node *pairingNode[V, P]) *pairingNode[V, P] {
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

// pop is an internal method that removes the root node and returns it.
// It handles the common logic of removing the root and merging children.
// Returns nil and an error if the heap is empty.
func (p *SimplePairingHeap[V, P]) pop() (SimpleNode[V, P], error) {
	if p.size == 0 {
		return nil, ErrHeapEmpty
	}

	rootNode := p.root
	p.root = p.merge(p.root.firstChild)
	p.size--
	return rootNode, nil
}

// Pop removes and returns a HeapNode containing the value and priority
// of the root node. The root's children are merged to form the new heap.
// Returns nil and an error if the heap is empty.
func (p *SimplePairingHeap[V, P]) Pop() (SimpleNode[V, P], error) {
	return p.pop()
}

// PopValue removes and returns just the value at the root.
// The root's children are merged to form the new heap.
// Returns zero value and an error if the heap is empty.
func (p *SimplePairingHeap[V, P]) PopValue() (V, error) {
	return valueFromNode(p.pop())
}

// PopPriority removes and returns just the priority at the root.
// The root's children are merged to form the new heap.
// Returns zero value and an error if the heap is empty.
func (p *SimplePairingHeap[V, P]) PopPriority() (P, error) {
	return priorityFromNode(p.pop())
}

// Push adds a new element with its priority by creating a single-node heap
// and melding it with the existing root. The new node becomes the root if
// its priority is higher than the current root's priority.
func (p *SimplePairingHeap[V, P]) Push(value V, priority P) {
	newNode := &pairingNode[V, P]{value: value, priority: priority}
	p.root = p.meld(newNode, p.root)
	p.size++
}
