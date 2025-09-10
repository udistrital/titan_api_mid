package models

type Desagregado struct {
	PrimaNavidad          float64
	Cesantias             float64
	PrimaServicios        float64
	PrimaVacaciones       float64
	Vacaciones            float64
	InteresesCesantias    float64
	BonificacionServicios float64
}

type PorcentajeDesagregado struct {
	PorcentajePrimaNavidad    float64
	PorcentajeCesantias       float64
	PorcentajePrimaServicios  float64
	PorcentajePrimaVacaciones float64
	PorcentajeVacaciones      float64
}
