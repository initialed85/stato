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

func (m *Machine) checkState(state *State) error {
	for _, possibleState := range m.states {
		if state != possibleState {
			continue
		}

		return nil
	}

	return fmt.Errorf("state %#+v not a known state", state.GetName())
}

func (m *Machine) checkTransition(transition *Transition) error {
	for _, possibleTransition := range m.transitions {
		if transition != possibleTransition {
			continue
		}

		return nil
	}

	return fmt.Errorf("transition %#+v not a known transition", transition.GetName())
}

func (m *Machine) GetName() string {
	return m.name
}

func (m *Machine) GetCurrentState() *State {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.currentState
}

func (m *Machine) DoTransition(transition *Transition) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	err := m.checkTransition(transition)
	if err != nil {
		return err
	}

	err = transition.transition(m.currentState)
	if err != nil {
		return err
	}

	m.currentState = transition.destinationState

	return nil
}

func (m *Machine) ChangeState(destinationState *State) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	err := m.checkState(destinationState)
	if err != nil {
		return err
	}

	possibleTransitions := make([]*Transition, 0)

	for _, possibleTransition := range m.transitions {
		if !(possibleTransition.sourceState == m.currentState && possibleTransition.destinationState == destinationState) {
			continue
		}

		possibleTransitions = append(possibleTransitions, possibleTransition)
	}

	if len(possibleTransitions) == 0 {
		return fmt.Errorf(
			"no transition from %#+v to %#+v (state change not permitted)",
			m.currentState.GetName(),
			destinationState.GetName(),
		)
	}

	if len(possibleTransitions) > 1 {
		return fmt.Errorf(
			"multiple transitions from %#+v to %#+v (state change permitted, but you have to tell me how)",
			m.currentState.GetName(),
			destinationState.GetName(),
		)
	}

	transition := possibleTransitions[0]

	err = transition.transition(m.currentState)
	if err != nil {
		return err
	}

	m.currentState = transition.destinationState

	return nil
}
