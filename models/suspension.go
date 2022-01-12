package models

import "time"

type Suspension struct {
	NumeroContrato string
	Vigencia       int
	Documento      string
	FechaInicio    time.Time
	FechaFin       time.Time
}
