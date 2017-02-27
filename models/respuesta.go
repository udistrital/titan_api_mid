package models

type Respuesta struct {
	Id           int
	Nombre_Cont  string
	NumDocumento float64
	Valor_bruto  string
	Valor_neto   string
	Conceptos    *[]ConceptosResumen
}
type FormatoPreliqu struct {
	//Contrato   *ContratoGeneral
	Respuesta *Respuesta
}
