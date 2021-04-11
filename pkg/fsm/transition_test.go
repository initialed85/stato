package fsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransition(t *testing.T) {
	sourceEnterCallCount := 0
	transitionEnterCallCount := 0
	sourceExitCallCount := 0
	destinationEnterCallCount := 0
	transitionExitCallCount := 0
	destinationExitCallCount := 0

	var err error

	sourceState := NewState(
		"source_state",
		[]func() error{
			func() error {
				sourceEnterCallCount++
				return nil
			},
		},
		[]func() error{
			func() error {
				sourceExitCallCount++
				return nil
			},
		},
	)

	destinationState := NewState(
		"destination_state",
		[]func() error{
			func() error {
				destinationEnterCallCount++
				return nil
			},
		},
		[]func() error{
			func() error {
				destinationExitCallCount++
				return nil
			},
		},
	)

	transition := NewTransition(
		"some_transition",
		[]func() error{
			func() error {
				transitionEnterCallCount++
				return nil
			},
		},
		[]func() error{
			func() error {
				transitionExitCallCount++
				return nil
			},
		},
		sourceState,
		destinationState,
	)

	sourceState.forceActive()
	assert.True(t, sourceState.GetActive())

	err = transition.transition(sourceState)
	assert.Nil(t, err)
	assert.Equal(t, 0, sourceEnterCallCount)
	assert.Equal(t, 1, transitionEnterCallCount)
	assert.Equal(t, 1, sourceExitCallCount)
	assert.Equal(t, 1, destinationEnterCallCount)
	assert.Equal(t, 1, transitionExitCallCount)
	assert.Equal(t, 0, destinationExitCallCount)
	assert.False(t, transition.GetSourceState().GetActive())
	assert.True(t, transition.GetDestinationState().GetActive())
}
