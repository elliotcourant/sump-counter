package sheets

import (
	"context"
	"github.com/elliotcourant/sump-counter/pkg/models"
	"github.com/elliotcourant/sump-counter/pkg/testutils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var (
	CredPath      = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	SpreadhseetId = os.Getenv("SPREADSHEET_ID")
)

func TestGoogleSheets(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		credFile := option.WithCredentialsFile(CredPath)
		service, err := sheets.NewService(context.Background(), credFile)
		require.NoError(t, err)

		sheetService := sheets.NewSpreadsheetsService(service)
		call := sheetService.Get(SpreadhseetId)
		spreadsheet, err := call.Do()
		require.NoError(t, err)
		require.NotNil(t, spreadsheet)

		cellRange := "Sheet1!A:A"
		result, err := sheetService.Values.Append(SpreadhseetId, cellRange, &sheets.ValueRange{
			MajorDimension: "",
			Range:          "",
			Values: [][]interface{}{
				{
					1,
				},
			},
			ForceSendFields: nil,
			NullFields:      nil,
		}).ValueInputOption("RAW").Do()
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("creds", func(t *testing.T) {
		json, err := ioutil.ReadFile(CredPath)
		require.NoError(t, err)

		creds, err := google.CredentialsFromJSON(context.Background(), json, "https://www.googleapis.com/auth/spreadsheets")
		require.NoError(t, err)
		require.NotNil(t, creds)
	})
}

func TestGoogleSheetsBase_AppendCycle(t *testing.T) {
	json, err := ioutil.ReadFile(CredPath)
	require.NoError(t, err)

	creds, err := google.CredentialsFromJSON(context.Background(), json, "https://www.googleapis.com/auth/spreadsheets")
	require.NoError(t, err)
	require.NotNil(t, creds)

	logger := logrus.New()

	config := Configuration{
		Log:              logrus.NewEntry(logger),
		Credentials:      creds,
		SpreadsheetId:    SpreadhseetId,
		SheetName:        "Sheet1",
		TimezoneName:     "America/Chicago",
		DateFormat:       "1/2/2006 15:04:05",
		ValueInputOption: "USER_ENTERED",
	}

	gSheets, err := NewGoogleSheets(config)
	require.NoError(t, err)
	require.NotNil(t, gSheets)

	t.Run("simple", func(t *testing.T) {
		testutils.NewTestWithContext(t, func(t *testing.T, ctx context.Context) {
			end := time.Now().Add(1 * time.Hour).UTC()
			cycle := models.Cycle{
				CycleId:   1243,
				PumpId:    543,
				StartTime: time.Now().UTC(),
				EndTime:   &end,
				Notified:  false,
			}
			pump := models.Pump{
				PumpId: 6543,
				Name:   "Ground Water",
				Pin:    13,
			}

			err := gSheets.AppendCycle(ctx, pump, cycle)
			assert.NoError(t, err)
		})
	})
}
