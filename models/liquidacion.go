package models

import (
	"time"
)

type Liquidacion struct {
	Id                int       `orm:"column(id);pk"`
	NombreLiquidacion string    `orm:"column(nombre_liquidacion)"`
	Nomina            *Nomina   `orm:"column(nomina);rel(fk)"`
	IdUsuario         int64     `orm:"column(id_usuario)"`
	EstadoLiquidacion string    `orm:"column(estado_liquidacion)"`
	FechaLiquidacion  time.Time `orm:"column(fecha_liquidacion);type(date)"`
	FechaInicio       time.Time `orm:"column(fecha_inicio);type(date)"`
	FechaFin          time.Time `orm:"column(fecha_fin);type(date)"`
}

type DatosLiquidacion struct {
	Preliquidacion *Preliquidacion
	Personas       []int
}
