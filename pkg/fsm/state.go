package fsm

type State struct {
	name             string
	contextCallbacks *ContextCallbacks
}

func NewState(
	name string,
	onEnterCallbacks []func() error,
	onExitCallbacks []func() error,
) *State {
	s := State{
		name: name,
		contextCallbacks: NewContextCallbacks(
			onEnterCallbacks,
			onExitCallbacks,
		),
	}

	return &s
}

func (c *State) forceActive() {
	c.contextCallbacks.forceActive()
}

func (c *State) enter() error {
	return c.contextCallbacks.enter()
}

func (c *State) exit() error {
	return c.contextCallbacks.exit()
}

func (c *State) GetName() string {
	return c.name
}

func (c *State) GetActive() bool {
	return c.contextCallbacks.GetActive()
}
