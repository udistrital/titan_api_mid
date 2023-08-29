package models

import "time"

type Reduccion struct {
	Vigencia            int
	Documento           string
	FechaReduccion      time.Time
	NivelAcademico      string
	Semanas             int
	ContratosOriginales []*ContratoReducir
	ContratoNuevo       *ContratoNuevo
}
