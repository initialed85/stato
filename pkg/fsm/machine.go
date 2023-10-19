package fsm

import (
	"context"
	"fmt"
	"sync"
)

type Machine struct {
	implySelfTransition      bool
	mu                       *sync.Mutex
	stateByName              map[string]*State
	transitionBySourceByName map[string]map[*State]*Transition
	initialState             *State
	currentState             *State
}

func NewMachine(
	states []*State,
	transitions []*Transition,
	initialState *State,
) (*Machine, error) {
	m := Machine{
		mu:                       new(sync.Mutex),
		stateByName:              make(map[string]*State),
		transitionBySourceByName: make(map[string]map[*State]*Transition),
	}

	for _, state := range states {
		_, ok := m.stateByName[state.Name()]
		if ok {
			return nil, fmt.Errorf("state %#+v already exists", state.Name())
		}
		m.stateByName[state.Name()] = state
	}

	for _, transition := range transitions {
		transitionBySource, ok := m.transitionBySourceByName[transition.Name()]
		if !ok {
			transitionBySource = make(map[*State]*Transition)
		}

		_, ok = transitionBySource[transition.source]
		if ok {
			return nil, fmt.Errorf(
				"transition %#+v already exists for source %#+v",
				transition.Name(), transition.GetSource().Name(),
			)
		}

		transitionBySource[transition.GetSource()] = transition
		m.transitionBySourceByName[transition.Name()] = transitionBySource
	}

	_, ok := m.stateByName[initialState.Name()]
	if !ok {
		return nil, fmt.Errorf("initial state %#+v not in states", initialState.Name())
	}

	m.currentState = initialState

	return &m, nil
}

func (m *Machine) State() string {
	return m.currentState.Name()
}

func (m *Machine) Transition(name string, ctx context.Context) (context.Context, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	transitionBySource, ok := m.transitionBySourceByName[name]
	if !ok {
		return nil, fmt.Errorf("transition %#+v not known", name)
	}

	transition, ok := transitionBySource[m.currentState]
	if !ok {
		return nil, fmt.Errorf(
			"transition %#+v not valid for current state %#+v",
			name, m.currentState.Name(),
		)
	}

	getScope := func(ctx context.Context) Scope {
		return Scope{
			Transition:  transition.Name(),
			Source:      transition.source.Name(),
			Destination: transition.destination.Name(),
			Context:     ctx,
		}
	}

	var err error

	ctx, err = transition.enter(getScope(ctx))
	if err != nil {
		return nil, err
	}

	if transition.destination != transition.source {
		ctx, err = transition.source.exit(getScope(ctx))
		if err != nil {
			return nil, err
		}

		m.currentState = transition.destination

		ctx, err = transition.destination.enter(getScope(ctx))
		if err != nil {
			m.currentState = transition.source
			return nil, err
		}
	}

	ctx, err = transition.exit(getScope(ctx))
	if err != nil {
		return nil, err
	}

	return ctx, nil
}
