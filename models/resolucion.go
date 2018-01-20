package models

import(
	"time"
)

type Resolucion struct {
	Id                            int       `orm:"column(id);pk"`
	NumResolucionPension          string    `orm:"column(num_resolucion_pension)"`
	FechaEmisionResolucionPension time.Time `orm:"column(fecha_emision_resolucion_pension);type(date)"`
}
