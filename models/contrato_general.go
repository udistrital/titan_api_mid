package models

import (
	"time"

)

type ContratoGeneral struct {
	Id                           string                 `orm:"column(numero_contrato);pk"`
	Vigencia                     int                 `orm:"column(vigencia)"`
	ObjetoContrato               string              `orm:"column(objeto_contrato);null"`
	PlazoEjecucion               int                 `orm:"column(plazo_ejecucion)"`
	FormaPago                    *Parametros         `orm:"column(forma_pago);rel(fk)"`
	OrdenadorGasto               *ArgoOrdenadores    `orm:"column(ordenador_gasto);rel(fk)"`
	ClausulaRegistroPresupuestal bool                `orm:"column(clausula_registro_presupuestal);null"`
	SedeSolicitante              string              `orm:"column(sede_solicitante);null"`
	DependenciaSolicitante       string              `orm:"column(dependencia_solicitante);null"`
	NumeroSolicitudNecesidad     int                 `orm:"column(numero_solicitud_necesidad)"`
	NumeroCdp                    int                 `orm:"column(numero_cdp)"`
	Contratista                  *InformacionProveedor   `orm:"column(contratista);rel(fk)"`
	UnidadEjecucion              *Parametros         `orm:"column(unidad_ejecucion);rel(fk)"`
	ValorContrato                float64             `orm:"column(valor_contrato)"`
	Justificacion                string              `orm:"column(justificacion)"`
	DescripcionFormaPago         string              `orm:"column(descripcion_forma_pago)"`
	Condiciones                  string              `orm:"column(condiciones)"`
	UnidadEjecutora              *UnidadEjecutora    `orm:"column(unidad_ejecutora);rel(fk)"`
	FechaRegistro                time.Time           `orm:"column(fecha_registro);type(date)"`
	TipologiaContrato            int                 `orm:"column(tipologia_contrato)"`
	TipoCompromiso               int                 `orm:"column(tipo_compromiso)"`
	ModalidadSeleccion           int                 `orm:"column(modalidad_seleccion)"`
	Procedimiento                int                 `orm:"column(procedimiento)"`
	RegimenContratacion          int                 `orm:"column(regimen_contratacion)"`
	TipoGasto                    int                 `orm:"column(tipo_gasto)"`
	TemaGastoInversion           int                 `orm:"column(tema_gasto_inversion)"`
	OrigenPresupueso             int                 `orm:"column(origen_presupueso)"`
	OrigenRecursos               int                 `orm:"column(origen_recursos)"`
	TipoMoneda                   int                 `orm:"column(tipo_moneda)"`
	ValorContratoMe              float64             `orm:"column(valor_contrato_me);null"`
	ValorTasaCambio              float64             `orm:"column(valor_tasa_cambio);null"`
	TipoControl                  int                 `orm:"column(tipo_control);null"`
	Observaciones                string              `orm:"column(observaciones);null"`
	Supervisor                   *SupervisorContrato `orm:"column(supervisor);rel(fk)"`
	ClaseContratista             int                 `orm:"column(clase_contratista)"`
	Convenio                     string              `orm:"column(convenio);null"`
	NumeroConstancia             int                 `orm:"column(numero_constancia);null"`
	Estado                       bool                `orm:"column(estado);null"`
	ResgistroPresupuestal        int                 `orm:"column(resgistro_presupuestal);null"`
	TipoContrato                 *TipoContrato       `orm:"column(tipo_contrato);rel(fk)"`
	LugarEjecucion               *LugarEjecucion     `orm:"column(lugar_ejecucion);rel(fk)"`
}
