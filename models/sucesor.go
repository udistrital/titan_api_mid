package models

import "time"

type Sucesor struct {
	NumeroContrato string
	Documento      string
	NombreCompleto string
	FechaInicio    time.Time
}
