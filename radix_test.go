package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSimpleRadixHeapPopOrder(t *testing.T) {
	raw := []Pair[string, uint]{
		{value: "value10", priority: 10},
		{value: "value3", priority: 3},
		{value: "value7", priority: 7},
		{value: "value1", priority: 1},
		{value: "value5", priority: 5},
		{value: "value2", priority: 2},
	}
	rh := NewSimpleRadixHeap(raw)
	assert.False(t, rh.IsEmpty())
	assert.Equal(t, len(raw), rh.Length())

	expected := []RadixPair[string, uint]{
		{ID: 4, value: "value1", priority: 1},
		{ID: 6, value: "value2", priority: 2},
		{ID: 2, value: "value3", priority: 3},
		{ID: 5, value: "value5", priority: 5},
		{ID: 3, value: "value7", priority: 7},
		{ID: 1, value: "value10", priority: 10},
	}
	actual := []RadixPair[string, uint]{}
	for !rh.IsEmpty() {
		vPtr, err := rh.Pop()
		assert.NoError(t, err)
		actual = append(actual, *vPtr)
	}
	assert.Equal(t, expected, actual)
	assert.True(t, rh.IsEmpty())

	_, err := rh.Pop()
	assert.Error(t, err)
}

func TestSimpleRadixHeapPushMonotonicity(t *testing.T) {
	rh := NewSimpleRadixHeap([]Pair[string, uint]{
		{value: "value2", priority: 2},
		{value: "value4", priority: 4},
		{value: "value6", priority: 6},
	})

	minPtr, err := rh.Pop()
	assert.NoError(t, err)
	assert.Equal(t, uint(2), (*minPtr).Priority())

	_, err = rh.Push("value3", 3)
	assert.NoError(t, err)
	assert.Equal(t, uint(3), (*rh.Peek()).Priority())

	_, err = rh.Push("value1", 1)
	assert.Error(t, err)
}

func TestSimpleRadixHeapPeek(t *testing.T) {
	rh := NewSimpleRadixHeap([]Pair[string, uint]{
		{value: "value8", priority: 8},
		{value: "value2", priority: 2},
		{value: "value5", priority: 5},
	})
	peekPtr := rh.Peek()
	assert.NotNil(t, peekPtr)
	assert.Equal(t, uint(2), (*peekPtr).Priority())

	// removes 2
	_, _ = rh.Pop()
	assert.Equal(t, uint(5), (*rh.Peek()).Priority())

	// clearing then Peek should return nil
	rh.Clear()
	assert.Nil(t, rh.Peek())
}

func TestSimpleRadixHeapClearCloneDeepClone(t *testing.T) {
	original := []Pair[string, uint]{
		{value: "value4", priority: 4},
		{value: "value1", priority: 1},
		{value: "value3", priority: 3},
	}
	rh := NewSimpleRadixHeap(original)
	assert.Equal(t, 3, rh.Length())

	clone := rh.Clone()
	assert.Equal(t, rh.Length(), clone.Length())

	// now last = 1, size = 2
	_, _ = rh.Pop()

	// valid since 2 >= last
	_, err := rh.Push("value2", 2)
	assert.NoError(t, err)

	cloneVals := []uint{}
	for !clone.IsEmpty() {
		vPtr, _ := clone.Pop()
		cloneVals = append(cloneVals, (*vPtr).Priority())
	}
	assert.Equal(t, []uint{1, 3, 4}, cloneVals)
}

func TestSimpleRadixHeapMerge(t *testing.T) {
	rh1 := NewSimpleRadixHeap([]Pair[string, uint]{
		{value: "value1", priority: 1},
		{value: "value4", priority: 4},
		{value: "value6", priority: 6},
	})
	rh2 := NewSimpleRadixHeap([]Pair[string, uint]{
		{value: "value2", priority: 2},
		{value: "value3", priority: 3},
		{value: "value5", priority: 5},
	})
	rh1.Merge(rh2)

	result := []uint{}
	for !rh1.IsEmpty() {
		vPtr, err := rh1.Pop()
		assert.NoError(t, err)
		result = append(result, (*vPtr).Priority())
	}
	assert.Equal(t, []uint{1, 2, 3, 4, 5, 6}, result)
}

func TestSimpleRadixHeapRemoveAndErrors(t *testing.T) {
	rh := NewSimpleRadixHeap([]Pair[string, uint]{})
	assert.True(t, rh.IsEmpty())
	_, err := rh.Pop()
	assert.Error(t, err)

	rh.Clear()
	_, err = rh.Push("value0", 0)
	assert.NoError(t, err)
	assert.Equal(t, uint(0), (*rh.Peek()).Priority())
}

func TestSimpleRadixHeapLengthIsEmpty(t *testing.T) {
	rh := NewSimpleRadixHeap([]Pair[string, uint]{})
	assert.True(t, rh.IsEmpty())
	assert.Equal(t, 0, rh.Length())

	_, _ = rh.Push("value7", 7)
	assert.False(t, rh.IsEmpty())
	assert.Equal(t, 1, rh.Length())
}

func TestNewRadixHeapContainsGet(t *testing.T) {
	data := []Pair[string, uint]{
		{value: "a", priority: 5},
		{value: "b", priority: 3},
		{value: "c", priority: 8},
	}
	rh := NewRadixHeap(data)

	assert.True(t, rh.Contains(1))
	assert.False(t, rh.Contains(4))
	assert.True(t, rh.Contains(2))
	assert.False(t, rh.Contains(5))

	elem, err := rh.GetElement(2)
	assert.NoError(t, err)
	assert.Equal(t, "b", elem.Value())
	assert.Equal(t, uint(3), elem.Priority())

	val, err := rh.GetValue(3)
	assert.NoError(t, err)
	assert.Equal(t, "c", *val)

	_, err = rh.GetPriority(10)
	assert.Error(t, err)

	_, err = rh.GetElement(99)
	assert.Error(t, err)
	_, err = rh.GetValue(99)
	assert.Error(t, err)
	_, err = rh.GetPriority(99)
	assert.Error(t, err)
}

func TestRadixHeapPushPopOrder(t *testing.T) {
	rh := NewRadixHeap([]Pair[string, uint]{})
	assert.True(t, rh.IsEmpty())
	assert.Equal(t, 0, rh.Length())

	_, err := rh.Push("x", 4)
	assert.NoError(t, err)
	_, err = rh.Push("y", 1)
	assert.NoError(t, err)
	_, err = rh.Push("z", 7)
	assert.NoError(t, err)
	_, err = rh.Push("w", 1)
	assert.NoError(t, err)

	assert.False(t, rh.IsEmpty())
	assert.Equal(t, 4, rh.Length())

	peeked := rh.Peek()
	assert.NotNil(t, peeked)
	assert.Equal(t, uint(1), peeked.Priority())

	var out []uint
	for !rh.IsEmpty() {
		item, err := rh.Pop()
		assert.NoError(t, err)
		out = append(out, item.Priority())
	}
	assert.Equal(t, 4, len(out))
	assert.Equal(t, uint(1), out[0])
	assert.Equal(t, uint(1), out[1])
	assert.Equal(t, uint(4), out[2])
	assert.Equal(t, uint(7), out[3])

	_, err = rh.Pop()
	assert.Error(t, err)
}

func TestRadixHeapRemoveAndErrors(t *testing.T) {
	data := []Pair[string, uint]{
		{value: "alpha", priority: 2},
		{value: "beta", priority: 5},
		{value: "gamma", priority: 4},
	}
	rh := NewRadixHeap(data)

	removed, err := rh.Remove(3) /// Issue
	assert.NoError(t, err)
	assert.Equal(t, "gamma", removed.Value())
	assert.Equal(t, uint(4), removed.Priority())
	assert.False(t, rh.Contains(3))
	assert.Equal(t, 2, rh.Length())

	_, err = rh.Remove(999)
	assert.Error(t, err)

	first, err := rh.Pop()
	assert.NoError(t, err)
	assert.Equal(t, uint(2), first.Priority())
	second, err := rh.Pop()
	assert.NoError(t, err)
	assert.Equal(t, uint(4), second.Priority())
	assert.True(t, rh.IsEmpty())
}

func TestRadixHeapUpdatePriority(t *testing.T) {
	data := []Pair[string, uint]{
		{value: "foo", priority: 10},
		{value: "bar", priority: 20},
	}
	rh := NewRadixHeap(data)

	first, err := rh.Pop()
	assert.NoError(t, err)
	assert.Equal(t, uint(10), first.Priority())

	newPair, err := rh.UpdatePriority(2, 5)
	assert.Error(t, err)
	assert.Nil(t, newPair)

	second, err := rh.Pop()
	assert.NoError(t, err)
	assert.Equal(t, uint(20), second.Priority())
	assert.True(t, rh.IsEmpty())
}

func TestRadixHeapClearPeekRebalance(t *testing.T) {
	rh := NewRadixHeap([]Pair[string, uint]{
		{value: "v7", priority: 7},
		{value: "v3", priority: 3},
		{value: "v9", priority: 9},
	})

	peeked := rh.Peek()
	assert.NotNil(t, peeked)
	assert.Equal(t, uint(3), peeked.Priority())

	_, _ = rh.Pop()
	_, err := rh.Push("v2", 2)
	assert.Error(t, err)

	err = rh.Rebalance()
	assert.NoError(t, err)

	rh.Clear()
	assert.True(t, rh.IsEmpty())
	assert.Nil(t, rh.Peek())
}
