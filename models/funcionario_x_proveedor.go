package models

type Funcionario_x_Proveedor struct {
	Id              int
	NombreProveedor string
	NumDocumento    int
	NumeroContrato  string
	VigenciaContrato  int
	//IdEPS                  int  							`xml:"id_eps"`
	//IdARL                  int  							`xml:"id_arl"`
	//IdFondoPension         int  							`xml:"id_fondo_pension"`
	//IdCajaCompensacion     int  							`xml:"id_caja_compensacion"`
}
