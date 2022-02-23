package models

type Resolucion struct {
	Id                int    `json:"Id"`
	FacultadId        int    `json:"IdFacultad"`
	Dedicacion        string `json:"Dedicacion"`
	NivelAcademico    string `json:"NivelAcademico"`
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
}
