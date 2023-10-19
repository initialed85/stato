package example

import (
	"context"
	"fmt"
	"github.com/initialed85/stato/pkg/fsm"
	"log"
	"sync"
)

type Connector struct {
	chargingStation *ChargingStation
	connectorID     int
	machine         *fsm.Machine
	mu              sync.Mutex
	transaction     *Transaction
}

func NewConnector(chargingStation *ChargingStation, connectorID int) (*Connector, error) {
	c := Connector{
		chargingStation: chargingStation,
		connectorID:     connectorID,
		transaction:     nil,
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
		nil,
	)

	preparing := fsm.NewState(
		Preparing,
		c.onOccupied,
		nil,
	)

	charging := fsm.NewState(
		Charging,
		c.onOccupied,
		nil,
	)

	suspendedEV := fsm.NewState(
		SuspendedEV,
		c.onOccupied,
		nil,
	)

	suspendedEVSE := fsm.NewState(
		SuspendedEVSE,
		c.onOccupied,
		nil,
	)

	finishing := fsm.NewState(
		Finishing,
		c.onOccupied,
		nil,
	)

	reserved := fsm.NewState(
		Reserved,
		c.onUnavailable,
		nil,
	)

	unavailable := fsm.NewState(
		Unavailable,
		c.onUnavailable,
		nil,
	)

	faulted := fsm.NewState(
		Faulted,
		c.onFaulted,
		nil,
	)

	transitions := make([]*fsm.Transition, 0)

	alpha := map[string]*fsm.State{
		"A": available,
		"B": preparing,
		"C": charging,
		"D": suspendedEV,
		"E": suspendedEVSE,
		"F": finishing,
		"G": reserved,
		"H": unavailable,
		"I": faulted,
	}

	numeric := map[string]*fsm.State{
		"1": available,
		"2": preparing,
		"3": charging,
		"4": suspendedEV,
		"5": suspendedEVSE,
		"6": finishing,
		"7": reserved,
		"8": unavailable,
		"9": faulted,
	}

	for _, destinationState := range alpha {
		transitions = append(
			transitions,
			fsm.NewTransition(
				fmt.Sprintf("%v%v", HandleStatusNotification, destinationState.Name()),
				initialised,
				destinationState,
				nil,
				nil,
			),
		)
	}

	for _, destinationState := range alpha {
		transitions = append(
			transitions,
			fsm.NewTransition(
				fmt.Sprintf("%v%v", HandleStatusNotification, destinationState.Name()),
				destinationState,
				destinationState,
				nil,
				nil,
			),
		)
	}

	permutations := []string{
		"__", "A2", "A3", "A4", "A5", "__", "A7", "A8", "A9",
		"B1", "__", "B3", "B4", "B5", "B6", "__", "__", "B9",
		"C1", "__", "__", "C4", "C5", "C6", "__", "C8", "C9",
		"D1", "__", "D3", "__", "D5", "D6", "__", "D8", "D9",
		"E1", "__", "E3", "E4", "__", "E6", "__", "E8", "E9",
		"F1", "F2", "__", "__", "__", "__", "__", "F8", "F9",
		"G1", "G2", "__", "__", "__", "__", "__", "G8", "G9",
		"H1", "H2", "H3", "H4", "H5", "__", "__", "__", "H9",
		"I1", "I2", "I3", "I4", "I5", "I6", "I7", "I8", "__",
	}

	for _, permutation := range permutations {
		source := string(permutation[0])
		destination := string(permutation[1])

		sourceState, alphaOk := alpha[source]
		destinationState, numericOk := numeric[destination]

		if !(alphaOk && numericOk) {
			continue
		}

		transitions = append(
			transitions,
			fsm.NewTransition(
				fmt.Sprintf("%v%v", HandleStatusNotification, destinationState.Name()),
				sourceState,
				destinationState,
				nil,
				nil,
			),
		)
	}

	configure := fsm.NewTransition(
		Configure,
		uninitialised,
		initialised,
		nil,
		nil,
	)

	remoteStart1 := fsm.NewTransition(
		RemoteStart,
		available,
		charging,
		c.onRemoteStart,
		c.onCharging,
	)

	remoteStart2 := fsm.NewTransition(
		RemoteStart,
		preparing,
		charging,
		c.onRemoteStart,
		c.onCharging,
	)

	remoteStop1 := fsm.NewTransition(
		RemoteStop,
		charging,
		finishing,
		nil,
		c.onFinishing,
	)

	remoteStop2 := fsm.NewTransition(
		RemoteStop,
		suspendedEV,
		finishing,
		nil,
		c.onFinishing,
	)

	remoteStop3 := fsm.NewTransition(
		RemoteStop,
		suspendedEVSE,
		finishing,
		nil,
		c.onFinishing,
	)

	transitions = append(
		transitions,
		[]*fsm.Transition{
			configure,
			remoteStart1,
			remoteStart2,
			remoteStop1,
			remoteStop2,
			remoteStop3,
		}...,
	)

	machine, err := fsm.NewMachine(
		[]*fsm.State{
			uninitialised,
			initialised,
			available,
			preparing,
			charging,
			suspendedEV,
			suspendedEVSE,
			finishing,
			reserved,
			unavailable,
			faulted,
		},
		transitions,
		uninitialised,
	)
	if err != nil {
		return nil, err
	}

	c.machine = machine

	return &c, nil
}

func (c *Connector) onInitialised(scope fsm.Scope) (context.Context, error) {
	log.Printf(
		"%v:%v created if it didn't exist",
		c.chargingStation.GetChargingStationID(), c.connectorID,
	)
	return scope.Context, nil
}

func (c *Connector) onAvailable(scope fsm.Scope) (context.Context, error) {
	log.Printf(
		"%v:%v available (%v)",
		c.chargingStation.GetChargingStationID(), c.connectorID, c.State(),
	)
	return scope.Context, nil
}

func (c *Connector) onOccupied(scope fsm.Scope) (context.Context, error) {
	log.Printf(
		"%v:%v occupied (%v)",
		c.chargingStation.GetChargingStationID(), c.connectorID, c.State(),
	)
	return scope.Context, nil
}

func (c *Connector) onUnavailable(scope fsm.Scope) (context.Context, error) {
	log.Printf(
		"%v:%v unavailable (%v)",
		c.chargingStation.GetChargingStationID(), c.connectorID, c.State(),
	)
	return scope.Context, nil
}

func (c *Connector) onFaulted(scope fsm.Scope) (context.Context, error) {
	log.Printf(
		"%v:%v faulted (%v)",
		c.chargingStation.GetChargingStationID(), c.connectorID, c.State(),
	)
	return scope.Context, nil
}

func (c *Connector) onRemoteStart(scope fsm.Scope) (context.Context, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.transaction != nil {
		return nil, fmt.Errorf(
			"connector %#+v already in transaction %#+v",
			c.connectorID, c.transaction.GetTransactionID(),
		)
	}

	transaction, err := NewTransaction(c.chargingStation, c)
	if err != nil {
		return nil, err
	}
	c.transaction = transaction

	scope.Context = context.WithValue(scope.Context, "transaction", transaction)

	return scope.Context, nil
}

func (c *Connector) onCharging(scope fsm.Scope) (context.Context, error) {
	log.Printf(
		"%v:%v EV charging underway (%v)",
		c.chargingStation.GetChargingStationID(), c.connectorID, c.State(),
	)
	return scope.Context, nil
}

func (c *Connector) onFinishing(scope fsm.Scope) (context.Context, error) {
	log.Printf(
		"%v:%v EV charging finished (%v)",
		c.chargingStation.GetChargingStationID(), c.connectorID, c.State(),
	)
	return scope.Context, nil
}

func (c *Connector) GetConnectorID() int {
	return c.connectorID
}

func (c *Connector) GetTransaction() (*Transaction, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.transaction == nil {
		return nil, fmt.Errorf(
			"%v:%v has no transaction",
			c.chargingStation.GetChargingStationID(), c.connectorID,
		)
	}

	return c.transaction, nil
}

func (c *Connector) State() string {
	return c.machine.State()
}

func (c *Connector) Configure(ctx context.Context) (err error) {
	ctx, err = c.machine.Transition(Configure, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connector) HandleStatusNotification(ctx context.Context, status string) (err error) {
	ctx, err = c.machine.Transition(
		fmt.Sprintf("%v%v", HandleStatusNotification, status),
		context.WithValue(ctx, "status", status),
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connector) RemoteStart(ctx context.Context) (transaction *Transaction, err error) {
	ctx, err = c.machine.Transition(RemoteStart, ctx)
	if err != nil {
		return nil, err
	}

	transaction, ok := ctx.Value("transaction").(*Transaction)
	if !ok {
		return nil, fmt.Errorf("failed to cast transaction %#+v", transaction)
	}

	err = transaction.Configure(ctx)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (c *Connector) HandleStart(ctx context.Context, transactionID int64) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.transaction == nil {
		return fmt.Errorf(
			"%v:%v:%v unexpectedly nil",
			c.chargingStation.GetChargingStationID(), c.connectorID, transactionID,
		)
	}

	err = c.transaction.HandleStart(ctx, transactionID)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connector) RemoteStop(ctx context.Context) (err error) {
	ctx, err = c.machine.Transition(RemoteStop, ctx)
	if err != nil {
		return err
	}

	return nil
}
