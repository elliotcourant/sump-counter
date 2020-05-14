package pio_mock

import (
	"fmt"
	"github.com/elliotcourant/sump-counter/pkg/pio"
	"sync"
)

var (
	_ pio.IO = &PIOMock{}
)

type (
	PIOMock struct {
		pinLock sync.RWMutex
		pins    map[pio.Pin]pio.State
	}
)

func NewPIOMock() *PIOMock {
	return &PIOMock{
		pinLock: sync.RWMutex{},
		pins:    make(map[pio.Pin]pio.State),
	}
}

func (p *PIOMock) Get(pin pio.Pin) pio.State {
	p.pinLock.RLock()
	defer p.pinLock.RUnlock()
	return p.pins[pin]
}

func (p *PIOMock) Set(pin pio.Pin, state pio.State) {
	fmt.Printf("[PIN: %d] Setting state: %s\n", pin, state)
	p.pinLock.Lock()
	defer p.pinLock.Unlock()
	p.pins[pin] = state
}
