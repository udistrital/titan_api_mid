package models

type Respuesta struct {
	Id           int
	NumeroContrato string
	VigenciaContrato string
	NumDocumento float64
	Saldo_RP     float64
	Valor_bruto  string
	Valor_neto   string
	Conceptos    *[]ConceptosResumen
}
type FormatoPreliqu struct {
	//Contrato   *ContratoGeneral
	Respuesta *Respuesta
}
