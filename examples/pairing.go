package examples

import (
	"fmt"

	"github.com/galactixx/heapcraft"
)

func PairingHeapExample() {
	// Create a pairing heap with min-heap ordering
	heap := heapcraft.NewPairingHeap[string](
		nil,
		func(a, b int) bool { return a < b },
		heapcraft.HeapConfig{UsePool: false},
	)

	// Push some elements with string values and integer priorities
	elements := []struct {
		value    string
		priority int
	}{
		{"task1", 5},
		{"task2", 2},
		{"task3", 8},
		{"task4", 1},
		{"task5", 3},
	}

	for _, elem := range elements {
		heap.Push(elem.value, elem.priority)
	}

	// Peek at the highest priority task
	if value, priority, err := heap.Peek(); err == nil {
		fmt.Printf("Highest priority task: %s (priority: %d)\n", value, priority)
	}

	// Pop tasks in priority order
	for !heap.IsEmpty() {
		value, priority, _ := heap.Pop()
		fmt.Printf("Processing: %s (priority: %d)\n", value, priority)
	}

	// Example with node tracking
	heap2 := heapcraft.NewPairingHeap[int](
		nil,
		func(a, b int) bool { return a < b },
		heapcraft.HeapConfig{UsePool: false},
	)

	// Push elements and store their IDs
	id1, _ := heap2.Push(10, 10)
	id2, _ := heap2.Push(5, 5)
	_, _ = heap2.Push(15, 15)

	// Update priority of a specific node
	if err := heap2.UpdatePriority(id1, 1); err == nil {
		fmt.Printf("Updated priority of node %s to 1\n", id1)
	}

	// Get value and priority of a specific node
	if value, priority, err := heap2.Get(id2); err == nil {
		fmt.Printf("Node %s: value=%d, priority=%d\n", id2, value, priority)
	}

	for !heap2.IsEmpty() {
		value, priority, _ := heap2.Pop()
		fmt.Printf("value=%d, priority=%d\n", value, priority)
	}
}
