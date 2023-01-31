package models

import "time"

type Reduccion struct {
	NumeroContratoOriginal  string
	NumeroContratoReduccion string
	Vigencia                int
	Documento               string
	ValorContratoReduccion  float64
	FechaReduccion          time.Time
	DesagregadoOriginal     *Desagregado
	DesagregadoReduccion    *Desagregado
}
