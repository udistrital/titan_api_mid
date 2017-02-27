package models



type Departamento struct {
	Id                  int    `orm:"column(id_departamento);pk"`
	IdPais              *Pais  `orm:"column(id_pais);rel(fk)"`
	Nombre              string `orm:"column(nombre)"`
	Abreviatura         string `orm:"column(abreviatura);null"`
	Descripcion         string `orm:"column(descripcion);null"`
	Estado              string `orm:"column(estado)"`
	AbPais              string `orm:"column(ab_pais);null"`
	Poblacion           int    `orm:"column(poblacion);null"`
	Area                int    `orm:"column(area);null"`
	CapitalDepartamento string `orm:"column(capital_departamento);null"`
	DepartamentoCapPais string `orm:"column(departamento_cap_pais);null"`
}
