package models

type VinculacionDocente struct {
	Id                             int
	NumeroContrato                 string
	Vigencia                       int
	PersonaId                      int
	NumeroHorasSemanales           int
	NumeroSemanas                  int
	PuntoSalarialId                int
	SalarioMinimoId                int
	ResolucionVinculacionDocenteId *Resolucion `json:"ResolucionVinculacionDocenteId"`
	DedicacionId                   int
	ProyectoCurricularId           int
	ValorContrato                  int
	Categoria                      string
	Emerito                        bool
	DependenciaAcademica           int
	NumeroRp                       int
	VigenciaRp                     int
	FechaInicio                    string
	Activo                         bool
	FechaCreacion                  string
	FechaModificacion              string
}
