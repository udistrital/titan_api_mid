package models

type ContratoNuevo struct {
	NumeroContratoReduccion string
	ValorContratoReduccion  float64
	NumeroResolucion        string
	IdResolucion            int
	DesagregadoReduccion    *Desagregado
}
