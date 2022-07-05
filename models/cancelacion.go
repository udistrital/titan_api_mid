package models

import "time"

type Cancelacion struct {
	Documento        string
	NumeroContrato   string
	Vigencia         int
	FechaCancelacion time.Time
}
