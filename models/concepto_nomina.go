package models

import "time"

type ConceptoNomina struct {
	Id                         int
	NombreConcepto             string
	AliasConcepto              string
	NaturalezaConceptoNominaId *NaturalezaConceptoNomina
	TipoConceptoNominaId       *TipoConceptoNomina
	EstadoConceptoNominaId     *EstadoConceptoNomina
	Activo                     bool
	FechaCreacion              time.Time
	FechaModificacion          time.Time
}
