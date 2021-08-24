package models

import "time"

type TipoNomina struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	NumeroOrden       float64
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
}
