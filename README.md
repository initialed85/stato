# stato

Named after the iconic Holden Statesman, `stato` tries to be a simple state machine library providing states and transitions with one or more on-enter and on-exit callbacks for both.

## Usage

Take a look at `pkg/machine_test.go` for some full code, but as a summary:

**Define two states**
```
stateA := NewState(
		"state_a",
		[]func() error{
			func() error {
				log.Printf("entering stateA")
				return nil
			},
		},
		[]func() error{
			func() error {
				log.Printf("exiting stateA")
				return nil
			},
		},
	)

stateB := NewState(
		"state_b",
		[]func() error{
			func() error {
				log.Printf("entering stateB")
				return nil
			},
		},
		[]func() error{
			func() error {
				log.Printf("exiting stateA")
				return nil
			},
		},
	)
```

**Define a transition**
```
transitionAB := NewTransition(
		"transition_a_b",
		[]func() error{
			func() error {
				log.Printf("entering transition A-B")
				return nil
			},
		},
		[]func() error{
			func() error {
				log.Printf("exiting transition A-B")
				return nil
			},
		},
		stateA,
		stateB,
	)
```

**Tie it together with a machine**
```
machine := NewMachine(
        "some_machine",     // name
        []*State{
            stateA,
            stateB,
        },                  // states
        []*Transition{
            transitionAB,
        },                  // transitions
        stateA,             // initial state
    )
```

At this point, a state machine exists in `stateA` but has fired no callbacks.

**Invoke a transition**
```
err = machine.Transition(transitionAB)
if err != nil {
    log.Fatal(err)
}
```

This should have worked without error- now the same state machine is in `stateB` and will have fired the callbacks causing the following output:

- `entering transitionA-B` 
- `exiting stateA` 
- `entering stateB` 
- `exiting transitionA-B` 

**Fail to invoke a transition**
```
err = machine.Transition(transitionAB)
if err != nil {
    log.Fatal(err)
}
```

This should explode because `transitionAB` requires that we're in `stateA`, but we're in `stateB`

## Suggestions

You may note that you need references to the `Transition` structs you define and the `Machine` you create in order to invoke a transition.

To save passing these around, one approach would be to define a struct that contains your states, transitions and the machine and then expose the means to invoke various transitions with methods.

To build on the earlier snippet:

```
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
```

This way the part of your application that needs to invoke a state change can just call that method and if the transition is permitted, the method will return nil.

It is of course critical to check for errors because that's the only way to know that the transition wasn't accepted!
