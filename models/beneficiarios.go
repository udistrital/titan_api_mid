package models

import(
	"time"
)


type Beneficiarios struct {
	Id                    int                    `orm:"column(id);pk"`
	InformacionPensionado int                    `orm:"column(informacion_pensionado);null"`
	InformacionProveedor  int  `orm:"column(informacion_proveedor);"`
	FechaNacBeneficiario  time.Time              `orm:"column(fecha_nac_beneficiario);type(date);null"`
	Tutor                 int                    `orm:"column(tutor);null"`
	SubFamiliar           string                 `orm:"column(sub_familiar);null"`
	CategoriaBeneficiario int `orm:"column(categoria_beneficiario);"`
	SubEstudios           string                 `orm:"column(aux_estudio);null"`
	Estado								string									`orm:"column(estado);null"`
}
