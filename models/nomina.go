package models



type Nomina struct {
	Id 							int               `orm:"auto;column(id);pk"`
	TipoVinculacion *TipoVinculacion `orm:"column(tipo_vinculacion);rel(fk)"`
	Descripcion     string           `orm:"column(descripcion)"`
	TipoNomina      *TipoNomina      `orm:"column(tipo_nomina);rel(fk)"`
	Activo          bool             `orm:"column(activo)"`
}
