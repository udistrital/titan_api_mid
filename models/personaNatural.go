package models

type PersonaNatural struct {
	Dependientes       string  `json:"PersonasACargo"`
	InteresViviendaAfc float64 `json:"InteresViviendaAfc"`
	ValorUvtPrepagada  float64 `json:"ValorUvtPrepagada"`
	Pensionado         string  `json:"Pensionado"`
	Afc                float64 `json:"ValorAfc"`
	PensionVoluntaria  float64 `json:"ValorPensionVoluntaria"`
	Reteiva            string  `json:"ResponsableIva"`
}
