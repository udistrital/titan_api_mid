package models


import (
	"time"

)
type PruebaGo struct {
	InformacionCargo []FuncionarioCargo
  Reglas string
  FechaPreliquidacion time.Time
  Valor_correcto_salario string
	Valor_correcto_Reteica				string
	Valor_correcto_EstampillaUD		string
	Valor_correcto_ProCultura			string
	Valor_correcto_AdultoMayor	  string										
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
	ValorReteica				string										`orm:"column(valor_reteica)"`
	ValorEstampillaUD		string										`orm:"column(valor_estampillaud)"`
	ValorProCultura			string										`orm:"column(valor_procultura)"`
	ValorAdultoMayor	  string										`orm:"column(valor_adultomayor)"`
}
