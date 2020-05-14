package service

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/elliotcourant/sump-counter/pkg/ioserver"
	"github.com/elliotcourant/sump-counter/pkg/models"
	"github.com/elliotcourant/sump-counter/pkg/sheets"
	"github.com/elliotcourant/sump-counter/pkg/storage"
	"github.com/elliotcourant/sump-counter/pkg/tracing"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Service struct {
	log           *logrus.Entry
	storage       storage.Storage
	sheets        sheets.GoogleSheets
	gpio          ioserver.GPIOServiceClient
	configuration Configuration
	quit          chan struct{}
}

func NewService(
	log *logrus.Entry,
	storage storage.Storage,
	sheets sheets.GoogleSheets,
	gpio ioserver.GPIOServiceClient,
	configuration Configuration,
) *Service {
	service := &Service{
		log:           log,
		storage:       storage,
		sheets:        sheets,
		gpio:          gpio,
		configuration: configuration,
		quit:          make(chan struct{}),
	}

	go service.doWork()

	return service
}

type pumpManagement struct {
	Pump       models.Pump
	Context    context.Context
	ctxCancel  context.CancelFunc
	Watch      ioserver.GPIOService_WatchPinStateClient
	cancelChan chan struct{}
	waitGroup  *sync.WaitGroup
}

func (p *pumpManagement) Cancel() {
	p.ctxCancel()
	p.cancelChan <- struct{}{}
}

func (s *Service) waitToDie() {
	<-s.quit
	return
}

func (s *Service) doWork() {
	s.log.Info("starting worker")

	var pumps []models.Pump
	for {
		ps, err := s.storage.GetAllPumps(context.Background())
		if err != nil {
			s.log.WithError(err).Errorf("failed to retrieve all pumps, waiting to die")
			s.waitToDie()
			return
		}

		if len(ps) == 0 {
			s.log.Infof("no pumps were found, checking again later")
			time.Sleep(5 * time.Second)
			continue
		}

		s.log.Infof("found %d pump(s), starting work", len(ps))
		pumps = ps
		break
	}

	management := map[int]*pumpManagement{}

	waitGroup := &sync.WaitGroup{}

	for _, pump := range pumps {
		log := s.log.WithFields(logrus.Fields{
			"pump":   pump.Name,
			"pumpId": pump.PumpId,
			"pin":    pump.Pin,
		})

		ctx, cancel := context.WithCancel(context.Background())
		watcher, err := s.gpio.WatchPinState(ctx, &ioserver.Pin{
			Pin: int32(pump.Pin),
		})
		if err != nil {
			log.WithError(err).Errorf("cannot watch pin")
			cancel()
			continue
		}

		waitGroup.Add(1)

		manager := &pumpManagement{
			Pump:       pump,
			Context:    ctx,
			ctxCancel:  cancel,
			Watch:      watcher,
			cancelChan: make(chan struct{}),
			waitGroup:  waitGroup,
		}

		management[pump.PumpId] = manager

		go s.monitorPump(manager)
	}

	s.log.Infof("watching %d pump(s)", len(management))

	s.googleSheetsWorker()

	s.log.Infof("stopping pump monitoring")

	for _, pump := range management {
		s.log.WithFields(logrus.Fields{
			"pump":   pump.Pump.Name,
			"pumpId": pump.Pump.PumpId,
			"pin":    pump.Pump.Pin,
		}).Debugf("canceling pump")
		pump.Cancel()
	}

	s.log.Infof("waiting for pump to finish")
	waitGroup.Wait()
	s.log.Infof("pump done working")
}

func (s *Service) googleSheetsWorker() {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-s.quit:
			ticker.Stop()
			return
		case <-ticker.C:
			break
		}

		if err := s.pushToGoogleSheets(); err != nil {
			s.log.WithError(err).Errorf("failed to push cycles to google sheets")
		}
	}
}

func (s *Service) pushToGoogleSheets() error {
	span, ctx := tracing.StartSpanFromContext(context.Background())
	defer span.Finish()

	cycles, err := s.storage.GetCyclesToNotify(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve cycles to be notified")
	}

	if len(cycles) == 0 {
		return nil
	}

	pumps, err := s.storage.GetAllPumps(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve pumps")
	}

	var cycleReports []models.CycleReport
	linq.From(cycles).JoinT(
		linq.From(pumps),
		func(cycle models.Cycle) int {
			return cycle.PumpId
		},
		func(pump models.Pump) int {
			return pump.PumpId
		},
		func(cycle models.Cycle, pump models.Pump) models.CycleReport {
			return models.CycleReport{
				Cycle: cycle,
				Name:  pump.Name,
			}
		}).
		ToSlice(&cycleReports)

	s.log.Infof("sending %d cycle(s) to google sheets", len(cycleReports))

	if err := s.sheets.AppendCycles(ctx, cycleReports...); err != nil {
		return errors.Wrap(err, "failed to push cycle reports to google sheets")
	}

	var cycleIds []int
	linq.From(cycles).
		SelectT(func(cycle models.Cycle) int {
			return cycle.CycleId
		}).
		ToSlice(&cycleIds)

	return s.storage.SetCyclesNotified(ctx, cycleIds)
}

func (s *Service) monitorPump(pump *pumpManagement) {
	defer pump.waitGroup.Done()
	log := s.log.WithFields(logrus.Fields{
		"pump":   pump.Pump.Name,
		"pumpId": pump.Pump.PumpId,
		"pin":    pump.Pump.Pin,
	})

	lastState, err := s.storage.GetLastPumpState(context.Background(), pump.Pump.PumpId)
	if err != nil {
		log.WithError(err).Errorf("failed to retrieve previous state of pump")
	}
	timeStamp := time.Now()

	for {
		select {
		case <-pump.cancelChan:
			return
		default:
		}

		pinState, err := pump.Watch.Recv()
		if err != nil {
			if err != context.Canceled {
				log.WithError(err).Errorf("failed to receive pump state")
			}
			continue // TODO (elliotcourant) after so many errors would we want to exit?
		}

		state := pinState.State.IOState()

		if state == lastState {
			log.Debugf("observed duplicate state")
			continue
		}

		log.Debugf("state changed %s -> %s, time as %s: %s", lastState, state, lastState, time.Since(timeStamp))

		if err = s.storage.PersistPumpState(context.Background(), pump.Pump.PumpId, state, time.Now().UTC()); err != nil {
			log.WithError(err).Errorf("failed to persist state change")
			continue
		}

		lastState = state
		timeStamp = time.Now()
	}
}

func (s *Service) Close() {
	s.quit <- struct{}{}
}
