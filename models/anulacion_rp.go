package models

import "time"

type AnulacionRp struct {
	NumeroContrato    string
	Vigencia          int
	Documento         string
	Rp                int
	FechaAnulacion    time.Time
	ContratosAnulados []*ContratoAnular
}
