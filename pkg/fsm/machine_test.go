package fsm

import (
	"context"
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
	transitionCycleEnterCallCount := 0
	transitionCycleExitCallCount := 0
	transitionSelfEnterCallCount := 0
	transitionSelfExitCallCount := 0

	stateA := NewState(
		"state_a",
		func(scope Scope) (context.Context, error) {
			stateAEnterCallCount++
			return scope.Context, nil
		},
		func(scope Scope) (context.Context, error) {
			stateAExitCallCount++
			return scope.Context, nil
		},
	)
	stateB := NewState(
		"state_b",
		func(scope Scope) (context.Context, error) {
			stateBEnterCallCount++
			return scope.Context, nil
		},
		func(scope Scope) (context.Context, error) {
			stateBExitCallCount++
			return scope.Context, nil
		},
	)
	stateC := NewState(
		"state_c",
		func(scope Scope) (context.Context, error) {
			stateCEnterCallCount++
			return scope.Context, nil
		},
		func(scope Scope) (context.Context, error) {
			stateCExitCallCount++
			return scope.Context, nil
		},
	)

	transitionAB := NewTransition(
		"transition_a_b",
		stateA,
		stateB,
		func(scope Scope) (context.Context, error) {
			transitionABEnterCallCount++
			return scope.Context, nil
		},
		func(scope Scope) (context.Context, error) {
			transitionABExitCallCount++
			return scope.Context, nil
		},
	)

	transitionBC := NewTransition(
		"transition_b_c",
		stateB,
		stateC,
		func(scope Scope) (context.Context, error) {
			transitionBCEnterCallCount++
			return scope.Context, nil
		},
		func(scope Scope) (context.Context, error) {
			transitionBCExitCallCount++
			return scope.Context, nil
		},
	)

	transitionCA := NewTransition(
		"transition_c_a",
		stateC,
		stateA,
		func(scope Scope) (context.Context, error) {
			transitionCAEnterCallCount++
			return scope.Context, nil
		},
		func(scope Scope) (context.Context, error) {
			transitionCAExitCallCount++
			return scope.Context, nil
		},
	)

	transitionCycleEnter := func(scope Scope) (context.Context, error) {
		transitionCycleEnterCallCount++
		data := scope.Context.Value("data").(int)
		scope.Context = context.WithValue(
			scope.Context,
			"data",
			data+1,
		)
		return scope.Context, nil
	}

	transitionCycleExit := func(scope Scope) (context.Context, error) {
		transitionCycleExitCallCount++
		data := scope.Context.Value("data").(int)
		scope.Context = context.WithValue(
			scope.Context,
			"data",
			data+1,
		)
		return scope.Context, nil
	}

	transitionCycle1 := NewTransition(
		"transition_cycle",
		stateA,
		stateB,
		transitionCycleEnter,
		transitionCycleExit,
	)

	transitionCycle2 := NewTransition(
		"transition_cycle",
		stateB,
		stateC,
		transitionCycleEnter,
		transitionCycleExit,
	)

	transitionSelfEnter := func(scope Scope) (context.Context, error) {
		transitionSelfEnterCallCount++
		data := scope.Context.Value("data").(int)
		scope.Context = context.WithValue(
			scope.Context,
			"data",
			data+1,
		)
		return scope.Context, nil
	}

	transitionSelfExit := func(scope Scope) (context.Context, error) {
		transitionSelfExitCallCount++
		data := scope.Context.Value("data").(int)
		scope.Context = context.WithValue(
			scope.Context,
			"data",
			data+1,
		)
		return scope.Context, nil
	}

	transitionSelf := NewTransition(
		"transition_self",
		stateA,
		stateA,
		transitionSelfEnter,
		transitionSelfExit,
	)

	m, err := NewMachine(
		[]*State{
			stateA,
			stateB,
			stateC,
		},
		[]*Transition{
			transitionAB,
			transitionBC,
			transitionCA,
			transitionCycle1,
			transitionCycle2,
			transitionSelf,
		},
		stateA,
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

	ctx := context.Background()

	for i := 0; i < 4; i++ {
		// permitted transition
		ctx, err = m.Transition(transitionAB.Name(), context.Background())
		require.NoError(t, err)
		require.Equal(t, stateB.Name(), m.State())

		// permitted transition
		ctx, err = m.Transition(transitionBC.Name(), context.Background())
		require.NoError(t, err)
		require.Equal(t, stateC.Name(), m.State())

		// permitted transition
		ctx, err = m.Transition(transitionCA.Name(), context.Background())
		require.NoError(t, err)
		require.Equal(t, stateA.Name(), m.State())

		// not permitted transition
		ctx, err = m.Transition(transitionCA.Name(), context.Background())
		require.Error(t, err)
		require.Equal(t, stateA.Name(), m.State())
	}

	require.Equal(t, 4, stateAEnterCallCount)
	require.Equal(t, 4, stateAExitCallCount)
	require.Equal(t, 4, stateBEnterCallCount)
	require.Equal(t, 4, stateBExitCallCount)
	require.Equal(t, 4, stateCEnterCallCount)
	require.Equal(t, 4, stateCExitCallCount)
	require.Equal(t, 4, transitionABEnterCallCount)
	require.Equal(t, 4, transitionABExitCallCount)
	require.Equal(t, 4, transitionBCEnterCallCount)
	require.Equal(t, 4, transitionBCExitCallCount)
	require.Equal(t, 4, transitionCAEnterCallCount)
	require.Equal(t, 4, transitionCAExitCallCount)

	ctx = context.WithValue(context.Background(), "data", 1)

	ctx, err = m.Transition("transition_cycle", ctx)
	require.NoError(t, err)
	require.Equal(t, stateB.Name(), m.State())
	require.Equal(t, 3, ctx.Value("data"))

	ctx, err = m.Transition("transition_cycle", ctx)
	require.NoError(t, err)
	require.Equal(t, stateC.Name(), m.State())
	require.Equal(t, 5, ctx.Value("data"))

	ctx, err = m.Transition("transition_cycle", ctx)
	require.Error(t, err)
	require.Equal(t, stateC.Name(), m.State())
	require.Nil(t, ctx)

	ctx = context.WithValue(context.Background(), "data", 1)

	ctx, err = m.Transition(transitionCA.Name(), ctx)
	require.NoError(t, err)
	require.Equal(t, stateA.Name(), m.State())
	require.Equal(t, 5, transitionCAEnterCallCount)
	require.Equal(t, 5, transitionCAExitCallCount)

	require.Equal(t, 5, stateAEnterCallCount)
	require.Equal(t, 5, stateAExitCallCount)

	for i := 0; i < 4; i++ {
		ctx, err = m.Transition(transitionSelf.Name(), ctx)
		require.NoError(t, err)
		require.Equal(t, stateA.Name(), m.State())
	}

	require.Equal(t, 5, stateAEnterCallCount)
	require.Equal(t, 5, stateAExitCallCount)

	require.Equal(t, 4, transitionSelfEnterCallCount)
	require.Equal(t, 4, transitionSelfExitCallCount)
}
