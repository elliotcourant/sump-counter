package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/elliotcourant/sump-counter/pkg/models"
	"github.com/elliotcourant/sump-counter/pkg/pio"
	"github.com/elliotcourant/sump-counter/pkg/tracing"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

const (
	maxAttempts       = 3
	sequenceBandwidth = 10
)

var (
	_ Storage = &sqlStorage{}
)

type (
	Configuration struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
		InMemory bool
		Log      *logrus.Entry
	}

	Storage interface {
		GetAllPumps(ctx context.Context) ([]models.Pump, error)
		GetLastPumpState(ctx context.Context, pumpId int) (pio.State, error)
		PersistPumpState(ctx context.Context, pumpId int, state pio.State, timestamp time.Time) error
		GetCyclesToNotify(ctx context.Context) ([]models.Cycle, error)
		SetCyclesNotified(ctx context.Context, cycleIds []int) error

		Close() error
	}

	sqlStorage struct {
		configuration Configuration
		log           *logrus.Entry
		db            *sql.DB
		orm           *goqu.Database
	}
)

func (s *sqlStorage) GetAllPumps(ctx context.Context) ([]models.Pump, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	return s.ListAllPumps(ctx)
}

func (s *sqlStorage) GetLastPumpState(ctx context.Context, pumpId int) (pio.State, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	lastCycle, ok, err := s.LatestCycleForPump(ctx, pumpId)
	if err != nil {
		return pio.Unknown, err
	}

	if !ok {
		return pio.Unknown, nil
	}

	if lastCycle.EndTime == nil {
		return pio.High, nil
	}

	return pio.Low, nil
}

func (s *sqlStorage) GetCyclesToNotify(ctx context.Context) ([]models.Cycle, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	return s.ListNonNotifiedCycles(ctx)
}

func (s *sqlStorage) SetCyclesNotified(ctx context.Context, cycleIds []int) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	return s.SetCyclesAsNotified(ctx, cycleIds...)
}

func (s *sqlStorage) PersistPumpState(ctx context.Context, pumpId int, state pio.State, timestamp time.Time) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	// TODO (elliotcourant) Add pump name and pin
	log := s.log.WithField("pump", pumpId)

	lastCycle, ok, err := s.LatestCycleForPump(ctx, pumpId)
	if err != nil {
		return err
	}

	switch state {
	case pio.Low:
		if lastCycle.EndTime != nil {
			// If we are observing a low state but the latest cycle is already closed then we want to do nothing.
			log.Warningf("observed low state for pump that was not running")
			return nil
		}
	case pio.High:
		if lastCycle.EndTime == nil && ok {
			// If we are observing a high state but the latest cycle has no end time then we are seeing a duplicate
			// state. We should not do anything.
			log.Warningf("observed high state for pump that was already running")
			return nil
		} else {
			// If the end time was not nil though, and we are seeing another high, that means a new cycle has begun.
			// We want to set ok to false so that we create a new cycle.
			ok = false
		}
	}

	if ok {
		log.Debugf("finishing pump cycle")
		// If we are ok, that means we can finish the current cycle. This should only happen when there is no end time
		// for the current cycle and we are observing a low state.
		return s.FinishCycle(ctx, lastCycle.CycleId, timestamp)
	}

	cycle := models.Cycle{
		PumpId:    pumpId,
		StartTime: timestamp,
		EndTime:   nil,
		Notified:  false,
	}

	log.Debugf("starting a new pump cycle")
	return s.StartCycle(ctx, &cycle)
}

func (s *sqlStorage) Close() error {
	return s.db.Close()
}

func NewStorage(configuration Configuration) (Storage, error) {
	return newSqlStorage(configuration)
}

func getConnectionString(configuration Configuration) string {
	props := []string{
		"sslmode=disable",
	}

	props = append(props, fmt.Sprintf("user=%s", configuration.User))
	props = append(props, fmt.Sprintf("dbname=%s", configuration.Database))
	props = append(props, fmt.Sprintf("password=%s", configuration.Password))
	props = append(props, fmt.Sprintf("host=%s", configuration.Host))
	props = append(props, fmt.Sprintf("port=%d", configuration.Port))

	return strings.Join(props, " ")
}

func newSqlStorage(configuration Configuration) (*sqlStorage, error) {

	db, err := sql.Open("postgres", getConnectionString(configuration))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &sqlStorage{
		configuration: configuration,
		log:           configuration.Log,
		db:            db,
		orm:           goqu.New("postgres", db),
	}, nil
}
