package models

import (
	"time"
)

type Preliquidacion struct {
	Nomina               *Nomina               `orm:"column(nomina);rel(fk)"`
	Id                   int                   `orm:"column(id);pk"`
	Descripcion          string                `orm:"column(descripcion);null"`
	Mes                  int                   `orm:"column(mes)"`
	Ano                  int                   `orm:"column(ano)"`
	FechaRegistro        time.Time             `orm:"column(fecha_registro);type(timestamp with time zone)"`
	EstadoPreliquidacion *EstadoPreliquidacion `orm:"column(estado_preliquidacion);rel(fk)"`
	Definitiva           bool
}
