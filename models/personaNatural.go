package models

type PersonaNatural struct {
	Dependientes       bool    `json:"PersonasACargo"`
	InteresViviendaAfc float64 `json:"InteresViviendaAfc"`
	ValorUvtPrepagada  float64 `json:"ValorUvtPrepagada"`
	Pensionado         string  `json:"Pensionado"`
	//Falta agregar el resto de alivios, como el AFC, la pensión voluntaria y responsable de iva
}
