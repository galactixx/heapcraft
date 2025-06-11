package heapcraft

// valueFromNode extracts the value from a SimpleNode and handles any error that occurred.
// If an error is present, it returns the zero value of type V and the error.
// Otherwise, it returns the node's value and nil error.
func valueFromNode[V any, P any](node SimpleNode[V, P], err error) (V, error) {
	if err != nil {
		var zero V
		return zero, err
	}
	return node.Value(), nil
}

// priorityFromNode extracts the priority from a SimpleNode and handles any error that occurred.
// If an error is present, it returns the zero value of type P and the error.
// Otherwise, it returns the node's priority and nil error.
func priorityFromNode[V any, P any](node SimpleNode[V, P], err error) (P, error) {
	if err != nil {
		var zero P
		return zero, err
	}
	return node.Priority(), nil
}
