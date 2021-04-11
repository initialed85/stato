package fsm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextCallbacks_Normal(t *testing.T) {
	enter1CallCount := 0
	enter2CallCount := 0
	exit1CallCount := 0
	exit2CallCount := 0
	var err error

	c := NewContextCallbacks(
		[]func() error{
			func() error {
				enter1CallCount++
				return nil
			},
			func() error {
				enter2CallCount++
				return nil
			},
		},
		[]func() error{
			func() error {
				exit1CallCount++
				return nil
			},
			func() error {
				exit2CallCount++
				return nil
			},
		},
	)

	// can't exit if not yet entered
	err = c.exit()
	assert.Nil(t, err)
	assert.Equal(t, 0, enter1CallCount)
	assert.Equal(t, 0, enter2CallCount)
	assert.Equal(t, 0, exit1CallCount)
	assert.Equal(t, 0, exit2CallCount)

	// can enter
	err = c.enter()
	assert.Nil(t, err)
	assert.Equal(t, 1, enter1CallCount)
	assert.Equal(t, 1, enter2CallCount)
	assert.Equal(t, 0, exit1CallCount)
	assert.Equal(t, 0, exit2CallCount)

	// can't enter if already entered
	err = c.enter()
	assert.Nil(t, err)
	assert.Equal(t, 1, enter1CallCount)
	assert.Equal(t, 1, enter2CallCount)
	assert.Equal(t, 0, exit1CallCount)
	assert.Equal(t, 0, exit2CallCount)

	// can exit
	err = c.exit()
	assert.Nil(t, err)
	assert.Equal(t, 1, enter1CallCount)
	assert.Equal(t, 1, enter2CallCount)
	assert.Equal(t, 1, exit1CallCount)
	assert.Equal(t, 1, exit2CallCount)

	// can't exit if already exited
	err = c.exit()
	assert.Nil(t, err)
	assert.Equal(t, 1, enter1CallCount)
	assert.Equal(t, 1, enter2CallCount)
	assert.Equal(t, 1, exit1CallCount)
	assert.Equal(t, 1, exit2CallCount)
}

func TestContextCallbacks_ForceActive(t *testing.T) {
	enter1CallCount := 0
	enter2CallCount := 0
	exit1CallCount := 0
	exit2CallCount := 0
	var err error

	c := NewContextCallbacks(
		[]func() error{
			func() error {
				enter1CallCount++
				return nil
			},
			func() error {
				enter2CallCount++
				return nil
			},
		},
		[]func() error{
			func() error {
				exit1CallCount++
				return nil
			},
			func() error {
				exit2CallCount++
				return nil
			},
		},
	)

	c.forceActive()

	// can't enter if already entered
	err = c.enter()
	assert.Nil(t, err)
	assert.Equal(t, 0, enter1CallCount)
	assert.Equal(t, 0, enter2CallCount)
	assert.Equal(t, 0, exit1CallCount)
	assert.Equal(t, 0, exit2CallCount)

	// can exit
	err = c.exit()
	assert.Nil(t, err)
	assert.Equal(t, 0, enter1CallCount)
	assert.Equal(t, 0, enter2CallCount)
	assert.Equal(t, 1, exit1CallCount)
	assert.Equal(t, 1, exit2CallCount)

	// can't exit if already exited
	err = c.exit()
	assert.Nil(t, err)
	assert.Equal(t, 0, enter1CallCount)
	assert.Equal(t, 0, enter2CallCount)
	assert.Equal(t, 1, exit1CallCount)
	assert.Equal(t, 1, exit2CallCount)

	// can enter once exited
	err = c.enter()
	assert.Nil(t, err)
	assert.Equal(t, 1, enter1CallCount)
	assert.Equal(t, 1, enter2CallCount)
	assert.Equal(t, 1, exit1CallCount)
	assert.Equal(t, 1, exit2CallCount)
}

func TestContextCallbacks_ErrorDuringEnterCallback(t *testing.T) {
	var err error

	c := NewContextCallbacks(
		[]func() error{
			func() error {
				return nil
			},
			func() error {
				return fmt.Errorf("oh no")
			},
		},
		[]func() error{
			func() error {
				return nil
			},
			func() error {
				return nil
			},
		},
	)

	// can't enter
	err = c.enter()
	assert.NotNil(t, err)
	assert.False(t, c.GetActive())
}

func TestContextCallbacks_ErrorDuringExitCallback(t *testing.T) {
	var err error

	c := NewContextCallbacks(
		[]func() error{
			func() error {
				return nil
			},
			func() error {
				return nil
			},
		},
		[]func() error{
			func() error {
				return nil
			},
			func() error {
				return fmt.Errorf("oh no")
			},
		},
	)

	// can enter
	err = c.enter()
	assert.Nil(t, err)

	// can't exit
	err = c.exit()
	assert.NotNil(t, err)
	assert.True(t, c.GetActive())
}
