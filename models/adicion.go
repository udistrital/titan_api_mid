package models

import "time"

type Adicion struct {
	NumeroContrato string
	Vigencia       int
	Documento      string
	Dedicacion     string
	Categoria      string
	NumeroSemanas  int
	HorasSemanales int
	NivelAcademico string
	FechaInicio    time.Time
	RpActual       int
	RpNuevo        int
	Cdp            int
}
