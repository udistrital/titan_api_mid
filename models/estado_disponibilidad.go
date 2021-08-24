package models

import "time"

type EstadoDisponibilidad struct {
	Id                int
	Nombre            string
	Descripcion       string
	CodigoAbreviacion string
	Activo            bool
	NumeroOrden       float64
	FechaCreacion     time.Time
	FechaModificacion time.Time
}
