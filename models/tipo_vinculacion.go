package models



type TipoVinculacion struct {
	Id     int    `orm:"column(id);pk"`
	Nombre string `orm:"column(nombre);null"`
}
