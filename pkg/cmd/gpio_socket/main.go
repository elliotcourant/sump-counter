package main

import (
	"flag"
	"fmt"
	"github.com/elliotcourant/sump-counter/pkg/ioserver"
	"github.com/elliotcourant/sump-counter/pkg/pio"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"syscall"
	"time"
)

var (
	socketPathInput = flag.String("socket", "/tmp/gpio", "the folder for the unix socket, socket will be named gpio.sock")
	logLevel        = flag.String("log-level", "trace", "log level for output")
)

func main() {
	flag.Parse()

	// sock, err := ioutil.TempDir("", "gpio_socket")
	// if err != nil {
	// 	panic(err)
	// }
	//
	// *socketPathInput = sock

	logger := logrus.New()
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

	io, err := pio.NewIO()
	if err != nil {
		logger.WithError(err).Fatal("failed to open pio")
		return
	} else {
		logger.Tracef("successfully setup gpio")
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	logger.Tracef("users home directory: %s", dir)

	folderPath, socketPath := *socketPathInput, *socketPathInput+"/gpio.sock"

	if err := os.RemoveAll(socketPath); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(folderPath, 777); err != nil {
		panic(err)
	}

	server, err := ioserver.NewIOServer(logger.WithField("component", "gpio-server"), socketPath, io)
	if err != nil {
		logger.WithError(err).Fatal("failed to create io server")
		return
	}

	logger.WithField("socket", socketPath).Info("starting server")
	server.Serve()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT)

	<-quit
	fmt.Println()
	logger.Warnf("quitting")

	server.Close()
	logger.Info("done")
}

func DebugPermissions(log *logrus.Logger, path string) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	log.WithFields(logrus.Fields{
		"path": path,
		"mode": fileInfo.Mode().String(),
	}).Debugf("permissions for path")
}
