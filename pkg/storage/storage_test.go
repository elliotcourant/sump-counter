package storage

import (
	"context"
	"github.com/elliotcourant/sump-counter/pkg/models"
	"github.com/elliotcourant/sump-counter/pkg/pio"
	"github.com/elliotcourant/sump-counter/pkg/testutils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestNewStorage(t *testing.T) {
	t.Run("get cyles", func(t *testing.T) {
		testutils.NewTestWithContext(t, func(t *testing.T, ctx context.Context) {
			tmp, err := ioutil.TempDir("", "sump-boi")
			require.NoError(t, err)
			defer func() {
				os.RemoveAll(tmp)
			}()

			logger := logrus.New()
			logger.SetLevel(logrus.TraceLevel)

			config := Configuration{
				InMemory: false,
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "password",
				Database: "sumpdata",
				Log:      logger.WithFields(logrus.Fields{}),
			}

			storage, err := newSqlStorage(config)
			assert.NoError(t, err)
			assert.NotNil(t, storage)

			pumps, err := storage.GetAllPumps(ctx)
			require.NoError(t, err)
			require.NotEmpty(t, pumps)

			t.Run("state base", func(t *testing.T) {
				pump := pumps[0]

				cycle := models.Cycle{
					PumpId:    pump.PumpId,
					StartTime: time.Now().UTC(),
					EndTime:   nil,
					Notified:  false,
				}
				err := storage.StartCycle(ctx, &cycle)
				require.NoError(t, err)

				state, err := storage.GetLastPumpState(ctx, pump.PumpId)
				require.NoError(t, err)
				require.Equal(t, pio.High, state)

				err = storage.FinishCycle(ctx, cycle.CycleId, time.Now().UTC())
				require.NoError(t, err)

				state, err = storage.GetLastPumpState(ctx, pump.PumpId)
				require.NoError(t, err)
				require.Equal(t, pio.Low, state)
			})

			t.Run("state persist", func(t *testing.T) {
				pump := pumps[0]

				state, err := storage.GetLastPumpState(ctx, pump.PumpId)
				require.NoError(t, err)
				require.Equal(t, pio.Low, state)

				err = storage.PersistPumpState(ctx, pump.PumpId, pio.High, time.Now().UTC())
				require.NoError(t, err)

				time.Sleep(1 * time.Second)

				state, err = storage.GetLastPumpState(ctx, pump.PumpId)
				require.NoError(t, err)
				require.Equal(t, pio.High, state)

				err = storage.PersistPumpState(ctx, pump.PumpId, pio.Low, time.Now().UTC())
				require.NoError(t, err)

				state, err = storage.GetLastPumpState(ctx, pump.PumpId)
				require.NoError(t, err)
				require.Equal(t, pio.Low, state)
			})
		})
	})
}
