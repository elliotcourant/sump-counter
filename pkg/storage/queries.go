package storage

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exec"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/elliotcourant/sump-counter/pkg/models"
	"github.com/elliotcourant/sump-counter/pkg/pio"
	"github.com/elliotcourant/sump-counter/pkg/tracing"
	"github.com/pkg/errors"
	"time"
)

type queryInterface interface {
	exp.SQLExpression
	Executor() exec.QueryExecutor
}

func (s *sqlStorage) ListAllPumps(ctx context.Context) ([]models.Pump, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	query := s.orm.
		From("pumps").
		Select("pump_id", "name", "pin")

	pumps := make([]models.Pump, 0)
	err := s.executeQuery(ctx, query, func(ctx context.Context, ex exec.QueryExecutor) error {
		return errors.WithStack(ex.ScanStructsContext(ctx, &pumps))
	})

	return pumps, err
}

func (s *sqlStorage) GetPump(ctx context.Context, name string, pin pio.Pin) (pump models.Pump, ok bool, err error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	query := s.orm.
		From("pumps").
		Select("pump_id", "name", "pin").
		Where(goqu.Ex{
			"name": name,
			"pin":  pin,
		}).
		Limit(1)

	err = s.executeQuery(ctx, query, func(ctx context.Context, ex exec.QueryExecutor) error {
		ok, err = ex.ScanStructContext(ctx, &pump)
		return errors.WithStack(err)
	})

	return pump, ok, errors.Wrap(err, "failed to retrieve pump")
}

func (s *sqlStorage) InsertPump(ctx context.Context, p *models.Pump) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	query := s.orm.
		Insert("pumps").
		Cols("name", "pin").
		Vals(
			goqu.Vals{p.Name, p.Pin},
		).
		Returning("pump_id", "name", "pin")

	err := s.executeQuery(ctx, query, func(ctx context.Context, ex exec.QueryExecutor) error {
		_, err := ex.ScanStructContext(ctx, p)
		return errors.WithStack(err)
	})

	return errors.Wrap(err, "failed to create pump")
}

func (s *sqlStorage) SelectPumpByPin(ctx context.Context, pin pio.Pin) (models.Pump, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	query := s.orm.
		From("pumps").
		Select("pump_id", "name", "pin").
		Where(goqu.Ex{
			"pin": pin,
		}).
		Limit(1)

	var pump models.Pump

	err := s.executeQuery(ctx, query, func(ctx context.Context, ex exec.QueryExecutor) error {
		_, err := ex.ScanStructContext(ctx, &pump)
		return errors.WithStack(err)
	})

	return pump, errors.Wrap(err, "failed to retrieve pump with pin")
}

func (s *sqlStorage) StartCycle(ctx context.Context, cycle *models.Cycle) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	query := s.orm.Insert("cycles").
		Cols("pump_id", "start_time", "end_time", "notified").
		Vals(
			goqu.Vals{cycle.PumpId, cycle.StartTime, cycle.EndTime, cycle.Notified},
		).
		Returning("cycle_id", "pump_id", "start_time", "end_time", "notified")

	return errors.Wrap(s.executeQuery(ctx, query, func(ctx context.Context, ex exec.QueryExecutor) error {
		_, err := ex.ScanStructContext(ctx, cycle)
		return errors.WithStack(err)
	}), "failed to start cycle")
}

func (s *sqlStorage) FinishCycle(ctx context.Context, cycleId int, timestamp time.Time) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	query := s.orm.Update("cycles").Set(goqu.Record{
		"end_time": timestamp.UTC(),
	}).Where(goqu.Ex{
		"cycle_id": cycleId,
	})

	return errors.Wrap(s.executeQuery(ctx, query, func(ctx context.Context, ex exec.QueryExecutor) error {
		_, err := ex.ExecContext(ctx)
		return errors.WithStack(err)
	}), "failed to finish cycle")
}

func (s *sqlStorage) LatestCycleForPump(ctx context.Context, pumpId int) (models.Cycle, bool, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	query := s.orm.
		From("cycles").
		Select(
			"cycles.cycle_id",
			"cycles.pump_id",
			"cycles.start_time",
			"cycles.end_time",
			"cycles.notified",
		).
		Where(goqu.Ex{
			"cycles.pump_id": pumpId,
		}).
		Order(
			goqu.T("cycles").Col("start_time").Desc(),
		).
		Limit(1)

	var cycle models.Cycle
	var valid bool
	err := s.executeQuery(ctx, query, func(ctx context.Context, ex exec.QueryExecutor) error {
		result, err := ex.ScanStructContext(ctx, &cycle)
		valid = result
		return errors.WithStack(err)
	})
	return cycle, valid, err
}

func (s *sqlStorage) ListNonNotifiedCycles(ctx context.Context) ([]models.Cycle, error) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	query := s.orm.
		From("cycles").
		Select(
			"cycles.cycle_id",
			"cycles.pump_id",
			"cycles.start_time",
			"cycles.end_time",
			"cycles.notified",
		).
		Where(
			goqu.Ex{
				"cycles.notified": false,
			},
			goqu.T("cycles").Col("end_time").IsNotNull(),
		).
		Order(
			goqu.T("cycles").Col("start_time").Asc(),
		)

	cycles := make([]models.Cycle, 0)
	err := s.executeQuery(ctx, query, func(ctx context.Context, ex exec.QueryExecutor) error {
		return errors.WithStack(ex.ScanStructsContext(ctx, &cycles))
	})

	return cycles, err
}

func (s *sqlStorage) SetCyclesAsNotified(ctx context.Context, cycleIds ...int) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	query := s.orm.
		Update("cycles").
		Set(goqu.Record{
			"notified": true,
		}).
		Where(goqu.Ex{
			"cycle_id": cycleIds,
		})

	return errors.Wrap(s.executeQuery(ctx, query, func(ctx context.Context, ex exec.QueryExecutor) error {
		_, err := ex.ExecContext(ctx)
		return errors.WithStack(err)
	}), "failed to update cycle(s)")
}

func (s *sqlStorage) executeQuery(ctx context.Context, query queryInterface, inner func(ctx context.Context, ex exec.QueryExecutor) error) error {
	sql, _, _ := query.ToSQL()
	// fmt.Println(sql)
	span, ctx := tracing.StartNamedSpanFromContext(ctx, sql)
	defer span.Finish()

	return inner(ctx, query.Executor())
}
