package models



type Ciudad struct {
	Id             int           `orm:"column(id_ciudad);pk"`
	IdDepartamento *Departamento `orm:"column(id_departamento);rel(fk)"`
	Nombre         string        `orm:"column(nombre)"`
	Abreviatura    string        `orm:"column(abreviatura);null"`
	Descripcion    string        `orm:"column(descripcion);null"`
	Estado         string        `orm:"column(estado)"`
	AbPais         string        `orm:"column(ab_pais);null"`
	Departamento   string        `orm:"column(departamento);null"`
	Poblacion      int           `orm:"column(poblacion);null"`
	Longitud       float64       `orm:"column(longitud);null"`
	Latitud        float64       `orm:"column(latitud);null"`
}
