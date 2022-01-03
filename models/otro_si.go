package models

import "time"

type OtroSi struct {
	NumeroContrato string
	Vigencia       int
	FechaFin       time.Time
}
