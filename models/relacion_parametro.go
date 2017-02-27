package models

import (
	
	"time"


)

type RelacionParametro struct {
	Id             int       `orm:"column(id_rel_parametro);pk"`
	Descripcion    string    `orm:"column(descripcion)"`
	EstadoRegistro bool      `orm:"column(estado_registro)"`
	FechaRegistro  time.Time `orm:"column(fecha_registro);type(date)"`
}
