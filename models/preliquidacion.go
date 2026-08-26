package models

type Preliquidacion struct {
	Id                     int    `json:"Id"`
	Descripcion            string `json:"Descripcion"`
	Mes                    int    `json:"Mes"`
	Ano                    int    `json:"Ano"`
	EstadoPreliquidacionId int    `json:"EstadoPreliquidacionId"`
	NominaId               int    `json:"NominaId"`
	Activo                 bool   `json:"Activo"`
	FechaCreacion          string `json:"FechaCreacion"`
	FechaModificacion      string `json:"FechaModificacion"`
}
