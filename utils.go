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

// valueFromNode extracts the value and handles any error that occurred.
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
// It uses a dynamic seed for the random number generator.
func generateRandomNumbers(b *testing.B, seed int64) []int {
	N := 10_000
	r := rand.New(rand.NewSource(seed))
	randomNumbers := make([]int, 0, b.N)
	for i := 0; i < b.N; i++ {
		randomNumbers = append(randomNumbers, r.Intn(N))
	}
	return randomNumbers
}

// generateRandomNumbersv1 generates a slice of random numbers for benchmarking.
// It uses a fixed seed of 42 for the random number generator.
func generateRandomNumbersv1(b *testing.B) []int {
	return generateRandomNumbers(b, 42)
}

// generateRandomNumbersv2 generates a slice of random numbers for benchmarking.
// It uses a fixed seed of 50 for the random number generator.
func generateRandomNumbersv2(b *testing.B) []int {
	return generateRandomNumbers(b, 50)
}
