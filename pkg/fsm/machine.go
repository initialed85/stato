package fsm

import (
	"fmt"
	"sync"
)

type Machine struct {
	mu                  sync.Mutex
	name                string
	states              []*State
	stateByName         map[string]*State
	transitions         []*Transition
	transitionByName    map[string]*Transition
	currentState        *State
	implySelfTransition bool
}

func NewMachine(
	name string,
	states []*State,
	transitions []*Transition,
	initialState *State,
	implySelfTransition bool,
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
					sourceState.name,
					destinationState.name,
				)
			} else if !foundSourceState {
				message = fmt.Sprintf(
					"source state of %#+v",
					sourceState.name,
				)
			} else if !foundDestinationState {
				message = fmt.Sprintf(
					"destination state of %#+v",
					destinationState.name,
				)
			}

			return &Machine{}, fmt.Errorf(
				"%v from transition %#+v must be in states",
				message,
				transition.name,
			)
		}
	}

	m := Machine{
		name:                name,
		states:              states,
		stateByName:         make(map[string]*State),
		transitionByName:    make(map[string]*Transition),
		transitions:         transitions,
		currentState:        initialState,
		implySelfTransition: implySelfTransition,
	}

	for _, state := range m.states {
		m.stateByName[state.GetName()] = state
	}

	for _, transition := range m.transitions {
		m.transitionByName[transition.GetName()] = transition
	}

	if implySelfTransition {
		for _, state := range m.states {
			transitionName := fmt.Sprintf("transition_%v_%v", state.name, state.name)
			_, ok := m.transitionByName[transitionName]
			if ok {
				return nil, fmt.Errorf(
					"implySelfTransition wants to create transition %#+v but it already exists",
					transitionName,
				)
			}

			transition := NewTransition(
				"some_transition",
				[]func() error{
					func() error {
						return nil
					},
				},
				[]func() error{
					func() error {
						return nil
					},
				},
				state,
				state,
			)

			m.transitions = append(m.transitions, transition)
			m.transitionByName[transitionName] = transition
		}
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

func (m *Machine) doTransition(transition *Transition) error {
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

func (m *Machine) DoTransition(transition *Transition) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.doTransition(transition)
}

func (m *Machine) changeState(destinationState *State) error {
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

func (m *Machine) ChangeState(destinationState *State) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.changeState(destinationState)
}

func (m *Machine) getState(name string) (*State, error) {
	state, ok := m.stateByName[name]
	if !ok {
		return nil, fmt.Errorf("state %#+v not a known state", name)
	}

	return state, nil
}

func (m *Machine) GetState(name string) (*State, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.getState(name)
}

func (m *Machine) getTransition(name string) (*Transition, error) {
	transition, ok := m.transitionByName[name]
	if !ok {
		return nil, fmt.Errorf("transition %#+v not a known transition", name)
	}

	return transition, nil
}

func (m *Machine) GetTransition(name string) (*Transition, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.getTransition(name)
}

func (m *Machine) DoTransitionByName(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	transition, err := m.getTransition(name)
	if err != nil {
		return err
	}

	return m.doTransition(transition)
}

func (m *Machine) ChangeStateByName(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	state, err := m.getState(name)
	if err != nil {
		return err
	}

	return m.changeState(state)
}
