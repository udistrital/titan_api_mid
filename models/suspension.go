package models

import "time"

type Suspension struct {
	NumeroContrato string
	Vigencia       int
	FechaInicio    time.Time
	FechaFin       time.Time
}
