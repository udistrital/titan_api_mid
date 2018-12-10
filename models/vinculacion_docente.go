package models


import (

	"database/sql"


)

type VinculacionDocente struct {
	Id                   int                           `orm:"column(id);pk;auto"`
  IdPersona            string                        `orm:"column(id_persona);"`
	NumeroContrato       sql.NullString                `orm:"column(numero_contrato);null"`
	Vigencia             sql.NullInt64                 `orm:"column(vigencia);null"`
	IdProyectoCurricular int16                         `orm:"column(id_proyecto_curricular)"`
	DependenciaAcademica int                           `orm:"column(dependencia_academica)"`
	IdResolucion         *Resolucion									 `orm:"column(id_resolucion)"`
	ValorContrato	float64
}
