# stato

Named after the iconic Holden Statesman, `stato` tries to be a simple state machine library providing states and
transitions with one or more on-enter and on-exit callbacks for both.

## Usage

The example below builds a state machine that's intended to be used via the `ChangeState()` approach (i.e. try to find
a unique for the source state to the destination state) but you can also explicitly use the `Transition()` approach
(i.e. try to invoke a transition that calls for a specific source state and destination state).

```golang
package state_machines

import (
	"fmt"
	"github.com/initialed85/stato/pkg/fsm"
)

type ConnectorFSM struct {
	Available     *fsm.State
	Preparing     *fsm.State
	Charging      *fsm.State
	SuspendedEV   *fsm.State
	SuspendedEVSE *fsm.State
	Finishing     *fsm.State
	Reserved      *fsm.State
	Unavailable   *fsm.State
	Faulted       *fsm.State
	states        []*fsm.State
	transitions   []*fsm.Transition
	machine       *fsm.Machine
}

func NewConnectorFSM() (*ConnectorFSM, error) {
	f := ConnectorFSM{
		Available: fsm.NewState(
			"Available",
			[]func() error{},
			[]func() error{},
		),
		Preparing: fsm.NewState(
			"Preparing",
			[]func() error{},
			[]func() error{},
		),
		Charging: fsm.NewState(
			"Charging",
			[]func() error{},
			[]func() error{},
		),
		SuspendedEV: fsm.NewState(
			"SuspendedEV",
			[]func() error{},
			[]func() error{},
		),
		SuspendedEVSE: fsm.NewState(
			"SuspendedEVSE",
			[]func() error{},
			[]func() error{},
		),
		Finishing: fsm.NewState(
			"Finishing",
			[]func() error{},
			[]func() error{},
		),
		Reserved: fsm.NewState(
			"Reserved",
			[]func() error{},
			[]func() error{},
		),
		Unavailable: fsm.NewState(
			"Unavailable",
			[]func() error{},
			[]func() error{},
		),
		Faulted: fsm.NewState(
			"Faulted",
			[]func() error{},
			[]func() error{},
		),
	}

	f.states = []*fsm.State{
		f.Available,
		f.Preparing,
		f.Charging,
		f.SuspendedEV,
		f.SuspendedEVSE,
		f.Finishing,
		f.Reserved,
		f.Unavailable,
		f.Faulted,
	}

	alpha := map[string]*fsm.State{
		"A": f.Available,
		"B": f.Preparing,
		"C": f.Charging,
		"D": f.SuspendedEV,
		"E": f.SuspendedEVSE,
		"F": f.Finishing,
		"G": f.Reserved,
		"H": f.Unavailable,
		"I": f.Faulted,
	}

	numeric := map[string]*fsm.State{
		"1": f.Available,
		"2": f.Preparing,
		"3": f.Charging,
		"4": f.SuspendedEV,
		"5": f.SuspendedEVSE,
		"6": f.Finishing,
		"7": f.Reserved,
		"8": f.Unavailable,
		"9": f.Faulted,
	}

	permutations := []string{
		"__", "A2", "A3", "A4", "A5", "__", "A7", "A8", "A9",
		"B1", "__", "B3", "B4", "B5", "__", "__", "__", "B9",
		"C1", "__", "__", "C4", "C5", "C6", "__", "C8", "C9",
		"D1", "__", "D3", "__", "D5", "D6", "__", "D8", "D9",
		"E1", "__", "E3", "E4", "__", "E6", "__", "E8", "E9",
		"F1", "F2", "__", "__", "__", "__", "__", "F8", "F9",
		"G1", "G2", "__", "__", "__", "__", "__", "G8", "G9",
		"H1", "H2", "H3", "H4", "H5", "__", "__", "__", "H9",
		"I1", "I2", "I3", "I4", "I5", "I6", "I7", "I8", "_",
	}

	for _, permutation := range permutations {
		source := permutation[0:1]
		destination := permutation[1:2]

		sourceState, alphaOk := alpha[source]
		destinationState, numericOk := numeric[destination]

		if !(alphaOk && numericOk) {
			continue
		}

		f.transitions = append(
			f.transitions,
			fsm.NewTransition(
				fmt.Sprintf("%vTo%v", source, destination),
				[]func() error{},
				[]func() error{},
				sourceState,
				destinationState,
			),
		)
	}

	machine, err := fsm.NewMachine(
		"Connector",
		f.states,
		f.transitions,
		f.Unavailable,
	)
	if err != nil {
		return nil, err
	}

	f.machine = machine

	return &f, nil
}

func (c *ConnectorFSM) ChangeState(destinationState *fsm.State) error {
	return c.machine.ChangeState(destinationState)
}
```
