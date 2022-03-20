package models

type PersonaNatural struct {
	PersonasACargo         string `json:"PersonasACargo"`
	InteresViviendaAfc     string `json:"InteresViviendaAfc"`
	ValorUvtPrepagada      string `json:"ValorUvtPrepagada"`
	Pensionado             string `json:"Pensionado"`
	ValorAfc               string `json:"ValorAfc"`
	ValorPensionVoluntaria string `json:"ValorPensionVoluntaria"`
	ResponsableIva         string `json:"ResponsableIva"`
}
