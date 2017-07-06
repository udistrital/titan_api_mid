package models



type DetallePreliquidacion struct {
	Id             int       `orm:"column(id);pk"`
	ValorCalculado int64     `orm:"column(valor_calculado)"`
	Preliquidacion int       `orm:"column(preliquidacion)"`
	Persona        int       `orm:"column(persona)"`
	Concepto       *Concepto `orm:"column(concepto);rel(fk)"`
	NumeroContrato *ContratoGeneral `orm:"column(numero_contrato);rel(fk)"`
	DiasLiquidados string       `orm:"column(dias_liquidados)"`
	TipoPreliquidacion string   `orm:"column(tipo_preliquidacion)"`
	VigenciaContrato int       `orm:"column(vigencia_contrato)"`
}
