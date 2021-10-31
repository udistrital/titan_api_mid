package models

import "time"

type Novedad struct {
	Id                int
	ContratoId        *Contrato
	ConceptoNominaId  *ConceptoNomina
	Valor             float64
	Cuotas            int
	FechaInicio       time.Time
	FechaFin          time.Time
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
}
