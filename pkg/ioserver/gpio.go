package ioserver

import (
	"context"
	"flag"
	"github.com/elliotcourant/sump-counter/pkg/pio"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

//go:generate protoc -I ../protos/ ../protos/gpio.proto --go_out=plugins=grpc:./
var (
	_ IOServer = &gpioServerBase{}
)

var (
	watchFrequency   = flag.Duration("watch-frequency", 50*time.Millisecond, "the frequency to check pin state")
	afterChangeDelay = flag.Duration("after-change-delay", 1*time.Second, "the delay after a change before another change can be observed")
)

type (
	IOServer interface {
		GPIOServiceServer
		Serve()
		Close() error
	}

	pinWatch struct {
		Pin        pio.Pin
		Result     chan State
		LastState  pio.State
		LastChange time.Time
	}

	gpioServerBase struct {
		log      *logrus.Entry
		listener net.Listener
		server   *grpc.Server
		io       pio.IO

		frequency        time.Duration
		afterChangeDelay time.Duration

		pinLock    sync.RWMutex
		pins       map[pio.Pin]*pinWatch
		stopWorker chan struct{}
	}
)

func NewIOServer(log *logrus.Entry, socketPath string, io pio.IO) (IOServer, error) {
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to listen at unix socket path: %s", socketPath)
	}

	ioServer := &gpioServerBase{
		log:              log,
		listener:         listener,
		server:           grpc.NewServer(),
		io:               io,
		frequency:        *watchFrequency,
		afterChangeDelay: *afterChangeDelay,
		pins:             make(map[pio.Pin]*pinWatch),
		stopWorker:       make(chan struct{}),
	}

	RegisterGPIOServiceServer(ioServer.server, ioServer)

	go ioServer.doWork()

	return ioServer, nil
}

func (g *gpioServerBase) doWork() {
	ticker := time.NewTicker(g.frequency)
	for {
		select {
		case <-g.stopWorker:
			ticker.Stop()
			return
		case <-ticker.C:
			g.checkPins()
		}
	}
}

func (g *gpioServerBase) checkPins() {
	g.pinLock.RLock()
	defer g.pinLock.RUnlock()
	var currentState pio.State
	for pin, watch := range g.pins {
		// If not enough time has passed since the last change was observed then just skip checking this pin.
		if time.Since(watch.LastChange) < g.afterChangeDelay {
			continue
		}

		// If the current state of the pin matches the previously observed state then nothing has changed.
		if currentState = g.io.Get(pin); currentState == watch.LastState {
			continue
		}

		g.log.WithField("pin", pin).Tracef("state changed: %s -> %s", watch.LastState, currentState)

		// If the current state does not match...
		// Update the last observed state.
		watch.LastState = currentState

		// Update the timestamp of the last state change.
		watch.LastChange = time.Now()

		// Push the new state to any watchers.
		switch currentState {
		case pio.Low:
			watch.Result <- State_Low
		case pio.High:
			watch.Result <- State_High
		}
	}
}

func (g *gpioServerBase) addPinToWatchlist(pin pio.Pin) *pinWatch {
	watch := &pinWatch{
		Pin:        pin,
		Result:     make(chan State),
		LastState:  pio.Unknown,
		LastChange: time.Now().Add(-1 * g.afterChangeDelay), // Make sure that we get the initial state right away.
	}

	g.pinLock.Lock()
	defer g.pinLock.Unlock()
	g.pins[pin] = watch

	g.log.WithField("pin", pin).Tracef("added pin to watch list")

	return watch
}

func (g *gpioServerBase) removePinFromWatchlist(pin pio.Pin) {
	g.pinLock.Lock()
	defer g.pinLock.Unlock()
	delete(g.pins, pin)

	g.log.WithField("pin", pin).Tracef("removed pin from watch list")
}

func (g *gpioServerBase) WatchPinState(pinRequest *Pin, server GPIOService_WatchPinStateServer) error {
	pin, ctx := pinRequest.IOPin(), server.Context()
	watch := g.addPinToWatchlist(pin)

	for {
		select {
		case <-ctx.Done():
			g.removePinFromWatchlist(pin)
			if err := ctx.Err(); err != context.Canceled {
				return err
			}
			return nil
		case state := <-watch.Result:
			if err := server.Send(&PinState{
				State: state,
			}); err != nil {
				return errors.Wrap(err, "failed to send new pin state")
			}
		}
	}
}

func (g *gpioServerBase) UpdatePinState(ctx context.Context, state *UpdatePinStateRequest) (*PinState, error) {
	g.log.
		WithField("pin", state.Pin.IOPin()).
		WithField("state", state.State.IOState()).
		Tracef("updating pin state")

	g.io.Set(state.Pin.IOPin(), state.State.IOState())

	return &PinState{
		State: state.State,
	}, nil
}

func (g *gpioServerBase) GetPinState(ctx context.Context, pinRequest *Pin) (*PinState, error) {
	g.log.
		WithField("pin", pinRequest.IOPin()).
		Tracef("retrieving pin state")

	state := g.io.Get(pinRequest.IOPin())
	switch state {
	case pio.Low:
		return &PinState{State: State_Low}, nil
	case pio.High:
		return &PinState{State: State_High}, nil
	}

	return nil, errors.Errorf("invalid pin state retrieved: %s", state)
}

func (g *gpioServerBase) Serve() {
	go func() {
		_ = g.server.Serve(g.listener)
	}()
}

func (g *gpioServerBase) Close() error {
	g.log.Info("closing gpio server")
	g.stopWorker <- struct{}{}
	g.server.Stop()
	return g.listener.Close()
}

func (m *Pin) IOPin() pio.Pin {
	return pio.Pin(m.Pin)
}

func (x State) IOState() pio.State {
	switch x {
	case State_Low:
		return pio.Low
	case State_High:
		return pio.High
	default:
		return pio.Unknown
	}
}
