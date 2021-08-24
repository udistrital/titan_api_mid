package models

import "time"

type DetallePreliquidacion struct {
	Id                     int
	ValorCalculado         float64
	NumeroContrato         string
	VigenciaContrato       int
	DiasLiquidados         float64
	TipoPreliquidacionId   *TipoPreliquidacion
	PreliquidacionId       *Preliquidacion
	ConceptoNominaId       *ConceptoNomina
	EstadoDisponibilidadId *EstadoDisponibilidad
	PersonaId              int
	DependenciaId          int
	FechaCreacion          time.Time
	FechaModificacion      time.Time
	Activo                 bool
	NombreCompleto         string
	Documento              string
}
