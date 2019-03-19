package models

type ObjetoInformacionContratista struct {
	InformacionContratista struct {
			Documento Documento `json:"Documento"`
			Contrato 	Contrato  `json:"contrato"`
			NombreCompleto  string  `json:"nombre_completo"`
			Dependencia Supervisor `json:"Supervisor"`
	} `json:"informacion_contratista"`
}


type Documento struct {
	Numero       string   `json:"numero"`
}

type Contrato struct {
	Numero       string   `json:"numero"`
}

type Supervisor struct {
	IdDependencia       string   `json:"id_dep"`
}
