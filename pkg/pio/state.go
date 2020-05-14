package pio

import (
	"github.com/stianeikeland/go-rpio"
	"math"
)

//go:generate stringer -type State -output state.strings.go
const (
	Low     State = State(rpio.Low)
	High    State = State(rpio.High)
	Unknown State = math.MaxUint8
)
