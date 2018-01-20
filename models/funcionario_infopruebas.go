package models

import (
	"time"

)

type FuncionarioInfoPruebas struct {

  InformacionCargo []FuncionarioCargo
  Reglas           string
  FechaPreliquidacion  time.Time
  Conceptos *[]ConceptosResumen
  Valor_correcto_salario string
	IdProveedor     int
	Dias_laborados float64
	Periodo 				string
	EsAnual         int
	PorcentajePT    int
	TipoNomina string
}
