package models

type ContratoReducir struct {
	NumeroContratoOriginal string
	ValorContratoReducido  float64
	DesagregadoOriginal    *Desagregado
}
