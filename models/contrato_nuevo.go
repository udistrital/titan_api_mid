package models

type ContratoNuevo struct {
	NumeroContratoReduccion string
	ValorContratoReduccion  float64
	DesagregadoReduccion    *Desagregado
}
