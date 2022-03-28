package models

import "time"

type Novedad struct {
	Id                int
	ContratoId        *Contrato
	ConceptoNominaId  *ConceptoNomina
	Valor             float64
	Cuotas            int
	FechaInicio       time.Time
	FechaFin          time.Time
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
}

type MensajeNovedad struct {
	Mensaje string
	Estado  int //1: Las cuotas Superan los meses //2: Se supera el 50% los descuentos  //3: No hay problema
}
