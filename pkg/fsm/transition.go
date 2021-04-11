package fsm

import (
	"fmt"
	"sync"
)

type Transition struct {
	mu               sync.Mutex
	name             string
	contextCallbacks *ContextCallbacks
	sourceState      *State
	destinationState *State
}

func NewTransition(
	name string,
	onEnterCallbacks []func() error,
	onExitCallbacks []func() error,
	sourceState *State,
	destinationState *State,
) *Transition {
	t := Transition{
		name: name,
		contextCallbacks: NewContextCallbacks(
			onEnterCallbacks,
			onExitCallbacks,
		),
		sourceState:      sourceState,
		destinationState: destinationState,
	}

	return &t
}

func (t *Transition) transition(currentState *State) error {
	if currentState != t.sourceState {
		return fmt.Errorf(
			"cannot call transition %#+v; want current state %#+v, got %#+v",
			t.GetName(),
			t.sourceState.GetName(),
			currentState.GetName(),
		)
	}

	if !currentState.GetActive() {
		return fmt.Errorf(
			"cannot call transition %#+v; current state %#+v is not active",
			t.GetName(),
			currentState.GetName(),
		)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	var err error

	// fire callbacks for start of transition
	err = t.contextCallbacks.enter()
	if err != nil {
		return err
	}

	// fire callbacks on the outgoing state
	err = t.sourceState.exit()
	if err != nil {
		return err
	}

	// fire callbacks on the incoming state
	err = t.destinationState.enter()
	if err != nil {
		return err
	}

	// fire callbacks for the end of the transition
	err = t.contextCallbacks.exit()
	if err != nil {
		return err
	}

	return nil
}

func (t *Transition) GetName() string {
	return t.name
}

func (t *Transition) GetActive() bool {
	return t.contextCallbacks.GetActive()
}

func (t *Transition) GetSourceState() *State {
	return t.sourceState
}

func (t *Transition) GetDestinationState() *State {
	return t.destinationState
}
