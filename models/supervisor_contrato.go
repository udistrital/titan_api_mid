package models


type SupervisorContrato struct {
	Id                    int    `orm:"column(id);pk"`
	Nombre                string `orm:"column(nombre)"`
	Documento             int    `orm:"column(documento)"`
	Cargo                 string `orm:"column(cargo)"`
	SedeSupervisor        string `orm:"column(sede_supervisor);null"`
	DependenciaSupervisor string `orm:"column(dependencia_supervisor);null"`
	Tipo                  int    `orm:"column(tipo);null"`
	Estado                bool   `orm:"column(estado);null"`
	DigitoVerificacion    int    `orm:"column(digito_verificacion);null"`
}
