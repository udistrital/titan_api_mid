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
	Vacaciones        float64
	DependenciaId     int
	ProyectoId        int
	Cdp               int
	Rp                int
	Unico             bool
	Completo          bool
	Activo            bool
	NumeroSemanas     int
	FechaCreacion     string
	FechaModificacion string
	ResolucionId      int
	Desagregado       *Desagregado
}

type ContratoOld struct {
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
	Vacaciones        float64
	DependenciaId     int
	ProyectoId        int
	Cdp               int
	Rp                int
	Unico             bool
	Completo          bool
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
}
