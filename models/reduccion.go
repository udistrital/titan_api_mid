package models

import "time"

type Reduccion struct {
	NumeroContratoReduccion string
	Vigencia                int
	Documento               string
	ValorContratoReduccion  float64
	FechaReduccion          time.Time
	ContratosOriginales     []*ContratoReducir
	DesagregadoReduccion    *Desagregado
}
