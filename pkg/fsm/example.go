package fsm

type FSM struct {
	stateA       *State
	stateB       *State
	transitionAB *Transition
	machine      *Machine
}

func NewFSM(
	onEnterStateA func() error,
	onExitStateA func() error,
	onEnterStateB func() error,
	onExitStateB func() error,
	onEnterTransitionAB func() error,
	onExitTransitionAB func() error,
) *FSM {
	stateA := NewState(
		"state_a",
		[]func() error{onEnterStateA},
		[]func() error{onExitStateA},
	)

	stateB := NewState(
		"state_b",
		[]func() error{onEnterStateB},
		[]func() error{onExitStateB},
	)

	f := FSM{
		stateA: stateA,
		stateB: stateB,
		transitionAB: NewTransition(
			"transition_a_b",
			[]func() error{onEnterTransitionAB},
			[]func() error{onExitTransitionAB},
			stateA,
			stateB,
		),
	}

	return &f
}

func (f *FSM) TransitionAB() error {
	return f.machine.Transition(f.transitionAB)
}
