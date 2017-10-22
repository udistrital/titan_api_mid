package models



type DetallePreliquidacion struct {
	Id                 int                   `orm:"column(id);pk"`
	ValorCalculado     float64               `orm:"column(valor_calculado)"`
	NumeroContrato     string                `orm:"column(numero_contrato);null"`
	VigenciaContrato   int                   `orm:"column(vigencia_contrato);null"`
	DiasLiquidados     float64               `orm:"column(dias_liquidados);null"`
	TipoPreliquidacion *TipoPreliquidacion   `orm:"column(tipo_preliquidacion);rel(fk)"`
	Preliquidacion     *Preliquidacion       `orm:"column(preliquidacion);rel(fk)"`
	Concepto           *ConceptoNomina       `orm:"column(concepto);rel(fk)"`
	EstadoDisponibilidad *EstadoDisponibilidad   `orm:"column(estado_disponibilidad);rel(fk)"`
}
