package ioserver

import (
	"github.com/elliotcourant/sump-counter/pkg/mocks/pio_mock"
	"github.com/elliotcourant/sump-counter/pkg/pio"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"sync"
	"testing"
	"time"
)

func TestGpioServerBase_WatchPinState(t *testing.T) {
	log := logrus.New()
	t.Run("simple", func(t *testing.T) {
		socketPath, socketPathCleanup := NewSocketPath(t)
		defer socketPathCleanup()

		io := pio_mock.NewPIOMock()

		server, err := NewIOServer(logrus.NewEntry(log), socketPath, io)
		assert.NoError(t, err, "io server should have been created successfully")

		defer server.Close()

		watchServer, watchClient, watchCancel := NewMockWatchServerClient(t)

		var wg sync.WaitGroup
		wg.Add(2)

		pinRequest := Pin{
			Pin: 13,
		}

		numberOfChanges := 10

		// Start the watch pin server
		go func(server IOServer) {
			defer wg.Done()
			err := server.WatchPinState(&pinRequest, watchServer)
			assert.NoError(t, err, "watch pin state should have finished without error")
		}(server)

		changesObserved := -1
		go func(watchClient GPIOService_WatchPinStateClient) {
			defer wg.Done()
			previousState := State(math.MaxUint8)
			for {
				state, err := watchClient.Recv()
				select {
				case <-watchClient.Context().Done():
					return
				default:
					assert.NoError(t, err, "received state should not have an error")
					require.NotEqual(t, previousState, state.State, "new state should not match previous state")

					// fmt.Printf("Received change: %s -> %s\n", previousState, state.State)

					changesObserved++
					previousState = state.State
				}
			}
		}(watchClient)

		previous := pio.Low
		pin := pinRequest.IOPin()
		for i := 0; i < numberOfChanges; i++ {
			time.Sleep(2 * time.Second)
			switch previous {
			case pio.Low:
				io.Set(pin, pio.High)
				previous = pio.High
			case pio.High:
				io.Set(pin, pio.Low)
				previous = pio.Low
			}
		}

		time.Sleep(2 * time.Second)

		watchCancel()

		wg.Wait()
		assert.Equal(t, numberOfChanges, changesObserved, "changes observed should equal expected")
	})
}
