package models

import "time"

type Reinicio struct {
	NumeroContrato string
	Vigencia       int
	Documento      string
	FechaReinicio  time.Time
}
