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
// @Success 201 {object} models.Detalle
// @Failure 403 body is empty
// @router /obtener_detalle_CT/:ano/:mes/:contrato/:vigencia/:documento [get]
func (c *DetallePreliquidacionController) ObtenerDetalleCT() {

	ano := c.Ctx.Input.Param(":ano")
	mes := c.Ctx.Input.Param(":mes")
	contrato := c.Ctx.Input.Param(":contrato")
	vigencia := c.Ctx.Input.Param(":vigencia")
	documento := c.Ctx.Input.Param(":documento")

	detalle, err := TraerDetalleMensual(ano, mes, contrato, vigencia, documento, true)

	if err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": detalle}
	} else {
		c.Data["mesaage"] = "Error al obtener el detalle, no hay "
		c.Abort("400")
	}

	c.ServeJSON()
}

func TraerDetalleMensual(ano, mes, contrato, vigencia, documento string, CPS bool) (detalle models.Detalle, err error) {
	var aux map[string]interface{}
	var tempDetalle []models.DetallePreliquidacion
	var query = "ContratoPreliquidacionId.PreliquidacionId.Ano:" + ano + ",ContratoPreliquidacionId.PreliquidacionId.Mes:" + mes + ",ContratoPreliquidacionId.ContratoId.NumeroContrato:" + contrato + ",ContratoPreliquidacionId.ContratoId.Vigencia:" + vigencia + ",ContratoPreliquidacionId.ContratoId.Documento:" + documento
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &tempDetalle)
		detalle.Contrato = tempDetalle[0].ContratoPreliquidacionId.ContratoId.NumeroContrato
		detalle.Vigencia = tempDetalle[0].ContratoPreliquidacionId.ContratoId.Vigencia
		for i := 0; i < len(tempDetalle); i++ {
			if tempDetalle[i].ConceptoNominaId.Id == 574 || tempDetalle[i].ConceptoNominaId.Id == 573 {
				fmt.Println("Salto")
			} else if tempDetalle[i].ConceptoNominaId.NaturalezaConceptoNominaId == 424 {
				if CPS && (tempDetalle[i].ConceptoNominaId.Id == 568 || tempDetalle[i].ConceptoNominaId.Id == 569 || tempDetalle[i].ConceptoNominaId.Id == 570) {
					detalle.Detalle = append(detalle.Detalle, tempDetalle[i])
				} else {
					detalle.Detalle = append(detalle.Detalle, tempDetalle[i])
					detalle.TotalDescuentos = detalle.TotalDescuentos + tempDetalle[i].ValorCalculado
				}
			} else if tempDetalle[i].ConceptoNominaId.NaturalezaConceptoNominaId == 423 {
				detalle.Detalle = append(detalle.Detalle, tempDetalle[i])
				detalle.TotalDevengado = detalle.TotalDevengado + tempDetalle[i].ValorCalculado
			} else {
				detalle.Detalle = append(detalle.Detalle, tempDetalle[i])
			}
		}
		detalle.TotalPago = detalle.TotalDevengado - detalle.TotalDescuentos
		return detalle, nil
	} else {
		fmt.Println("Error al obtener detalle ", err)
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
			auxDetalle, err := TraerDetalleMensual(ano, mes, contratoPreliquidacion[i].ContratoId.NumeroContrato, strconv.Itoa(contratoPreliquidacion[i].ContratoId.Vigencia), documento, false)
			if err == nil {
				detallesDVE = append(detallesDVE, auxDetalle)

				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": detallesDVE}
			} else {
				c.Data["mesaage"] = "Error al obtener detalle de 1 o más contratos"
				c.Abort("404")
				break
			}
		}

	} else {
		fmt.Println("Error al obtener detalle ", err)
		c.Data["message"] = "No existe un contrato asociado al documento que sea vigente para la preliquidación " + err.Error()
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
	var res map[string]interface{}
	var mesIterativo = int(novedad.FechaInicio.Month())
	var anoIterativo = novedad.FechaInicio.Year()
	var auxCuotas int
	var contratoPreliquidacion []models.ContratoPreliquidacion
	var descuentos []models.DetallePreliquidacion
	var totalAPagar []models.DetallePreliquidacion
	var honorarios []models.DetallePreliquidacion
	var detalleNuevo models.DetallePreliquidacion
	auxCuotas = novedad.Cuotas

	for { //itera desde el mes en el que se aplicó la novedad hasta el fin del numero de cuotas

		fmt.Println("Mes: ", mesIterativo)
		fmt.Println("Ano: ", anoIterativo)
		fmt.Println("Cuotas: ", auxCuotas)
		var query = "ContratoId.Id:" + strconv.Itoa(novedad.ContratoId.Id) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo)
		fmt.Println(beego.AppConfig.String("UrlTitanCrud") + "/contrato_preliquidacion?limit=-1&query=" + query)
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &res); err == nil { //obtiene el contrato_preliquidacion de ese mes y año
			LimpiezaRespuestaRefactor(res, &contratoPreliquidacion)
			//Para descuentos
			if novedad.ConceptoNominaId.NaturalezaConceptoNominaId == 424 {
				//Obtener el valor de los honorarios para ese mes
				query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:87"
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &res); err == nil {
					LimpiezaRespuestaRefactor(res, &honorarios)
					//Actualizar los descuentos
					query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:573"
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &res); err == nil {
						LimpiezaRespuestaRefactor(res, &descuentos)
						if novedad.ConceptoNominaId.TipoConceptoNominaId == 419 {
							descuentos[0].ValorCalculado = descuentos[0].ValorCalculado + novedad.Valor
							detalleNuevo.ValorCalculado = novedad.Valor
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(descuentos[0].Id), "PUT", &res, descuentos[0]); err == nil {
								fmt.Println("Descuentos actualizados")
							} else {
								fmt.Println("Error al actualizar descuentos ", err)
								return "Error al actualizar descuentos ", err
							}
						} else if novedad.ConceptoNominaId.TipoConceptoNominaId == 420 {
							descuentos[0].ValorCalculado = (descuentos[0].ValorCalculado + (honorarios[0].ValorCalculado * (novedad.Valor / 100)))
							detalleNuevo.ValorCalculado = math.Round(honorarios[0].ValorCalculado * (novedad.Valor / 100))
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(descuentos[0].Id), "PUT", &res, descuentos[0]); err == nil {
								fmt.Println("Descuentos actualizados")
							} else {
								fmt.Println("Error al actualizar descuentos ", err)
								return "Error al actualizar descuentos ", err
							}
						}
					} else {
						fmt.Println("Error al obtener el valor de los descuentos ", err)
					}
					//Obtener y actualizar el total a pagar
					query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:574"
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &res); err == nil {
						LimpiezaRespuestaRefactor(res, &totalAPagar)
						if novedad.ConceptoNominaId.TipoConceptoNominaId == 419 {
							totalAPagar[0].ValorCalculado = totalAPagar[0].ValorCalculado - novedad.Valor
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(totalAPagar[0].Id), "PUT", &res, totalAPagar[0]); err == nil {
								fmt.Println("Total a pagar actualizado")
							} else {
								fmt.Println("Error al actualizar total a pagar ", err)
								return "Error al actualizar total a pagar ", err
							}
						} else if novedad.ConceptoNominaId.TipoConceptoNominaId == 420 {
							totalAPagar[0].ValorCalculado = totalAPagar[0].ValorCalculado - math.Round(honorarios[0].ValorCalculado*(novedad.Valor/100))
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(totalAPagar[0].Id), "PUT", &res, totalAPagar[0]); err == nil {
								fmt.Println("Total a pagar actualizado")
							} else {
								fmt.Println("Error al actualizar total a pagar ", err)
								return "Error al actualizar total a pagar ", err
							}
						}
					} else {
						fmt.Println("Error al obtener el total a pagar ", err)
						return "Error al obtener el total a pagar ", err
					}

					//Agregar la novedad a los detalles de esa preliquidacion
					detalleNuevo.Id = 0
					detalleNuevo.TipoPreliquidacionId = 397
					detalleNuevo.Activo = true
					detalleNuevo.EstadoDisponibilidadId = 426
					detalleNuevo.ConceptoNominaId = novedad.ConceptoNominaId
					detalleNuevo.DiasEspecificos = honorarios[0].DiasEspecificos
					detalleNuevo.DiasLiquidados = honorarios[0].DiasLiquidados
					detalleNuevo.ContratoPreliquidacionId = &contratoPreliquidacion[0]

					if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/", "POST", &res, detalleNuevo); err == nil {
						fmt.Println("Concepto Añadido")
					} else {
						fmt.Println("Error al agregar concepto", err)
						return "Error al agregar concepto", err
					}

				} else {
					fmt.Println("Error al obtener el valor de los honorarios ", err)
					return "Error al obtener el valor de los honorarios ", err
				}
				//Para devengos
			} else if novedad.ConceptoNominaId.NaturalezaConceptoNominaId == 423 {
				//Obtener el valor de los honorarios para ese mes
				query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:87"
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &res); err == nil {
					LimpiezaRespuestaRefactor(res, &honorarios)
					//Actualizar el total a pagar
					query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:574"
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &res); err == nil {
						LimpiezaRespuestaRefactor(res, &totalAPagar)
						if novedad.ConceptoNominaId.TipoConceptoNominaId == 419 {
							totalAPagar[0].ValorCalculado = totalAPagar[0].ValorCalculado + novedad.Valor
							detalleNuevo.ValorCalculado = novedad.Valor
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(totalAPagar[0].Id), "PUT", &res, totalAPagar[0]); err == nil {
								fmt.Println("Total a pagar actualizado")
							} else {
								fmt.Println("Error al actualizar total a pagar ", err)
								return "Error al actualizar total a pagar ", err
							}
						} else if novedad.ConceptoNominaId.TipoConceptoNominaId == 420 {
							totalAPagar[0].ValorCalculado = totalAPagar[0].ValorCalculado + (honorarios[0].ValorCalculado * (novedad.Valor / 100))
							detalleNuevo.ValorCalculado = (honorarios[0].ValorCalculado * (novedad.Valor / 100))
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(totalAPagar[0].Id), "PUT", &res, totalAPagar[0]); err == nil {
								fmt.Println("Total a pagar actualizado")
							} else {
								fmt.Println("Error al actualizar total a pagar ", err)
								return "Error al actualizar total a pagar ", err
							}
						}
					} else {
						fmt.Println("Error al obtener el total a pagar ", err)
						return "Error al actualizar total a pagar ", err
					}

					//Agregar la novedad a los detalles de esa preliquidacion
					detalleNuevo.Id = 0
					detalleNuevo.TipoPreliquidacionId = 397
					detalleNuevo.Activo = true
					detalleNuevo.EstadoDisponibilidadId = 426
					detalleNuevo.ConceptoNominaId = novedad.ConceptoNominaId
					detalleNuevo.DiasEspecificos = honorarios[0].DiasEspecificos
					detalleNuevo.DiasLiquidados = honorarios[0].DiasLiquidados
					detalleNuevo.ContratoPreliquidacionId = &contratoPreliquidacion[0]

					if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/", "POST", &res, detalleNuevo); err == nil {
						fmt.Println("Concepto Añadido")
					} else {
						fmt.Println("Error al agregar concepto", err)
						return "Error al agregar concepto", err
					}

				} else {
					fmt.Println("Error al obtener el valor de los honorarios ", err)
					return "Error al obtener el valor de los honorarios ", err
				}
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
	return "", nil
}

func EliminarValorNovedad(novedad models.Novedad, fecha_actual time.Time) (mensaje string, err error) {

	var aux map[string]interface{}
	var contratoPreliquidacion []models.ContratoPreliquidacion
	var descuentos []models.DetallePreliquidacion
	var totalAPagar []models.DetallePreliquidacion
	var honorarios []models.DetallePreliquidacion
	var detalle []models.DetallePreliquidacion

	var mesFin int
	if fecha_actual.Year() == novedad.FechaFin.Year() {
		mesFin = int(novedad.FechaFin.Month())
	} else {
		mesFin = 12
	}

	for i := int(fecha_actual.Month()); i <= mesFin; i++ {
		var query = "ContratoId.NumeroContrato:" + novedad.ContratoId.NumeroContrato + ",PreliquidacionId.Ano:" + strconv.Itoa(fecha_actual.Year()) + ",PreliquidacionId.Mes:" + strconv.Itoa(i)
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil { //obtiene el contrato_preliquidacion de ese mes y año
			LimpiezaRespuestaRefactor(aux, &contratoPreliquidacion)
			if novedad.ConceptoNominaId.NaturalezaConceptoNominaId == 423 {
				//Obtener el valor de los honorarios para ese mes
				query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:87"
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil { //obtiene los descuentos aplicados para ese mes para sumarle el de la novedad
					LimpiezaRespuestaRefactor(aux, &honorarios)
					//Actualizar el total a pagar
					query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:574"
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil { //obtiene los descuentos aplicados para ese mes para sumarle el de la novedad
						LimpiezaRespuestaRefactor(aux, &totalAPagar)
						if novedad.ConceptoNominaId.TipoConceptoNominaId == 419 {
							totalAPagar[0].ValorCalculado = totalAPagar[0].ValorCalculado - novedad.Valor
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(totalAPagar[0].Id), "PUT", &aux, totalAPagar[0]); err == nil {
								fmt.Println("Total a pagar actualizado")
							} else {
								fmt.Println("Error al actualizar total a pagar ", err)
								return "Error al actualizar total a pagar  ", err
							}
						} else if novedad.ConceptoNominaId.TipoConceptoNominaId == 420 {
							totalAPagar[0].ValorCalculado = totalAPagar[0].ValorCalculado - (honorarios[0].ValorCalculado * (novedad.Valor / 100))
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(totalAPagar[0].Id), "PUT", &aux, totalAPagar[0]); err == nil {
								fmt.Println("Total a pagar actualizado")
							} else {
								fmt.Println("Error al actualizar total a pagar ", err)
								return "Error al actualizar total a pagar  ", err
							}
						}
					} else {
						fmt.Println("Error al obtener el total a pagar ", err)
						return "Error al obtener el total a pagar ", err
					}
				} else {
					fmt.Println("Error al obtener el valor de los honorarios ", err)
					return "Error al obtener el valor de los honorarios ", err
				}
			} else if novedad.ConceptoNominaId.NaturalezaConceptoNominaId == 424 {
				//Obtener el valor de los honorarios para ese mes
				query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:87"
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil { //obtiene los descuentos aplicados para ese mes para sumarle el de la novedad
					LimpiezaRespuestaRefactor(aux, &honorarios)
					//Actualizar los descuentos
					query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:573"
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil { //obtiene los descuentos aplicados para ese mes para sumarle el de la novedad
						LimpiezaRespuestaRefactor(aux, &descuentos)
						if novedad.ConceptoNominaId.TipoConceptoNominaId == 419 {
							descuentos[0].ValorCalculado = descuentos[0].ValorCalculado - novedad.Valor
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(descuentos[0].Id), "PUT", &aux, descuentos[0]); err == nil {
								fmt.Println("Descuentos actualizados")
							} else {
								fmt.Println("Error al actualizar descuentos ", err)
								return "Error al actualizar descuentos ", err
							}
						} else if novedad.ConceptoNominaId.TipoConceptoNominaId == 420 {
							descuentos[0].ValorCalculado = (descuentos[0].ValorCalculado - (honorarios[0].ValorCalculado * (novedad.Valor / 100)))
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(descuentos[0].Id), "PUT", &aux, descuentos[0]); err == nil {
								fmt.Println("Descuentos actualizados")
							} else {
								fmt.Println("Error al actualizar descuentos ", err)
								return "Error al actualizar descuentos ", err
							}
						}
					} else {
						fmt.Println("Error al obtener el valor de los descuentos ", err)
						return "Error al obtener el valor de los descuentos ", err
					}

					query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:574"
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil { //obtiene los descuentos aplicados para ese mes para sumarle el de la novedad

						LimpiezaRespuestaRefactor(aux, &totalAPagar)
						if novedad.ConceptoNominaId.TipoConceptoNominaId == 419 {
							totalAPagar[0].ValorCalculado = totalAPagar[0].ValorCalculado + novedad.Valor
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(totalAPagar[0].Id), "PUT", &aux, totalAPagar[0]); err == nil {
								fmt.Println("Total a pagar actualizado")
							} else {
								fmt.Println("Error al actualizar total a pagar ", err)
								return "Error al actualizar total a pagar ", err
							}
						} else if novedad.ConceptoNominaId.TipoConceptoNominaId == 420 {
							totalAPagar[0].ValorCalculado = totalAPagar[0].ValorCalculado + (honorarios[0].ValorCalculado * (novedad.Valor / 100))
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(totalAPagar[0].Id), "PUT", &aux, totalAPagar[0]); err == nil {
								fmt.Println("Total a pagar actualizado")
							} else {
								fmt.Println("Error al actualizar total a pagar ", err)
								return "Error al actualizar total a pagar ", err
							}
						}

					} else {
						fmt.Println("Error al obtener el total a pagar ", err)
						return "Error al obtener el total a pagar ", err
					}

					//Eliminar el Detalle del concepto
					query := "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:" + strconv.Itoa(novedad.ConceptoNominaId.Id)
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &detalle)
						if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalle[0].Id), "DELETE", &aux, nil); err == nil {
							fmt.Println("Detalle de novedad eliminado con éxtio")
						} else {
							fmt.Println("Error al eliminar Detalle de novedad: ", err)
							return "Error al eliminar Detalle de novedad: ", err
						}
					} else {
						fmt.Println("Error al obtener el detalle de la novedad ", err)
						return "Error al obtener el detalle de la novedad ", err
					}
					//Actualizar fecha de finalización de la novedad
					novedad.FechaFin = time.Now()
					novedad.Activo = false
					fmt.Println("novedad: ", novedad)
					if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/novedad/"+strconv.Itoa(novedad.Id), "PUT", &aux, novedad); err == nil {
						fmt.Println("Novedad Actualizada")
					} else {
						fmt.Println("Error al actualizar novedad: ", err)
						return "Error al actualizar novedad: ", err
					}
				} else {
					fmt.Println("Error al obtener el valor de los honorarios ", err)
					return "Error al obtener el valor de los honorarios ", err
				}
			}

		} else {
			fmt.Println("Error al intentar obtener el id del contrato_preliquidación ", err)
			return "Error al intentar obtener el id del contrato_preliquidación ", err
		}
	}
	return "", err
}
