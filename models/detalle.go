package models

type Detalle struct {
	Contrato        string
	Vigencia        int
	TotalDevengado  float64
	TotalDescuentos float64
	TotalPago       float64
	Detalle         []DetallePreliquidacion
}

type DetalleDVE struct {
	//Resolucion         *Resolucion
	//ResolucionCompleta *ResolucionCompleta
	Detalle []Detalle
}

type DetalleMensual struct {
	TotalDevengado  float64
	TotalDescuentos float64
	Detalle         []DetallePreliquidacion
}
