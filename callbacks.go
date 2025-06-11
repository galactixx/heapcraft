package heapcraft

import (
	"fmt"
	"sync"
)

// Callbacks maintains a registry of callback functions (ID → function).
type Callbacks struct {
	callbacks map[int]Callback
	curId     int
	lock      sync.RWMutex
}

// Callback stores a unique ID and the function to invoke when swaps occur.
type Callback struct {
	ID       int
	Function func(x, y int)
}

// run invokes each registered callback function with the provided indices x and y.
func (c *Callbacks) run(x, y int) {
	c.lock.RLock()
	for _, callback := range c.callbacks {
		callback.Function(x, y)
	}
	c.lock.RUnlock()
}

// Deregister removes the callback with the specified ID, returning an error
// if it does not exist.
func (c *Callbacks) deregister(id int) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, exists := c.callbacks[id]; !exists {
		return fmt.Errorf("%d is not an ID of a callback", id)
	}
	delete(c.callbacks, id)
	return nil
}

// Register adds a callback function to be called on each swap and returns a
// Callback struct containing the function and its unique ID.
func (c *Callbacks) register(fn func(x, y int)) Callback {
	c.lock.Lock()
	defer c.lock.Unlock()
	newId := c.curId + 1
	callback := Callback{ID: newId, Function: fn}
	c.callbacks[newId] = callback
	c.curId = newId
	return callback
}
