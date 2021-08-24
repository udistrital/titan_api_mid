package models

type Resumen struct {
	NombreConcepto       string
	NaturalezaConcepto   string
	NaturalezaConceptoId string
	Total                string
}

type ResumentCompleto struct {
	TotalDevengos         int
	TotalDescuentos       int
	ResumenTotalConceptos []Resumen
}
