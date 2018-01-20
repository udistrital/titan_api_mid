package models

type Contratista struct {
	Id                int `orm:"column(id_proveedor)"`
	NumeroContrato    int `orm:"column(numero_contrato)"`
	Asignacion_basica int `orm:"column(asignacion_basica)"`
}

//informacionproveedor.id_proveedor, contratos.contratista,contratos.numero_contrato,informacionproveedor.nom_proveedor

// last inserted Id on success.
