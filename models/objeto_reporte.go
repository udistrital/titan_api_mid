package models

type ObjetoReporte struct {
	ProyectoCurricular int
	Facultad int
	Preliquidacion  *Preliquidacion
	TotalDev 			float64
	TotalDesc       float64
	TotalDocentes	int
}
