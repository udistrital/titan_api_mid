package models

type InformacionProveedor struct {
	Id           int    `orm:"column(id);pk"`
	NumDocumento string `orm:"column(num_documento);"`
	NomProveedor string `orm:"column(nom_proveedor);pk"`
}
