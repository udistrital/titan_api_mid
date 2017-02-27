package models

type TipoPension struct {
	Id            int    `orm:"column(id);pk"`
	NombrePension string `orm:"column(nombre_pension)"`
}
