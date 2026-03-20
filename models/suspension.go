package models

import "time"

type Suspension struct {
	Id             int       `json:"Id"`
	NumeroContrato string    `json:"NumeroContrato"`
	Vigencia       int       `json:"Vigencia"`
	Documento      string    `json:"Documento"`
	FechaInicio    time.Time `json:"FechaInicio,omitempty"`
	FechaFin       time.Time `json:"FechaFin,omitempty"`
}
