package models

import "time"

type Suspension struct {
	NumeroContrato string
	FechaInicio    time.Time
	FechaFin       time.Time
}
