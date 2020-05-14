package models

import (
	"github.com/elliotcourant/sump-counter/pkg/pio"
)

type Pump struct {
	PumpId int     `db:"pump_id"`
	Name   string  `db:"name"`
	Pin    pio.Pin `db:"pin"`
}
