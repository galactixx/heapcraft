package heapcraft

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testNode is a simple implementation of SimpleNode for testing
type testNode[V any, P any] struct {
	val V
	pri P
}

func (m testNode[V, P]) Value() V    { return m.val }
func (m testNode[V, P]) Priority() P { return m.pri }

func TestValueFromNode(t *testing.T) {
	tests := []struct {
		name    string
		node    SimpleNode[string, int]
		err     error
		wantVal string
		wantErr bool
	}{
		{
			name:    "successful value extraction",
			node:    testNode[string, int]{val: "test", pri: 1},
			err:     nil,
			wantVal: "test",
			wantErr: false,
		},
		{
			name:    "error case",
			node:    nil,
			err:     errors.New("test error"),
			wantVal: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := valueFromNode(tt.node, tt.err)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantVal, got)
			}
		})
	}
}

func TestPriorityFromNode(t *testing.T) {
	tests := []struct {
		name    string
		node    SimpleNode[string, int]
		err     error
		wantPri int
		wantErr bool
	}{
		{
			name:    "successful priority extraction",
			node:    testNode[string, int]{val: "test", pri: 42},
			err:     nil,
			wantPri: 42,
			wantErr: false,
		},
		{
			name:    "error case",
			node:    nil,
			err:     errors.New("test error"),
			wantPri: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := priorityFromNode(tt.node, tt.err)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Zero(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPri, got)
			}
		})
	}
}
