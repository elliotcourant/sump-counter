package pio

import (
	"github.com/pkg/errors"
	"github.com/stianeikeland/go-rpio"
	"sync"
)

var (
	_ IO = &ioBase{}
)

var (
	setupLock sync.Mutex
	setup     bool
)

type (
	Pin   rpio.Pin
	State rpio.State

	IO interface {
		Get(pin Pin) State
		Set(pin Pin, state State)
	}

	ioBase struct {
	}
)

func NewIO() (IO, error) {
	return &ioBase{}, setupPi()
}

func setupPi() error {
	setupLock.Lock()
	defer setupLock.Unlock()
	if setup {
		return nil
	}

	if err := rpio.Open(); err != nil {
		return errors.Wrap(err, "failed to open rpio")
	}

	setup = true

	return nil
}

func (i *ioBase) Get(pin Pin) State {
	return State(rpio.ReadPin(rpio.Pin(pin)))
}

func (i *ioBase) Set(pin Pin, state State) {
	rpio.WritePin(rpio.Pin(pin), rpio.State(state))
}
