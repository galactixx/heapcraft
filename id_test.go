package heapcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegerIDGenerator(t *testing.T) {
	generator := &IntegerIDGenerator{NextID: 0}

	// Test first few IDs
	assert.Equal(t, "0", generator.Next())
	assert.Equal(t, "1", generator.Next())
	assert.Equal(t, "2", generator.Next())

	// Test with custom starting ID
	generator2 := &IntegerIDGenerator{NextID: 100}
	assert.Equal(t, "100", generator2.Next())
	assert.Equal(t, "101", generator2.Next())
}

func TestUUIDGenerator(t *testing.T) {
	generator := &UUIDGenerator{}

	// Test that UUIDs are generated and are different
	id1 := generator.Next()
	id2 := generator.Next()
	id3 := generator.Next()

	// UUIDs should be different
	assert.NotEqual(t, id1, id2)
	assert.NotEqual(t, id2, id3)
	assert.NotEqual(t, id1, id3)

	// UUIDs should be valid UUID format (36 characters with hyphens)
	assert.Len(t, id1, 36)
	assert.Len(t, id2, 36)
	assert.Len(t, id3, 36)

	// Check UUID format (8-4-4-4-12 pattern)
	assert.Contains(t, id1, "-")
	assert.Contains(t, id2, "-")
	assert.Contains(t, id3, "-")
}

func TestIDGeneratorInterface(t *testing.T) {
	var generator IDGenerator

	// Test IntegerIDGenerator implements interface
	generator = &IntegerIDGenerator{NextID: 0}
	assert.Equal(t, "0", generator.Next())

	// Test UUIDGenerator implements interface
	generator = &UUIDGenerator{}
	id := generator.Next()
	assert.Len(t, id, 36)
}
