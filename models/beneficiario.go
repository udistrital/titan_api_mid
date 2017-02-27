package models

import(
	"time"
)


type Beneficiario struct {
	Id                    int                    `orm:"column(id);pk"`
	InformacionPensionado int                    `orm:"column(informacion_pensionado);null"`
	InformacionProveedor  *InformacionProveedor  `orm:"column(informacion_proveedor);rel(fk)"`
	FechaNacBeneficiario  time.Time              `orm:"column(fecha_nac_beneficiario);type(date);null"`
	Tutor                 int                    `orm:"column(tutor);null"`
	SubFamiliar           string                 `orm:"column(sub_familiar);null"`
	CategoriaBeneficiario *CategoriaBeneficiario `orm:"column(categoria_beneficiario);rel(fk)"`
	SubEstudios           string                 `orm:"column(sub_estudios);null"`
}
