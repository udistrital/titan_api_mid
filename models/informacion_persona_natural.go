package models

/*type InformacionPersonaNatural struct {
	Id string 
	PersonasACargo bool
	DeclaranteRenta bool
	InteresViviendaAfc float64
	MedicinaPrepagada bool
	Pensionado bool

}*/
type InformacionPersonaNatural struct {
	InformacionPersonaNatural struct {
		Id 			string `json:"Id"`
		PersonasACargo		bool `json:"PersonasACargo"`
		DeclaranteRenta      	bool `json:"DeclaranteRenta"`
		InteresViviendaAfc      float64 `json:"InteresViviendaAfc"`
		MedicinaPrepagada      	bool `json:"MedicinaPrepagada"`
		Pensionado      	bool `json:"Pensionado"`
	} `json:"informacion_persona_natural"`
}

//informacionproveedor.id_proveedor, contratos.contratista,contratos.numero_contrato,informacionproveedor.nom_proveedor

// last inserted Id on success.
