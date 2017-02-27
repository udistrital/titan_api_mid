package models


type Sustituto struct {
	Id                              int           `orm:"column(id);pk"`
	ParentescoInformacionPensionado int           `orm:"column(parentesco_informacion_pensionado)"`
	BeneficiarioSustituto           *Beneficiario `orm:"column(beneficiario_sustituto);rel(fk)"`
	Estado                          string        `orm:"column(estado);null"`
}
