package models

type ResolucionCompleta struct {
	Id                      int
	NumeroResolucion        string
	FechaExpedicion         string
	Vigencia                int
	DependenciaId           int
	TipoResolucionId        int
	PreambuloResolucion     string
	ConsideracionResolucion string
	NumeroSemanas           int
	Periodo                 int
	Titulo                  int
	DependenciaFirmaId      int
	VigenciaCarga           int
	PeriodoCarga            int
	CuadroResponsabilidades string
	NuxeoUid                string
	Activo                  bool
	FechaCreacion           string
	FechaModificacion       string
}
