package models


import (
	"time"

)

type FuncionarioCargo struct {
  Id                      int                `orm:"column(id)"`
  Asignacion_basica         int             `orm:"column(asignacion_basica)"`
  FechaInicio    time.Time `orm:"column(emp_desde);type(date);null"`
	FechaFin       time.Time `orm:"column(emp_hasta);type(date);null"`

}
