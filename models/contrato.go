package models

import "time"

type Contrato struct {
	Id                int
	NumeroContrato    string
	Vigencia          int
	NombreCompleto    string
	Documento         string
	PersonaId         int
	TipoNominaId      int
	FechaInicio       time.Time
	FechaFin          time.Time
	ValorContrato     float64
	DependenciaId     int
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
}
