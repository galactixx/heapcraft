package heapcraft

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallbacksRegister(t *testing.T) {
	callbacks := &Callbacks{
		callbacks: make(map[int]Callback),
		curId:     0,
	}

	// Test registering a callback
	called := false
	var capturedX, capturedY int
	fn := func(x, y int) {
		called = true
		capturedX = x
		capturedY = y
	}

	callback := callbacks.register(fn)
	assert.Equal(t, 1, callback.ID)
	assert.NotNil(t, callback.Function)
	assert.Equal(t, 1, len(callbacks.callbacks))
	assert.Equal(t, 1, callbacks.curId)

	// Test the callback works
	callbacks.run(5, 10)
	assert.True(t, called)
	assert.Equal(t, 5, capturedX)
	assert.Equal(t, 10, capturedY)

	// Test registering another callback
	called2 := false
	fn2 := func(x, y int) {
		called2 = true
	}

	callback2 := callbacks.register(fn2)
	assert.Equal(t, 2, callback2.ID)
	assert.NotNil(t, callback2.Function)
	assert.Equal(t, 2, len(callbacks.callbacks))
	assert.Equal(t, 2, callbacks.curId)

	// Test the second callback works
	callbacks.run(15, 20)
	assert.True(t, called2)
}

func TestCallbacksRun(t *testing.T) {
	callbacks := &Callbacks{
		callbacks: make(map[int]Callback),
		curId:     0,
	}

	// Register multiple callbacks
	called1 := false
	called2 := false
	var capturedX1, capturedY1, capturedX2, capturedY2 int

	fn1 := func(x, y int) {
		called1 = true
		capturedX1 = x
		capturedY1 = y
	}

	fn2 := func(x, y int) {
		called2 = true
		capturedX2 = x
		capturedY2 = y
	}

	callbacks.register(fn1)
	callbacks.register(fn2)

	// Run callbacks
	callbacks.run(5, 10)

	// Verify both callbacks were called with correct parameters
	assert.True(t, called1)
	assert.True(t, called2)
	assert.Equal(t, 5, capturedX1)
	assert.Equal(t, 10, capturedY1)
	assert.Equal(t, 5, capturedX2)
	assert.Equal(t, 10, capturedY2)
}

func TestCallbacksDeregister(t *testing.T) {
	callbacks := &Callbacks{
		callbacks: make(map[int]Callback),
		curId:     0,
	}

	// Register callbacks
	called1 := false
	called2 := false

	fn1 := func(x, y int) {
		called1 = true
	}

	fn2 := func(x, y int) {
		called2 = true
	}

	callback1 := callbacks.register(fn1)
	callback2 := callbacks.register(fn2)

	assert.Equal(t, 2, len(callbacks.callbacks))

	// Deregister first callback
	err := callbacks.deregister(callback1.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(callbacks.callbacks))

	// Verify deregistered callback is not in the map
	_, exists := callbacks.callbacks[callback1.ID]
	assert.False(t, exists)

	// Verify remaining callback is still there
	_, exists = callbacks.callbacks[callback2.ID]
	assert.True(t, exists)

	// Run callbacks - only the second one should be called
	callbacks.run(1, 2)
	assert.False(t, called1)
	assert.True(t, called2)
}

func TestCallbacksDeregisterNonExistent(t *testing.T) {
	callbacks := &Callbacks{
		callbacks: make(map[int]Callback),
		curId:     0,
	}

	// Try to deregister non-existent callback
	err := callbacks.deregister(999)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "999 is not an ID of a callback")
}

func TestCallbacksThreadSafety(t *testing.T) {
	callbacks := &Callbacks{
		callbacks: make(map[int]Callback),
		curId:     0,
	}

	var res int = 0

	// Register a callback
	callbacks.register(func(x, y int) { res += x * y })
	var wg sync.WaitGroup

	for i := range 50 {
		wg.Add(1)
		go func(num int) {
			callbacks.run(num, num)
			wg.Done()
		}(i)
	}

	wg.Wait()

	var exp int = 0
	for i := range 50 {
		exp += i * i
	}
	assert.Equal(t, exp, res)
}

func TestCallbacksEmptyRun(t *testing.T) {
	callbacks := &Callbacks{
		callbacks: make(map[int]Callback),
		curId:     0,
	}

	// Running callbacks on empty registry should not panic
	assert.NotPanics(t, func() {
		callbacks.run(1, 2)
	})
}

func TestCallbacksSequentialIDs(t *testing.T) {
	callbacks := &Callbacks{
		callbacks: make(map[int]Callback),
		curId:     0,
	}

	// Register multiple callbacks and verify sequential IDs
	fn := func(x, y int) {}

	callback1 := callbacks.register(fn)
	callback2 := callbacks.register(fn)
	callback3 := callbacks.register(fn)

	assert.Equal(t, 1, callback1.ID)
	assert.Equal(t, 2, callback2.ID)
	assert.Equal(t, 3, callback3.ID)
	assert.Equal(t, 3, callbacks.curId)

	// Deregister middle callback and register new one
	callbacks.deregister(callback2.ID)
	callback4 := callbacks.register(fn)

	assert.Equal(t, 4, callback4.ID)
	assert.Equal(t, 4, callbacks.curId)
}
