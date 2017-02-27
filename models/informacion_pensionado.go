package models

import (
	"time"
)

type InformacionPensionado struct {
	Id                   int          `orm:"column(id);pk"`
	InformacionProveedor int          `orm:"column(informacion_proveedor);rel(fk)"`
	Estado               string       `orm:"column(estado)"`
	FechaNacEmpleado     time.Time    `orm:"column(fecha_nac_empleado);type(date);null"`
	EstadoCivil          *EstadoCivil `orm:"column(estado_civil);rel(fk)"`
	FechaRetiro          time.Time    `orm:"column(fecha_retiro);type(date);null"`
	PersonaFallecido     string       `orm:"column(persona_fallecido)"`
	PensionadoEnExterior string       `orm:"column(pensionado_en_exterior);null"`
	TipoPensionado       int          `orm:"column(tipo_pensionado)"`
	TipoPension          *TipoPension `orm:"column(tipo_pension);rel(fk)"`
	ValorPensionAsignada int          `orm:"column(valor_pension_asignada);null"`
	FechaPension         time.Time    `orm:"column(fecha_pension);type(date);null"`
	Resolucion           *Resolucion  `orm:"column(resolucion);rel(fk)"`
	EstadoPension        string       `orm:"column(estado_pension);null"`
}
