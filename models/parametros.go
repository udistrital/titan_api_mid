package models

type Periodo struct {
	Id                int    `json:"Id"`
	Nombre            string `json:"Nombre"`
	Descripcion       string `json:"Descripcion"`
	Year              int    `json:"Year"`
	Ciclo             string `json:"Ciclo"`
	CodigoAbreviacion string `json:"CodigoAbreviacion"`
	Activo            bool   `json:"Activo"`
	AplicacionId      int    `json:"AplicacionId"`
	InicioVigencia    string `json:"InicioVigencia"`
	FinVigencia       string `json:"FinVigencia"`
	FechaCreacion     string `json:"FechaCreacion"`
	FechaModificacion string `json:"FechaModificacion"`
}

type Parametro struct {
	Id                int    `json:"Id"`
	Nombre            string `json:"Nombre"`
	Descripcion       string `json:"Descripcion"`
	CodigoAbreviacion string `json:"CodigoAbreviacion"`
	Activo            bool   `json:"Activo"`
	NumeroOrden       int    `json:"NumeroOrden"`
	FechaCreacion     string `json:"FechaCreacion"`
	FechaModificacion string `json:"FechaModificacion"`
	TipoParametroId   int    `json:"TipoParametroId"`
	ParametroPadreId  int    `json:"ParametroPadreId"`
}

type ParametroPeriodo struct {
	Id                int    `json:"Id"`
	ParametroId       int    `json:"ParametroId"`
	PeriodoId         int    `json:"PeriodoId"`
	Valor             string `json:"Valor"`
	FechaCreacion     string `json:"FechaCreacion"`
	FechaModificacion string `json:"FechaModificacion"`
	Activo            bool   `json:"Activo"`
}
