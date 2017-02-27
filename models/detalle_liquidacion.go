package models

type DetalleLiquidacion struct {
	Id             int              `orm:"column(id);pk"`
	ValorCalculado int64            `orm:"column(valor_calculado)"`
	EstadoConcepto string           `orm:"column(estado_concepto)"`
	Liquidacion    *Liquidacion     `orm:"column(liquidacion);rel(fk)"`
	Persona        int              `orm:"column(persona)"`
	Concepto       *Concepto        `orm:"column(concepto);rel(fk)"`
	NumeroContrato *ContratoGeneral `orm:"column(numero_contrato);rel(fk)"`
}
