package models



type Nomina struct {
	Id          int              `orm:"column(id);pk"`
	Vinculacion *TipoVinculacion `orm:"column(vinculacion);rel(fk)"`
	Nombre      string           `orm:"column(nombre)"`
	Descripcion string           `orm:"column(descripcion)"`
	Estado      string           `orm:"column(estado)"`
	Periodo     string           `orm:"column(periodo);null"`
	TipoNomina  *TipoNomina      `orm:"column(tipo_nomina);rel(fk)"`
}
