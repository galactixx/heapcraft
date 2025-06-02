package heapcraft

import (
	"github.com/mohae/deepcopy"
)

// leftistQueue is a simple FIFO queue used for building the
// heap via pairwise merging.
type leftistQueue[T any] struct {
	data []*LeftistNode[T]
	head int
	size int
}

// push adds a node to the end of the queue.
func (l *leftistQueue[T]) push(element *LeftistNode[T]) {
	l.data = append(l.data, element)
	l.size++
}

// remainingElements returns the number of remaining elements.
func (l leftistQueue[T]) remainingElements() int {
	return l.size
}

// length returns the current number of elements in the queue.
func (l leftistQueue[T]) length() int { return len(l.data) }

// pop removes and returns the node at the head of the queue,
// advancing head efficiently.
func (l *leftistQueue[T]) pop() *LeftistNode[T] {
	if l.remainingElements() == 0 {
		return nil
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

// NewLeftistHeap constructs a leftist heap from an initial
// slice of values. It uses a queue to iteratively merge singleton
// nodes until one root remains.
func NewLeftistHeap[T any](data []T, cmp func(a, b T) bool) LeftistHeap[T] {
	if len(data) == 0 {
		return LeftistHeap[T]{cmp: cmp, size: 0}
	}

	n := len(data)
	queueData := make([]*LeftistNode[T], 0, n)
	initQueue := leftistQueue[T]{data: queueData, head: 0, size: 0}

	heap := LeftistHeap[T]{cmp: cmp, size: n}

	for i := range data {
		initQueue.push(&LeftistNode[T]{val: data[i], s: 0})
	}

	for initQueue.remainingElements() > 1 {
		merged := heap.merge(initQueue.pop(), initQueue.pop())
		initQueue.push(merged)
	}

	heap.root = initQueue.pop()
	return heap
}

// LeftistNode represents a node within a leftist heap, storing its
// value, children, and null-path length (s).
type LeftistNode[T any] struct {
	val   T
	left  *LeftistNode[T]
	right *LeftistNode[T]
	s     int
}

// LeftistHeap is the main heap structure, holding a pointer to the
// root node, a comparison function, and the total size.
type LeftistHeap[T any] struct {
	root *LeftistNode[T]
	cmp  func(a, b T) bool
	size int
}

// deepCloner recursively copies a subtree, duplicating values via
// deep copy and preserving structure and s-field.
func (l LeftistHeap[T]) deepCloner(node *LeftistNode[T]) *LeftistNode[T] {
	if node == nil {
		return node
	}

	newNode := LeftistNode[T]{}
	newNode.val = deepcopy.Copy(node.val).(T)
	newNode.left = l.deepCloner(node.left)
	newNode.right = l.deepCloner(node.right)
	newNode.s = node.s
	return &newNode
}

// DeepClone returns a completely independent copy of the heap,
// including all nodes and values.
func (l LeftistHeap[T]) DeepClone() LeftistHeap[T] {
	newHeap := LeftistHeap[T]{cmp: l.cmp, size: l.size}
	newHeap.root = l.deepCloner(l.root)
	return newHeap
}

// Clone returns a shallow copy of the heap, sharing the same nodes
// and underlying structure.
func (l LeftistHeap[T]) Clone() LeftistHeap[T] {
	return LeftistHeap[T]{root: l.root, cmp: l.cmp, size: l.size}
}

// Clear removes all elements from the heap, resetting root and size.
func (l *LeftistHeap[T]) Clear() { l.root = nil; l.size = 0 }

// Length returns the number of elements currently in the heap.
func (l LeftistHeap[T]) Length() int { return l.size }

// IsEmpty returns true if the heap has no elements.
func (l LeftistHeap[T]) IsEmpty() bool { return l.Length() == 0 }

// Peek returns a pointer to the minimum value without removing it,
// or nil if the heap is empty.
func (l *LeftistHeap[T]) Peek() *T {
	if l.IsEmpty() {
		return nil
	}

	return &l.root.val
}

// Pop removes and returns the minimum value from the heap, restoring
// the leftist property via merge.
func (l *LeftistHeap[T]) Pop() *T {
	if l.IsEmpty() {
		return nil
	}

	root := l.root
	l.root = l.merge(l.root.right, l.root.left)
	l.size--
	return &root.val
}

func (l *LeftistHeap[T]) mergeLarger(a, b *LeftistNode[T]) *LeftistNode[T] {
	if a.left == nil {
		a.left = b
	} else {
		a.right = l.merge(a.right, b)
		if a.left.s < a.right.s {
			tempNode := a.left
			a.left = a.right
			a.right = tempNode
		}
		a.s = a.right.s + 1
	}
	return a
}

// merge combines two leftist subheaps into one, ensuring the
// resulting tree satisfies heap order and leftist property.
func (l *LeftistHeap[T]) merge(a, b *LeftistNode[T]) *LeftistNode[T] {
	if a == nil {
		return b
	}

	if b == nil {
		return a
	}

	if l.cmp(a.val, b.val) {
		return l.mergeLarger(a, b)
	} else {
		return l.mergeLarger(b, a)
	}
}

// Insert adds a new element into the heap by merging a
// singleton node with the existing tree.
func (l *LeftistHeap[T]) Insert(element T) {
	newNode := &LeftistNode[T]{val: element, s: 0}
	l.root = l.merge(newNode, l.root)
	l.size++
}
