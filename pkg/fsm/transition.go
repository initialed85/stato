package fsm

type Transition struct {
	source      *State
	destination *State
	*Named
	*Callbacks
}

func NewTransition(
	name string,
	source *State,
	destination *State,
	enterCallback Callback,
	exitCallback Callback,
) *Transition {
	t := Transition{
		source:      source,
		destination: destination,
		Named:       NewNamed(name),
		Callbacks: NewCallbacks(
			enterCallback,
			exitCallback,
		),
	}

	if source == nil {
		panic("source state unexpectedly nil")
	}

	if destination == nil {
		panic("destination state unexpectedly nil")
	}

	return &t
}

func (t *Transition) GetSource() *State {
	return t.source
}

func (t *Transition) GetDestination() *State {
	return t.destination
}
