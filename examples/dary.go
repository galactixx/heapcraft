package examples

import (
	"fmt"

	"github.com/galactixx/heapcraft"
)

func DaryHeapExample() {
	// Create a binary heap (d=2) with min-heap ordering
	heap := heapcraft.NewBinaryHeap[int](nil, func(a, b int) bool { return a < b }, false)

	// Push some elements
	elements := []struct {
		value    int
		priority int
	}{
		{10, 10},
		{5, 5},
		{15, 15},
		{3, 3},
		{8, 8},
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

	// Example with a 3-ary heap
	heap3 := heapcraft.NewDaryHeap[int](3, nil, func(a, b int) bool { return a > b }, false)

	// Push elements
	for _, elem := range elements {
		heap3.Push(elem.value, elem.priority)
	}

	// Pop elements
	for !heap3.IsEmpty() {
		value, priority, _ := heap3.Pop()
		fmt.Printf("Popped: value=%d, priority=%d\n", value, priority)
	}
}
