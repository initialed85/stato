package fsm

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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
		false,
	)
	require.NoError(t, err)

	require.Equal(t, 0, stateAEnterCallCount)
	require.Equal(t, 0, stateAExitCallCount)
	require.Equal(t, 0, stateBEnterCallCount)
	require.Equal(t, 0, stateBExitCallCount)
	require.Equal(t, 0, stateCEnterCallCount)
	require.Equal(t, 0, stateCExitCallCount)
	require.Equal(t, 0, transitionABEnterCallCount)
	require.Equal(t, 0, transitionABExitCallCount)
	require.Equal(t, 0, transitionBCEnterCallCount)
	require.Equal(t, 0, transitionBCExitCallCount)
	require.Equal(t, 0, transitionCAEnterCallCount)
	require.Equal(t, 0, transitionCAExitCallCount)

	// permitted transition
	err = m.DoTransition(transitionAB)
	require.NoError(t, err)

	// permitted transition
	err = m.DoTransition(transitionBC)
	require.NoError(t, err)

	// permitted transition
	err = m.DoTransition(transitionCA)
	require.NoError(t, err)

	// not permitted transition
	err = m.DoTransition(transitionCA)
	require.Error(t, err)

	// permitted state change
	err = m.ChangeState(stateB)
	require.NoError(t, err)

	// permitted state change
	err = m.ChangeState(stateC)
	require.NoError(t, err)

	// permitted state change
	err = m.ChangeState(stateA)
	require.NoError(t, err)

	// not permitted state change
	err = m.ChangeState(stateA)
	require.Error(t, err)

	require.Equal(t, 2, stateAEnterCallCount)
	require.Equal(t, 2, stateAExitCallCount)
	require.Equal(t, 2, stateBEnterCallCount)
	require.Equal(t, 2, stateBExitCallCount)
	require.Equal(t, 2, stateCEnterCallCount)
	require.Equal(t, 2, stateCExitCallCount)
	require.Equal(t, 2, transitionABEnterCallCount)
	require.Equal(t, 2, transitionABExitCallCount)
	require.Equal(t, 2, transitionBCEnterCallCount)
	require.Equal(t, 2, transitionBCExitCallCount)
	require.Equal(t, 2, transitionCAEnterCallCount)
	require.Equal(t, 2, transitionCAExitCallCount)

	possibleStateA, err := m.GetState(stateA.GetName())
	require.NoError(t, err)
	assert.Equal(t, stateA, possibleStateA)

	possibleTransitionAB, err := m.GetTransition(transitionAB.GetName())
	require.NoError(t, err)
	assert.Equal(t, possibleTransitionAB, possibleTransitionAB)
}

func TestMachineByName(t *testing.T) {
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
		false,
	)
	require.NoError(t, err)

	require.Equal(t, 0, stateAEnterCallCount)
	require.Equal(t, 0, stateAExitCallCount)
	require.Equal(t, 0, stateBEnterCallCount)
	require.Equal(t, 0, stateBExitCallCount)
	require.Equal(t, 0, stateCEnterCallCount)
	require.Equal(t, 0, stateCExitCallCount)
	require.Equal(t, 0, transitionABEnterCallCount)
	require.Equal(t, 0, transitionABExitCallCount)
	require.Equal(t, 0, transitionBCEnterCallCount)
	require.Equal(t, 0, transitionBCExitCallCount)
	require.Equal(t, 0, transitionCAEnterCallCount)
	require.Equal(t, 0, transitionCAExitCallCount)

	// permitted transition
	err = m.DoTransitionByName(transitionAB.GetName())
	require.NoError(t, err)

	// permitted transition
	err = m.DoTransitionByName(transitionBC.GetName())
	require.NoError(t, err)

	// permitted transition
	err = m.DoTransitionByName(transitionCA.GetName())
	require.NoError(t, err)

	// not permitted transition
	err = m.DoTransitionByName(transitionCA.GetName())
	require.Error(t, err)

	// permitted state change
	err = m.ChangeStateByName(stateB.GetName())
	require.NoError(t, err)

	// permitted state change
	err = m.ChangeStateByName(stateC.GetName())
	require.NoError(t, err)

	// permitted state change
	err = m.ChangeStateByName(stateA.GetName())
	require.NoError(t, err)

	// not permitted state change (implySelfTransition: false)
	err = m.ChangeStateByName(stateA.GetName())
	require.Error(t, err)

	// self transition not permitted
	for i := 0; i < 10; i++ {
		err = m.ChangeStateByName(stateA.GetName())
		require.Error(t, err)
	}

	require.Equal(t, 2, stateAEnterCallCount)
	require.Equal(t, 2, stateAExitCallCount)
	require.Equal(t, 2, stateBEnterCallCount)
	require.Equal(t, 2, stateBExitCallCount)
	require.Equal(t, 2, stateCEnterCallCount)
	require.Equal(t, 2, stateCExitCallCount)
	require.Equal(t, 2, transitionABEnterCallCount)
	require.Equal(t, 2, transitionABExitCallCount)
	require.Equal(t, 2, transitionBCEnterCallCount)
	require.Equal(t, 2, transitionBCExitCallCount)
	require.Equal(t, 2, transitionCAEnterCallCount)
	require.Equal(t, 2, transitionCAExitCallCount)

	possibleStateA, err := m.GetState(stateA.GetName())
	require.NoError(t, err)
	assert.Equal(t, stateA, possibleStateA)

	possibleTransitionAB, err := m.GetTransition(transitionAB.GetName())
	require.NoError(t, err)
	assert.Equal(t, possibleTransitionAB, possibleTransitionAB)
}

func TestMachineByNameImplySelfTransition(t *testing.T) {
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
		true,
	)
	require.NoError(t, err)

	require.Equal(t, 0, stateAEnterCallCount)
	require.Equal(t, 0, stateAExitCallCount)
	require.Equal(t, 0, stateBEnterCallCount)
	require.Equal(t, 0, stateBExitCallCount)
	require.Equal(t, 0, stateCEnterCallCount)
	require.Equal(t, 0, stateCExitCallCount)
	require.Equal(t, 0, transitionABEnterCallCount)
	require.Equal(t, 0, transitionABExitCallCount)
	require.Equal(t, 0, transitionBCEnterCallCount)
	require.Equal(t, 0, transitionBCExitCallCount)
	require.Equal(t, 0, transitionCAEnterCallCount)
	require.Equal(t, 0, transitionCAExitCallCount)

	// permitted transition
	err = m.DoTransitionByName(transitionAB.GetName())
	require.NoError(t, err)

	// permitted transition
	err = m.DoTransitionByName(transitionBC.GetName())
	require.NoError(t, err)

	// permitted transition
	err = m.DoTransitionByName(transitionCA.GetName())
	require.NoError(t, err)

	// not permitted transition
	err = m.DoTransitionByName(transitionCA.GetName())
	require.Error(t, err)

	// permitted state change
	err = m.ChangeStateByName(stateB.GetName())
	require.NoError(t, err)

	// permitted state change
	err = m.ChangeStateByName(stateC.GetName())
	require.NoError(t, err)

	// permitted state change
	err = m.ChangeStateByName(stateA.GetName())
	require.NoError(t, err)

	// permitted state change (implySelfTransition: true)
	err = m.ChangeStateByName(stateA.GetName())
	require.NoError(t, err)

	// self transition not permitted
	for i := 0; i < 10; i++ {
		err = m.ChangeStateByName(stateA.GetName())
		require.NoError(t, err)
	}

	require.Equal(t, 13, stateAEnterCallCount)
	require.Equal(t, 13, stateAExitCallCount)
	require.Equal(t, 2, stateBEnterCallCount)
	require.Equal(t, 2, stateBExitCallCount)
	require.Equal(t, 2, stateCEnterCallCount)
	require.Equal(t, 2, stateCExitCallCount)
	require.Equal(t, 2, transitionABEnterCallCount)
	require.Equal(t, 2, transitionABExitCallCount)
	require.Equal(t, 2, transitionBCEnterCallCount)
	require.Equal(t, 2, transitionBCExitCallCount)
	require.Equal(t, 2, transitionCAEnterCallCount)
	require.Equal(t, 2, transitionCAExitCallCount)

	possibleStateA, err := m.GetState(stateA.GetName())
	require.NoError(t, err)
	assert.Equal(t, stateA, possibleStateA)

	possibleTransitionAB, err := m.GetTransition(transitionAB.GetName())
	require.NoError(t, err)
	assert.Equal(t, possibleTransitionAB, possibleTransitionAB)
}
