package models

import (
	"time"
)

type Cycle struct {
	CycleId   int        `db:"cycle_id"`
	PumpId    int        `db:"pump_id"`
	StartTime time.Time  `db:"start_time"`
	EndTime   *time.Time `db:"end_time"`
	Notified  bool       `db:"notified"`
}

type CycleReport struct {
	Cycle
	Name string
}
