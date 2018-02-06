package models

type ObjetoActaInicio struct {
	ActaInicio struct {
		FechaInicioTemp 	string `json:"fechaInicio"`
		FechaFinTemp      string				 `json:"fechaFin"`
	} `json:"actaInicio"`
}
