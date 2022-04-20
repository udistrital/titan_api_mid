package models

type Resolucion struct {
	Id                int
	FacultadId        int
	Dedicacion        string
	NivelAcademico    string
	Activo            bool
	FechaCreacion     string
	FechaModificacion string
}
