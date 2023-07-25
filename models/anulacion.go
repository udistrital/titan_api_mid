package models

import "time"

type Anulacion struct {
	NumeroContrato string
	Vigencia       int
	Documento      string
	ValorContrato  float64
	NivelAcademico string
	FechaAnulacion time.Time
	Desagregado    *Desagregado
}

type AnulacionOld struct {
	NumeroContrato string
	Vigencia       int
	Documento      string
	FechaAnulacion time.Time
	Rp             int
}
