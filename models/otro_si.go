package models

import "time"

type OtroSi struct {
	NumeroContrato string
	Vigencia       int
	Documento      string
	Valor          float64
	Rp             int
	Cdp            int
	FechaFin       time.Time
}
