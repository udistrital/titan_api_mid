package models

type Resumen struct {
	NombreConcepto  string
	NaturalezaConcepto string
  Total           string
}

type ResumentCompleto struct {
	TotalDevengos    int
	TotalDescuentos  int
	ResumenTotalConceptos []Resumen
}
