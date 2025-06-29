package heapcraft

// HeapConfig is a struct that contains the configuration for a heap.
type HeapConfig struct {
	// UsePool is a boolean that indicates whether to use a pool for the heap.
	UsePool bool
	// IDGenerator is a pointer to an IDGenerator that is used to generate
	// unique IDs for the heap. If nil, the default IDGenerator is used.
	IDGenerator IDGenerator
}

// GetGenerator returns the IDGenerator from the HeapConfig.
// If the IDGenerator is nil, the default IDGenerator is returned.
func (h *HeapConfig) GetGenerator() IDGenerator {
	if h.IDGenerator == nil {
		return &UUIDGenerator{}
	}
	return h.IDGenerator
}
