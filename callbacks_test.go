package heapcraft

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBaseCallbacksRegister tests registering callbacks with baseCallbacks.
func TestBaseCallbacksRegister(t *testing.T) {
	callbacks := make(baseCallbacks, 0)

	// Test registering a callback
	called := false
	var capturedX, capturedY int
	fn := func(x, y int) {
		called = true
		capturedX = x
		capturedY = y
	}

	callback := callbacks.register(fn)
	assert.NotEmpty(t, callback.ID)
	assert.NotNil(t, callback.Function)
	assert.Equal(t, 1, callbacks.count())

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
	assert.NotEmpty(t, callback2.ID)
	assert.NotNil(t, callback2.Function)
	assert.Equal(t, 2, callbacks.count())

	// Test the second callback works
	callbacks.run(15, 20)
	assert.True(t, called2)
}

// TestSyncCallbacksRegister tests registering callbacks with syncCallbacks.
func TestSyncCallbacksRegister(t *testing.T) {
	callbacks := NewSyncCallbacks()

	// Test registering a callback
	called := false
	var capturedX, capturedY int
	fn := func(x, y int) {
		called = true
		capturedX = x
		capturedY = y
	}

	callback := callbacks.register(fn)
	assert.NotEmpty(t, callback.ID)
	assert.NotNil(t, callback.Function)
	assert.Equal(t, 1, callbacks.count())

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
	assert.NotEmpty(t, callback2.ID)
	assert.NotNil(t, callback2.Function)
	assert.Equal(t, 2, callbacks.count())

	// Test the second callback works
	callbacks.run(15, 20)
	assert.True(t, called2)
}

// TestBaseCallbacksRun tests running callbacks with baseCallbacks.
func TestBaseCallbacksRun(t *testing.T) {
	callbacks := make(baseCallbacks, 0)

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

// TestSyncCallbacksRun tests running callbacks with syncCallbacks.
func TestSyncCallbacksRun(t *testing.T) {
	callbacks := NewSyncCallbacks()

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

// TestBaseCallbacksDeregisterNonExistent tests deregistering non-existent callbacks with baseCallbacks.
func TestBaseCallbacksDeregisterNonExistent(t *testing.T) {
	callbacks := make(baseCallbacks, 0)

	// Try to deregister non-existent callback
	err := callbacks.deregister("999")
	assert.Error(t, err)
	assert.Equal(t, ErrCallbackNotFound, err)
}

// TestSyncCallbacksDeregisterNonExistent tests deregistering non-existent callbacks with syncCallbacks.
func TestSyncCallbacksDeregisterNonExistent(t *testing.T) {
	callbacks := NewSyncCallbacks()

	// Try to deregister non-existent callback
	err := callbacks.deregister("999")
	assert.Error(t, err)
	assert.Equal(t, ErrCallbackNotFound, err)
}

// TestBaseCallbacksCount tests the count method with baseCallbacks.
func TestBaseCallbacksCount(t *testing.T) {
	callbacks := make(baseCallbacks, 0)

	assert.Equal(t, 0, callbacks.count())

	// Register a callback
	fn := func(x, y int) {}
	callbacks.register(fn)
	assert.Equal(t, 1, callbacks.count())

	// Register another callback
	callbacks.register(fn)
	assert.Equal(t, 2, callbacks.count())

	// Deregister a callback
	callback := callbacks.register(fn)
	assert.Equal(t, 3, callbacks.count())

	callbacks.deregister(callback.ID)
	assert.Equal(t, 2, callbacks.count())
}

// TestSyncCallbacksCount tests the count method with syncCallbacks.
func TestSyncCallbacksCount(t *testing.T) {
	callbacks := NewSyncCallbacks()

	assert.Equal(t, 0, callbacks.count())

	// Register a callback
	fn := func(x, y int) {}
	callbacks.register(fn)
	assert.Equal(t, 1, callbacks.count())

	// Register another callback
	callbacks.register(fn)
	assert.Equal(t, 2, callbacks.count())

	// Deregister a callback
	callback := callbacks.register(fn)
	assert.Equal(t, 3, callbacks.count())

	callbacks.deregister(callback.ID)
	assert.Equal(t, 2, callbacks.count())
}

// TestBaseCallbacksGetCallbacks tests the getCallbacks method with baseCallbacks.
func TestBaseCallbacksGetCallbacks(t *testing.T) {
	callbacks := make(baseCallbacks, 0)

	// Register some callbacks
	fn := func(x, y int) {}
	callback1 := callbacks.register(fn)
	callbacks.register(fn)

	// Get a copy of the callbacks
	copied := callbacks.getCallbacks()
	assert.Equal(t, 2, copied.count())

	// Verify the copy is independent
	callbacks.deregister(callback1.ID)
	assert.Equal(t, 1, callbacks.count())
	assert.Equal(t, 2, copied.count()) // Copy should be unchanged
}

// TestSyncCallbacksGetCallbacks tests the getCallbacks method with syncCallbacks.
func TestSyncCallbacksGetCallbacks(t *testing.T) {
	callbacks := NewSyncCallbacks()

	// Register some callbacks
	fn := func(x, y int) {}
	callback1 := callbacks.register(fn)
	callbacks.register(fn)

	// Get a copy of the callbacks
	copied := callbacks.getCallbacks()
	assert.Equal(t, 2, copied.count())

	// Verify the copy is independent
	callbacks.deregister(callback1.ID)
	assert.Equal(t, 1, callbacks.count())
	assert.Equal(t, 2, copied.count()) // Copy should be unchanged
}

// TestBaseCallbacksThreadSafety tests thread safety of baseCallbacks (should not be thread-safe).
func TestBaseCallbacksThreadSafety(t *testing.T) {
	callbacks := make(baseCallbacks, 0)

	var res int32 = 0

	// Register a callback
	callbacks.register(func(x, y int) {
		atomic.AddInt32(&res, int32(x*y))
	})
	var wg sync.WaitGroup

	for i := range 50 {
		wg.Add(1)
		go func(num int) {
			callbacks.run(num, num)
			wg.Done()
		}(i)
	}

	wg.Wait()

	var exp int32 = 0
	for i := range 50 {
		exp += int32(i * i)
	}
	assert.Equal(t, exp, res)
}

// TestSyncCallbacksThreadSafety tests thread safety of syncCallbacks.
func TestSyncCallbacksThreadSafety(t *testing.T) {
	callbacks := NewSyncCallbacks()

	var res int32 = 0

	// Register a callback
	callbacks.register(func(x, y int) {
		atomic.AddInt32(&res, int32(x*y))
	})
	var wg sync.WaitGroup

	for i := range 50 {
		wg.Add(1)
		go func(num int) {
			callbacks.run(num, num)
			wg.Done()
		}(i)
	}

	wg.Wait()

	var exp int32 = 0
	for i := range 50 {
		exp += int32(i * i)
	}
	assert.Equal(t, exp, res)
}

// TestSyncCallbacksConcurrentRegistration tests concurrent registration with syncCallbacks.
func TestSyncCallbacksConcurrentRegistration(t *testing.T) {
	callbacks := NewSyncCallbacks()
	var wg sync.WaitGroup
	numGoroutines := 10

	// Start multiple goroutines that register callbacks concurrently
	for range numGoroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn := func(x, y int) {}
			callbacks.register(fn)
		}()
	}

	wg.Wait()

	// Verify all callbacks were registered
	assert.Equal(t, numGoroutines, callbacks.count())
}

// TestBaseCallbacksEmptyRun tests running callbacks on empty baseCallbacks.
func TestBaseCallbacksEmptyRun(t *testing.T) {
	callbacks := make(baseCallbacks, 0)

	// Running callbacks on empty registry should not panic
	assert.NotPanics(t, func() {
		callbacks.run(1, 2)
	})
}

// TestSyncCallbacksEmptyRun tests running callbacks on empty syncCallbacks.
func TestSyncCallbacksEmptyRun(t *testing.T) {
	callbacks := NewSyncCallbacks()

	// Running callbacks on empty registry should not panic
	assert.NotPanics(t, func() {
		callbacks.run(1, 2)
	})
}

// TestBaseCallbacksUniqueIDs tests that baseCallbacks generates unique IDs.
func TestBaseCallbacksUniqueIDs(t *testing.T) {
	callbacks := make(baseCallbacks, 0)

	// Register multiple callbacks and verify unique IDs
	fn := func(x, y int) {}

	callback1 := callbacks.register(fn)
	callbacks.register(fn)
	callback3 := callbacks.register(fn)

	assert.NotEmpty(t, callback1.ID)
	assert.NotEmpty(t, callback3.ID)
	assert.NotEqual(t, callback1.ID, callback3.ID)
}

// TestSyncCallbacksUniqueIDs tests that syncCallbacks generates unique IDs.
func TestSyncCallbacksUniqueIDs(t *testing.T) {
	callbacks := NewSyncCallbacks()

	// Register multiple callbacks and verify unique IDs
	fn := func(x, y int) {}

	callback1 := callbacks.register(fn)
	callbacks.register(fn)
	callback3 := callbacks.register(fn)

	assert.NotEmpty(t, callback1.ID)
	assert.NotEmpty(t, callback3.ID)
	assert.NotEqual(t, callback1.ID, callback3.ID)
}
