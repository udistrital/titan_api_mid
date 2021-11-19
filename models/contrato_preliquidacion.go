package models

type ContratoPreliquidacion struct {
	Id                   int
	ContratoId           *Contrato
	PreliquidacionId     *Preliquidacion
	Cumplido             bool
	Preliquidado         bool
	ResponsableIva       bool
	Dependientes         bool
	Pensionado           bool
	InteresesVivienda    float64
	MedicinaPrepagadaUvt float64
	PensionVoluntaria    float64
	Afc                  float64
	Activo               bool
	FechaCreacion        string
	FechaModificacion    string
}
