package fsm

import (
	"sync"
)

type ContextCallbacks struct {
	mu               sync.Mutex
	active           bool
	onEnterCallbacks []func() error
	onExitCallbacks  []func() error
}

func NewContextCallbacks(
	onEnterCallbacks []func() error,
	onExitCallbacks []func() error,
) *ContextCallbacks {
	c := ContextCallbacks{
		onEnterCallbacks: onEnterCallbacks,
		onExitCallbacks:  onExitCallbacks,
		active:           false,
	}

	return &c
}

// used to force an active state without firing the onEnterCallbacks
func (c *ContextCallbacks) forceActive() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.active = true
}

// fire the onEnterCallbacks if we're not active and mark ourselves as active
func (c *ContextCallbacks) enter() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// if we're already active, do nothing
	if c.active {
		return nil
	}

	var err error
	for _, callback := range c.onEnterCallbacks {
		err = callback()

		// failed callback will mark a failure to enter
		if err != nil {
			return err
		}
	}

	// we're active now
	c.active = true

	return nil
}

// fire the onExitCallbacks if we're active and mark ourselves as not active
func (c *ContextCallbacks) exit() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// if we're not active, do nothing
	if !c.active {
		return nil
	}

	var err error
	for _, callback := range c.onExitCallbacks {
		err = callback()

		// failed callback will mark a failure to exit
		if err != nil {
			return err
		}
	}

	// we're not active now
	c.active = false

	return nil
}

func (c *ContextCallbacks) GetActive() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.active
}
