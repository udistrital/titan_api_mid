package models


import (
	"time"

)

type ConceptoNominaPorPersona struct {
	ValorNovedad  float64               `orm:"column(valor_novedad)"`
	NumCuotas     int                   `orm:"column(num_cuotas)"`
	Id 						int                    `orm:"auto;column(id);pk"`
	FechaDesde    time.Time             `orm:"column(fecha_desde);type(timestamp with time zone);null"`
	FechaHasta    time.Time             `orm:"column(fecha_hasta);type(timestamp with time zone);null"`
	FechaRegistro time.Time             `orm:"column(fecha_registro);type(timestamp with time zone);null"`
	Persona       int										 `orm:"column(persona)"`
	Concepto      *ConceptoNomina       `orm:"column(concepto);rel(fk)"`
	Nomina        *Nomina               `orm:"column(nomina);rel(fk)"`
	Activo        bool                  `orm:"column(activo)"`
}
