package heapcraft

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSimplePairingHeapPopOrder(t *testing.T) {
	data := []*HeapNode[int, int]{
		CreateHeapNodePtr(9, 9),
		CreateHeapNodePtr(4, 4),
		CreateHeapNodePtr(6, 6),
		CreateHeapNodePtr(1, 1),
		CreateHeapNodePtr(7, 7),
		CreateHeapNodePtr(3, 3),
	}

	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap(data, cmp)

	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(data), h.Length())

	var values []int
	for !h.IsEmpty() {
		pair, err := h.Pop()
		if err == nil {
			values = append(values, pair.Value())
		}
	}

	expected := []int{1, 3, 4, 6, 7, 9}
	assert.Equal(t, expected, values)
	assert.True(t, h.IsEmpty())
	_, err := h.Pop()
	assert.NotNil(t, err)
}

func TestInsertPopPeekLenIsEmptySimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap([]*HeapNode[int, int]{}, cmp)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())
	_, err := h.Peek()
	assert.NotNil(t, err)

	input := []*HeapNode[int, int]{
		CreateHeapNodePtr(5, 5),
		CreateHeapNodePtr(2, 2),
		CreateHeapNodePtr(8, 8),
		CreateHeapNodePtr(3, 3),
		CreateHeapNodePtr(6, 6),
	}
	expectedOrder := []int{2, 3, 5, 6, 8}

	for _, pair := range input {
		h.Push(pair.Value(), pair.Priority())
	}

	assert.False(t, h.IsEmpty())
	assert.Equal(t, len(input), h.Length())
	peekValue, _ := h.PeekValue()
	assert.Equal(t, 2, peekValue)

	for i, expected := range expectedOrder {
		popped, err := h.Pop()
		assert.Nil(t, err)
		assert.Equal(t, expected, popped.Value())
		assert.Equal(t, len(input)-(i+1), h.Length())
	}

	assert.True(t, h.IsEmpty())
	_, err = h.Peek()
	assert.NotNil(t, err)
}

func TestClearCloneSimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap([]*HeapNode[int, int]{
		CreateHeapNodePtr(1, 1),
		CreateHeapNodePtr(2, 2),
		CreateHeapNodePtr(3, 3),
	}, cmp)

	// Test basic cloning
	clone := h.Clone()
	assert.Equal(t, h.Length(), clone.Length())
	hPeekValue, _ := h.PeekValue()
	clonePeekValue, _ := clone.PeekValue()
	assert.Equal(t, hPeekValue, clonePeekValue)
	hPeekPriority, _ := h.PeekPriority()
	clonePeekPriority, _ := clone.PeekPriority()
	assert.Equal(t, hPeekPriority, clonePeekPriority)

	// Test independence of clone
	h.Push(0, 0)
	hPeekValueAfterInsert, _ := h.PeekValue()
	assert.Equal(t, 0, hPeekValueAfterInsert)
	clonePeekValueAfterInsert, _ := clone.PeekValue()
	assert.Equal(t, 1, clonePeekValueAfterInsert)

	// Test that clone maintains its own state
	clone.Push(5, 5)
	assert.Equal(t, 4, clone.Length())
	assert.Equal(t, 4, h.Length())

	// Test that clearing original doesn't affect clone
	h.Clear()
	assert.True(t, h.IsEmpty())
	assert.False(t, clone.IsEmpty())
	assert.Equal(t, 4, clone.Length())
}

func TestSimplePairingHeapDeepClone(t *testing.T) {
	// Create a heap with a complex structure
	h := NewSimplePairingHeap([]*HeapNode[int, int]{}, lt)
	h.Push(5, 5)
	h.Push(3, 3)
	h.Push(7, 7)
	h.Push(1, 1)
	h.Push(9, 9)

	// Create a clone
	clone := h.Clone()

	// Test that all elements are in the same order
	originalElements := make([]int, 0)
	cloneElements := make([]int, 0)

	for !h.IsEmpty() {
		val, _ := h.PopValue()
		originalElements = append(originalElements, val)
	}

	for !clone.IsEmpty() {
		val, _ := clone.PopValue()
		cloneElements = append(cloneElements, val)
	}

	assert.Equal(t, originalElements, cloneElements)

	// Test that modifying clone doesn't affect original
	h = NewSimplePairingHeap([]*HeapNode[int, int]{}, lt)
	h.Push(5, 5)
	h.Push(3, 3)
	clone = h.Clone()

	clone.Push(1, 1)
	assert.Equal(t, 2, h.Length())
	assert.Equal(t, 3, clone.Length())

	// Test that clone maintains heap property
	val, _ := clone.PopValue()
	assert.Equal(t, 1, val)
}

func TestPairingHeapDeepClone(t *testing.T) {
	// Create a heap with a complex structure
	h := NewPairingHeap([]*HeapNode[int, int]{}, lt)
	id1 := h.Push(5, 5)
	id2 := h.Push(3, 3)
	id3 := h.Push(7, 7)
	id4 := h.Push(1, 1)
	h.Push(9, 9)

	// Create a clone
	clone := h.Clone()

	// Test that all elements are preserved with their IDs
	for _, id := range []uint{id1, id2, id3, id4} {
		val1, err1 := h.GetValue(id)
		val2, err2 := clone.GetValue(id)
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, val1, val2)
	}

	// Test that modifying clone doesn't affect original
	h.Clear()
	h.Push(5, 5)
	h.Push(3, 3)
	clone = h.Clone()

	newID := clone.Push(1, 1)
	assert.Equal(t, 2, h.Length())
	assert.Equal(t, 3, clone.Length())

	// Test that clone maintains heap property and node tracking
	val, _ := clone.PopValue()
	assert.Equal(t, 1, val)

	// Test that new nodes in clone have unique IDs
	_, err := h.Get(newID)
	assert.Error(t, err)
	_, err = clone.Get(newID)
	assert.Error(t, err)

	// Test that clone maintains independent node tracking
	h.Push(10, 10)
	clone.Push(20, 20)

	hVal, _ := h.PeekValue()
	cloneVal, _ := clone.PeekValue()
	assert.Equal(t, hVal, cloneVal)
}

func TestPairingHeapCloneWithUpdates(t *testing.T) {
	// Create a heap with a complex structure
	h := NewPairingHeap([]*HeapNode[int, int]{}, lt)
	id1 := h.Push(5, 5)
	id2 := h.Push(3, 3)
	id3 := h.Push(7, 7)
	id4 := h.Push(1, 1)

	// Create a clone
	clone := h.Clone()

	// Update values in original
	err := h.UpdateValue(id1, 50)
	assert.NoError(t, err)
	err = h.UpdatePriority(id2, 30)
	assert.NoError(t, err)

	// Verify clone remains unchanged
	val1, _ := clone.GetValue(id1)
	val2, _ := clone.GetValue(id2)
	assert.Equal(t, 5, val1)
	assert.Equal(t, 3, val2)

	// Update values in clone
	err = clone.UpdateValue(id3, 70)
	assert.NoError(t, err)
	err = clone.UpdatePriority(id4, 10)
	assert.NoError(t, err)

	// Verify original remains unchanged
	val3, _ := h.GetValue(id3)
	val4, _ := h.GetValue(id4)
	assert.Equal(t, 7, val3)
	assert.Equal(t, 1, val4)

	// Test that both heaps maintain correct order after updates
	hVal, _ := h.PeekValue()
	cloneVal, _ := clone.PeekValue()
	assert.Equal(t, 1, hVal)
	assert.Equal(t, 3, cloneVal)
}

func TestPeekPopEmptySimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap([]*HeapNode[int, int]{}, cmp)
	_, err := h.Peek()
	assert.NotNil(t, err)
	_, err = h.Pop()
	assert.NotNil(t, err)
	_, err = h.PopValue()
	assert.NotNil(t, err)
	_, err = h.PopPriority()
	assert.NotNil(t, err)
}

func TestLengthIsEmptySimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap([]*HeapNode[int, int]{}, cmp)
	assert.True(t, h.IsEmpty())
	assert.Equal(t, 0, h.Length())

	h.Push(10, 10)
	assert.False(t, h.IsEmpty())
	assert.Equal(t, 1, h.Length())
}

func TestPeekValueAndPrioritySimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }

	h := NewSimplePairingHeap([]*HeapNode[int, int]{}, cmp)
	peekValueEmpty, _ := h.PeekValue()
	assert.Equal(t, 0, peekValueEmpty)
	peekPriorityEmpty, _ := h.PeekPriority()
	assert.Equal(t, 0, peekPriorityEmpty)

	h.Push(42, 10)
	peekValue42, _ := h.PeekValue()
	assert.Equal(t, 42, peekValue42)
	peekPriority10, _ := h.PeekPriority()
	assert.Equal(t, 10, peekPriority10)

	h.Push(15, 5)
	peekValue15, _ := h.PeekValue()
	assert.Equal(t, 15, peekValue15)
	peekPriority5, _ := h.PeekPriority()
	assert.Equal(t, 5, peekPriority5)

	h.Push(100, 1)
	peekValue100, _ := h.PeekValue()
	assert.Equal(t, 100, peekValue100)
	peekPriority1, _ := h.PeekPriority()
	assert.Equal(t, 1, peekPriority1)

	h.Pop()
	peekValueAfterPop, _ := h.PeekValue()
	assert.Equal(t, 15, peekValueAfterPop)
	peekPriorityAfterPop, _ := h.PeekPriority()
	assert.Equal(t, 5, peekPriorityAfterPop)

	h.Clear()
	peekValueAfterClear, _ := h.PeekValue()
	assert.Equal(t, 0, peekValueAfterClear)
	peekPriorityAfterClear, _ := h.PeekPriority()
	assert.Equal(t, 0, peekPriorityAfterClear)
}

func TestPopValueAndPrioritySimplePairing(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewSimplePairingHeap([]*HeapNode[int, int]{
		CreateHeapNodePtr(42, 10),
		CreateHeapNodePtr(15, 5),
		CreateHeapNodePtr(100, 1),
	}, cmp)

	val, err := h.PopValue()
	assert.Nil(t, err)
	assert.Equal(t, 100, val)
	peekValue15AfterPop, _ := h.PeekValue()
	assert.Equal(t, 15, peekValue15AfterPop)

	pri, err := h.PopPriority()
	assert.Nil(t, err)
	assert.Equal(t, 5, pri)
	peekValue42AfterPop, _ := h.PeekValue()
	assert.Equal(t, 42, peekValue42AfterPop)

	h.Clear()
	_, err = h.PopValue()
	assert.NotNil(t, err)
	_, err = h.PopPriority()
	assert.NotNil(t, err)
}

func TestPairingHeapIDTracking(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewPairingHeap([]*HeapNode[int, int]{}, cmp)
	assert.NotNil(t, h.elements)
	assert.Equal(t, 0, len(h.elements))

	h.Push(1, 10)
	h.Push(2, 20)
	h.Push(3, 30)

	assert.Equal(t, 3, len(h.elements))
	assert.Equal(t, uint(1), h.curID-3)

	for i := uint(1); i < h.curID; i++ {
		node, exists := h.elements[i]
		assert.True(t, exists)
		assert.Equal(t, i, node.ID())
	}

	popped, err := h.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(h.elements))
	assert.Equal(t, 1, popped.Value())

	h.Clear()
	assert.Equal(t, 0, len(h.elements))
	assert.Equal(t, uint(1), h.curID)
}

func TestPairingHeapUpdateValue(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewPairingHeap([]*HeapNode[int, int]{}, cmp)

	h.Push(1, 10)
	h.Push(2, 20)
	h.Push(3, 30)

	err := h.UpdateValue(1, 100)
	assert.Nil(t, err)
	node, exists := h.elements[1]
	assert.True(t, exists)
	assert.Equal(t, 100, node.Value())

	err = h.UpdateValue(999, 100)
	assert.NotNil(t, err)
	assert.Equal(t, "id does not link to existing node", err.Error())

	popped, err := h.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 100, popped.Value())
}

func TestPairingHeapUpdatePriority(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewPairingHeap([]*HeapNode[int, int]{}, cmp)

	h.Push(1, 10)
	h.Push(2, 20)
	h.Push(3, 30)

	err := h.UpdatePriority(2, 5)
	assert.Nil(t, err)

	popped, err := h.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 2, popped.Value())
	assert.Equal(t, 5, popped.Priority())

	err = h.UpdatePriority(1, 15)
	assert.Nil(t, err)

	popped, err = h.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 1, popped.Value())
}

func TestPairingHeapUpdatePriorityEdgeCases(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewPairingHeap([]*HeapNode[int, int]{}, cmp)

	h.Push(1, 10)
	err := h.UpdatePriority(1, 20)
	assert.Nil(t, err)
	popped, err := h.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 1, popped.Value())
	assert.Equal(t, 20, popped.Priority())

	h.Push(1, 10)
	h.Push(2, 20)
	h.Push(3, 30)
	err = h.UpdatePriority(2, 5)
	assert.Nil(t, err)
	popped, err = h.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 1, popped.Value())
	assert.Equal(t, 5, popped.Priority())

	h.Clear()
	h.Push(1, 10)
	h.Push(2, 20)
	h.Push(3, 30)
	err = h.UpdatePriority(3, 5)
	assert.Nil(t, err)
	popped, err = h.Pop()
	assert.Nil(t, err)
	assert.Equal(t, 3, popped.Value())
	assert.Equal(t, 5, popped.Priority())
}

func TestPairingHeapClone(t *testing.T) {
	cmp := func(a, b int) bool { return a < b }
	h := NewPairingHeap([]*HeapNode[int, int]{}, cmp)

	h.Push(1, 10)
	h.Push(2, 20)
	h.Push(3, 30)

	clone := h.Clone()
	assert.Equal(t, h.size, clone.size)
	assert.Equal(t, h.curID, clone.curID)
	assert.Equal(t, len(h.elements), len(clone.elements))

	for id, node := range h.elements {
		cloneNode, exists := clone.elements[id]
		assert.True(t, exists)
		assert.Equal(t, node.Value(), cloneNode.Value())
		assert.Equal(t, node.Priority(), cloneNode.Priority())
	}
}

func TestComplexHeapStructure(t *testing.T) {
	h := NewPairingHeap[int](nil, func(a, b int) bool { return a < b })

	h.Push(1, 1)
	h.Push(2, 2)
	h.Push(3, 3)
	h.Push(4, 4)
	h.Push(5, 5)
	h.Push(6, 6)
	h.Push(7, 7)

	nodeIDs := make(map[int]uint)
	for id, node := range h.elements {
		nodeIDs[node.value] = id
	}

	assert.Equal(t, 7, h.Length())
	peekValueComplex, _ := h.PeekValue()
	assert.Equal(t, 1, peekValueComplex)
	peekPriorityComplex, _ := h.PeekPriority()
	assert.Equal(t, 1, peekPriorityComplex)

	values := make([]int, 0)
	for !h.IsEmpty() {
		val, err := h.PopValue()
		assert.Nil(t, err)
		values = append(values, val)
	}
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, values)
}

func TestLeafNodeUpdate(t *testing.T) {
	h := NewPairingHeap[int](nil, func(a, b int) bool { return a < b })

	h.Push(1, 1)
	h.Push(2, 2)
	h.Push(3, 3)
	h.Push(4, 4)
	h.Push(5, 5)
	h.Push(6, 6)
	h.Push(7, 7)

	nodeIDs := make(map[int]uint)
	for id, node := range h.elements {
		nodeIDs[node.value] = id
	}

	err := h.UpdatePriority(nodeIDs[7], 0)
	assert.Nil(t, err)
	peekValueLeaf, _ := h.PeekValue()
	assert.Equal(t, 7, peekValueLeaf)
	peekPriorityLeaf, _ := h.PeekPriority()
	assert.Equal(t, 0, peekPriorityLeaf)

	values := make([]int, 0)
	for !h.IsEmpty() {
		val, err := h.PopValue()
		assert.Nil(t, err)
		values = append(values, val)
	}
	assert.Equal(t, []int{7, 1, 2, 3, 4, 5, 6}, values)
}

func TestMiddleNodeUpdate(t *testing.T) {
	h := NewPairingHeap[int](nil, func(a, b int) bool { return a < b })

	h.Push(1, 1)
	h.Push(2, 2)
	h.Push(3, 3)
	h.Push(4, 4)
	h.Push(5, 5)
	h.Push(6, 6)
	h.Push(7, 7)

	nodeIDs := make(map[int]uint)
	for id, node := range h.elements {
		nodeIDs[node.value] = id
	}

	err := h.UpdatePriority(nodeIDs[3], 0)
	assert.Nil(t, err)
	peekValueMiddle, _ := h.PeekValue()
	assert.Equal(t, 3, peekValueMiddle)
	peekPriorityMiddle, _ := h.PeekPriority()
	assert.Equal(t, 0, peekPriorityMiddle)

	values := make([]int, 0)
	for !h.IsEmpty() {
		val, err := h.PopValue()
		assert.Nil(t, err)
		values = append(values, val)
	}
	assert.Equal(t, []int{3, 1, 2, 4, 5, 6, 7}, values)
}

func TestMultipleNodeUpdates(t *testing.T) {
	h := NewPairingHeap[int](nil, func(a, b int) bool { return a < b })

	h.Push(1, 1)
	h.Push(2, 2)
	h.Push(3, 3)
	h.Push(4, 4)
	h.Push(5, 5)
	h.Push(6, 6)
	h.Push(7, 7)

	nodeIDs := make(map[int]uint)
	for id, node := range h.elements {
		nodeIDs[node.value] = id
	}

	err := h.UpdatePriority(nodeIDs[4], 0)
	assert.Nil(t, err)
	peekValueMultiple1, _ := h.PeekValue()
	assert.Equal(t, 4, peekValueMultiple1)
	peekPriorityMultiple1, _ := h.PeekPriority()
	assert.Equal(t, 0, peekPriorityMultiple1)

	err = h.UpdatePriority(nodeIDs[2], 1)
	assert.Nil(t, err)
	peekValueMultiple2, _ := h.PeekValue()
	assert.Equal(t, 4, peekValueMultiple2)
	peekPriorityMultiple2, _ := h.PeekPriority()
	assert.Equal(t, 0, peekPriorityMultiple2)

	err = h.UpdatePriority(nodeIDs[6], -1)
	assert.Nil(t, err)
	peekValueMultiple3, _ := h.PeekValue()
	assert.Equal(t, 6, peekValueMultiple3)
	peekPriorityMultiple3, _ := h.PeekPriority()
	assert.Equal(t, -1, peekPriorityMultiple3)

	values := make([]int, 0)
	for !h.IsEmpty() {
		val, err := h.PopValue()
		assert.Nil(t, err)
		values = append(values, val)
	}
	assert.Equal(t, []int{6, 4, 1, 2, 3, 5, 7}, values)
}

func TestReversePriorityUpdates(t *testing.T) {
	h := NewPairingHeap[int](nil, func(a, b int) bool { return a < b })

	h.Push(1, 10)
	h.Push(2, 20)
	h.Push(3, 30)
	h.Push(4, 40)
	h.Push(5, 50)
	h.Push(6, 60)
	h.Push(7, 70)

	nodeIDs := make(map[int]uint)
	for id, node := range h.elements {
		nodeIDs[node.value] = id
	}

	err := h.UpdatePriority(nodeIDs[7], 1)
	assert.Nil(t, err)
	err = h.UpdatePriority(nodeIDs[6], 2)
	assert.Nil(t, err)
	err = h.UpdatePriority(nodeIDs[5], 3)
	assert.Nil(t, err)
	err = h.UpdatePriority(nodeIDs[4], 4)
	assert.Nil(t, err)
	err = h.UpdatePriority(nodeIDs[3], 5)
	assert.Nil(t, err)
	err = h.UpdatePriority(nodeIDs[2], 6)
	assert.Nil(t, err)
	err = h.UpdatePriority(nodeIDs[1], 7)
	assert.Nil(t, err)

	values := make([]int, 0)
	for !h.IsEmpty() {
		val, err := h.PopValue()
		assert.Nil(t, err)
		values = append(values, val)
	}
	assert.Equal(t, []int{7, 6, 5, 4, 3, 2, 1}, values)
}

func TestPairingHeapGetters(t *testing.T) {
	h := NewPairingHeap[int, int](nil, func(a, b int) bool { return a < b })
	h.Push(42, 10)
	h.Push(15, 5)
	h.Push(100, 1)

	nodeIDs := make(map[int]uint)
	for id, node := range h.elements {
		nodeIDs[node.value] = id
	}

	pair, _ := h.Get(nodeIDs[42])
	assert.Equal(t, 42, pair.Value())
	assert.Equal(t, 10, pair.Priority())

	val, _ := h.GetValue(nodeIDs[15])
	assert.Equal(t, 15, val)

	pri, _ := h.GetPriority(nodeIDs[100])
	assert.Equal(t, 1, pri)

	_, err := h.Get(999)
	assert.NotNil(t, err)
	_, err = h.GetValue(999)
	assert.NotNil(t, err)
	_, err = h.GetPriority(999)
	assert.NotNil(t, err)

	h.Pop()
	_, err = h.Get(nodeIDs[100])
	assert.NotNil(t, err)
}

// Pairing Heap Benchmarks
func BenchmarkPairingHeapInsertion(b *testing.B) {
	N := 10_000
	data := make([]*HeapNode[int, int], 0)
	heap := NewPairingHeap(data, func(a, b int) bool { return a < b })
	b.ReportAllocs()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var num int
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			num = r.Intn(N)
			heap.Push(num, num)
		}
	})
}

func BenchmarkPairingHeapDeletion(b *testing.B) {
	data := make([]*HeapNode[int, int], 0)
	heap := NewPairingHeap(data, func(a, b int) bool { return a < b })

	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			heap.Pop()
		}
	})
}

func BenchmarkSimplePairingHeapInsertion(b *testing.B) {
	N := 10_000
	data := make([]*HeapNode[int, int], 0)
	heap := NewSimplePairingHeap(data, func(a, b int) bool { return a < b })
	b.ReportAllocs()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var num int
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			num = r.Intn(N)
			heap.Push(num, num)
		}
	})
}

func BenchmarkSimplePairingHeapDeletion(b *testing.B) {
	data := make([]*HeapNode[int, int], 0)
	heap := NewSimplePairingHeap(data, func(a, b int) bool { return a < b })

	for i := 0; i < b.N; i++ {
		heap.Push(i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			heap.Pop()
		}
	})
}

func TestPairingHeapInsertReturnsID(t *testing.T) {
	h := NewPairingHeap([]*HeapNode[int, int]{}, lt)

	// Test sequential ID assignment
	id1 := h.Push(10, 10)
	id2 := h.Push(20, 20)
	id3 := h.Push(30, 30)

	assert.Equal(t, uint(1), id1)
	assert.Equal(t, uint(2), id2)
	assert.Equal(t, uint(3), id3)

	// Verify elements can be retrieved using IDs
	val1, _ := h.GetValue(id1)
	val2, _ := h.GetValue(id2)
	val3, _ := h.GetValue(id3)
	assert.Equal(t, 10, val1)
	assert.Equal(t, 20, val2)
	assert.Equal(t, 30, val3)

	// Test ID continues after operations
	h.Pop()
	id4 := h.Push(40, 40)
	assert.Equal(t, uint(4), id4)
}

func TestPairingHeapInsertIDAfterClear(t *testing.T) {
	h := NewPairingHeap([]*HeapNode[int, int]{}, lt)

	id1 := h.Push(10, 10)
	h.Clear()
	id2 := h.Push(20, 20)

	assert.Equal(t, uint(1), id1)
	assert.Equal(t, uint(1), id2) // Should reset to 1
}

func TestSimplePairingHeapInsertNoID(t *testing.T) {
	h := NewSimplePairingHeap([]*HeapNode[int, int]{}, lt)

	// SimplePairingHeap Push should not return ID
	h.Push(10, 10)
	h.Push(20, 20)

	assert.Equal(t, 2, h.Length())
	val1, _ := h.PopValue()
	val2, _ := h.PopValue()
	assert.Equal(t, 10, val1)
	assert.Equal(t, 20, val2)
}
