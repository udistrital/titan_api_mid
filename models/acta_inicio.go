package models


import (
	"time"

)

type ActaInicio struct {
	Id             int       `orm:"column(id);pk"`
	NumeroContrato *ContratoGeneral    `orm:"rel(one);column(numero_contrato);null"`
	Vigencia       int       `orm:"column(vigencia);null"`
	FechaInicio    time.Time `orm:"column(fecha_inicio);type(date);null"`
	FechaFin       time.Time `orm:"column(fecha_fin);type(date);null"`
	Descripcion    string    `orm:"column(descripcion);null"`
	Usuario        string    `orm:"column(usuario);null"`
}
