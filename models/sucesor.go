package models

import "time"

type Sucesor struct {
	NumeroContrato  string
	Vigencia        int
	DocumentoNuevo  string
	NombreCompleto  string
	DocumentoActual string
	FechaInicio     time.Time
}
