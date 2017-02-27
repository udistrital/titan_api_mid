package models



type ArgoOrdenadores struct {
	ORGTIPOORDENADOR  string  `orm:"column(ORG_TIPO_ORDENADOR)"`
	ORGIDENTIFICADOR  float64 `orm:"column(ORG_IDENTIFICADOR)"`
	ORGORDENADORGASTO string  `orm:"column(ORG_ORDENADOR_GASTO)"`
	ORGNOMBRE         string  `orm:"column(ORG_NOMBRE)"`
	ORGIDENTIFICACION float64 `orm:"column(ORG_IDENTIFICACION)"`
	ORGESTADO         string  `orm:"column(ORG_ESTADO)"`
	Id                int     `orm:"column(ORG_IDENTIFICADOR_UNICO);pk"`
}
