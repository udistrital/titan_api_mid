package models

type ObjetoDocentePlanta struct {
	DocenteCollection struct {
		Docente []struct {
			Planta string `json:"planta"`
		} `json:"docentes"`
	} `json:"docentesCollection"`
}
