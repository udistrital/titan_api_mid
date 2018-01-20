package models



type Pais struct {
	Id            int    `orm:"column(id_pais);pk"`
	NombrePais    string `orm:"column(nombre_pais)"`
	Abreviatura   string `orm:"column(abreviatura);null"`
	Estado        string `orm:"column(estado)"`
	CapitalPais   string `orm:"column(capital_pais);null"`
	ProvinciaPais string `orm:"column(provincia_pais);null"`
	AreaPais      int    `orm:"column(area_pais);null"`
	PoblacionPais int    `orm:"column(poblacion_pais);null"`
	CodigoPais    int    `orm:"column(codigo_pais);null"`
}
