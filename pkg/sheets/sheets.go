package sheets

import (
	"context"
	"fmt"
	"github.com/elliotcourant/sump-counter/pkg/models"
	"github.com/elliotcourant/sump-counter/pkg/tracing"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"time"
)

type (
	Configuration struct {
		Log              *logrus.Entry
		Credentials      *google.Credentials
		SpreadsheetId    string
		SheetName        string
		TimezoneName     string
		DateFormat       string
		ValueInputOption string
	}

	GoogleSheets interface {
		AppendCycle(ctx context.Context, pump models.Pump, cycle models.Cycle) error
		AppendCycles(ctx context.Context, cycles ...models.CycleReport) error
	}

	googleSheetsBase struct {
		log      *logrus.Entry
		config   Configuration
		service  *sheets.SpreadsheetsService
		timeZone *time.Location
	}
)

func NewGoogleSheets(configuration Configuration) (GoogleSheets, error) {
	timeZone, err := time.LoadLocation(configuration.TimezoneName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load time zone")
	}

	service, err := sheets.NewService(context.Background(), option.WithCredentials(configuration.Credentials))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create google service")
	}

	googleSheets := &googleSheetsBase{
		config:   configuration,
		service:  sheets.NewSpreadsheetsService(service),
		log:      configuration.Log.WithField("component", "sheets"),
		timeZone: timeZone,
	}

	return googleSheets, nil
}

func (g *googleSheetsBase) AppendCycles(ctx context.Context, cycles ...models.CycleReport) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	values := make([][]interface{}, len(cycles), len(cycles))
	for i, cycle := range cycles {
		values[i] = []interface{}{
			cycle.CycleId,
			cycle.PumpId,
			cycle.Name,
			cycle.StartTime.In(g.timeZone).Format(g.config.DateFormat),
			cycle.EndTime.In(g.timeZone).Format(g.config.DateFormat),
			cycle.EndTime.Truncate(time.Millisecond).Sub(cycle.StartTime.Truncate(time.Millisecond)).String(),
		}
	}

	cellRange := fmt.Sprintf("%s!A:A", g.config.SheetName)
	_, err := g.service.Values.
		Append(g.config.SpreadsheetId, cellRange, &sheets.ValueRange{
			Values: values,
		}).
		Context(ctx).
		ValueInputOption(g.config.ValueInputOption).
		InsertDataOption("INSERT_ROWS").
		Do()

	return errors.Wrap(err, "failed to append cycle(s) to google sheet")
}

func (g *googleSheetsBase) AppendCycle(ctx context.Context, pump models.Pump, cycle models.Cycle) error {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	return g.AppendCycles(ctx, models.CycleReport{
		Cycle: cycle,
		Name:  pump.Name,
	})
}
