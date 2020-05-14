package main

import (
	"flag"
	"fmt"
	"github.com/elliotcourant/sump-counter/pkg/authentication"
	"github.com/elliotcourant/sump-counter/pkg/service"
	"github.com/elliotcourant/sump-counter/pkg/sheets"
	"github.com/elliotcourant/sump-counter/pkg/storage"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"syscall"
	"time"
)

var (
	logLevel = flag.String("log-level", "trace", "log level for output")
)

var (
	gpioSocketPath = os.Getenv("GPIO_SOCKET_PATH")
)

func main() {
	flag.Parse()

	logger := logrus.New()
	logger.ExitFunc = os.Exit
	logger.Formatter = &logrus.TextFormatter{
		DisableTimestamp:       true,
		TimestampFormat:        time.RFC3339,
		DisableSorting:         false,
		SortingFunc:            sort.Strings,
		DisableLevelTruncation: false,
		QuoteEmptyFields:       false,
		FieldMap:               nil,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return frame.Function, fmt.Sprintf("%s:%d", frame.File, frame.Line)
		},
	}

	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logger.WithError(err).WithField("log-level", logLevel).Warnf("invalid log level provided")
	}

	logger.SetLevel(level)

	googleCredentials, err := authentication.GetGoogleCredentialsFromEnv()
	if err != nil {
		logger.WithError(err).Fatalf("failed to read google credentials")
		return
	}

	configuration, err := service.GetConfigurationFromEnv()
	if err != nil {
		logger.WithError(err).Fatalf("failed to read sump-boi configuration")
		return
	}

	googleSheets, err := sheets.NewGoogleSheets(sheets.Configuration{
		Log:              logrus.NewEntry(logger),
		Credentials:      googleCredentials,
		SpreadsheetId:    configuration.SpreadsheetId,
		SheetName:        configuration.SheetName,
		TimezoneName:     configuration.TimeZone,
		DateFormat:       configuration.DateFormat,
		ValueInputOption: configuration.ValueInputOption,
	})
	if err != nil {
		logger.WithError(err).Fatalf("failed to setup google sheets client")
		return
	}

	logger.Info("running")

	db, err := storage.NewStorage(storage.Configuration{
		Host:     os.Getenv("POSTGRESQL_SERVICE_HOST"),
		Port:     5432,
		User:     os.Getenv("POSTGRES_USER"),
		Database: os.Getenv("POSTGRES_DB"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Log:      logger.WithField("component", "storage"),
		InMemory: false,
	})
	if err != nil {
		logger.WithError(err).Fatalf("failed to setup storage")
		return
	}

	gpio, err := NewGPIO()
	if err != nil {
		logger.WithError(err).Fatalf("failed to setup gpio")
		return
	}

	workService := service.NewService(
		logger.WithField("component", "worker"),
		db,
		googleSheets,
		gpio.Client(),
		*configuration,
	)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-quit

	logger.Info("stopping service")

	workService.Close()

	if err := db.Close(); err != nil {
		logger.WithError(err).Errorf("failed to gracefully close storage")
	}

	if err := gpio.Close(); err != nil {
		logger.WithError(err).Errorf("failed to gracefully close grpc connection")
	}

	logger.Info("done...")
}
