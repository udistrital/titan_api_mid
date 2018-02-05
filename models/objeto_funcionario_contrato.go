package models

type ObjetoFuncionarioContrato struct {
	ContratosTipo struct {
		ContratoTipo[]struct {
			Id              string     `json:"id_proveedor"`
			NombreProveedor string  `json:"nom_proveedor"`
			NumDocumento    string `json:"num_documento"`
			NumeroContrato  string  `json:"numero_contrato"`
			VigenciaContrato  string  `json:"vigencia"`

		} `json:"contrato_tipo"`
	} `json:"contratos_tipo"`
}
