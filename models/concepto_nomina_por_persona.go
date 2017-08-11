package models


import (
	"time"

)

type ConceptoNominaPorPersona struct {
	ValorNovedad  float64               `orm:"column(valor_novedad)"`
	NumCuotas     int                   `orm:"column(num_cuotas)"`
	Id            int                   `orm:"column(id);pk"`
	FechaDesde    time.Time             `orm:"column(fecha_desde);type(timestamp with time zone);null"`
	FechaHasta    time.Time             `orm:"column(fecha_hasta);type(timestamp with time zone);null"`
	FechaRegistro time.Time             `orm:"column(fecha_registro);type(timestamp with time zone);null"`
	Persona       *InformacionProveedor `orm:"column(persona);rel(fk)"`
	Concepto      *ConceptoNomina       `orm:"column(concepto);rel(fk)"`
	Nomina        int                   `orm:"column(nomina)"`
	Activo        bool                  `orm:"column(activo)"`
}
