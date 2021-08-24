package models

type PersonasPreliquidacion struct {
	IdPersona            int
	NombreCompleto       string
	NumDocumento         int
	NumeroContrato       string
	VigenciaContrato     int
	Preliquidacion       int
	EstadoDisponibilidad int
	TipoContrato         string
	FechaInicio          string
	FechaFin             string
	ValorContrato        string
}
