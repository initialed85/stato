package fsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	s := NewState(
		"some_state",
		[]func() error{},
		[]func() error{},
	)

	assert.False(t, s.GetActive())

	s.forceActive()
	assert.Equal(t, "some_state", s.GetName())
	assert.True(t, s.GetActive())
}
