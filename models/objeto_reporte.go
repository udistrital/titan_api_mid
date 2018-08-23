package models

type ObjetoReporte struct {
	ProyectoCurricular int
	Preliquidacion  *Preliquidacion
	TotalDev 			float64
	TotalDesc       float64
	TotalDocentes	int
}
