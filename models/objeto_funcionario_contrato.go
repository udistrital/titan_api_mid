package models

type ObjetoFuncionarioContrato struct {
	ContratosTipo struct {
		ContratoTipo[]struct {
			Id              string     `json:"id_proveedor"`
			NombreProveedor string  `json:"nom_proveedor"`
			NumDocumento    string `json:"num_documento"`
			NumeroContrato  string  `json:"numero_contrato"` //Puede borrarse
			VigenciaContrato  string  `json:"vigencia"`      //Puede borrarse
			Preliquidado     string
		} `json:"contrato_tipo"`
	} `json:"contratos_tipo"`
}
