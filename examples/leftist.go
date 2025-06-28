package examples

import (
	"fmt"

	"github.com/galactixx/heapcraft"
)

func LeftistHeapExample() {
	// Create a leftist heap with min-heap ordering
	heap := heapcraft.NewLeftistHeap[int](nil, func(a, b int) bool { return a < b }, false)

	// Push some elements
	elements := []struct {
		value    int
		priority int
	}{
		{20, 20},
		{10, 10},
		{30, 30},
		{5, 5},
		{15, 15},
		{25, 25},
	}

	for _, elem := range elements {
		heap.Push(elem.value, elem.priority)
	}

	// Peek at the root
	if value, priority, err := heap.Peek(); err == nil {
		fmt.Printf("Peek: value=%d, priority=%d\n", value, priority)
	}

	// Pop elements in order
	for !heap.IsEmpty() {
		value, priority, _ := heap.Pop()
		fmt.Printf("Popped: value=%d, priority=%d\n", value, priority)
	}

	// Example with node updates
	heap2 := heapcraft.NewLeftistHeap[string](nil, func(a, b int) bool { return a < b }, false)

	// Push elements and store their IDs
	id1 := heap2.Push("apple", 5)
	id2 := heap2.Push("banana", 3)
	_ = heap2.Push("cherry", 7)

	// Update priority of a specific node
	if err := heap2.UpdatePriority(id1, 1); err == nil {
		fmt.Printf("Updated priority of %s to 1\n", id1)
	}

	// Get value and priority of a specific node
	if value, priority, err := heap2.Get(id2); err == nil {
		fmt.Printf("Node %s: value=%s, priority=%d\n", id2, value, priority)
	}

	for !heap2.IsEmpty() {
		value, priority, _ := heap2.Pop()
		fmt.Printf("value=%s, priority=%d\n", value, priority)
	}
}
