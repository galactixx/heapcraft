package heapcraft

import "sync"

type pool[T any] interface {
	Get() T
	Put(node T)
}

// syncPool is a pool that uses a sync.Pool to store the nodes.
type syncPool[T any] struct{ pool sync.Pool }

// Get returns a node from the pool.
func (p *syncPool[T]) Get() T { return p.pool.Get().(T) }

// Put returns a node to the pool
func (p *syncPool[T]) Put(node T) { p.pool.Put(node) }

// defaultPool is a pool that uses a constructor function to create a new node.
// this is the default pool used by the heapcraft package, where the nodes are
// created on the fly.
type defaultPool[T any] struct{ constructor func() T }

// Get just generates a new node based on the constructor function.
func (p *defaultPool[T]) Get() T { return p.constructor() }

// Put is a no-op for the default pool.
func (p *defaultPool[T]) Put(node T) {}

// newDefaultPool creates a new default pool with the given constructor function.
func newDefaultPool[T any](constructor func() T) pool[T] {
	return &defaultPool[T]{constructor: constructor}
}

// newSyncPool creates a new sync pool with the given constructor function.
func newSyncPool[T any](constructor func() T) pool[T] {
	return &syncPool[T]{
		pool: sync.Pool{
			New: func() any { return constructor() },
		},
	}
}

// newPool creates a new pool based on the usePool flag.
func newPool[T any](usePool bool, constructor func() T) pool[T] {
	if usePool {
		return newSyncPool(constructor)
	}
	return newDefaultPool(constructor)
}
