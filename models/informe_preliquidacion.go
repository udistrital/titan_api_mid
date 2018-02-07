package models

type InformePreliquidacion struct {
	//IdPersona      int     `orm:"column(id)"`
	NombreCompleto  string
	Documento 			string
	NumeroContrato  string
	Vigencia       int
	Conceptos      []ConceptosInforme
	Disponibilidad int
}


type ConceptosInforme struct {
	Id         int    `orm:"column(id)"`
	Nombre     string `orm:"column(nombre)"`
	Naturaleza string `orm:"column(naturaleza)"`
	Valor      string `orm:"column(valor)"`
	TipoPreliquidacion string `orm:"column(tipo)"`
	EstadoDisponibilidad int `orm:"column(id_disp)"`
}
