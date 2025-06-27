package heapcraft

import (
	"sync"

	"github.com/google/uuid"
)

// callbacks is an interface that defines the methods for managing
// callbacks.
type callbacks interface {
	run(x, y int)
	register(fn func(x, y int)) callback
	deregister(id string) error
	count() int
	getCallbacks() callbacks
}

// callbacks maintains a registry of callback functions (ID → function).
type baseCallbacks map[string]callback

// callback stores a unique ID and the function to invoke when swaps occur.
type callback struct {
	ID       string
	Function func(x, y int)
}

// run invokes each registered callback function with the provided indices x and y.
func (c baseCallbacks) run(x, y int) {
	for _, callback := range c {
		callback.Function(x, y)
	}
}

// register adds a callback function to be called on each swap and returns a
// callback struct containing the function and its unique ID.
func (c baseCallbacks) register(fn func(x, y int)) callback {
	newId := uuid.New().String()
	callback := callback{ID: newId, Function: fn}
	c[newId] = callback
	return callback
}

// deregister removes the callback with the specified ID, returning an error
// if it does not exist.
func (c baseCallbacks) deregister(id string) error {
	if _, exists := c[id]; !exists {
		return ErrCallbackNotFound
	}
	delete(c, id)
	return nil
}

// count returns the number of registered callbacks.
func (c baseCallbacks) count() int { return len(c) }

// getCallbacks returns a copy of the callbacks map.
func (c baseCallbacks) getCallbacks() callbacks {
	callbacksMap := make(baseCallbacks, len(c))
	for k, v := range c {
		callbacksMap[k] = v
	}
	return callbacksMap
}

// NewSyncCallbacks creates a new thread-safe callbacks instance.
func NewSyncCallbacks() *syncCallbacks {
	return &syncCallbacks{callbacks: make(baseCallbacks, 0)}
}

// syncCallbacks represents a thread-safe wrapper around callbacks.
// It provides the same interface as callbacks but with mutex-protected operations.
type syncCallbacks struct {
	callbacks baseCallbacks
	lock      sync.RWMutex
}

// Run invokes each registered callback function with the provided indices x and y.
// This is the thread-safe version of run.
func (c *syncCallbacks) run(x, y int) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	c.callbacks.run(x, y)
}

// Register adds a callback function to be called on each swap and returns a
// callback struct containing the function and its unique ID.
// This is the thread-safe version of register.
func (c *syncCallbacks) register(fn func(x, y int)) callback {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.callbacks.register(fn)
}

// Deregister removes the callback with the specified ID, returning an error
// if it does not exist. This is the thread-safe version of deregister.
func (c *syncCallbacks) deregister(id string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.callbacks.deregister(id)
}

// Count returns the number of registered callbacks.
// This is the thread-safe version of count.
func (c *syncCallbacks) count() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.callbacks.count()
}

// getCallbacks returns a copy of the callbacks map.
// This is the thread-safe version of getCallbacks.
func (c *syncCallbacks) getCallbacks() callbacks {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.callbacks.getCallbacks()
}
