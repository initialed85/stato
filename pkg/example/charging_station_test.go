package example

import (
	"context"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestChargingStation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// us provisioning the charger

	chargingStation, err := NewChargingStation("test_001")
	require.NoError(t, err)
	require.Equal(t, Unitialised, chargingStation.State())

	err = chargingStation.Configure(ctx)
	require.NoError(t, err)
	require.Equal(t, Initialised, chargingStation.State())

	// boot notification udpating the charger status

	err = chargingStation.HandleBootNotification(ctx, "ACME Charger 1")
	require.NoError(t, err)
	require.Equal(t, Available, chargingStation.State())

	// status notification udpating the connector status

	connector, err := chargingStation.AddConnector(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, Initialised, connector.State())

	err = connector.HandleStatusNotification(ctx, Available)
	require.NoError(t, err)
	require.Equal(t, Available, connector.State())

	// user plugging in their ev

	err = connector.HandleStatusNotification(ctx, Preparing)
	require.NoError(t, err)
	require.Equal(t, Preparing, connector.State())

	// us starting charging

	var transaction *Transaction

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		transaction, err = connector.RemoteStart(ctx)
		require.NoError(t, err)
		require.Equal(t, Charging, connector.State())
		wg.Done()
	}()
	wg.Wait()

	require.NotNil(t, transaction)
	require.Equal(t, int64(1), transaction.GetTransactionID())
	require.Equal(t, Initialised, transaction.State())

	// the charger starting charging

	err = transaction.HandleStart(ctx, transaction.GetTransactionID())
	require.NoError(t, err)
	require.Equal(t, Charging, transaction.State())
	require.Equal(t, Charging, connector.State())

	err = connector.HandleStatusNotification(ctx, Charging)
	require.NoError(t, err)
	require.Equal(t, Charging, connector.State())

}
