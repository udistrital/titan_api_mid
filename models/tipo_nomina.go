package models



type TipoNomina struct {
	Id     int    `orm:"column(id);pk"`
	Nombre string `orm:"column(Nombre);null"`
}
