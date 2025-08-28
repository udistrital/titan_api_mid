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
	PuntoSalarial  float64
	ObjetoNovedad  *ObjetoNovedad
}

type ObjetoNovedad struct {
	SemanasNuevas               int
	TipoResolucion              string
	VinculacionOriginal         string
	VigenciaVinculacionOriginal int
}
