package fsm

type State struct {
	*Named
	*Callbacks
}

func NewState(
	name string,
	enterCallback Callback,
	exitCallback Callback,
) *State {
	s := State{
		Named: NewNamed(name),
		Callbacks: NewCallbacks(
			enterCallback,
			exitCallback,
		),
	}

	return &s
}
