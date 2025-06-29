package heapcraft

import (
	"strconv"

	"github.com/google/uuid"
)

// IDGenerator is an interface that details a structure
// that can generate unique IDs.
type IDGenerator interface{ Next() string }

// IntegerIDGenerator is a generator that uses integers.
type IntegerIDGenerator struct{ NextID int }

// Next returns the next integer ID as a string.
func (g *IntegerIDGenerator) Next() string {
	intID := strconv.Itoa(g.NextID)
	g.NextID++
	return intID
}

// UUIDGenerator is a generator that uses UUIDs.
type UUIDGenerator struct{}

// Next returns a new UUID as a string (UUIDv4).
func (g *UUIDGenerator) Next() string {
	return uuid.New().String()
}
