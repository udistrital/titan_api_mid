package models

type Mensaje struct {
	ArnTopic       string
	Asunto         string
	Atributos      map[string]interface{}
	DestinatarioId []int
	Mensaje        string
	RemitenteId    string
}
