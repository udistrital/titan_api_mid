package models

import "time"

type NaturalezaConceptoNomina struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	NumeroOrden       float64
	FechaCreacion     time.Time
	FechaModificacion time.Time
}
