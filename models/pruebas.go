package models


import (
	"time"

)
type PruebaGo struct {
	InformacionCargo []FuncionarioCargo
  Reglas string
  FechaPreliquidacion time.Time
  Valor_correcto_salario string
  IdProveedor int
  Dias_laborados float64
  Periodo string
  EsAnual int
  PorcentajePT int
  TipoNomina int
}
