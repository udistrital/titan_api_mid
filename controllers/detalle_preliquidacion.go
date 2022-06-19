package controllers

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

type DetallePreliquidacionController struct {
	beego.Controller
}

func (c *DetallePreliquidacionController) URLMapping() {
	c.Mapping("ObtenerDetalleCT", c.ObtenerDetalleCT)
	c.Mapping("ObtenerDetalleDVE", c.ObtenerDetalleDVE)
}

// Get ...
// @Title Obtener Resumen CT
// @Description Obtener el detalle de la preliquidación para CPS
// @Param	ano		path 	string	true		"Año de la preliquidación"
// @Param	mes		path 	string	true		"Mes de la preliquidación"
// @Param	contrato		path 	string	true		"Contrato a buscar"
// @Param	vigencia		path 	string	true		"vigencia del contrato"
// @Param	documento		path 	string	true		"Documento del contratista"
// @Success 201 {object} []models.Detalle
// @Failure 403 body is empty
// @router /obtener_detalle_CT/:ano/:mes/:contrato/:vigencia/:documento [get]
func (c *DetallePreliquidacionController) ObtenerDetalleCT() {

	ano := c.Ctx.Input.Param(":ano")
	mes := c.Ctx.Input.Param(":mes")
	contrato := c.Ctx.Input.Param(":contrato")
	vigencia := c.Ctx.Input.Param(":vigencia")
	documento := c.Ctx.Input.Param(":documento")

	detalle, err := TraerDetalleMensual(ano, mes, contrato, vigencia, documento, true, false)

	if err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": detalle}
	} else {
		c.Data["mesaage"] = "Error al obtener el detalle, no hay detalle para el contrato en el mes seleccionado"
		c.Abort("400")
	}

	c.ServeJSON()
}

func TraerDetalleMensual(ano, mes, contrato, vigencia, documento string, CPS bool, HCS bool) (detalle []models.Detalle, err error) {
	var aux map[string]interface{}
	var tempDetalle []models.DetallePreliquidacion
	var auxDetalle models.Detalle
	var auxContratos []models.Contrato
	var query = "NumeroContrato:" + contrato + ",Vigencia:" + vigencia + ",Documento:" + documento
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &auxContratos)
		if auxContratos[0].Id != 0 {
			for j := 0; j < len(auxContratos); j++ {
				query = "ContratoPreliquidacionId.PreliquidacionId.Ano:" + ano + ",ContratoPreliquidacionId.PreliquidacionId.Mes:" + mes + ",ContratoPreliquidacionId.ContratoId.Id:" + strconv.Itoa(auxContratos[j].Id)
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &tempDetalle)
					auxDetalle.Contrato = tempDetalle[0].ContratoPreliquidacionId.ContratoId.NumeroContrato
					auxDetalle.Vigencia = tempDetalle[0].ContratoPreliquidacionId.ContratoId.Vigencia
					for i := 0; i < len(tempDetalle); i++ {
						if tempDetalle[i].ConceptoNominaId.Id == 574 || tempDetalle[i].ConceptoNominaId.Id == 573 {
							fmt.Println("Salto")
						} else if tempDetalle[i].ConceptoNominaId.NaturalezaConceptoNominaId == 424 {
							if HCS && tempDetalle[i].ConceptoNominaId.Id == 570 {
								auxDetalle.Detalle = append(auxDetalle.Detalle, tempDetalle[i])
							} else if CPS && (tempDetalle[i].ConceptoNominaId.Id == 568 || tempDetalle[i].ConceptoNominaId.Id == 569 || tempDetalle[i].ConceptoNominaId.Id == 570) {
								auxDetalle.Detalle = append(auxDetalle.Detalle, tempDetalle[i])
							} else {
								auxDetalle.Detalle = append(auxDetalle.Detalle, tempDetalle[i])
								auxDetalle.TotalDescuentos = auxDetalle.TotalDescuentos + tempDetalle[i].ValorCalculado
							}
						} else if tempDetalle[i].ConceptoNominaId.NaturalezaConceptoNominaId == 423 {
							auxDetalle.Detalle = append(auxDetalle.Detalle, tempDetalle[i])
							auxDetalle.TotalDevengado = auxDetalle.TotalDevengado + tempDetalle[i].ValorCalculado
						} else {
							auxDetalle.Detalle = append(auxDetalle.Detalle, tempDetalle[i])
						}
					}
					auxDetalle.TotalPago = auxDetalle.TotalDevengado - auxDetalle.TotalDescuentos
					detalle = append(detalle, auxDetalle)
				} else {
					fmt.Println("Error al obtener detalle ", err)
					return detalle, err
				}
			}
			return detalle, nil
		} else {
			fmt.Println("Error al obtener contrato o contratos ", err)
			return detalle, err
		}
	} else {
		fmt.Println("Error al obtener contrato o contratos ", err)
		return detalle, err
	}
}

// Get ...
// @Title Obtener Resumen DVE
// @Description Obtener el detalle de la preliquidación para DVE sea HCS o HCH
// @Param	ano		path 	string	true		"Año de la preliquidación"
// @Param	mes		path 	string	true		"Mes de la preliquidación"
// @Param	documento		path 	string	true		"Documento a buscar"
// @Param	nomina		path 	string	true		"Nomina, si es HCH o HCS"
// @Success 201 {object} []models.DetalleDVE
// @Failure 400 the request contains incorrect syntax
// @router /obtener_detalle_DVE/:ano/:mes/:documento/:nomina [get]
func (c *DetallePreliquidacionController) ObtenerDetalleDVE() {
	ano := c.Ctx.Input.Param(":ano")
	mes := c.Ctx.Input.Param(":mes")
	documento := c.Ctx.Input.Param(":documento")
	nomina := c.Ctx.Input.Param(":nomina")
	var auxDetalle []models.Detalle

	var aux map[string]interface{}
	//var vinculacion []models.VinculacionDocente
	var contratoPreliquidacion []models.ContratoPreliquidacion
	//var resolucionCompleta []models.ResolucionCompleta
	var detallesDVE []models.Detalle
	//var tempDetalle models.DetalleDVE

	//obtener los contratos asociados a la persona
	var query = "ContratoId.Documento:" + documento + ",ContratoId.Vigencia:" + ano + ",PreliquidacionId.Ano:" + ano + ",PreliquidacionId.Mes:" + mes + ",PreliquidacionId.NominaId:" + nomina
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contratoPreliquidacion)
		//Recorrer los contratos para obtener las resoluciones
		/*
			for i := 0; i < len(tempContrato); i++ {
				//Encontrar el id de la resolucion
				query := "NumeroContrato:" + tempContrato[i].ContratoId.NumeroContrato + ",Vigencia:" + strconv.Itoa(tempContrato[i].ContratoId.Vigencia)
				if err := request.GetJson(beego.AppConfig.String("UrlAdministrativaCrud")+"/vinculacion_docente?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &vinculacion)
					//Encontrar el número de resolución
					if err := request.GetJson(beego.AppConfig.String("UrlAdministrativaCrud")+"/resolucion?limit=-1&query=Id:"+strconv.Itoa(vinculacion[0].ResolucionVinculacionDocenteId.Id), &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &resolucionCompleta)
						if len(detallesDVE) == 0 {
							tempDetalle.Resolucion = vinculacion[0].ResolucionVinculacionDocenteId
							tempDetalle.ResolucionCompleta = &resolucionCompleta[0]
							tempDetalle.Detalle = append(tempDetalle.Detalle, TraerDetalleMensual(ano, mes, tempContrato[i].ContratoId.NumeroContrato, strconv.Itoa(tempContrato[i].ContratoId.Vigencia), documento))
							detallesDVE = append(detallesDVE, tempDetalle)
						} else {
							res, pos := encontrarResolucion(resolucionCompleta[0].Id, detallesDVE)
							if res {
								detallesDVE[pos].Detalle = append(detallesDVE[pos].Detalle, TraerDetalleMensual(ano, mes, tempContrato[i].ContratoId.NumeroContrato, strconv.Itoa(tempContrato[i].ContratoId.Vigencia), documento))
							} else {
								tempDetalle.Resolucion = vinculacion[0].ResolucionVinculacionDocenteId
								tempDetalle.Detalle = append(tempDetalle.Detalle, TraerDetalleMensual(ano, mes, tempContrato[i].ContratoId.NumeroContrato, strconv.Itoa(tempContrato[i].ContratoId.Vigencia), documento))
								detallesDVE = append(detallesDVE, tempDetalle)
							}
						}
					} else {
						fmt.Println("Error al obtener el número de la resolución: ", err)
					}
				} else {
					fmt.Println("Error al obtener resolución: ", err)
				}
			}
		*/
		//Agregar los detalles de todos los contratos
		for i := 0; i < len(contratoPreliquidacion); i++ {
			if nomina == "416" {
				auxDetalle, _ = TraerDetalleMensual(ano, mes, contratoPreliquidacion[i].ContratoId.NumeroContrato, strconv.Itoa(contratoPreliquidacion[i].ContratoId.Vigencia), documento, false, true)
			} else {
				auxDetalle, _ = TraerDetalleMensual(ano, mes, contratoPreliquidacion[i].ContratoId.NumeroContrato, strconv.Itoa(contratoPreliquidacion[i].ContratoId.Vigencia), documento, false, false)
			}

			if err == nil {
				detallesDVE = append(detallesDVE, auxDetalle[0])
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": detallesDVE}
			} else {
				c.Data["mesaage"] = "Error al obtener detalle de 1 o más contratos"
				c.Abort("404")
				break
			}
		}

	} else {
		fmt.Println("Error al obtener detalle ", err)
		c.Data["mesaage"] = "No existe un contrato asociado al documento que sea vigente para la preliquidación " + err.Error()
		c.Abort("404")
	}
	c.ServeJSON()
}

/*
func encontrarResolucion(idResolucion int, resoluciones []models.DetalleDVE) (res bool, pos int) {
	for i := 0; i < len(resoluciones); i++ {
		if idResolucion == resoluciones[i].ResolucionCompleta.Id {
			return true, i
		}
	}
	return false, 0
}
*/

func AgregarValorNovedad(novedad models.Novedad) (mensaje string, err error) {

	var aux map[string]interface{}
	var mesIterativo = int(novedad.FechaInicio.Month())
	var anoIterativo = novedad.FechaInicio.Year()
	var auxCuotas int
	var contratoPreliquidacion []models.ContratoPreliquidacion
	var honorarios []models.DetallePreliquidacion
	var detalleNuevo models.DetallePreliquidacion
	var idHonorarios int
	auxCuotas = novedad.Cuotas
	fmt.Println(novedad)
	if novedad.ContratoId.TipoNominaId == 410 {
		idHonorarios = 152 //Salario Base
	} else {
		idHonorarios = 87 //Honorarios
	}

	for { //itera desde el mes en el que se aplicó la novedad hasta el fin del numero de cuotas

		fmt.Println("Mes: ", mesIterativo)
		fmt.Println("Ano: ", anoIterativo)
		fmt.Println("Cuotas: ", auxCuotas)
		var query = "ContratoId.Id:" + strconv.Itoa(novedad.ContratoId.Id) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo)
		fmt.Println(beego.AppConfig.String("UrlTitanCrud") + "/contrato_preliquidacion?limit=-1&query=" + query)
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil { //obtiene el contrato_preliquidacion de ese mes y año
			LimpiezaRespuestaRefactor(aux, &contratoPreliquidacion)
			//Obtener el valor de los honorarios para ese mes
			query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:" + strconv.Itoa(idHonorarios)
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &honorarios)
				detalleNuevo.Id = 0
				detalleNuevo.TipoPreliquidacionId = 397
				detalleNuevo.Activo = true
				detalleNuevo.EstadoDisponibilidadId = 426
				detalleNuevo.ConceptoNominaId = novedad.ConceptoNominaId
				detalleNuevo.DiasEspecificos = honorarios[0].DiasEspecificos
				detalleNuevo.DiasLiquidados = honorarios[0].DiasLiquidados
				detalleNuevo.ContratoPreliquidacionId = &contratoPreliquidacion[0]

				//Dependiendo de si es fijo o porcentual calcula el valor de la novedad
				if novedad.ConceptoNominaId.TipoConceptoNominaId == 419 {
					detalleNuevo.ValorCalculado = novedad.Valor
				} else if novedad.ConceptoNominaId.TipoConceptoNominaId == 420 {
					detalleNuevo.ValorCalculado = math.Round((honorarios[0].ValorCalculado * (novedad.Valor / 100)))
				}

				if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/", "POST", &aux, detalleNuevo); err == nil {
					fmt.Println("Concepto Añadido")
				} else {
					fmt.Println("Error al agregar concepto", err)
					return "Error al agregar concepto", err
				}
				//Agregar la novedad a los detalles de esa preliquidacion
			} else {
				fmt.Println("Error al obtener el valor de los honorarios ", err)
				return "Error al obtener el valor de los honorarios ", err
			}
		} else {
			fmt.Println("Error al intentar obtener el id del contrato_preliquidación ", err)
			return "Error al intentar obtener el id del contrato_preliquidación ", err
		}
		auxCuotas = auxCuotas - 1

		if auxCuotas == 0 {
			break
		} else {
			if mesIterativo == 12 {
				mesIterativo = 1
				anoIterativo = anoIterativo + 1
			} else {
				mesIterativo = mesIterativo + 1
			}
		}
	}
	return "No se generó ningún error", nil
}

func EliminarValorNovedad(novedad models.Novedad, fecha_actual time.Time) (mensaje string, err error) {

	var aux map[string]interface{}
	var contratoPreliquidacion []models.ContratoPreliquidacion
	var detalle []models.DetallePreliquidacion

	mesFin := int(novedad.FechaFin.Month())

	for i := int(fecha_actual.Month()); i <= mesFin; i++ {
		var query = "ContratoId.NumeroContrato:" + novedad.ContratoId.NumeroContrato + ",PreliquidacionId.Ano:" + strconv.Itoa(fecha_actual.Year()) + ",PreliquidacionId.Mes:" + strconv.Itoa(i)
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil { //obtiene el contrato_preliquidacion de ese mes y año
			LimpiezaRespuestaRefactor(aux, &contratoPreliquidacion)
			//Eliminar el Concepto de la liquidacion de ese mes
			query := "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId.Id:" + strconv.Itoa(novedad.ConceptoNominaId.Id)
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &detalle)
				if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalle[0].Id), "DELETE", &aux, nil); err == nil {
					fmt.Println("Detalle de novedad eliminado con éxtio")
					//Actualizar fecha de finalización de la novedad
					novedad.FechaFin = fecha_actual
					novedad.Activo = false
					if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/novedad/"+strconv.Itoa(novedad.Id), "PUT", &aux, novedad); err == nil {
						fmt.Println("Novedad Actualizada")
					} else {
						fmt.Println("Error al actualizar novedad: ", err)
						return "Error al actualizar novedad: ", err
					}
				} else {
					fmt.Println("Error al eliminar Detalle de novedad: ", err)
					return "Error al eliminar Detalle de novedad: ", err
				}
			} else {
				fmt.Println("Error al obtener el detalle de la novedad ", err)
				return "Error al obtener el detalle de la novedad ", err
			}
		} else {
			fmt.Println("Error al intentar obtener el id del contrato_preliquidación ", err)
			return "Error al intentar obtener el id del contrato_preliquidación ", err
		}
	}
	return "", err
}

//Agregar novedades al contrato general del docente

func AgregarNovedadGeneral(documento string, mes int, ano int) (mensaje string, err error) {
	var aux map[string]interface{}
	var contratoGeneral []models.Contrato
	var detalleGeneral []models.DetallePreliquidacion
	var detallesTotales []models.DetallePreliquidacion

	//Traer el contrato general del mes a iterar
	query := "NumeroContrato:GENERAL" + strconv.Itoa(mes) + ",Vigencia:" + strconv.Itoa(ano) + ",Documento:" + documento
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1"+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contratoGeneral)
		if contratoGeneral[0].Id != 0 {
			//Traer los detalles de ese contrato:
			query := "ContratoPreliquidacionId.ContratoId.Id:" + strconv.Itoa(contratoGeneral[0].Id) + "," + documento
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1"+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &detalleGeneral)
				if detalleGeneral[0].Id != 0 {
					//Traer los detalles de todos los contratos de ese docente para ese mes en específico
					query := "ContratoPreliquidacionId.ContratoId.Documento:" + documento + ",ContratoPreliquidacionId.PreliquidacionId.Mes:" + strconv.Itoa(mes) + ",ContratoPreliquidacionId.PreliquidacionId.Ano:" + strconv.Itoa(ano)
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1"+query, &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &detallesTotales)
						if detallesTotales[0].Id != 0 {
							//Verificar cuáles conceptos son novedades y agregar su valor al contrato

						} else {
							fmt.Println("No hay detalles para ese contrato")
							return "No hay detalles para ese contrato", nil
						}
					} else {
						fmt.Println("Error al obtener detalles del contrato general", err)
						return "Error al obtener detalles del contrato general", err
					}

				} else {
					fmt.Println("No hay detalles para ese contrato")
					return "No hay detalles para ese contrato", nil
				}
			} else {
				fmt.Println("Error al obtener detalles del contrato general", err)
				return "Error al obtener detalles del contrato general", err
			}

		} else {
			fmt.Println("Error al agregar el concepto al contrato general, no hay conrato general para este docente")
			return "Error al agregar el concepto al contrato general, no hay conrato general para este docente", nil
		}
	} else {
		fmt.Println("Error al obtener el contrato generla", err)
		return "Error al obtener el contrato general", err
	}

	return "Novedades generales actualizadas", nil
}
