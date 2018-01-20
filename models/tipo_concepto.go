package models


type TipoConcepto struct {
	Id     int    `orm:"column(id);pk"`
	Nombre string `orm:"column(nombre)"`
}
