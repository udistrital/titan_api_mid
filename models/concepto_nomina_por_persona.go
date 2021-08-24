package models

import (
	"time"
)

type ConceptoNominaPorPersona struct {
	Id                int             `orm:"auto;column(id);pk"`
	ValorNovedad      float64         `orm:"column(valor_novedad)"`
	NumCuotas         int             `orm:"column(num_cuotas)"`
	FechaDesde        time.Time       `orm:"column(fecha_desde);type(timestamp with time zone);null"`
	FechaHasta        time.Time       `orm:"column(fecha_hasta);type(timestamp with time zone);null"`
	NominaId          *Nomina         `orm:"column(nomina_id);rel(fk)"`
	ConceptoNominaId  *ConceptoNomina `orm:"column(concepto_nomina_id);rel(fk)"`
	NumeroContrato    string          `orm:"column(numero_contrato)"`
	VigenciaContrato  int             `orm:"column(vigencia_contrato)"`
	PersonaId         int             `orm:"column(persona_id)"`
	FechaCreacion     time.Time       `orm:"column(fecha_creacion);type(timestamp with time zone);null"`
	FechaModificacion time.Time       `orm:"column(fecha_modificacion);type(timestamp with time zone);null"`
	Activo            bool            `orm:"column(activo)"`
}
