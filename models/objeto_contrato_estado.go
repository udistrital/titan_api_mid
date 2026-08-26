package models

type ObjetoContratoEstado struct {
	ContratoEstado struct {
		ValorContrato  string `json:"valorContrato"`
		Estado         Estado `json:"estado"`
		Vigencia       string `json:"vigencia"`
		NumeroContrato string `json:"numeroContrato"`
	} `json:"contratoEstado"`
}

type Estado struct {
	Id           int
	NombreEstado string
}

type ListaContratos struct {
	VigenciaContrato string
	Total            string
	NumeroContrato   string
}
