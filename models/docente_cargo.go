package models

import "time"

type DocenteCargo struct {
	Id                int       `orm:"column(id)"`
	Asignacion_basica int       `orm:"column(asignacion_basica)"`
	FechaInicio       time.Time `orm:"column(emp_desde);type(date);null"`
	FechaFin          time.Time `orm:"column(emp_hasta);type(date);null"`
	Puntos            float64   `orm:"column(puntos)"`
	Regimen           string    `orm:"column(regimen)"`
}

type Puntos struct {
	Puntos_salariales   float64
	Puntos_bonificacion float64
}

// last inserted Id on success.
