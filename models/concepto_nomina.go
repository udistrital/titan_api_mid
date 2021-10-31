package models

type ConceptoNomina struct {
	Id                         int
	NombreConcepto             string
	AliasConcepto              string
	NaturalezaConceptoNominaId int
	TipoConceptoNominaId       int
	EstadoConceptoNominaId     int
	Activo                     bool
	FechaCreacion              string
	FechaModificacion          string
}
