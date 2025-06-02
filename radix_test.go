package heapcraft

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRadixHeapPopOrder(t *testing.T) {
	raw := []uint{10, 3, 7, 1, 5, 2}
	rh := NewRadixHeap(raw)
	assert.False(t, rh.IsEmpty())
	assert.Equal(t, len(raw), rh.Length())

	expected := []uint{1, 2, 3, 5, 7, 10}
	actual := []uint{}
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

func TestRadixHeapPushMonotonicity(t *testing.T) {
	rh := NewRadixHeap([]uint{2, 4, 6})

	minPtr, err := rh.Pop()
	assert.NoError(t, err)
	assert.Equal(t, uint(2), *minPtr)

	err = rh.Push(3)
	assert.NoError(t, err)
	assert.Equal(t, 3, int(*rh.Peek()))

	err = rh.Push(1)
	assert.Error(t, err)
}

func TestRadixHeapPeek(t *testing.T) {
	rh := NewRadixHeap([]uint{8, 2, 5})
	peekPtr := rh.Peek()
	assert.NotNil(t, peekPtr)
	assert.Equal(t, uint(2), *peekPtr)

	// removes 2
	_, _ = rh.Pop()
	assert.Equal(t, uint(5), *rh.Peek())

	// clearing then Peek should return nil
	rh.Clear()
	assert.Nil(t, rh.Peek())
}

func TestRadixHeapClearCloneDeepClone(t *testing.T) {
	original := []uint{4, 1, 3}
	rh := NewRadixHeap(original)
	assert.Equal(t, 3, rh.Length())

	clone := rh.Clone()
	assert.Equal(t, rh.Length(), clone.Length())

	// now last = 1, size = 2
	_, _ = rh.Pop()

	// valid since 2 >= last
	err := rh.Push(2)
	assert.NoError(t, err)

	fmt.Println(rh.buckets)
	fmt.Println(clone.buckets)
	cloneVals := []uint{}
	for !clone.IsEmpty() {
		vPtr, _ := clone.Pop()
		cloneVals = append(cloneVals, *vPtr)
	}
	assert.Equal(t, []uint{1, 3, 4}, cloneVals)

	rh2 := NewRadixHeap([]uint{9, 7, 5})
	deep := rh2.DeepClone()
	assert.Equal(t, rh2.Length(), deep.Length())
	_, _ = rh2.Pop()
	err = rh2.Push(6)
	assert.NoError(t, err)

	deepVals := []uint{}
	for !deep.IsEmpty() {
		vPtr, _ := deep.Pop()
		deepVals = append(deepVals, *vPtr)
	}
	assert.Equal(t, []uint{5, 7, 9}, deepVals)
}

func TestRadixHeapMerge(t *testing.T) {
	rh1 := NewRadixHeap([]uint{1, 4, 6})
	rh2 := NewRadixHeap([]uint{2, 3, 5})

	rh1.Merge(rh2)

	result := []uint{}
	for !rh1.IsEmpty() {
		vPtr, err := rh1.Pop()
		assert.NoError(t, err)
		result = append(result, *vPtr)
	}
	assert.Equal(t, []uint{1, 2, 3, 4, 5, 6}, result)
}

func TestRadixHeapRemoveAndErrors(t *testing.T) {
	rh := NewRadixHeap([]uint{})
	assert.True(t, rh.IsEmpty())
	_, err := rh.Pop()
	assert.Error(t, err)

	rh.Clear()
	err = rh.Push(0)
	assert.NoError(t, err)
	assert.Equal(t, uint(0), *rh.Peek())
}

func TestRadixHeapLengthIsEmpty(t *testing.T) {
	rh := NewRadixHeap([]uint{})
	assert.True(t, rh.IsEmpty())
	assert.Equal(t, 0, rh.Length())

	_ = rh.Push(7)
	assert.False(t, rh.IsEmpty())
	assert.Equal(t, 1, rh.Length())
}
