package models

type VinculacionDocente struct {
	Id                   int
	IdPersona            string
	IdProyectoCurricular int16
	DependenciaAcademica int
	IdResolucion         *Resolucion `json:"IdResolucion"`
	ValorContrato        float64
}
