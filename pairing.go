package heapcraft

import (
	"github.com/mohae/deepcopy"
)

// NewPairingHeap constructs a pairing heap from an existing slice.
// It initializes the root with the first element and inserts the rest.
func NewPairingHeap[T any](data []T, cmp func(a, b T) bool) PairingHeap[T] {
	if len(data) == 0 {
		return PairingHeap[T]{cmp: cmp, size: 0}
	}

	heap := PairingHeap[T]{cmp: cmp, size: 0}
	for i := range data {
		heap.Insert(data[i])
	}

	return heap
}

// PairingNode represents a node in the pairing heap, with pointers to its
// first child and next sibling
type PairingNode[T any] struct {
	val         T
	firstChild  *PairingNode[T]
	nextSibling *PairingNode[T]
}

// PairingHeap holds the root node pointer, comparison function, and size
// of the heap
type PairingHeap[T any] struct {
	root *PairingNode[T]
	cmp  func(a, b T) bool
	size int
}

// deepCloner recursively creates a deep copy of a subtree, including values
// and child/sibling pointers
func (p *PairingHeap[T]) deepCloner(node *PairingNode[T]) *PairingNode[T] {
	if node == nil {
		return node
	}

	newNode := PairingNode[T]{}
	newNode.val = deepcopy.Copy(node.val).(T)
	newNode.firstChild = p.deepCloner(node.firstChild)
	newNode.nextSibling = p.deepCloner(node.nextSibling)
	return &newNode
}

// DeepClone produces a deep copy of the entire heap structure, including
// all nodes
func (p PairingHeap[T]) DeepClone() PairingHeap[T] {
	newHeap := PairingHeap[T]{cmp: p.cmp, size: p.size}
	newHeap.root = p.deepCloner(p.root)
	return newHeap
}

// Clone creates a shallow copy of the heap, sharing the same nodes (no
// duplication)
func (p PairingHeap[T]) Clone() PairingHeap[T] {
	return PairingHeap[T]{root: p.root, cmp: p.cmp, size: p.size}
}

// Clear removes all elements by resetting root and size
func (p *PairingHeap[T]) Clear() { p.root = nil; p.size = 0 }

// Length returns the number of elements in the heap
func (p PairingHeap[T]) Length() int { return p.size }

// IsEmpty checks whether the heap has no elements
func (p PairingHeap[T]) IsEmpty() bool { return p.Length() == 0 }

// Peek returns the minimum value at the root without removing it;
// returns nil if empty
func (p *PairingHeap[T]) Peek() *T {
	if p.IsEmpty() {
		return nil
	}

	return &p.root.val
}

// meld links two pairing-heap trees and returns the new root.
// The smaller key becomes the new root, and the other tree becomes
// its first child or a sibling.
func (p *PairingHeap[T]) meld(new *PairingNode[T], root *PairingNode[T]) *PairingNode[T] {
	if root == nil {
		return new
	}

	if new == nil {
		return root
	}

	newRoot := root

	if p.cmp(new.val, newRoot.val) {
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
// remaining siblings.
func (p *PairingHeap[T]) merge(node *PairingNode[T]) *PairingNode[T] {
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

// Pop removes and returns the root value (minimum). It then merges the
// root's children to form the new heap.
func (p *PairingHeap[T]) Pop() *T {
	if p.IsEmpty() {
		return nil
	}

	rootNode := p.root
	p.root = p.merge(p.root.firstChild)

	p.size--
	return &rootNode.val
}

// Insert adds a new element by creating a single-node heap and melding
// it with the existing root.
func (p *PairingHeap[T]) Insert(element T) {
	newNode := &PairingNode[T]{val: element}
	p.root = p.meld(newNode, p.root)
	p.size++
}

// MergeWith combines another pairing heap into this one by melding their
// roots and updating the size.
func (p *PairingHeap[T]) MergeWith(heap PairingHeap[T]) {
	p.root = p.meld(heap.root, p.root)
	p.size = p.size + heap.size
}
