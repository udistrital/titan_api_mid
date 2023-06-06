package models

type Predicado struct {
	Id            int            `orm:"column(id);pk;auto"`
	Nombre        string         `orm:"column(nombre)"`
	Descripcion   string         `orm:"column(descripcion)"`
	Dominio       *Dominio       `orm:"column(dominio)"`
	TipoPredicado *TipoPredicado `orm:"column(tipo_predicado)"`
}
