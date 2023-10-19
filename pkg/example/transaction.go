package example

import (
	"context"
	"fmt"
	"github.com/initialed85/stato/pkg/fsm"
	"log"
	"sync"
)

var (
	mu                    sync.Mutex
	transactionIDSequence int64 = 0
)

type Transaction struct {
	chargingStation *ChargingStation
	connector       *Connector
	transactionID   int64
	machine         *fsm.Machine
}

func NewTransaction(
	chargingStation *ChargingStation,
	connector *Connector,
) (*Transaction, error) {
	mu.Lock()
	defer mu.Unlock()

	transactionIDSequence++

	t := Transaction{
		chargingStation: chargingStation,
		connector:       connector,
		transactionID:   transactionIDSequence,
		machine:         nil,
	}

	uninitialised := fsm.NewState(
		Unitialised,
		nil,
		nil,
	)

	initialised := fsm.NewState(
		Initialised,
		t.onInitialised,
		nil,
	)

	charging := fsm.NewState(
		Charging,
		t.onCharging,
		nil,
	)

	parking := fsm.NewState(
		Parking,
		t.onParking,
		nil,
	)

	done := fsm.NewState(
		Done,
		t.onDone,
		nil,
	)

	failed := fsm.NewState(
		Failed,
		t.onFailed,
		nil,
	)

	configure := fsm.NewTransition(
		Configure,
		uninitialised,
		initialised,
		nil,
		nil,
	)

	handleStart := fsm.NewTransition(
		HandleStart,
		initialised,
		charging,
		t.onCharging,
		nil,
	)

	meterValues := fsm.NewTransition(
		HandleMeterValues,
		charging,
		charging,
		t.onMeterValues,
		nil,
	)

	handleStop := fsm.NewTransition(
		HandleStop,
		charging,
		done,
		t.onParking,
		nil,
	)

	machine, err := fsm.NewMachine(
		[]*fsm.State{
			uninitialised,
			initialised,
			charging,
			parking,
			done,
			failed,
		},
		[]*fsm.Transition{
			configure,
			handleStart,
			meterValues,
			handleStop,
		},
		uninitialised,
	)
	if err != nil {
		return nil, err
	}

	t.machine = machine

	return &t, nil
}

func (t *Transaction) onInitialised(scope fsm.Scope) (context.Context, error) {
	log.Printf(
		"%v:%v:%v created",
		t.chargingStation.GetChargingStationID(),
		t.connector.GetConnectorID(),
		t.transactionID,
	)

	return scope.Context, nil
}

func (t *Transaction) onCharging(scope fsm.Scope) (context.Context, error) {
	transactionID, ok := scope.Context.Value("transactionID").(int64)
	if !ok {
		return nil, fmt.Errorf(
			"transaction %v failed to cast transactionID %#+v",
			t.transactionID, transactionID,
		)
	}

	if transactionID != t.transactionID {
		return nil, fmt.Errorf(
			"transaction %v got unexpected transactionID %v",
			t.transactionID, transactionID,
		)
	}

	log.Printf(
		"%v:%v:%v charging",
		t.chargingStation.GetChargingStationID(),
		t.connector.GetConnectorID(),
		t.transactionID,
	)
	return scope.Context, nil
}

func (t *Transaction) onMeterValues(scope fsm.Scope) (context.Context, error) {
	log.Printf(
		"%v:%v:%v meterValues",
		t.chargingStation.GetChargingStationID(),
		t.connector.GetConnectorID(),
		t.transactionID,
	)
	return scope.Context, nil
}

func (t *Transaction) onParking(scope fsm.Scope) (context.Context, error) {
	log.Printf(
		"%v:%v:%v parking",
		t.chargingStation.GetChargingStationID(),
		t.connector.GetConnectorID(),
		t.transactionID,
	)
	return scope.Context, nil
}

func (t *Transaction) onDone(scope fsm.Scope) (context.Context, error) {
	log.Printf(
		"%v:%v:%v done",
		t.chargingStation.GetChargingStationID(),
		t.connector.GetConnectorID(),
		t.transactionID,
	)
	return scope.Context, nil
}

func (t *Transaction) onFailed(scope fsm.Scope) (context.Context, error) {
	log.Printf(
		"%v:%v:%v failed",
		t.chargingStation.GetChargingStationID(),
		t.connector.GetConnectorID(),
		t.transactionID,
	)
	return scope.Context, nil
}

func (t *Transaction) GetTransactionID() int64 {
	return t.transactionID
}

func (t *Transaction) State() string {
	return t.machine.State()
}

func (t *Transaction) Configure(ctx context.Context) (err error) {
	ctx, err = t.machine.Transition(Configure, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (t *Transaction) HandleStart(ctx context.Context, transactionID int64) (err error) {
	ctx, err = t.machine.Transition(HandleStart, context.WithValue(ctx, "transactionID", transactionID))
	if err != nil {
		return err
	}

	return nil

}

func (t *Transaction) HandleMeterValues(ctx context.Context) (err error) {
	ctx, err = t.machine.Transition(HandleMeterValues, ctx)
	if err != nil {
		return err
	}

	return nil

}

func (t *Transaction) HandleStop(ctx context.Context) (err error) {
	ctx, err = t.machine.Transition(HandleStop, ctx)
	if err != nil {
		return err
	}

	return nil

}
