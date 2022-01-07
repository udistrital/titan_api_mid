package models

import "time"

type OtroSi struct {
	NumeroContrato string
	Vigencia       int
	Documento      string
	FechaFin       time.Time
}
