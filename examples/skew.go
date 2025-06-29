package examples

import (
	"fmt"

	"github.com/galactixx/heapcraft"
)

func SkewHeapExample() {
	// Create a skew heap with min-heap ordering
	heap := heapcraft.NewFullSkewHeap[int](
		nil,
		func(a, b int) bool { return a < b },
		heapcraft.HeapConfig{UsePool: false},
	)

	// Push some elements
	elements := []struct {
		value    int
		priority int
	}{
		{40, 40},
		{20, 20},
		{60, 60},
		{10, 10},
		{30, 30},
		{50, 50},
		{70, 70},
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
	heap2 := heapcraft.NewFullSkewHeap[string](
		nil,
		func(a, b int) bool { return a < b },
		heapcraft.HeapConfig{UsePool: false},
	)

	// Push elements and store their IDs
	id1, _ := heap2.Push("red", 5)
	id2, _ := heap2.Push("green", 3)
	_, _ = heap2.Push("blue", 7)

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
