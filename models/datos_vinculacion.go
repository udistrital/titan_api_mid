package models

type DatosVinculacion struct {
	NumeroContrato string
	Vigencia       int
	Documento      string
	Dedicacion     string
	Categoria      string
	NumeroSemanas  int
	HorasSemanales int
	NivelAcademico string
	Cancelacion    bool
}
