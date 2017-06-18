package dto

import (
	"encoding/gob"
	"time"
)

type SensorMessage struct {
	Name      string
	value     float64
	Timestamp time.Time
}

func init() {
	gob.Register(SensorMessage)
}
