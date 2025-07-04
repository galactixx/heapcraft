package heapcraft

// clearNodeLinks resets all the linking pointers of a node to nil.
// This is used when removing a node from its current position in the heap
// before reinserting it elsewhere.
func clearNodeLinks[V any, P any](node *pairingHeapNode[V, P]) {
	node.nextSibling = nil
	node.parent = nil
	node.prevSibling = nil
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

// Value returns the value stored in the node.
func (n *pairingHeapNode[V, P]) Value() V { return n.value }

// Priority returns the priority of the node.
func (n *pairingHeapNode[V, P]) Priority() P { return n.priority }

// FullPairingHeap implements a pairing heap data structure with node tracking.
// It maintains a multi-way tree structure where each node can have multiple children.
// The heap supports efficient insertion, deletion, and priority updates of nodes.
// Nodes are tracked by unique IDs, allowing for O(1) access and updates.
type FullPairingHeap[V any, P any] struct {
	root     *pairingHeapNode[V, P]
	cmp      func(a, b P) bool
	size     int
	elements map[string]*pairingHeapNode[V, P]
	pool     pool[*pairingHeapNode[V, P]]
	idGen    IDGenerator
}

// UpdateValue updates the value of a node with the given ID.
// Returns an error if the ID does not exist in the heap.
// The heap structure remains unchanged as this operation only modifies the value.
func (p *FullPairingHeap[V, P]) UpdateValue(id string, value V) error {
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
func (p *FullPairingHeap[V, P]) UpdatePriority(id string, priority P) error {
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
func (p *FullPairingHeap[V, P]) Clone() *FullPairingHeap[V, P] {
	elements := make(map[string]*pairingHeapNode[V, P], len(p.elements))
	for _, node := range p.elements {
		cloned := p.pool.Get()
		cloned.id = node.id
		cloned.value = node.value
		cloned.priority = node.priority
		cloned.parent = node.parent
		cloned.firstChild = node.firstChild
		cloned.nextSibling = node.nextSibling
		elements[node.id] = cloned
	}

	// Re-link the nodes to the new heap structure to avoid reference
	// issues
	for _, node := range elements {
		// Re-link the parent pointer to the new heap
		if node.parent != nil {
			node.parent = elements[node.parent.id]
		}

		// Re-link the first child pointer to the new heap
		if node.firstChild != nil {
			node.firstChild = elements[node.firstChild.id]
		}

		// Re-link the next sibling pointer to the new heap
		if node.nextSibling != nil {
			node.nextSibling = elements[node.nextSibling.id]
		}

		// Re-link the previous sibling pointer to the new heap
		if node.prevSibling != nil {
			node.prevSibling = elements[node.prevSibling.id]
		}
	}

	return &FullPairingHeap[V, P]{
		root:     elements[p.root.id],
		cmp:      p.cmp,
		size:     p.size,
		elements: elements,
		pool:     p.pool,
		idGen:    p.idGen,
	}
}

// Clear removes all elements from the heap.
// Resets the root to nil, size to zero, and initializes a new empty element map.
// The next node ID is reset to 1.
func (p *FullPairingHeap[V, P]) Clear() {
	p.root = nil
	p.size = 0
	p.elements = make(map[string]*pairingHeapNode[V, P], 0)
}

// Length returns the current number of elements in the heap.
func (p *FullPairingHeap[V, P]) Length() int { return p.size }

// IsEmpty returns true if the heap contains no elements.
func (p *FullPairingHeap[V, P]) IsEmpty() bool { return p.size == 0 }

// peek is an internal method that returns the root node's value and priority without removing it.
// Returns nil and an error if the heap is empty.
func (p *FullPairingHeap[V, P]) peek() (V, P, error) {
	if p.size == 0 {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrHeapEmpty
	}
	v, pr := p.root.value, p.root.priority
	return v, pr, nil
}

// Peek returns a HeapNode containing the value and priority
// of the root node without removing it. Returns nil and an error if the heap is empty.
func (p *FullPairingHeap[V, P]) Peek() (V, P, error) { return p.peek() }

// PeekValue returns the value at the root without removing it.
// Returns zero value and an error if the heap is empty.
func (p *FullPairingHeap[V, P]) PeekValue() (V, error) {
	return valueFromNode(p.peek())
}

// PeekPriority returns the priority at the root without removing it.
// Returns zero value and an error if the heap is empty.
func (p *FullPairingHeap[V, P]) PeekPriority() (P, error) {
	return priorityFromNode(p.peek())
}

// get is an internal method that retrieves a HeapNode for the node with the given ID.
// Returns an error if the ID does not exist in the heap.
func (p *FullPairingHeap[V, P]) get(id string) (V, P, error) {
	node, exists := p.elements[id]
	if !exists {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrNodeNotFound
	}
	v, pr := node.value, node.priority
	return v, pr, nil
}

// Get retrieves a HeapNode for the node with the given ID.
// Returns an error if the ID does not exist in the heap.
func (p *FullPairingHeap[V, P]) Get(id string) (V, P, error) { return p.get(id) }

// GetValue retrieves the value of the node with the given ID.
// Returns zero value and an error if the ID does not exist in the heap.
func (p *FullPairingHeap[V, P]) GetValue(id string) (V, error) {
	return valueFromNode(p.get(id))
}

// GetPriority retrieves the priority of the node with the given ID.
// Returns zero value and an error if the ID does not exist in the heap.
func (p *FullPairingHeap[V, P]) GetPriority(id string) (P, error) {
	return priorityFromNode(p.get(id))
}

// meld combines two pairing heap trees into a single tree.
// The tree with the higher priority (according to cmp) becomes the root,
// and the other tree becomes its first child. The operation maintains
// the doubly-linked sibling list structure.
// Returns the new root of the combined tree.
func (p *FullPairingHeap[V, P]) meld(new *pairingHeapNode[V, P], root *pairingHeapNode[V, P]) *pairingHeapNode[V, P] {
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
func (p *FullPairingHeap[V, P]) merge(node *pairingHeapNode[V, P]) *pairingHeapNode[V, P] {
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
func (p *FullPairingHeap[V, P]) pop() (V, P, error) {
	if p.size == 0 {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrHeapEmpty
	}

	removed := p.root
	p.root = p.merge(p.root.firstChild)
	p.size--
	removed.firstChild = nil
	removed.nextSibling = nil
	removed.parent = nil
	removed.prevSibling = nil
	delete(p.elements, removed.id)
	v, pr := removed.value, removed.priority
	p.pool.Put(removed)
	return v, pr, nil
}

// Pop removes and returns a HeapNode containing the value and priority
// of the root node. The root's children are merged to form the new heap.
// Returns nil and an error if the heap is empty.
func (p *FullPairingHeap[V, P]) Pop() (V, P, error) { return p.pop() }

// PopValue removes and returns just the value at the root.
// The root's children are merged to form the new heap.
// Returns zero value and an error if the heap is empty.
func (p *FullPairingHeap[V, P]) PopValue() (V, error) {
	return valueFromNode(p.pop())
}

// PopPriority removes and returns just the priority at the root.
// The root's children are merged to form the new heap.
// Returns zero value and an error if the heap is empty.
func (p *FullPairingHeap[V, P]) PopPriority() (P, error) {
	return priorityFromNode(p.pop())
}

// Push adds a new element with the given value and priority to the heap.
// A new node is created with a unique ID and melded with the existing root.
// The new node becomes the root if its priority is higher than the current root's.
// Returns the ID of the inserted node.
func (p *FullPairingHeap[V, P]) Push(value V, priority P) (string, error) {
	newNode := p.pool.Get()
	newNode.id = p.idGen.Next()
	if _, exists := p.elements[newNode.id]; exists {
		return "", ErrIDGenerationFailed
	}

	newNode.value = value
	newNode.priority = priority
	p.elements[newNode.id] = newNode
	p.root = p.meld(newNode, p.root)
	p.size++
	return newNode.id, nil
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

// PairingHeap implements a basic pairing heap without node tracking.
// It maintains a multi-way tree structure but does not support node updates
// or removal of arbitrary nodes. This implementation is simpler but less
// feature-rich than FullPairingHeap.
type PairingHeap[V any, P any] struct {
	root *pairingNode[V, P]
	cmp  func(a, b P) bool
	size int
	pool pool[*pairingNode[V, P]]
}

// cloneNode creates a deep copy of a pairing node.
// It recursively clones the first child and next sibling.
func (p *PairingHeap[V, P]) cloneNode(node *pairingNode[V, P]) *pairingNode[V, P] {
	if node == nil {
		return nil
	}

	cloned := p.pool.Get()
	cloned.value = node.value
	cloned.priority = node.priority
	cloned.firstChild = p.cloneNode(node.firstChild)
	cloned.nextSibling = p.cloneNode(node.nextSibling)
	return cloned
}

// Clone creates a deep copy of the heap structure and nodes. If values or
// priorities are reference types, those reference values are shared between the
// original and cloned heaps.
func (p *PairingHeap[V, P]) Clone() *PairingHeap[V, P] {
	return &PairingHeap[V, P]{
		root: p.cloneNode(p.root),
		cmp:  p.cmp,
		size: p.size,
		pool: p.pool,
	}
}

// Clear removes all elements from the simple heap.
// The heap is ready for new insertions after clearing.
func (p *PairingHeap[V, P]) Clear() {
	p.root = nil
	p.size = 0
}

// Length returns the current number of elements in the heap.
func (p *PairingHeap[V, P]) Length() int { return p.size }

// IsEmpty returns true if the simple heap contains no elements.
func (p *PairingHeap[V, P]) IsEmpty() bool { return p.size == 0 }

// peek is an internal method that returns the root node's value and priority without removing it.
// Returns nil and an error if the heap is empty.
func (p *PairingHeap[V, P]) peek() (V, P, error) {
	if p.size == 0 {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrHeapEmpty
	}
	v, pr := p.root.value, p.root.priority
	return v, pr, nil
}

// Peek returns a HeapNode containing the value and priority
// of the root node without removing it. Returns nil and an error if the heap is empty.
func (p *PairingHeap[V, P]) Peek() (V, P, error) { return p.peek() }

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

// meld links two pairing-heap trees and returns the new root.
// The tree with the higher priority (according to cmp) becomes the new root,
// and the other tree becomes its first child. The nextSibling pointer of the
// new child is set to the original first child of the new root.
func (p *PairingHeap[V, P]) meld(new *pairingNode[V, P], root *pairingNode[V, P]) *pairingNode[V, P] {
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
func (p *PairingHeap[V, P]) merge(node *pairingNode[V, P]) *pairingNode[V, P] {
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
func (p *PairingHeap[V, P]) pop() (V, P, error) {
	if p.size == 0 {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrHeapEmpty
	}

	removed := p.root
	p.root = p.merge(p.root.firstChild)
	removed.firstChild = nil
	removed.nextSibling = nil
	v, pr := removed.value, removed.priority
	p.pool.Put(removed)
	p.size--
	return v, pr, nil
}

// Pop removes and returns a HeapNode containing the value and priority
// of the root node. The root's children are merged to form the new heap.
// Returns nil and an error if the heap is empty.
func (p *PairingHeap[V, P]) Pop() (V, P, error) { return p.pop() }

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

// Push adds a new element with its priority by creating a single-node heap
// and melding it with the existing root. The new node becomes the root if
// its priority is higher than the current root's priority.
func (p *PairingHeap[V, P]) Push(value V, priority P) {
	newNode := p.pool.Get()
	newNode.value = value
	newNode.priority = priority
	p.root = p.meld(newNode, p.root)
	p.size++
}
