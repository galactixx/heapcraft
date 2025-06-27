package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSyncRadixHeap(t *testing.T) {
	tests := []struct {
		name     string
		data     []HeapNode[int, uint]
		expected int
	}{
		{
			name:     "empty heap",
			data:     []HeapNode[int, uint]{},
			expected: 0,
		},
		{
			name: "single element",
			data: []HeapNode[int, uint]{
				{value: 42, priority: 10},
			},
			expected: 1,
		},
		{
			name: "multiple elements",
			data: []HeapNode[int, uint]{
				{value: 42, priority: 10},
				{value: 24, priority: 5},
				{value: 100, priority: 15},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			heap := NewSyncRadixHeap(tt.data)
			assert.NotNil(t, heap)
			assert.Equal(t, tt.expected, heap.Length())
		})
	}
}

func TestSyncRadixHeap_Clone(t *testing.T) {
	data := []HeapNode[int, uint]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
		{value: 100, priority: 15},
	}

	original := NewSyncRadixHeap(data)
	cloned := original.Clone()

	assert.Equal(t, original.Length(), cloned.Length())
	assert.Equal(t, original.IsEmpty(), cloned.IsEmpty())

	original.Push(999, 17)
	assert.NotEqual(t, original.Length(), cloned.Length())
}

func TestSyncRadixHeap_Push(t *testing.T) {
	heap := NewSyncRadixHeap([]HeapNode[int, uint]{})

	t.Run("push to empty heap", func(t *testing.T) {
		err := heap.Push(42, 10)
		assert.NoError(t, err)
		assert.Equal(t, 1, heap.Length())
		assert.False(t, heap.IsEmpty())
	})

	t.Run("push with higher priority", func(t *testing.T) {
		err := heap.Push(24, 15)
		assert.NoError(t, err)
		assert.Equal(t, 2, heap.Length())
	})

	t.Run("push with lower priority should fail", func(t *testing.T) {
		err := heap.Push(100, 5)
		assert.Error(t, err)
		assert.Equal(t, ErrPriorityLessThanLast, err)
		assert.Equal(t, 2, heap.Length())
	})
}

func TestSyncRadixHeap_Pop(t *testing.T) {
	data := []HeapNode[int, uint]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
		{value: 100, priority: 15},
	}

	heap := NewSyncRadixHeap(data)

	t.Run("pop from non-empty heap", func(t *testing.T) {
		node, err := heap.Pop()
		require.NoError(t, err)
		assert.Equal(t, 24, node.Value())
		assert.Equal(t, uint(5), node.Priority())
		assert.Equal(t, 2, heap.Length())
	})

	t.Run("pop until empty", func(t *testing.T) {
		node, err := heap.Pop()
		require.NoError(t, err)
		assert.Equal(t, 42, node.Value())
		assert.Equal(t, uint(10), node.Priority())

		node, err = heap.Pop()
		require.NoError(t, err)
		assert.Equal(t, 100, node.Value())
		assert.Equal(t, uint(15), node.Priority())

		_, err = heap.Pop()
		assert.Error(t, err)
		assert.Equal(t, ErrHeapEmpty, err)
		assert.True(t, heap.IsEmpty())
	})
}

func TestSyncRadixHeap_Peek(t *testing.T) {
	data := []HeapNode[int, uint]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
		{value: 100, priority: 15},
	}

	heap := NewSyncRadixHeap(data)

	t.Run("peek from non-empty heap", func(t *testing.T) {
		node, err := heap.Peek()
		require.NoError(t, err)
		assert.Equal(t, 24, node.Value())
		assert.Equal(t, uint(5), node.Priority())
		assert.Equal(t, 3, heap.Length())
	})

	t.Run("peek from empty heap", func(t *testing.T) {
		heap.Clear()
		_, err := heap.Peek()
		assert.Error(t, err)
		assert.Equal(t, ErrHeapEmpty, err)
	})
}

func TestSyncRadixHeap_PopValue(t *testing.T) {
	data := []HeapNode[int, uint]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncRadixHeap(data)

	t.Run("pop value from non-empty heap", func(t *testing.T) {
		value, err := heap.PopValue()
		require.NoError(t, err)
		assert.Equal(t, 24, value)
		assert.Equal(t, 1, heap.Length())
	})

	t.Run("pop value from empty heap", func(t *testing.T) {
		heap.Clear()
		value, err := heap.PopValue()
		assert.Error(t, err)
		assert.Equal(t, ErrHeapEmpty, err)
		assert.Equal(t, 0, value)
	})
}

func TestSyncRadixHeap_PopPriority(t *testing.T) {
	data := []HeapNode[int, uint]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncRadixHeap(data)

	t.Run("pop priority from non-empty heap", func(t *testing.T) {
		priority, err := heap.PopPriority()
		require.NoError(t, err)
		assert.Equal(t, uint(5), priority)
		assert.Equal(t, 1, heap.Length())
	})

	t.Run("pop priority from empty heap", func(t *testing.T) {
		heap.Clear()
		priority, err := heap.PopPriority()
		assert.Error(t, err)
		assert.Equal(t, ErrHeapEmpty, err)
		assert.Equal(t, uint(0), priority)
	})
}

func TestSyncRadixHeap_PeekValue(t *testing.T) {
	data := []HeapNode[int, uint]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncRadixHeap(data)

	t.Run("peek value from non-empty heap", func(t *testing.T) {
		value, err := heap.PeekValue()
		require.NoError(t, err)
		assert.Equal(t, 24, value)
		assert.Equal(t, 2, heap.Length())
	})

	t.Run("peek value from empty heap", func(t *testing.T) {
		heap.Clear()
		value, err := heap.PeekValue()
		assert.Error(t, err)
		assert.Equal(t, ErrHeapEmpty, err)
		assert.Equal(t, 0, value)
	})
}

func TestSyncRadixHeap_PeekPriority(t *testing.T) {
	data := []HeapNode[int, uint]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncRadixHeap(data)

	t.Run("peek priority from non-empty heap", func(t *testing.T) {
		priority, err := heap.PeekPriority()
		require.NoError(t, err)
		assert.Equal(t, uint(5), priority)
		assert.Equal(t, 2, heap.Length())
	})

	t.Run("peek priority from empty heap", func(t *testing.T) {
		heap.Clear()
		priority, err := heap.PeekPriority()
		assert.Error(t, err)
		assert.Equal(t, ErrHeapEmpty, err)
		assert.Equal(t, uint(0), priority)
	})
}

func TestSyncRadixHeap_Clear(t *testing.T) {
	data := []HeapNode[int, uint]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncRadixHeap(data)
	assert.Equal(t, 2, heap.Length())
	assert.False(t, heap.IsEmpty())

	heap.Clear()
	assert.Equal(t, 0, heap.Length())
	assert.True(t, heap.IsEmpty())
}

func TestSyncRadixHeap_Rebalance(t *testing.T) {
	t.Run("rebalance when bucket 0 is empty", func(t *testing.T) {
		data := []HeapNode[int, uint]{
			{value: 42, priority: 10},
			{value: 24, priority: 5},
			{value: 100, priority: 15},
		}
		heap := NewSyncRadixHeap(data)

		_, err := heap.Pop()
		require.NoError(t, err)

		err = heap.Rebalance()
		assert.NoError(t, err)
		assert.Equal(t, 2, heap.Length())
	})

	t.Run("rebalance when bucket 0 has elements", func(t *testing.T) {
		data := []HeapNode[int, uint]{
			{value: 42, priority: 10},
		}
		heap := NewSyncRadixHeap(data)

		err := heap.Rebalance()
		assert.Error(t, err)
		assert.Equal(t, ErrNoRebalancingNeeded, err)
	})

	t.Run("rebalance empty heap", func(t *testing.T) {
		heap := NewSyncRadixHeap([]HeapNode[int, uint]{})
		err := heap.Rebalance()
		assert.Error(t, err)
		assert.Equal(t, ErrHeapEmpty, err)
	})
}

func TestSyncRadixHeap_Length(t *testing.T) {
	heap := NewSyncRadixHeap([]HeapNode[int, uint]{})
	assert.Equal(t, 0, heap.Length())

	heap.Push(42, 10)
	assert.Equal(t, 1, heap.Length())

	heap.Push(24, 15)
	assert.Equal(t, 2, heap.Length())
}

func TestSyncRadixHeap_IsEmpty(t *testing.T) {
	heap := NewSyncRadixHeap([]HeapNode[int, uint]{})
	assert.True(t, heap.IsEmpty())

	heap.Push(42, 10)
	assert.False(t, heap.IsEmpty())

	heap.Clear()
	assert.True(t, heap.IsEmpty())
}

func TestSyncRadixHeap_Merge(t *testing.T) {
	data1 := []HeapNode[int, uint]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}
	data2 := []HeapNode[int, uint]{
		{value: 100, priority: 15},
		{value: 50, priority: 8},
	}

	heap1 := NewSyncRadixHeap(data1)
	heap2 := NewSyncRadixHeap(data2)

	originalLength1 := heap1.Length()
	originalLength2 := heap2.Length()

	heap1.Merge(heap2)

	assert.Equal(t, originalLength1+originalLength2, heap1.Length())
	assert.Equal(t, originalLength2, heap2.Length())

	allValues := make([]int, 0)
	for !heap1.IsEmpty() {
		value, err := heap1.PopValue()
		require.NoError(t, err)
		allValues = append(allValues, value)
	}

	expectedValues := []int{24, 50, 42, 100}
	assert.ElementsMatch(t, expectedValues, allValues)
}
