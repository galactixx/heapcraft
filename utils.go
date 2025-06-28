package heapcraft

import (
	"math/rand"
	"testing"
)

// zeroValuePair returns the zero value of type V and P.
func zeroValuePair[V any, P any]() (V, P) {
	var zeroV V
	var zeroP P
	return zeroV, zeroP
}

// valueFromNode extracts the value from a SimpleNode and handles any error that occurred.
// If an error is present, it returns the zero value of type V and the error.
// Otherwise, it returns the node's value and nil error.
func valueFromNode[V any, P any](v V, _ P, err error) (V, error) {
	if err != nil {
		var zero V
		return zero, err
	}
	return v, nil
}

// priorityFromNode extracts the priority from a Node.
func priorityFromNode[V any, P any](_ V, p P, err error) (P, error) {
	if err != nil {
		var zero P
		return zero, err
	}
	return p, nil
}

// generateRandomNumbers generates a slice of random numbers for benchmarking.
// It uses the current time as the seed for the random number generator.
// The numbers are generated using the rand package.
// The numbers are between 0 and N-1.
func generateRandomNumbers(b *testing.B) []int {
	N := 10_000
	r := rand.New(rand.NewSource(42))
	randomNumbers := make([]int, 0, b.N)
	for i := 0; i < b.N; i++ {
		randomNumbers = append(randomNumbers, r.Intn(N))
	}
	return randomNumbers
}
