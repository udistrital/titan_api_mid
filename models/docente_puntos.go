package models

type Docente_puntos struct {
	Id        int    `orm:"column(id);pk;auto"`
	Documento string `orm:"column(num_documento)"`
	Puntos    int    `orm:"column(puntos);null"`
}
type Foo struct {
	Bar string
}
