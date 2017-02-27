package models


import (
	"time"

)

type ConceptoPorPersona struct {
	ValorNovedad  float64     `orm:"column(valor_novedad)"`
	EstadoNovedad string     `orm:"column(estado_novedad)"`
	FechaDesde    time.Time `orm:"column(fecha_desde);type(date)"`
	FechaHasta    time.Time `orm:"column(fecha_hasta);type(date)"`
	NumCuotas     int64     `orm:"column(num_cuotas)"`
	Persona       *InformacionProveedor       `orm:"column(persona)"`
	Concepto      *Concepto `orm:"column(concepto);rel(fk)"`
	Nomina        *Nomina       `orm:"column(nomina)"`
	Id            int       `orm:"auto;column(id);pk"`
	Tipo          string    `orm:"column(tipo);null"`
}
