package models

import (
	"time"
)

type Preliquidacion struct {
	Nombre      string    `orm:"column(nombre)"`
	Nomina      *Nomina   `orm:"column(nomina);rel(fk)"`
	Estado      string    `orm:"column(estado)"`
	Fecha       time.Time `orm:"column(fecha);type(date)"`
	Descripcion string    `orm:"column(descripcion);null"`
	FechaInicio time.Time `orm:"column(fecha_inicio);type(date)"`
	FechaFin    time.Time `orm:"column(fecha_fin);type(date)"`
	Id          int       `orm:"column(id);pk"`
	Liquidada   string    `orm:"column(liquidada)"`
}
