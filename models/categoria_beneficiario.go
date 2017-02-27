package models


type CategoriaBeneficiario struct {
	Id                          int    `orm:"column(id);pk"`
	NombreTipoBeneficiario      string `orm:"column(nombre_tipo_beneficiario)"`
	DescripcionTipoBeneficiario string `orm:"column(descripcion_tipo_beneficiario)"`
}
