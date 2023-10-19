package fsm

type Named struct {
	name string
}

func NewNamed(name string) *Named {
	n := Named{
		name: name,
	}

	return &n
}

func (n *Named) Name() string {
	return n.name
}
