package models



type ParametroEstandar struct {
	Id                   int    `orm:"column(id_parametro);pk"`
	ClaseParametro       string `orm:"column(clase_parametro);null"`
	ValorParametro       string `orm:"column(valor_parametro)"`
	DescripcionParametro string `orm:"column(descripcion_parametro);null"`
}
