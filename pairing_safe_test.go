package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSyncFullPairingHeap(t *testing.T) {
	tests := []struct {
		name     string
		data     []HeapNode[int, int]
		expected int
	}{
		{
			name:     "empty heap",
			data:     []HeapNode[int, int]{},
			expected: 0,
		},
		{
			name: "single element",
			data: []HeapNode[int, int]{
				{value: 42, priority: 10},
			},
			expected: 1,
		},
		{
			name: "multiple elements",
			data: []HeapNode[int, int]{
				{value: 42, priority: 10},
				{value: 24, priority: 5},
				{value: 100, priority: 15},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			heap := NewSyncFullPairingHeap(tt.data, lt, HeapConfig{UsePool: false})
			assert.NotNil(t, heap)
			assert.Equal(t, tt.expected, heap.Length())
		})
	}
}

func TestNewSyncPairingHeap(t *testing.T) {
	tests := []struct {
		name     string
		data     []HeapNode[int, int]
		expected int
	}{
		{
			name:     "empty heap",
			data:     []HeapNode[int, int]{},
			expected: 0,
		},
		{
			name: "single element",
			data: []HeapNode[int, int]{
				{value: 42, priority: 10},
			},
			expected: 1,
		},
		{
			name: "multiple elements",
			data: []HeapNode[int, int]{
				{value: 42, priority: 10},
				{value: 24, priority: 5},
				{value: 100, priority: 15},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			heap := NewSyncPairingHeap(tt.data, lt, false)
			assert.NotNil(t, heap)
			assert.Equal(t, tt.expected, heap.Length())
		})
	}
}

func TestSyncFullPairingHeap_Clone(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
		{value: 100, priority: 15},
	}

	original := NewSyncFullPairingHeap(data, lt, HeapConfig{UsePool: false})
	cloned := original.Clone()

	assert.Equal(t, original.Length(), cloned.Length())
	assert.Equal(t, original.IsEmpty(), cloned.IsEmpty())

	original.Push(999, 1)
	assert.NotEqual(t, original.Length(), cloned.Length())
}

func TestSyncPairingHeap_Clone(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
		{value: 100, priority: 15},
	}

	original := NewSyncPairingHeap(data, lt, false)
	cloned := original.Clone()

	assert.Equal(t, original.Length(), cloned.Length())
	assert.Equal(t, original.IsEmpty(), cloned.IsEmpty())

	original.Push(999, 1)
	assert.NotEqual(t, original.Length(), cloned.Length())
}

func TestSyncFullPairingHeap_Push(t *testing.T) {
	heap := NewSyncFullPairingHeap([]HeapNode[int, int]{}, lt, HeapConfig{UsePool: false})

	id, _ := heap.Push(42, 10)
	assert.NotEmpty(t, id)
	assert.Equal(t, 1, heap.Length())
	assert.False(t, heap.IsEmpty())

	id2, _ := heap.Push(24, 5)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id, id2)
	assert.Equal(t, 2, heap.Length())
}

func TestSyncPairingHeap_Push(t *testing.T) {
	heap := NewSyncPairingHeap([]HeapNode[int, int]{}, lt, false)

	heap.Push(42, 10)
	assert.Equal(t, 1, heap.Length())
	assert.False(t, heap.IsEmpty())

	heap.Push(24, 5)
	assert.Equal(t, 2, heap.Length())
}

func TestSyncFullPairingHeap_Pop(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
		{value: 100, priority: 15},
	}

	heap := NewSyncFullPairingHeap(data, lt, HeapConfig{UsePool: false})

	value, priority, err := heap.Pop()
	require.NoError(t, err)
	assert.Equal(t, 24, value)
	assert.Equal(t, 5, priority)
	assert.Equal(t, 2, heap.Length())

	value, _, err = heap.Pop()
	require.NoError(t, err)
	assert.Equal(t, 42, value)

	value, _, err = heap.Pop()
	require.NoError(t, err)
	assert.Equal(t, 100, value)

	_, _, err = heap.Pop()
	assert.Error(t, err)
	assert.True(t, heap.IsEmpty())
}

func TestSyncPairingHeap_Pop(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
		{value: 100, priority: 15},
	}

	heap := NewSyncPairingHeap(data, lt, false)

	value, priority, err := heap.Pop()
	require.NoError(t, err)
	assert.Equal(t, 24, value)
	assert.Equal(t, 5, priority)
	assert.Equal(t, 2, heap.Length())

	value, _, err = heap.Pop()
	require.NoError(t, err)
	assert.Equal(t, 42, value)

	value, _, err = heap.Pop()
	require.NoError(t, err)
	assert.Equal(t, 100, value)

	_, _, err = heap.Pop()
	assert.Error(t, err)
	assert.True(t, heap.IsEmpty())
}

func TestSyncFullPairingHeap_Peek(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
		{value: 100, priority: 15},
	}

	heap := NewSyncFullPairingHeap(data, lt, HeapConfig{UsePool: false})

	value, priority, err := heap.Peek()
	require.NoError(t, err)
	assert.Equal(t, 24, value)
	assert.Equal(t, 5, priority)
	assert.Equal(t, 3, heap.Length())

	heap.Clear()
	_, _, err = heap.Peek()
	assert.Error(t, err)
}

func TestSyncPairingHeap_Peek(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
		{value: 100, priority: 15},
	}

	heap := NewSyncPairingHeap(data, lt, false)

	value, priority, err := heap.Peek()
	require.NoError(t, err)
	assert.Equal(t, 24, value)
	assert.Equal(t, 5, priority)
	assert.Equal(t, 3, heap.Length())

	heap.Clear()
	_, _, err = heap.Peek()
	assert.Error(t, err)
}

func TestSyncFullPairingHeap_PopValue(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncFullPairingHeap(data, lt, HeapConfig{UsePool: false})

	value, err := heap.PopValue()
	require.NoError(t, err)
	assert.Equal(t, 24, value)
	assert.Equal(t, 1, heap.Length())

	heap.Clear()
	value, err = heap.PopValue()
	assert.Error(t, err)
	assert.Equal(t, 0, value)
}

func TestSyncPairingHeap_PopValue(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncPairingHeap(data, lt, false)

	value, err := heap.PopValue()
	require.NoError(t, err)
	assert.Equal(t, 24, value)
	assert.Equal(t, 1, heap.Length())

	heap.Clear()
	value, err = heap.PopValue()
	assert.Error(t, err)
	assert.Equal(t, 0, value)
}

func TestSyncFullPairingHeap_PopPriority(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncFullPairingHeap(data, lt, HeapConfig{UsePool: false})

	priority, err := heap.PopPriority()
	require.NoError(t, err)
	assert.Equal(t, 5, priority)
	assert.Equal(t, 1, heap.Length())

	heap.Clear()
	priority, err = heap.PopPriority()
	assert.Error(t, err)
	assert.Equal(t, 0, priority)
}

func TestSyncPairingHeap_PopPriority(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncPairingHeap(data, lt, false)

	priority, err := heap.PopPriority()
	require.NoError(t, err)
	assert.Equal(t, 5, priority)
	assert.Equal(t, 1, heap.Length())

	heap.Clear()
	priority, err = heap.PopPriority()
	assert.Error(t, err)
	assert.Equal(t, 0, priority)
}

func TestSyncFullPairingHeap_PeekValue(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncFullPairingHeap(data, lt, HeapConfig{UsePool: false})

	value, err := heap.PeekValue()
	require.NoError(t, err)
	assert.Equal(t, 24, value)
	assert.Equal(t, 2, heap.Length())

	heap.Clear()
	value, err = heap.PeekValue()
	assert.Error(t, err)
	assert.Equal(t, 0, value)
}

func TestSyncPairingHeap_PeekValue(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncPairingHeap(data, lt, false)

	value, err := heap.PeekValue()
	require.NoError(t, err)
	assert.Equal(t, 24, value)
	assert.Equal(t, 2, heap.Length())

	heap.Clear()
	value, err = heap.PeekValue()
	assert.Error(t, err)
	assert.Equal(t, 0, value)
}

func TestSyncFullPairingHeap_PeekPriority(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncFullPairingHeap(data, lt, HeapConfig{UsePool: false})

	priority, err := heap.PeekPriority()
	require.NoError(t, err)
	assert.Equal(t, 5, priority)
	assert.Equal(t, 2, heap.Length())

	heap.Clear()
	priority, err = heap.PeekPriority()
	assert.Error(t, err)
	assert.Equal(t, 0, priority)
}

func TestSyncPairingHeap_PeekPriority(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncPairingHeap(data, lt, false)

	priority, err := heap.PeekPriority()
	require.NoError(t, err)
	assert.Equal(t, 5, priority)
	assert.Equal(t, 2, heap.Length())

	heap.Clear()
	priority, err = heap.PeekPriority()
	assert.Error(t, err)
	assert.Equal(t, 0, priority)
}

func TestSyncFullPairingHeap_UpdateValue(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncFullPairingHeap(data, lt, HeapConfig{UsePool: false})

	id, _ := heap.Push(100, 15)
	err := heap.UpdateValue(id, 999)
	assert.NoError(t, err)

	value, err := heap.GetValue(id)
	require.NoError(t, err)
	assert.Equal(t, 999, value)

	err = heap.UpdateValue("non-existent-id", 123)
	assert.Error(t, err)
}

func TestSyncFullPairingHeap_UpdatePriority(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncFullPairingHeap(data, lt, HeapConfig{UsePool: false})

	id, _ := heap.Push(100, 15)
	err := heap.UpdatePriority(id, 1)
	assert.NoError(t, err)

	priority, err := heap.GetPriority(id)
	require.NoError(t, err)
	assert.Equal(t, 1, priority)

	value, err := heap.PeekValue()
	require.NoError(t, err)
	assert.Equal(t, 100, value)

	err = heap.UpdatePriority("non-existent-id", 123)
	assert.Error(t, err)
}

func TestSyncFullPairingHeap_Get(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncFullPairingHeap(data, lt, HeapConfig{UsePool: false})

	id, _ := heap.Push(100, 15)
	value, priority, err := heap.Get(id)
	require.NoError(t, err)
	assert.Equal(t, 100, value)
	assert.Equal(t, 15, priority)

	_, _, err = heap.Get("non-existent-id")
	assert.Error(t, err)
}

func TestSyncFullPairingHeap_GetValue(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncFullPairingHeap(data, lt, HeapConfig{UsePool: false})

	id, _ := heap.Push(100, 15)
	value, err := heap.GetValue(id)
	require.NoError(t, err)
	assert.Equal(t, 100, value)

	value, err = heap.GetValue("non-existent-id")
	assert.Error(t, err)
	assert.Equal(t, 0, value)
}

func TestSyncFullPairingHeap_GetPriority(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncFullPairingHeap(data, lt, HeapConfig{UsePool: false})

	id, _ := heap.Push(100, 15)
	priority, err := heap.GetPriority(id)
	require.NoError(t, err)
	assert.Equal(t, 15, priority)

	priority, err = heap.GetPriority("non-existent-id")
	assert.Error(t, err)
	assert.Equal(t, 0, priority)
}

func TestSyncFullPairingHeap_Clear(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncFullPairingHeap(data, lt, HeapConfig{UsePool: false})
	assert.Equal(t, 2, heap.Length())
	assert.False(t, heap.IsEmpty())

	heap.Clear()
	assert.Equal(t, 0, heap.Length())
	assert.True(t, heap.IsEmpty())
}

func TestSyncPairingHeap_Clear(t *testing.T) {
	data := []HeapNode[int, int]{
		{value: 42, priority: 10},
		{value: 24, priority: 5},
	}

	heap := NewSyncPairingHeap(data, lt, false)
	assert.Equal(t, 2, heap.Length())
	assert.False(t, heap.IsEmpty())

	heap.Clear()
	assert.Equal(t, 0, heap.Length())
	assert.True(t, heap.IsEmpty())
}

func TestSyncFullPairingHeap_Length(t *testing.T) {
	heap := NewSyncFullPairingHeap([]HeapNode[int, int]{}, lt, HeapConfig{UsePool: false})
	assert.Equal(t, 0, heap.Length())

	heap.Push(42, 10)
	assert.Equal(t, 1, heap.Length())

	heap.Push(24, 5)
	assert.Equal(t, 2, heap.Length())
}

func TestSyncPairingHeap_Length(t *testing.T) {
	heap := NewSyncPairingHeap([]HeapNode[int, int]{}, lt, false)
	assert.Equal(t, 0, heap.Length())

	heap.Push(42, 10)
	assert.Equal(t, 1, heap.Length())

	heap.Push(24, 5)
	assert.Equal(t, 2, heap.Length())
}

func TestSyncFullPairingHeap_IsEmpty(t *testing.T) {
	heap := NewSyncFullPairingHeap([]HeapNode[int, int]{}, lt, HeapConfig{UsePool: false})
	assert.True(t, heap.IsEmpty())

	heap.Push(42, 10)
	assert.False(t, heap.IsEmpty())

	heap.Clear()
	assert.True(t, heap.IsEmpty())
}

func TestSyncPairingHeap_IsEmpty(t *testing.T) {
	heap := NewSyncPairingHeap([]HeapNode[int, int]{}, lt, false)
	assert.True(t, heap.IsEmpty())

	heap.Push(42, 10)
	assert.False(t, heap.IsEmpty())

	heap.Clear()
	assert.True(t, heap.IsEmpty())
}
