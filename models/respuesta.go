package models

type Respuesta struct {
	Id           int
	NumeroContrato string
	VigenciaContrato string
	NumDocumento float64
	Saldo_RP     float64
	TotalDevengos int
	TotalDescuentos  int
	TotalAPagar    int
	Conceptos    *[]ConceptosResumen
}
type FormatoPreliqu struct {
	//Contrato   *ContratoGeneral
	Respuesta *Respuesta
}

type TotalPersona struct {
	Id        int
	Total     string
}
