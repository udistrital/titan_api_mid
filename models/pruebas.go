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
	NumDocumento int
  Dias_laborados float64
	Mes int
  Ano int
  EsAnual int
  PorcentajePT int
  TipoNomina int
}


type DatosPruebas struct {
	Id 									int                      `orm:"auto;column(id);pk"`
	NumDocumento				int 											`orm:"column(num_documento)"`
	MesPreliq						int 											`orm:"column(mes_preliquidacion)"`
	AnoPreliq						int 											`orm:"column(ano_preliquidacion)"`
	ValorSalario				string										`orm:"column(valor_salario)"`
}
