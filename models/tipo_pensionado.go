package models

type TipoPensionado struct {
	Id                int    `orm:"column(id);pk"`
	DesTipoPensionado string `orm:"column(des_tipo_pensionado)"`
}
