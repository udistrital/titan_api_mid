package models



type Nomina struct {
	Id 							int               `orm:"auto;column(id);pk"`
	Descripcion     string           `orm:"column(descripcion)"`
	TipoNomina      *TipoNomina      `orm:"column(tipo_nomina);rel(fk)"`
	Activo          bool             `orm:"column(activo)"`
}
