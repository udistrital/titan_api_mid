package models

type ObjetoInformacionContratista struct {
	InformacionContratista struct {
			Documento Documento `json:"Documento"`
			Contrato 	Contrato  `json:"contrato"`
			NombreCompleto  string  `json:"nombre_completo"`
	} `json:"informacion_contratista"`
}


type Documento struct {
	Numero       string   `json:"numero"`
}

type Contrato struct {
	Numero       string   `json:"numero"`
}
