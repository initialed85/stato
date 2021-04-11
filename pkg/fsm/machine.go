package fsm

import (
	"fmt"
	"sync"
)

type Machine struct {
	mu           sync.Mutex
	name         string
	states       []*State
	transitions  []*Transition
	currentState *State
}

func NewMachine(
	name string,
	states []*State,
	transitions []*Transition,
	initialState *State,
) (*Machine, error) {
	if len(states) == 0 {
		return &Machine{}, fmt.Errorf("states cannot be empty")
	}

	if len(transitions) == 0 {
		return &Machine{}, fmt.Errorf("transitions cannot be empty")
	}

	if initialState == nil {
		return &Machine{}, fmt.Errorf("initialState cannot be nil")
	}

	found := false
	for _, state := range states {
		if state == initialState {
			found = true
			break
		}
	}

	if !found {
		return &Machine{}, fmt.Errorf("initialState must be in states")
	}

	for _, transition := range transitions {
		sourceState := transition.GetSourceState()
		destinationState := transition.GetDestinationState()

		foundSourceState := false
		foundDestinationState := false
		for _, state := range states {
			if sourceState == state {
				foundSourceState = true
			}

			if destinationState == state {
				foundDestinationState = true
			}
		}

		if !(foundSourceState || foundDestinationState) {
			message := ""

			if !foundSourceState && !foundDestinationState {
				message = fmt.Sprintf(
					"source state of %#+v and destination state of %#+v",
					sourceState.GetName(),
					destinationState.GetName(),
				)
			} else if !foundSourceState {
				message = fmt.Sprintf(
					"source state of %#+v",
					sourceState.GetName(),
				)
			} else if !foundDestinationState {
				message = fmt.Sprintf(
					"destination state of %#+v",
					destinationState.GetName(),
				)
			}

			return &Machine{}, fmt.Errorf(
				"%v from transition %#+v must be in states",
				message,
				transition.GetName(),
			)
		}
	}

	m := Machine{
		name:         name,
		states:       states,
		transitions:  transitions,
		currentState: initialState,
	}

	// activate the state but don't fire the enter callbacks
	initialState.forceActive()

	return &m, nil
}

func (m *Machine) GetName() string {
	return m.name
}

func (m *Machine) GetCurrentState() *State {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.currentState
}

func (m *Machine) Transition(transition *Transition) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	found := false
	for _, possibleTransition := range m.transitions {
		if transition == possibleTransition {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf(
			"transition %#+v not in transitions",
			transition.GetName(),
		)
	}

	err := transition.transition(m.currentState)
	if err != nil {
		return err
	}

	m.currentState = transition.destinationState

	return nil
}
