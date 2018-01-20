package models


type EstadoCivil struct {
	Id                int    `orm:"column(id);pk"`
	NombreEstadoCivil string `orm:"column(nombre_estado_civil);null"`
}
