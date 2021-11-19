package models

type Resolucion struct {
	Id             int    `json:"Id"`
	IdFacultad     int    `json:"IdFacultad"`
	Dedicacion     string `json:"Dedicacion"`
	NivelAcademico string `json:"NivelAcademico"`
}
