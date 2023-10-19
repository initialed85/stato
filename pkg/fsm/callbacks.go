package fsm

import "context"

type Scope struct {
	Transition  string
	Source      string
	Destination string
	Context     context.Context
}

type Callback func(scope Scope) (context.Context, error)

type Callbacks struct {
	enterCallback Callback
	exitCallback  Callback
}

func NewCallbacks(
	enterCallback Callback,
	exitCallback Callback,
) *Callbacks {
	c := Callbacks{
		enterCallback: enterCallback,
		exitCallback:  exitCallback,
	}

	return &c
}

func (c *Callbacks) enter(scope Scope) (context.Context, error) {
	if c.enterCallback == nil {
		return scope.Context, nil
	}

	return c.enterCallback(scope)
}

func (c *Callbacks) exit(scope Scope) (context.Context, error) {
	if c.exitCallback == nil {
		return scope.Context, nil
	}

	return c.exitCallback(scope)
}
