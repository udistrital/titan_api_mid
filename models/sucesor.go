package models

import "time"

type Sucesor struct {
	NumeroContrato string
	Vigencia       int
	Documento      string
	NombreCompleto string
	FechaInicio    time.Time
}
