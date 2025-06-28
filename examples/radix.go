package examples

import (
	"fmt"

	"github.com/galactixx/heapcraft"
)

func RadixHeapExample() {
	// Create a radix heap (only works with unsigned integers)
	heap := heapcraft.NewRadixHeap[int, uint](nil, false)

	// Push some elements with unsigned integer priorities
	elements := []struct {
		value    int
		priority uint
	}{
		{25, 25},
		{50, 50},
		{75, 75},
		{100, 100},
		{150, 150},
		{200, 200},
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
}
