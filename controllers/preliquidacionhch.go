package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

// PreliquidacionhchController operations for Preliquidacioncthch
type PreliquidacionhchController struct {
	beego.Controller
}

// Preliquidar ...
// @Title Preliquidar
// @Description Preliquidacion para Honorarios
func (c *PreliquidacionhchController) Preliquidar(datos models.DatosPreliquidacion, reglasbase string) (res []models.Respuesta) {
	//declaracion de variables

	var resumenPreliqu []models.Respuesta
	var tempDocentes models.ObjetoFuncionarioContrato
	var controlError error

	//var al, ml, dl int
	//-----------------------

	//carga de informacion de los empleados a partir del id de persona Natural (en este momento id proveedor)

	//var wg sync.WaitGroup
	//wg.Add(len(datos.PersonasPreLiquidacion))

	for i := 0; i < len(datos.PersonasPreLiquidacion); i++ {
		//go func(i int) {
		//defer wg.Done()
		if datos.PersonasPreLiquidacion[i].EstadoDisponibilidad == 1 {

			var respuesta string
			var verificacionPagoPendientes = 2

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

			tempDocentes, controlError = GetContratosPorPersonaHCH(datos, datos.PersonasPreLiquidacion[i])

			//AGRUPAR PARA CALCULAR SOBRE VALORES TOTALES
			if controlError == nil {

				//BORRAR LO YA PRELIQUIDADO ANTERIORMENTE
				if datos.Preliquidacion.Definitiva == true {
					var d []models.DetallePreliquidacion
					query := "Preliquidacion.Id:" + strconv.Itoa(datos.Preliquidacion.Id) + ",Persona:" + strconv.Itoa(datos.PersonasPreLiquidacion[i].IdPersona)

					if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion?limit=-1&query="+query, &d); err == nil {
						if len(d) != 0 || d[0].Id != 0 {
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

				//tempAgrupar := make(map[string]interface{}) // este mapa tiene la siguiente estructura: tempAgrupar[numero_cedula_docente][id_resolucion][valor_total] (cada resolucion tiene un único tipo de nivel académico, por lo tanto los valores totales se van sumando de acuerdo a la resolución )
				infoResolucion := make(map[string]string)
				infoResoluciones := make(map[string]interface{})

				for _, dato := range tempDocentes.ContratosTipo.ContratoTipo {
					fmt.Println("dato contratos", dato)
					var vinculaciones []models.VinculacionDocente
					query := "NumeroContrato:" + dato.NumeroContrato + ",Vigencia:" + dato.VigenciaContrato
					if err := request.GetJson("http://"+beego.AppConfig.String("Urlargocrud")+":"+beego.AppConfig.String("Portargocrud")+"/"+beego.AppConfig.String("Nsargocrud")+"/vinculacion_docente?limit=-1&query="+query, &vinculaciones); err == nil {

						_, ok := infoResolucion[strconv.Itoa(vinculaciones[0].IdResolucion.Id)]
						if ok {

							infoResolucionTemp := make(map[string]string)
							tempValor, _ := strconv.Atoi(infoResolucion[strconv.Itoa(vinculaciones[0].IdResolucion.Id)])
							tempValor = tempValor + int(vinculaciones[0].ValorContrato)
							infoResolucion[strconv.Itoa(vinculaciones[0].IdResolucion.Id)] = strconv.Itoa(tempValor)
							infoResolucionTemp["NumeroContrato"] = dato.NumeroContrato
							infoResolucionTemp["VigenciaContrato"] = dato.VigenciaContrato
							infoResolucionTemp["Total"] = infoResolucion[strconv.Itoa(vinculaciones[0].IdResolucion.Id)]
							infoResoluciones[strconv.Itoa(vinculaciones[0].IdResolucion.Id)] = infoResolucionTemp

						} else {

							infoResolucionTemp := make(map[string]string)
							tempValor := int(vinculaciones[0].ValorContrato)

							infoResolucion[strconv.Itoa(vinculaciones[0].IdResolucion.Id)] = strconv.Itoa(tempValor)
							infoResolucionTemp["NumeroContrato"] = dato.NumeroContrato
							infoResolucionTemp["VigenciaContrato"] = dato.VigenciaContrato
							infoResolucionTemp["Total"] = infoResolucion[strconv.Itoa(vinculaciones[0].IdResolucion.Id)]
							infoResoluciones[strconv.Itoa(vinculaciones[0].IdResolucion.Id)] = infoResolucionTemp

						}
					}

				}

				//CALCULAR PRELIQUIDACIÓN PARA CADA VALOR AGRUPADO
				for key := range infoResoluciones {
					aux := models.ListaContratos{}
					if err := formatdata.FillStruct(infoResoluciones[key], &aux); err == nil {

						resumenPreliqu = append(resumenPreliqu, liquidarContratoHCH(reglasbase, datos.PersonasPreLiquidacion[i].NumDocumento, datos.PersonasPreLiquidacion[i].IdPersona, datos.Preliquidacion, aux)...)
					} else {
						fmt.Println("error al guardar información agrupada", err)
					}
				}

			}

		}
		//}(i)
	}
	//wg.Wait()
	//CALCULAR FONDO DE SOLIDARIDAD Y RETEFUENTE
	resultadoDesc := CalcularDescuentosTotales(reglasbase, datos.Preliquidacion, resumenPreliqu)
	var idDetaPre interface{}
	if datos.Preliquidacion.Definitiva == true {
		fmt.Println("resultado fondo", resultadoDesc)
		for _, descuentos := range resultadoDesc {
			valor, _ := strconv.ParseFloat(descuentos.Valor, 64)
			diasLiquidados, _ := strconv.ParseFloat(descuentos.DiasLiquidados, 64)
			tipoPreliquidacion, _ := strconv.Atoi(descuentos.TipoPreliquidacion)
			detallepreliqu := models.DetallePreliquidacion{Concepto: &models.ConceptoNomina{Id: descuentos.Id}, Preliquidacion: &models.Preliquidacion{Id: datos.Preliquidacion.Id}, ValorCalculado: valor, Persona: descuentos.IdPersona, DiasLiquidados: diasLiquidados, TipoPreliquidacion: &models.TipoPreliquidacion{Id: tipoPreliquidacion}}

			if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion", "POST", &idDetaPre, &detallepreliqu); err == nil {

			} else {
				fmt.Println("error1: ", err)
			}
		}
	} else {
		*resumenPreliqu[0].Conceptos = append(*resumenPreliqu[0].Conceptos, resultadoDesc...)
	}

	return resumenPreliqu
}

func liquidarContratoHCH(reglasbase string, NumDocumento, Persona int, preliquidacion models.Preliquidacion, informacionContrato models.ListaContratos) (res []models.Respuesta) {
	var predicados []models.Predicado //variable para inyectar reglas
	var objetoDatosActa models.ObjetoActaInicio
	var errorConsultaActa error
	var predicadosRetefuente string

	var resumenPreliqu []models.Respuesta
	var reglasinyectadas string
	var reglas string
	var FechaInicio time.Time
	var FechaFin time.Time
	var dispo int

	var idDetaPre interface{}

	objetoDatosActa, errorConsultaActa = ActaInicioDVE(informacionContrato.NumeroContrato, informacionContrato.VigenciaContrato)

	if errorConsultaActa == nil {

		datosActa := objetoDatosActa
		layout := "2006-01-02"
		vigenciaContratoString, _ := strconv.Atoi(informacionContrato.VigenciaContrato)

		FechaInicio, _ = time.Parse(layout, datosActa.ActaInicio.FechaInicioTemp)
		FechaFin, _ = time.Parse(layout, datosActa.ActaInicio.FechaFinTemp)

		diasContrato := CalcularDias(FechaInicio, FechaFin) + 1 //Suma uno para día inclusive
		vigenciaContrato := informacionContrato.VigenciaContrato
		fmt.Println("valor_contrato->", informacionContrato.Total)
		predicados = append(predicados, models.Predicado{Nombre: "valor_contrato(" + strconv.Itoa(Persona) + "," + informacionContrato.Total + "). "})
		predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + strconv.Itoa(Persona) + "," + strconv.FormatFloat(diasContrato, 'f', -1, 64) + "," + informacionContrato.VigenciaContrato + "). "})

		reglasinyectadas = FormatoReglas(predicados)

		predicadosRetefuente = CargarDatosRetefuente(NumDocumento)
		dispo = verificacionPago(Persona, preliquidacion.Ano, preliquidacion.Mes, informacionContrato.NumeroContrato, informacionContrato.VigenciaContrato)
		reglas = reglasinyectadas + reglasbase + predicadosRetefuente + "estado_pago(" + strconv.Itoa(dispo) + ")."

		temp := golog.CargarReglasCT(Persona, reglas, preliquidacion, vigenciaContrato, datosActa)

		resultado := temp[len(temp)-1]
		resultado.Id = Persona
		resultado.NumDocumento = float64(NumDocumento)
		resultado.NumeroContrato = informacionContrato.NumeroContrato
		resultado.VigenciaContrato = informacionContrato.VigenciaContrato
		resultado.TotalDevengos, resultado.TotalDescuentos, resultado.TotalAPagar = CalcularTotalesPorPersona(*resultado.Conceptos)

		dispo = verificacionPago(NumDocumento, preliquidacion.Ano, preliquidacion.Mes, informacionContrato.NumeroContrato, informacionContrato.VigenciaContrato)

		//se guardan los conceptos calculados en la nomina

		if preliquidacion.Definitiva == true {

			for _, descuentos := range *resultado.Conceptos {
				valor, _ := strconv.ParseFloat(descuentos.Valor, 64)
				diasLiquidados, _ := strconv.ParseFloat(descuentos.DiasLiquidados, 64)
				tipoPreliquidacion, _ := strconv.Atoi(descuentos.TipoPreliquidacion)
				detallepreliqu := models.DetallePreliquidacion{Concepto: &models.ConceptoNomina{Id: descuentos.Id}, Preliquidacion: &models.Preliquidacion{Id: preliquidacion.Id}, ValorCalculado: valor, NumeroContrato: informacionContrato.NumeroContrato, VigenciaContrato: vigenciaContratoString, Persona: Persona, DiasLiquidados: diasLiquidados, TipoPreliquidacion: &models.TipoPreliquidacion{Id: tipoPreliquidacion}, EstadoDisponibilidad: &models.EstadoDisponibilidad{Id: dispo}}

				if err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/detalle_preliquidacion", "POST", &idDetaPre, &detallepreliqu); err == nil {

				} else {
					fmt.Println("error1: ", err)
				}
			}
		}

		//
		resumenPreliqu = append(resumenPreliqu, resultado)
		predicados = nil
		reglas = ""
		reglasinyectadas = ""

	} else {
		fmt.Println("error al traer acta de inicio")
	}

	return resumenPreliqu

}
