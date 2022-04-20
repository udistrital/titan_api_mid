package models

type VinculacionDocente struct {
	Id                             int
	PersonaId                      int
	NumeroHorasSemanales           int
	NumeroSemanas                  int
	ResolucionVinculacionDocenteId *Resolucion
	DedicacionId                   int
	ValorContrato                  int
	Categoria                      string
}
