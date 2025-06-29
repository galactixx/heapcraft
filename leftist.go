package heapcraft

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
func (l leftistQueue[N]) remainingElements() int { return l.size }

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
	id       string
	value    V
	priority P
	parent   *leftistHeapNode[V, P]
	left     *leftistHeapNode[V, P]
	right    *leftistHeapNode[V, P]
	s        int
}

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
	elements map[string]*leftistHeapNode[V, P]
	pool     pool[*leftistHeapNode[V, P]]
	idGen    IDGenerator
}

// UpdateValue changes the value of the node with the given ID.
// Returns an error if the ID doesn't exist in the heap.
func (l *LeftistHeap[V, P]) UpdateValue(id string, value V) error {
	if _, exists := l.elements[id]; !exists {
		return ErrNodeNotFound
	}

	l.elements[id].value = value
	return nil
}

// UpdatePriority changes the priority of the node with the given ID and
// restructures the heap to maintain the heap property.
// Returns an error if the ID doesn't exist in the heap.
func (l *LeftistHeap[V, P]) UpdatePriority(id string, priority P) error {
	if _, exists := l.elements[id]; !exists {
		return ErrNodeNotFound
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

// Clone creates a deep copy of the heap structure and nodes. If values or
// priorities are reference types, those reference values are shared between the
// original and cloned heaps.
func (l *LeftistHeap[V, P]) Clone() *LeftistHeap[V, P] {
	elements := make(map[string]*leftistHeapNode[V, P], len(l.elements))
	for _, node := range l.elements {
		cloned := l.pool.Get()
		cloned.id = node.id
		cloned.value = node.value
		cloned.priority = node.priority
		cloned.parent = node.parent
		cloned.left = node.left
		cloned.right = node.right
		cloned.s = node.s
		elements[node.id] = cloned
	}

	// Re-assign parent, left, and right pointers to the cloned nodes.
	for _, node := range elements {
		// Re-assign parent, left, and right pointers to the cloned nodes.
		if node.parent != nil {
			node.parent = elements[node.parent.id]
		}

		// Re-assign left pointer to the cloned node if it exists.
		if node.left != nil {
			node.left = elements[node.left.id]
		}

		// Re-assign right pointer to the cloned node if it exists.
		if node.right != nil {
			node.right = elements[node.right.id]
		}
	}

	return &LeftistHeap[V, P]{
		root:     elements[l.root.id],
		cmp:      l.cmp,
		size:     l.size,
		elements: elements,
		pool:     l.pool,
		idGen:    l.idGen,
	}
}

// Clear removes all elements from the heap and resets its state.
// The heap is ready for new insertions after clearing.
func (l *LeftistHeap[V, P]) Clear() {
	l.root = nil
	l.size = 0
	l.elements = make(map[string]*leftistHeapNode[V, P])
}

// Length returns the current number of elements in the heap.
func (l *LeftistHeap[V, P]) Length() int { return l.size }

// IsEmpty returns true if the heap contains no elements.
func (l *LeftistHeap[V, P]) IsEmpty() bool { return l.size == 0 }

// peek is an internal method that returns the root node without removing it.
// Returns nil and an error if the heap is empty.
func (l *LeftistHeap[V, P]) peek() (V, P, error) {
	if l.size == 0 {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrHeapEmpty
	}
	v, p := l.root.value, l.root.priority
	return v, p, nil
}

// Peek returns the minimum element without removing it.
// Returns nil and an error if the heap is empty.
func (l *LeftistHeap[V, P]) Peek() (V, P, error) { return l.peek() }

// PeekValue returns the value at the root without removing it.
// Returns zero value and an error if the heap is empty.
func (l *LeftistHeap[V, P]) PeekValue() (V, error) {
	return valueFromNode(l.peek())
}

// PeekPriority returns the priority at the root without removing it.
// Returns zero value and an error if the heap is empty.
func (l *LeftistHeap[V, P]) PeekPriority() (P, error) {
	return priorityFromNode(l.peek())
}

// get is an internal method that retrieves a node with the given ID.
// Returns an error if the ID doesn't exist in the heap.
func (l *LeftistHeap[V, P]) get(id string) (V, P, error) {
	if node, exists := l.elements[id]; exists {
		return node.value, node.priority, nil
	}
	v, p := zeroValuePair[V, P]()
	return v, p, ErrNodeNotFound
}

// Get returns the element associated with the given ID.
// Returns an error if the ID doesn't exist in the heap.
func (l *LeftistHeap[V, P]) Get(id string) (V, P, error) { return l.get(id) }

// GetValue returns the value associated with the given ID.
// Returns zero value and an error if the ID doesn't exist in the heap.
func (l *LeftistHeap[V, P]) GetValue(id string) (V, error) {
	return valueFromNode(l.get(id))
}

// GetPriority returns the priority associated with the given ID.
// Returns zero value and an error if the ID doesn't exist in the heap.
func (l *LeftistHeap[V, P]) GetPriority(id string) (P, error) {
	return priorityFromNode(l.get(id))
}

// Pop removes and returns the minimum element from the heap.
// The heap property is restored through merging the root's children.
// Returns nil and an error if the heap is empty.
func (l *LeftistHeap[V, P]) Pop() (V, P, error) { return l.pop() }

// PopValue removes and returns just the value at the root.
// The heap property is restored through merging the root's children.
// Returns zero value and an error if the heap is empty.
func (l *LeftistHeap[V, P]) PopValue() (V, error) {
	return valueFromNode(l.pop())
}

// PopPriority removes and returns just the priority at the root.
// The heap property is restored through merging the root's children.
// Returns zero value and an error if the heap is empty.
func (l *LeftistHeap[V, P]) PopPriority() (P, error) {
	return priorityFromNode(l.pop())
}

// pop is an internal method that removes the root node and returns it.
// Handles the common logic of removing the root and merging its children.
// Returns nil and an error if the heap is empty.
func (l *LeftistHeap[V, P]) pop() (V, P, error) {
	if l.size == 0 {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrHeapEmpty
	}

	rootNode := l.root
	l.root = l.merge(l.root.right, l.root.left)
	if l.root != nil {
		l.root.parent = nil
	}
	delete(l.elements, rootNode.id)
	rootNode.left, rootNode.right, rootNode.parent = nil, nil, nil
	l.size--
	v, p := rootNode.value, rootNode.priority
	l.pool.Put(rootNode)
	return v, p, nil
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
		b.s = 1
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
func (l *LeftistHeap[V, P]) Push(value V, priority P) (string, error) {
	newNode := l.pool.Get()
	newNode.id = l.idGen.Next()
	if _, exists := l.elements[newNode.id]; exists {
		return "", ErrIDGenerationFailed
	}

	newNode.value = value
	newNode.priority = priority
	newNode.s = 1
	l.root = l.merge(newNode, l.root)
	l.elements[newNode.id] = newNode
	l.size++
	return newNode.id, nil
}

// SimpleLeftistHeap implements a basic leftist heap without node tracking.
// Maintains the heap property through the comparison function and
// the leftist property through s-values.
type SimpleLeftistHeap[V any, P any] struct {
	root *leftistNode[V, P]
	cmp  func(a, b P) bool
	size int
	pool pool[*leftistNode[V, P]]
}

// cloneNode creates a deep copy of a leftist node.
// It recursively clones the left and right children.
func (l *SimpleLeftistHeap[V, P]) cloneNode(node *leftistNode[V, P]) *leftistNode[V, P] {
	if node == nil {
		return nil
	}

	cloned := l.pool.Get()
	cloned.value = node.value
	cloned.priority = node.priority
	cloned.s = node.s
	cloned.left = l.cloneNode(node.left)
	cloned.right = l.cloneNode(node.right)
	return cloned
}

// Clone creates a deep copy of the heap structure and nodes. If values or
// priorities are reference types, those reference values are shared between the
// original and cloned heaps.
func (l *SimpleLeftistHeap[V, P]) Clone() *SimpleLeftistHeap[V, P] {
	return &SimpleLeftistHeap[V, P]{
		root: l.cloneNode(l.root),
		cmp:  l.cmp,
		size: l.size,
		pool: l.pool,
	}
}

// Clear removes all elements from the simple heap.
// The heap is ready for new insertions after clearing.
func (l *SimpleLeftistHeap[V, P]) Clear() {
	l.root = nil
	l.size = 0
}

// Length returns the current number of elements in the simple heap.
func (l *SimpleLeftistHeap[V, P]) Length() int { return l.size }

// IsEmpty returns true if the simple heap contains no elements.
func (l *SimpleLeftistHeap[V, P]) IsEmpty() bool { return l.size == 0 }

// peek is an internal method that returns the root node without removing it.
// Returns nil and an error if the heap is empty.
func (l *SimpleLeftistHeap[V, P]) peek() (V, P, error) {
	if l.size == 0 {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrHeapEmpty
	}
	v, p := l.root.value, l.root.priority
	return v, p, nil
}

// Peek returns the minimum element without removing it.
// Returns nil and an error if the heap is empty.
func (l *SimpleLeftistHeap[V, P]) Peek() (V, P, error) { return l.peek() }

// PeekValue returns the value at the root without removing it.
// Returns zero value and an error if the heap is empty.
func (l *SimpleLeftistHeap[V, P]) PeekValue() (V, error) {
	return valueFromNode(l.peek())
}

// PeekPriority returns the priority at the root without removing it.
// Returns zero value and an error if the heap is empty.
func (l *SimpleLeftistHeap[V, P]) PeekPriority() (P, error) {
	return priorityFromNode(l.peek())
}

// pop is an internal method that removes the root node and returns it.
// Handles the common logic of removing the root and merging its children.
// Returns nil and an error if the heap is empty.
func (l *SimpleLeftistHeap[V, P]) pop() (V, P, error) {
	if l.size == 0 {
		v, p := zeroValuePair[V, P]()
		return v, p, ErrHeapEmpty
	}

	removed := l.root
	l.root = l.merge(l.root.right, l.root.left)
	removed.left, removed.right = nil, nil
	l.size--
	v, p := removed.value, removed.priority
	l.pool.Put(removed)
	return v, p, nil
}

// Pop removes and returns the minimum element from the simple heap.
// The heap property is restored through merging the root's children.
// Returns nil and an error if the heap is empty.
func (l *SimpleLeftistHeap[V, P]) Pop() (V, P, error) { return l.pop() }

// PopValue removes and returns just the value at the root.
// The heap property is restored through merging the root's children.
// Returns zero value and an error if the heap is empty.
func (l *SimpleLeftistHeap[V, P]) PopValue() (V, error) {
	return valueFromNode(l.pop())
}

// PopPriority removes and returns just the priority at the root.
// The heap property is restored through merging the root's children.
// Returns zero value and an error if the heap is empty.
func (l *SimpleLeftistHeap[V, P]) PopPriority() (P, error) {
	return priorityFromNode(l.pop())
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
		b.s = 1
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
	newNode := l.pool.Get()
	newNode.value = value
	newNode.priority = priority
	newNode.s = 1
	l.root = l.merge(newNode, l.root)
	l.size++
}
