package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeapConfigDefaultGenerator(t *testing.T) {
	config := &HeapConfig{
		UsePool:     false,
		IDGenerator: nil,
	}

	generator := config.GetGenerator()
	assert.IsType(t, &UUIDGenerator{}, generator)
}

func TestHeapConfigCustomGenerator(t *testing.T) {
	customGenerator := &IntegerIDGenerator{NextID: 0}
	config := &HeapConfig{
		UsePool:     true,
		IDGenerator: customGenerator,
	}

	generator := config.GetGenerator()
	assert.Equal(t, customGenerator, generator)
	assert.IsType(t, &IntegerIDGenerator{}, generator)
}

func TestHeapConfigUsePool(t *testing.T) {
	config := &HeapConfig{
		UsePool:     true,
		IDGenerator: nil,
	}

	assert.True(t, config.UsePool)

	config.UsePool = false
	assert.False(t, config.UsePool)
}
