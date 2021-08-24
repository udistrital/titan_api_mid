package models

import (
	"time"
)

type Preliquidacion struct {
	Id                     int
	Descripcion            string
	Mes                    int
	Ano                    int
	Activo                 bool
	EstadoPreliquidacionId *EstadoPreliquidacion
	NominaId               *Nomina
	FechaCreacion          time.Time
	FechaModificacion      time.Time
	Definitiva             bool
}
