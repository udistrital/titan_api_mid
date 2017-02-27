package models

type TipoContrato struct {
	Id           int    `orm:"column(id);pk"`
	TipoContrato string `orm:"column(tipo_contrato);null"`
}
