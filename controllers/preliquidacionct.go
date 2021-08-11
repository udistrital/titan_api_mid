package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// PreliquidacionctController operations for Preliquidacionct
type PreliquidacionctController struct {
	beego.Controller
}

// GetIBCPorNovedad ...
// @Title GetIBCPorNovedad
// @Description Funcion para calcular IBC para una novedad específica
func (c *PreliquidacionctController) GetIBCPorNovedad(ano, mes, numDocumento, idPersona int, reglasbase, novedad string) (res int) {

	var resumenPreliqu []models.Respuesta
	var ibcNovedad int

	preliquidacion := models.Preliquidacion{Ano: ano, Mes: mes}
	datos := models.DatosPreliquidacion{Preliquidacion: preliquidacion}

	arreglo_contratos, _ := GetContratosPorPersonaCT(datos, numDocumento)

	for _, info := range arreglo_contratos.ContratosTipo.ContratoTipo {
		vigencia, _ := strconv.Atoi(info.VigenciaContrato)
		persona := models.PersonasPreliquidacion{NumDocumento: numDocumento, IdPersona: idPersona, NumeroContrato: strings.Replace(info.NumeroContrato, "c", "", -1), VigenciaContrato: vigencia}
		resumenPreliqu = append(resumenPreliqu, liquidarContratoCT(persona, preliquidacion, reglasbase, novedad)...)
	}

	for _, res := range resumenPreliqu {
		for _, concepto := range *res.Conceptos {
			if concepto.Id == 311 {
				temp, _ := strconv.Atoi(concepto.Valor)
				ibcNovedad = ibcNovedad + temp
			}
		}
	}

	return ibcNovedad

}

// Preliquidar ...
// @Title Preliquidar
// @Description Preliquidacion para contratistas
func (c *PreliquidacionctController) Preliquidar(datos models.DatosPreliquidacion, reglasbase string) (res []models.Respuesta) {

	var resumenPreliqu []models.Respuesta
	//-----------------------

	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {

		if datos.PersonasPreLiquidacion[i].IdPersona != 0 {
			
			if datos.PersonasPreLiquidacion[i].EstadoDisponibilidad == 1 {

				var respuesta string
				var verificacionPagoPendientes = 2
				datos.PersonasPreLiquidacion[i].NumeroContrato = strings.Replace(datos.PersonasPreLiquidacion[i].NumeroContrato, "c", "", -1)
				datos.PersonasPreLiquidacion[i].ValorContrato = strings.Replace(datos.PersonasPreLiquidacion[i].ValorContrato, "v", "", -1)
				detallesAMod := ConsultarDetalleAModificar(datos.PersonasPreLiquidacion[i].NumeroContrato, datos.PersonasPreLiquidacion[i].VigenciaContrato, datos.PersonasPreLiquidacion[i].Preliquidacion)
				for _, pos := range detallesAMod {

					verificacionPagoPendientes = verificacionPago(0, datos.Preliquidacion.Ano, datos.Preliquidacion.Mes, pos.NumeroContrato, strconv.Itoa(pos.VigenciaContrato))
					pos.EstadoDisponibilidad = &models.EstadoDisponibilidad{Id: verificacionPagoPendientes}
					if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion/"+strconv.Itoa(pos.Id), "PUT", &respuesta, pos); err == nil {
						fmt.Println("preliquidaciones actualizadas")
					} else {
						fmt.Println("error al actualizar detalle de preliquidación: ", err)
					}
				}

			} else {

				//eliminar los registros ya existentes en caso de ser definitiva y no solo consulta
				if datos.Preliquidacion.Definitiva == true {
					var d []models.DetallePreliquidacion
					datos.PersonasPreLiquidacion[i].NumeroContrato = strings.Replace(datos.PersonasPreLiquidacion[i].NumeroContrato, "c", "", -1)
					query := "Preliquidacion.Id:" + strconv.Itoa(datos.Preliquidacion.Id) + ",NumeroContrato:" + datos.PersonasPreLiquidacion[i].NumeroContrato + ",VigenciaContrato:" + strconv.Itoa(datos.PersonasPreLiquidacion[i].VigenciaContrato)

					if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
						if len(d) != 0 {
							for _, dato := range d {
								urlcrud := "http://" + beego.AppConfig.String("Urlcrud") + ":" + beego.AppConfig.String("Portcrud") + "/" + beego.AppConfig.String("Nscrud") + "/detalle_preliquidacion/" + strconv.Itoa(dato.Id)
								var res string
								if err := request.SendJson(urlcrud, "DELETE", &res, nil); err == nil {
									fmt.Println("borrado correctamente")
								} else {
									fmt.Println("error", err)
								}
							}
						}

					} else {
						fmt.Println("error de detalle", err)
					}
				}

				resumenPreliqu = liquidarContratoCT(datos.PersonasPreLiquidacion[i], datos.Preliquidacion, reglasbase, "")

			}
		}
	}

	return resumenPreliqu
}

func liquidarContratoCT(persona models.PersonasPreliquidacion, preliquidacion models.Preliquidacion, reglasbase, novedadInyectada string) (res []models.Respuesta) {

	var predicados []models.Predicado //variable para inyectar reglas
	var predicadosRetefuente string
	var resumenPreliqu []models.Respuesta

	var reglasinyectadas string
	var reglas string
	var disp int
        var pensionado bool
	var dependientes bool
	var idDetaPre interface{}

	var objetoDatosContrato models.ObjetoContratoEstado
	var objetoDatosActa models.ObjetoActaInicio
	var errorConsultaContrato error
	var errorConsultaActa error

	var FechaInicio time.Time
	var FechaFin time.Time
	fmt.Println("valor del contrato ", persona.NumeroContrato)
	fmt.Println("valor del contrato ", persona.VigenciaContrato)
	objetoDatosContrato, errorConsultaContrato = ContratosContratistas(persona.NumeroContrato, persona.VigenciaContrato)

	objetoDatosActa, errorConsultaActa = ActaInicioContratistas(persona.NumeroContrato, persona.VigenciaContrato)

	if errorConsultaContrato == nil {
		if errorConsultaActa == nil {
			datosContrato := objetoDatosContrato.ContratoEstado
			datosActa := objetoDatosActa.ActaInicio

			layout := "2006-01-02"
			//se modifica por fechas listadas
			FechaInicio, _ = time.Parse(layout, datosActa.FechaInicioTemp)
			FechaFin, _ = time.Parse(layout, datosActa.FechaFinTemp)

			diasContrato := CalcularDias(FechaInicio, FechaFin) + 1 //Suma uno para día inclusive
			fmt.Println("valor del contrato ", datosContrato.ValorContrato)
			fmt.Println("días contrato ", diasContrato)

			vigenciaContrato := strconv.Itoa(persona.VigenciaContrato)
			predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + strconv.Itoa(persona.IdPersona) + "," + datosContrato.ValorContrato + "). "})
			predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + strconv.Itoa(persona.IdPersona) + "," + strconv.FormatFloat(diasContrato, 'f', -1, 64) + "," + vigenciaContrato + "). "})

			reglasinyectadas = FormatoReglas(predicados)

			if novedadInyectada == "" {
				reglasinyectadas = reglasinyectadas + CargarNovedadesPersona(persona.IdPersona, persona.NumeroContrato, strconv.Itoa(persona.VigenciaContrato), preliquidacion)
			} else {
				reglasinyectadas = reglasinyectadas + novedadInyectada
			}
                        
			predicadosRetefuente, pensionado, dependientes = CargarDatosRetefuente(persona.NumDocumento)
			disp = verificacionPago(persona.IdPersona, preliquidacion.Ano, preliquidacion.Mes, persona.NumeroContrato, strconv.Itoa(persona.VigenciaContrato))
			reglas = reglasinyectadas + reglasbase + predicadosRetefuente + "estado_pago(" + strconv.Itoa(disp) + ")."
			//reglas = reglasinyectadas + reglasbase + predicadosRetefuente + "estado_pago(2)."
			fmt.Println("dep", dependientes)
			temp := golog.CargarReglasCT(persona.IdPersona, reglas, preliquidacion, vigenciaContrato, objetoDatosActa, pensionado, dependientes)
			resultado := temp[len(temp)-1]
			fmt.Println("resultado", resultado)
			resultado.NumDocumento = float64(persona.NumDocumento)
			resultado.NumeroContrato = persona.NumeroContrato
			resultado.VigenciaContrato = strconv.Itoa(persona.VigenciaContrato)
			resultado.TotalDevengos, resultado.TotalDescuentos, resultado.TotalAPagar = CalcularTotalesPorPersona(*resultado.Conceptos)

			if disp == 1 {
				resultado.EstadoPago = "Pendiente"
			} else if disp == 2 {
				resultado.EstadoPago = "Listo para pago"
			}

			//INSERTAR LOS REGISTROS SI LA PRELIQUIDACIÓN ES DEFINITIVA
			if preliquidacion.Definitiva == true {
				for _, descuentos := range *resultado.Conceptos {
					valor, _ := strconv.ParseFloat(descuentos.Valor, 64)
					diasLiquidados, _ := strconv.ParseFloat(descuentos.DiasLiquidados, 64)
					tipoPreliquidacion, _ := strconv.Atoi(descuentos.TipoPreliquidacion)
					detallepreliqu := models.DetallePreliquidacion{Concepto: &models.ConceptoNomina{Id: descuentos.Id}, Preliquidacion: &models.Preliquidacion{Id: preliquidacion.Id}, ValorCalculado: valor, NumeroContrato: persona.NumeroContrato, VigenciaContrato: persona.VigenciaContrato, Persona: persona.IdPersona, DiasLiquidados: diasLiquidados, TipoPreliquidacion: &models.TipoPreliquidacion{Id: tipoPreliquidacion}, EstadoDisponibilidad: &models.EstadoDisponibilidad{Id: disp}}

					if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion", "POST", &idDetaPre, &detallepreliqu); err == nil {

					} else {
						fmt.Println("error1: ", err)
					}
				}
			}

			//
			resumenPreliqu = append(resumenPreliqu, resultado)
			predicados = nil
			objetoDatosContrato = models.ObjetoContratoEstado{}
			reglas = ""
			reglasinyectadas = ""

		} else {
			fmt.Println("error al traer acta de inicio")
		}

	} else {
		fmt.Println("error al traer valor del contrato")
	}

	return resumenPreliqu
}

// ConsultarDetalleAModificar ...
// @Title ConsultarDetalleAModificar
// @Description Consultar los detalles por contrato, vigencia y preliquidacion
func ConsultarDetalleAModificar(id_contrato string, vigencia, preliquidacion int) (det []models.DetallePreliquidacion) {
	var v []models.DetallePreliquidacion
	query := "NumeroContrato:" + id_contrato + ",VigenciaContrato:" + strconv.Itoa(vigencia) + ",Preliquidacion.Id:" + strconv.Itoa(preliquidacion)
	if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?query="+query, &v); err == nil && v != nil {

	} else {
		fmt.Println("error al consultar preliquidacion a modificar ", err)
	}

	return v

}
