package models



type Concepto struct {
	Id             int           `orm:"column(id);pk"`
	NombreConcepto string        `orm:"column(nombre_concepto)"`
	Naturaleza    string    `orm:"column(naturaleza);null"`
	AliasConcepto string        `orm:"column(alias_concepto)"`
}
