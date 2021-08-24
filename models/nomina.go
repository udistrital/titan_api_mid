package models

import "time"

type Nomina struct {
	Id                int
	Descripcion       string
	Activo            bool
	FechaCreacion     time.Time
	FechaModificacion time.Time
	TipoNominaId      *TipoNomina
}
