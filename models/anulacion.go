package models

import "time"

type Anulacion struct {
	NumeroContrato string
	Vigencia       int
	Documento      string
	FechaAnulacion time.Time
	Desagregado    *Desagregado
}
