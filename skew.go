package heapcraft

import "github.com/mohae/deepcopy"

// NewSkewHeap constructs a SkewHeap from the given slice by
// inserting elements one by one using the provided comparison
// function. If the slice is empty, it returns an empty heap.
func NewSkewHeap[T any](data []T, cmp func(a, b T) bool) SkewHeap[T] {
	if len(data) == 0 {
		return SkewHeap[T]{cmp: cmp, size: 0}
	}

	root := &SkewNode[T]{val: data[0]}
	heap := SkewHeap[T]{root: root, cmp: cmp, size: 0}
	for i := 1; i < len(data); i++ {
		heap.Insert(data[i])
	}
	return heap
}

type SkewNode[T any] struct {
	val   T
	left  *SkewNode[T]
	right *SkewNode[T]
}

type SkewHeap[T any] struct {
	root *SkewNode[T]
	cmp  func(a, b T) bool
	size int
}

// deepCloner recursively creates a deep copy of the subtree
// rooted at the given node by copying each node’s value (via deepcopy.Copy)
// and cloning its left and right children.
func (s SkewHeap[T]) deepCloner(node *SkewNode[T]) *SkewNode[T] {
	if node == nil {
		return node
	}

	newNode := SkewNode[T]{}
	newNode.val = deepcopy.Copy(node.val).(T)
	newNode.left = s.deepCloner(node.left)
	newNode.right = s.deepCloner(node.right)
	return &newNode
}

// DeepClone returns a new SkewHeap whose structure and stored values are
// deep-copied from the receiver, ensuring that modifying the clone does
// not affect the original.
func (s SkewHeap[T]) DeepClone() SkewHeap[T] {
	newHeap := SkewHeap[T]{cmp: s.cmp, size: s.size}
	newHeap.root = s.deepCloner(s.root)
	return newHeap
}

// Clone returns a shallow copy of the SkewHeap, sharing the same nodes
// and comparison function. Any structural modifications to one heap will
// affect the other.
func (s SkewHeap[T]) Clone() SkewHeap[T] {
	return SkewHeap[T]{root: s.root, cmp: s.cmp, size: s.size}
}

// Clear removes all elements from the heap by dropping its root and
// resetting size to zero.
func (s *SkewHeap[T]) Clear() { s.root = nil; s.size = 0 }

// Length returns the number of elements currently stored in the heap.
func (s SkewHeap[T]) Length() int { return s.size }

// IsEmpty reports whether the heap contains no elements.
func (s SkewHeap[T]) IsEmpty() bool { return s.Length() == 0 }

// Peek returns a pointer to the minimum element (root) without removing it.
// If the heap is empty, it returns nil.
func (s *SkewHeap[T]) Peek() *T {
	if s.IsEmpty() {
		return nil
	}

	return &s.root.val
}

// Pop removes and returns a pointer to the minimum element (root) from the heap.
// It merges the left and right subtrees to restore the heap structure.
// Returns nil if the heap is empty.
func (s *SkewHeap[T]) Pop() *T {
	if s.IsEmpty() {
		return nil
	}

	rootNode := s.root
	s.root = s.merge(s.root.left, s.root.right)
	s.size--
	return &rootNode.val
}

// merge combines two skew-heap subtrees rooted at 'new' and 'root' into
// a single heap. It ensures the heap-order property by recursively merging
// the smaller root with the opposite subtree and swapping children to
// maintain the self-adjusting structure.
func (s *SkewHeap[T]) merge(new *SkewNode[T], root *SkewNode[T]) *SkewNode[T] {
	if new == nil {
		return root
	}

	if root == nil {
		return new
	}

	first := new
	second := root

	if s.cmp(first.val, second.val) {
		tempNode := first.right
		first.right = first.left
		first.left = s.merge(second, tempNode)
		return first
	} else {
		return s.merge(second, first)
	}
}

// Insert adds a new element to the heap by creating a single-node tree
// and merging it with the existing root. The size is incremented
// accordingly.
func (s *SkewHeap[T]) Insert(element T) {
	newNode := SkewNode[T]{val: element}
	s.root = s.merge(&newNode, s.root)
	s.size++
}
