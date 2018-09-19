package models

type ObjetoReporte struct {
	ProyectoCurricular int
	Facultad int
	Dependencia string
	Preliquidacion  *Preliquidacion
	TotalDev 			float64
	TotalDesc       float64
	TotalDocentes	int
}
