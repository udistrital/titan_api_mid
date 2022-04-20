package models

type DesagregadoContratoHCS struct {
	NumeroContrato        string
	Vigencia              int
	SueldoBasico          float64
	PrimaServicios        float64
	PrimaVacaciones       float64
	InteresesCesantias    float64
	Cesantias             float64
	Vacaciones            float64
	PrimaNavidad          float64
	BonificacionServicios float64
}
