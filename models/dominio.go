package models



type Dominio struct {
	Id          int    `orm:"column(id);pk;auto"`
	Nombre      string `orm:"column(nombre)"`
	Descripcion string `orm:"column(descripcion);null"`
}
