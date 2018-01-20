package models



type ConceptoNomina struct {
	Id                 int                       `orm:"column(id);pk"`
	NombreConcepto     string                    `orm:"column(nombre_concepto)"`
	AliasConcepto      string                    `orm:"column(alias_concepto);null"`
	TipoConcepto       *TipoConceptoNomina       `orm:"column(tipo_concepto);rel(fk)"`
	NaturalezaConcepto *NaturalezaConceptoNomina `orm:"column(naturaleza_concepto);rel(fk)"`
}
