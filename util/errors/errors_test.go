package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var errorsTests = map[error]bool{}

func TestIsHashable(t *testing.T) {
	assert.False(t, IsHashable(NewMapError()))
	assert.False(t, IsHashable(NewSliceError()))
	assert.False(t, IsHashable(NewNonHashableStruct()))
}

type MapError map[string]string

func NewMapError() error {
	return MapError(make(map[string]string))
}

// Error implements error.
func (m MapError) Error() string {
	return "map error"
}

type SliceError []string

func NewSliceError() error {
	return SliceError([]string{})
}

// Error implements error.
func (s SliceError) Error() string {
	return "slice error"
}

type StructNonHashable struct {
	m map[string]string
}

// Error implements error.
func (s StructNonHashable) Error() string {
	return "nonhashable struct"
}

func NewNonHashableStruct() error {
	return StructNonHashable{
		m: make(map[string]string),
	}
}
