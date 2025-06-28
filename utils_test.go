package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZeroValuePair(t *testing.T) {
	v, p := zeroValuePair[int, string]()
	assert.Equal(t, 0, v)
	assert.Equal(t, "", p)

	v2, p2 := zeroValuePair[string, float64]()
	assert.Equal(t, "", v2)
	assert.Equal(t, 0.0, p2)
}

func TestValueFromNode(t *testing.T) {
	value, err := valueFromNode(42, "high", nil)
	require.NoError(t, err)
	assert.Equal(t, 42, value)

	value, err = valueFromNode(0, "", ErrHeapEmpty)
	assert.Error(t, err)
	assert.Equal(t, 0, value)
	assert.Equal(t, ErrHeapEmpty, err)
}

func TestPriorityFromNode(t *testing.T) {
	priority, err := priorityFromNode(42, "high", nil)
	require.NoError(t, err)
	assert.Equal(t, "high", priority)

	priority, err = priorityFromNode(0, "", ErrNodeNotFound)
	assert.Error(t, err)
	assert.Equal(t, "", priority)
	assert.Equal(t, ErrNodeNotFound, err)
}
