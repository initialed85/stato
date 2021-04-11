package fsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMachine(t *testing.T) {
	stateAEnterCallCount := 0
	stateAExitCallCount := 0
	stateBEnterCallCount := 0
	stateBExitCallCount := 0
	stateCEnterCallCount := 0
	stateCExitCallCount := 0
	transitionABEnterCallCount := 0
	transitionABExitCallCount := 0
	transitionBCEnterCallCount := 0
	transitionBCExitCallCount := 0
	transitionCAEnterCallCount := 0
	transitionCAExitCallCount := 0

	stateA := NewState(
		"state_a",
		[]func() error{
			func() error {
				stateAEnterCallCount++
				return nil
			},
		},
		[]func() error{
			func() error {
				stateAExitCallCount++
				return nil
			},
		},
	)
	stateB := NewState(
		"state_b",
		[]func() error{
			func() error {
				stateBEnterCallCount++
				return nil
			},
		},
		[]func() error{
			func() error {
				stateBExitCallCount++
				return nil
			},
		},
	)
	stateC := NewState(
		"state_c",
		[]func() error{
			func() error {
				stateCEnterCallCount++
				return nil
			},
		},
		[]func() error{
			func() error {
				stateCExitCallCount++
				return nil
			},
		},
	)

	transitionAB := NewTransition(
		"transition_a_b",
		[]func() error{
			func() error {
				transitionABEnterCallCount++
				return nil
			},
		},
		[]func() error{
			func() error {
				transitionABExitCallCount++
				return nil
			},
		},
		stateA,
		stateB,
	)

	transitionBC := NewTransition(
		"transition_b_c",
		[]func() error{
			func() error {
				transitionBCEnterCallCount++
				return nil
			},
		},
		[]func() error{
			func() error {
				transitionBCExitCallCount++
				return nil
			},
		},
		stateB,
		stateC,
	)

	transitionCA := NewTransition(
		"transition_c_a",
		[]func() error{
			func() error {
				transitionCAEnterCallCount++
				return nil
			},
		},
		[]func() error{
			func() error {
				transitionCAExitCallCount++
				return nil
			},
		},
		stateC,
		stateA,
	)

	m, err := NewMachine(
		"some_machine",
		[]*State{
			stateA,
			stateB,
			stateC,
		},
		[]*Transition{
			transitionAB,
			transitionBC,
			transitionCA,
		},
		stateA,
	)
	assert.Nil(t, err)

	assert.Equal(t, 0, stateAEnterCallCount)
	assert.Equal(t, 0, stateAExitCallCount)
	assert.Equal(t, 0, stateBEnterCallCount)
	assert.Equal(t, 0, stateBExitCallCount)
	assert.Equal(t, 0, stateCEnterCallCount)
	assert.Equal(t, 0, stateCExitCallCount)
	assert.Equal(t, 0, transitionABEnterCallCount)
	assert.Equal(t, 0, transitionABExitCallCount)
	assert.Equal(t, 0, transitionBCEnterCallCount)
	assert.Equal(t, 0, transitionBCExitCallCount)
	assert.Equal(t, 0, transitionCAEnterCallCount)
	assert.Equal(t, 0, transitionCAExitCallCount)

	// permitted transition
	err = m.Transition(transitionAB)
	assert.Nil(t, err)

	// permitted transition
	err = m.Transition(transitionBC)
	assert.Nil(t, err)

	// permitted transition
	err = m.Transition(transitionCA)
	assert.Nil(t, err)

	// not permitted transition
	err = m.Transition(transitionCA)
	assert.NotNil(t, err)
}
