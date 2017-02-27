package models


type LugarEjecucion struct {
	Id          int     `orm:"column(id);pk"`
	Direccion   string  `orm:"column(direccion)"`
	Sede        string  `orm:"column(sede);null"`
	Dependencia string  `orm:"column(dependencia);null"`
	Ciudad      float64 `orm:"column(ciudad)"`
}
