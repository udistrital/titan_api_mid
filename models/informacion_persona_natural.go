package models

type InformacionPersonaNatural struct {
	Id                 string
	PersonasACargo     bool
	DeclaranteRenta    bool
	InteresViviendaAfc float64
	MedicinaPrepagada  bool
	IdFondoPension     int
	ValorUvtPrepagada  int
}

//informacionproveedor.id_proveedor, contratos.contratista,contratos.numero_contrato,informacionproveedor.nom_proveedor

// last inserted Id on success.
