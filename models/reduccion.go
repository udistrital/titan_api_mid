package models

import "time"

type Reduccion struct {
	Vigencia            int
	Documento           string
	FechaReduccion      time.Time
	ContratosOriginales []*ContratoReducir
	ContratoNuevo       *ContratoNuevo
}
