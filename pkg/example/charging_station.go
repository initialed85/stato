package example

import (
	"context"
	"fmt"
	"github.com/initialed85/stato/pkg/fsm"
	"log"
	"sync"
)

type ChargingStation struct {
	chargingStationID string
	machine           *fsm.Machine
	mu                sync.Mutex
	connectorByID     map[int]*Connector
}

func NewChargingStation(
	chargingStationID string,
) (*ChargingStation, error) {
	c := ChargingStation{
		chargingStationID: chargingStationID,
		connectorByID:     make(map[int]*Connector),
	}

	uninitialised := fsm.NewState(
		Unitialised,
		nil,
		nil,
	)

	initialised := fsm.NewState(
		Initialised,
		c.onInitialised,
		nil,
	)

	available := fsm.NewState(
		Available,
		c.onAvailable,
		c.onUnavailable,
	)

	unavailable := fsm.NewState(
		Unavailable,
		nil,
		nil,
	)

	configure := fsm.NewTransition(
		Configure,
		uninitialised,
		initialised,
		c.onConfiguring,
		nil,
	)

	handleBootNotification := fsm.NewTransition(
		HandleBootNotification,
		initialised,
		available,
		c.onHandleBootNotification,
		nil,
	)

	shutdown := fsm.NewTransition(
		Shutdown,
		available,
		unavailable,
		c.onShutdown,
		nil,
	)

	machine, err := fsm.NewMachine(
		[]*fsm.State{
			uninitialised,
			initialised,
			available,
			unavailable,
		},
		[]*fsm.Transition{
			configure,
			handleBootNotification,
			shutdown,
		},
		uninitialised,
	)
	if err != nil {
		return nil, err
	}

	c.machine = machine

	return &c, nil
}

func (c *ChargingStation) onInitialised(scope fsm.Scope) (context.Context, error) {
	log.Printf("%v created if it didn't exist", c.chargingStationID)
	return scope.Context, nil
}

func (c *ChargingStation) onAvailable(scope fsm.Scope) (context.Context, error) {
	log.Printf("%v available", c.chargingStationID)
	return scope.Context, nil
}

func (c *ChargingStation) onUnavailable(scope fsm.Scope) (context.Context, error) {
	log.Printf("%v unavailable", c.chargingStationID)
	return scope.Context, nil
}

func (c *ChargingStation) onConfiguring(scope fsm.Scope) (context.Context, error) {
	log.Printf("charging station  %v configuring...", c.chargingStationID)
	return scope.Context, nil
}

func (c *ChargingStation) onHandleBootNotification(scope fsm.Scope) (context.Context, error) {
	model := scope.Context.Value("model").(string)
	log.Printf("%v booted, model: %v...", c.chargingStationID, model)
	return scope.Context, nil
}

func (c *ChargingStation) onShutdown(scope fsm.Scope) (context.Context, error) {
	log.Printf("%v shutting down...", c.chargingStationID)
	return scope.Context, nil
}

func (c *ChargingStation) AddConnector(ctx context.Context, connectorID int) (*Connector, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.connectorByID[connectorID]
	if ok {
		return nil, fmt.Errorf("connectorID %#+v already exists", connectorID)
	}

	connector, err := NewConnector(c, connectorID)
	if err != nil {
		return nil, err
	}

	err = connector.Configure(ctx)
	if err != nil {
		return nil, err
	}

	c.connectorByID[connectorID] = connector

	return connector, nil
}

func (c *ChargingStation) GetChargingStationID() string {
	return c.chargingStationID
}

func (c *ChargingStation) GetConnector(connectorID int) (*Connector, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	connector, ok := c.connectorByID[connectorID]
	if !ok {
		return nil, fmt.Errorf("connectorID %#+v doesn't exist", connectorID)
	}

	return connector, nil
}

func (c *ChargingStation) State() string {
	return c.machine.State()
}

func (c *ChargingStation) Configure(ctx context.Context) (err error) {
	ctx, err = c.machine.Transition(Configure, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChargingStation) HandleBootNotification(ctx context.Context, model string) (err error) {
	ctx, err = c.machine.Transition(
		HandleBootNotification,
		context.WithValue(ctx, "model", model),
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChargingStation) Shutdown(ctx context.Context) (err error) {
	ctx, err = c.machine.Transition(Shutdown, ctx)
	if err != nil {
		return err
	}

	return nil
}
