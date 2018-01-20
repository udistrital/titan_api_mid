package models


type Sustituto struct {
	Id                              int           `orm:"column(id);pk"`
	Proveedor											int								`orm:"column(informacion_proveedor)"`
	Beneficiario           				int `orm:"column(beneficiario)"`
	Porcentaje										int					`orm:"column(porcentaje)"`
	Estado                          string        `orm:"column(estado);null"`
	NumeroContrato								string						`orm:"column(numero_contrato);null"`
	Tutor	 												int								`orm:"column(tutor)"`
}
