package models
import (

	"time"


)

type Parametros struct {
	Id                int                `orm:"column(id_parametro);pk"`
	Descripcion       string             `orm:"column(descripcion)"`
	CodigoContraloria string             `orm:"column(codigo_contraloria);null"`
	RelParametro      *RelacionParametro `orm:"column(rel_parametro);rel(fk)"`
	EstadoRegistro    bool               `orm:"column(estado_registro)"`
	FechaRegistro     time.Time          `orm:"column(fecha_registro);type(date)"`
}
