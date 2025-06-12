package heapcraft

import (
	"fmt"
	"sync"
)

// callbacks maintains a registry of callback functions (ID → function).
type callbacks struct {
	callbacks map[int]callback
	curId     int
	lock      sync.RWMutex
}

// callback stores a unique ID and the function to invoke when swaps occur.
type callback struct {
	ID       int
	Function func(x, y int)
}

// run invokes each registered callback function with the provided indices x and y.
func (c *callbacks) run(x, y int) {
	c.lock.RLock()
	for _, callback := range c.callbacks {
		callback.Function(x, y)
	}
	c.lock.RUnlock()
}

// Deregister removes the callback with the specified ID, returning an error
// if it does not exist.
func (c *callbacks) deregister(id int) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, exists := c.callbacks[id]; !exists {
		return fmt.Errorf("%d is not an ID of a callback", id)
	}
	delete(c.callbacks, id)
	return nil
}

// Register adds a callback function to be called on each swap and returns a
// callback struct containing the function and its unique ID.
func (c *callbacks) register(fn func(x, y int)) callback {
	c.lock.Lock()
	defer c.lock.Unlock()
	newId := c.curId + 1
	callback := callback{ID: newId, Function: fn}
	c.callbacks[newId] = callback
	c.curId = newId
	return callback
}
