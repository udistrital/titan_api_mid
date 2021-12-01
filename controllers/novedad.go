package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

type NovedadController struct {
	beego.Controller
}

func (c *NovedadController) URLMapping() {
	c.Mapping("AgregarNovedad", c.AgregarNovedad)
	c.Mapping("EliminarNovedad", c.EliminarNovedad)
	c.Mapping("CancelarContrato", c.CancelarContrato)
	c.Mapping("CederContrato", c.CederContrato)
}

// Post ...
// @Title Agregar Novedad
// @Description Agregar Novedad a un contrato y liquidar nuevamente
// @Param	novedad		body 	models.Novedad 	true	"Cuerpo de la novedad a guardar"
// @Success 201 {object} models.Novedad
// @Failure 403 body is empty
// @router /agregar_novedad [post]
func (c *NovedadController) AgregarNovedad() {
	var aux map[string]interface{}
	var auxNovedad []models.Novedad
	var contratoPreliquidacion []models.ContratoPreliquidacion
	var contrato []models.Contrato
	var auxDetalle []models.DetallePreliquidacion
	var auxConcepto []models.ConceptoNomina
	var honorarios float64
	var descuentos float64
	var novedad models.Novedad
	var fecha_actual = time.Now()
	var anoFin int
	var mesFin int
	var iteracion int
	var posible bool
	posible = true
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &novedad); err == nil {
		fmt.Println(novedad)
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/concepto_nomina?limit=-1&query=Id:"+strconv.Itoa(novedad.ConceptoNominaId.Id), &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &auxConcepto)
			fmt.Println(auxConcepto[0])
		} else {
			fmt.Println("Error al obtener el concepto: ", err)
		}

		//Si es devengo
		if auxConcepto[0].NaturalezaConceptoNominaId == 423 {
			if int(fecha_actual.Month())+novedad.Cuotas-1 > 12 {
				fmt.Println(int(fecha_actual.Month()))
				mesFin = int(fecha_actual.Month()) + novedad.Cuotas - 13
				anoFin = fecha_actual.Year() + 1
			} else {
				mesFin = int(fecha_actual.Month()) + novedad.Cuotas - 1
				anoFin = fecha_actual.Year()
			}
			novedad.FechaInicio = fecha_actual
			novedad.FechaFin = time.Date(anoFin, time.Month(mesFin), 30, 0, 0, 0, 0, time.UTC)
			if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/novedad", "POST", &aux, novedad); err == nil {
				LimpiezaRespuestaRefactor(aux, &novedad)
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/novedad?limit=-1&query=Id:"+strconv.Itoa(novedad.Id), &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &auxNovedad)
					novedad = auxNovedad[0]
					if novedad.ConceptoNominaId.TipoConceptoNominaId == 419 || novedad.ConceptoNominaId.TipoConceptoNominaId == 420 {
						//Agregar el Valor al detalle
						AgregarValorNovedad(novedad)
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": novedad}
					}
				}
			} else {
				fmt.Println("No se pudo guardar la novedad", err)
				c.Data["message"] = "Error service Delete: Request contains incorrect parameter"
			}
			//si es desuento
		} else if auxConcepto[0].NaturalezaConceptoNominaId == 424 {
			if int(fecha_actual.Month())+novedad.Cuotas-1 > 12 {
				mesFin = int(fecha_actual.Month()) + novedad.Cuotas - 13
				anoFin = fecha_actual.Year() + 1
			} else {
				mesFin = int(fecha_actual.Month()) + novedad.Cuotas - 1
				anoFin = fecha_actual.Year()
			}
			novedad.FechaInicio = fecha_actual
			novedad.FechaFin = time.Date(anoFin, time.Month(mesFin), 30, 0, 0, 0, 0, time.Local)
			//Verificar que las cuotas no se pasen del tiempo restante del contrato
			//Obtener el contrato

			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=Id:"+strconv.Itoa(novedad.ContratoId.Id), &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &contrato)
			} else {
				fmt.Println("Error al obtener el contrato: ", err)
			}
			fmt.Println(novedad.Cuotas)

			if contrato[0].FechaFin.Year() == novedad.FechaInicio.Year() {
				if int(contrato[0].FechaFin.Month())-int(fecha_actual.Month()) < novedad.Cuotas {
					c.Data["message"] = "Imposible agregar la novedad, las cuotas superan los meses restantes"
					fmt.Println("Supera los meses año igual")
					posible = false
				}
			} else {

				if int(contrato[0].FechaFin.Month())+12-int(fecha_actual.Month()) < novedad.Cuotas {
					c.Data["message"] = "Imposible agregar la novedad, las cuotas superan los meses restantes"
					fmt.Println("Supera los meses año diferente")
					posible = false
				}
			}

			iteracion = int(novedad.FechaInicio.Month())
			for iteracion < int(novedad.FechaInicio.Month())+novedad.Cuotas {
				//Obtener el valor de los honorarios de ese mes
				var query = "ContratoId:" + strconv.Itoa(novedad.ContratoId.Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(iteracion) + ",PreliquidacionId.Ano:" + strconv.Itoa(fecha_actual.Year())
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contratoPreliquidacion)
					query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:87"
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &auxDetalle)
						honorarios = auxDetalle[0].ValorCalculado
					} else {
						fmt.Println("Error al obtener el valor de los honorarios: ", err)
					}
				} else {
					fmt.Println("Error al obtener el contrato peliquidacion: ", err)
				}
				//Obtener el Valor de los descuentos de ese mes
				query = "ContratoId:" + strconv.Itoa(novedad.ContratoId.Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(iteracion) + ",PreliquidacionId.Ano:" + strconv.Itoa(fecha_actual.Year())
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contratoPreliquidacion)
					query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:573"
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &auxDetalle)
						descuentos = auxDetalle[0].ValorCalculado
					} else {
						fmt.Println("Error al obtener el valor de los descuentos: ", err)
					}
				} else {
					fmt.Println("Error al obtener el contrato peliquidacion: ", err)
				}
				//Verificar que el valor de los descuentos no supera la mitad de los honorarios
				if auxConcepto[0].TipoConceptoNominaId == 419 {
					if descuentos+novedad.Valor > (honorarios / 2) {
						posible = false
					}
				} else if auxConcepto[0].TipoConceptoNominaId == 420 {
					if descuentos+((novedad.Valor/100)*honorarios) > (honorarios / 2) {
						posible = false
					}
				}

				if posible {
					if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/novedad", "POST", &aux, novedad); err == nil {
						fmt.Println("Novedad Registrada ")
						LimpiezaRespuestaRefactor(aux, &novedad)
						if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/novedad?limit=-1&query=Id:"+strconv.Itoa(novedad.Id), &aux); err == nil {
							LimpiezaRespuestaRefactor(aux, &auxNovedad)
							novedad = auxNovedad[0]
							if novedad.ConceptoNominaId.TipoConceptoNominaId == 419 || novedad.ConceptoNominaId.TipoConceptoNominaId == 420 {
								//Agregar el Valor al detalle
								AgregarValorNovedad(novedad)
								c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": novedad}
							}
						}
					} else {
						fmt.Println("No se pudo guardar la novedad", err)
						c.Data["message"] = "Error service Delete: Request contains incorrect parameter"
					}
				} else {
					fmt.Println("No se pudo guardar la novedad", err)
					c.Data["message"] = "Imposible agregar la novedad, los descuentos sobrepasan la mitad de los honorarios"
				}
				if iteracion == 12 {
					break
				} else {
					iteracion = iteracion + 1
				}
			}
		}
	} else {
		fmt.Println("Error al Unmarshal de novedad: ", err)
		c.Data["message"] = "Error service Delete: Request contains incorrect parameter"
	}
	c.ServeJSON()
}

// Get ...
// @Title Eliminar Novedad
// @Description Eliminar Novedad Novedad a un contrato y liquidar nuevamente
// @Param	id		path 	true	"Id de la novedad que se va retirar"
// @Success 201 models.Novedad.Id
// @Failure 403 body is empty
// @router /eliminar_novedad/:id [get]
func (c *NovedadController) EliminarNovedad() {
	var id = c.Ctx.Input.Param(":id")
	var fecha_actual = time.Now()
	var aux map[string]interface{}
	var novedad []models.Novedad
	//Verificar cómo se van a enviar los datos del contrato al trabajar el front
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/novedad?limit=-1&query=Id:"+id, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &novedad)
		if novedad[0].ConceptoNominaId.TipoConceptoNominaId == 419 || novedad[0].ConceptoNominaId.TipoConceptoNominaId == 420 {
			EliminarValorNovedad(novedad[0], fecha_actual)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": novedad[0].Id}
		} else {
			c.Data["mesaage"] = "Error service Delete: Request contains incorrect parameter"
			c.Abort("404")
		}
	} else {
		fmt.Println("Error al unmarshal de la novedad ", err)
		c.Data["mesaage"] = "Error service Delete: Request contains incorrect parameter"
		c.Abort("404")
	}
}

// Get ...
// @Title Registrar Cancelacion
// @Description Maneja la novedad contractual de cancelación
// @Param	NumeroContrato		path 	true	"Numero del contrato que se va a cancelar"
// @Success 201 {object} models.Contrato
// @Failure 403 body is empty
// @router /cancelar_contrato/:NumeroContrato [get]
func (c *NovedadController) CancelarContrato() {
	var numero = c.Ctx.Input.Param(":NumeroContrato")
	var fecha_actual = time.Now()
	var aux map[string]interface{}
	var contrato []models.Contrato
	var contrato_preliquidacion []models.ContratoPreliquidacion
	var valorDia float64
	var detalles []models.DetallePreliquidacion
	var mesInicio = int(fecha_actual.Month())
	var mesFin int

	//Traer el contrato a cancelar
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=NumeroContrato:"+numero, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contrato)
	} else {
		fmt.Println("Error al obtener el contrato: ", err)
	}

	if contrato[0].FechaInicio.Year() != contrato[0].FechaFin.Year() {
		mesFin = 12
	} else {
		mesFin = int(contrato[0].FechaFin.Month())
	}

	//Eliminar los detalles y los contratos_preliquidacion
	for i := mesInicio; i <= mesFin; i++ {
		//Obtener contrat_preliquidacion para ese mes
		query := "ContratoId:" + strconv.Itoa(contrato[0].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(i) + ",PreliquidacionId.Ano:" + strconv.Itoa(fecha_actual.Year())
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &detalles)
				for j := 0; j < len(detalles); j++ {
					if detalles[j].ConceptoNominaId.Id == 87 && i == mesInicio {
						valorDia = detalles[j].ValorCalculado / detalles[j].DiasLiquidados
					}
					if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalles[j].Id), "DELETE", &aux, nil); err == nil {
						fmt.Println("Detalle eliminado con exito")
					} else {
						fmt.Println("Error al eliminar detalle: ", err)
					}
				}
				if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[0].Id), "DELETE", &aux, nil); err == nil {
					fmt.Println("contrato preliquidacion eliminado con exito")
				} else {
					fmt.Println("Error al eliminar contrato preliquidacion: ", err)
				}
			} else {
				fmt.Println("Error al obtener detalles")
			}
		} else {
			fmt.Println("Error al obtener contrato_preliquidacion")
		}
	}

	contrato[0].FechaFin = time.Date(fecha_actual.Year(), fecha_actual.Month(), fecha_actual.Day(), 12, 0, 0, 0, time.UTC)
	//Actualizar feca de finalización del contrato
	if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato[0].Id), "PUT", &aux, contrato[0]); err == nil {
		fmt.Println("Contrato Actualizado")
	} else {
		fmt.Println("Error al ctualizar el contrato: ", err)
	}

	if contrato[0].FechaInicio.Month() != fecha_actual.Month() {
		contrato[0].FechaInicio = time.Date(fecha_actual.Year(), fecha_actual.Month(), 1, 12, 0, 0, 0, time.UTC)
	}

	if fecha_actual.Day() == 31 {
		contrato[0].ValorContrato = (valorDia * 30)
	} else {
		contrato[0].ValorContrato = valorDia * float64(fecha_actual.Day())
	}

	if contrato[0].TipoNominaId == 411 {
		liquidarCPS(contrato[0], 0)
	} else if contrato[0].TipoNominaId == 409 {
		liquidarHCH(contrato[0])
	} else if contrato[0].TipoNominaId == 410 {

	}
}

// Post ...
// @Title Ceder Contrato
// @Description Maneja la novedad contractual de cesión de contrato
// @Param	Sucesor		body 	models.Sucesor 	true	"Datos del sucesor del contrato"
// @Success 201 {object} models.Contrato
// @Failure 403 body is empty
// @router /ceder_contrato [post]
func (c *NovedadController) CederContrato() {
	var sucesor models.Sucesor
	var aux map[string]interface{}
	var contrato []models.Contrato
	var contratoNuevo models.Contrato
	var contrato_preliquidacion []models.ContratoPreliquidacion
	var detalles []models.DetallePreliquidacion
	var valorDia float64
	var valorNuevo float64
	var valorViejo float64
	var mesInicio int
	var mesFin int
	var diasRestantes int
	var diasCompletos int
	var diasLiquidados int

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &sucesor); err == nil {
		mesInicio = int(sucesor.FechaInicio.Month())

		//Traer el contrato a cancelar
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=NumeroContrato:"+sucesor.NumeroContrato, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contrato)
			contratoNuevo = contrato[0]
			valorViejo = contrato[0].ValorContrato

			if contrato[0].TipoNominaId == 411 {
				diasCompletos, _ = calcularDiasContratoCPS(contrato[0].FechaInicio, contrato[0].FechaFin)
			} else if contrato[0].TipoNominaId == 409 {
				//liquidarHCH(contrato[0])
			} else if contrato[0].TipoNominaId == 410 {

			}
		} else {
			fmt.Println("Error al obtener el contrato: ", err)
		}

		if contrato[0].FechaInicio.Year() != contrato[0].FechaFin.Year() {
			mesFin = 12
		} else {
			mesFin = int(contrato[0].FechaFin.Month())
		}

		//Eliminar los detalles y los contratos_preliquidacion
		for i := mesInicio; i <= mesFin; i++ {
			//Obtener contrat_preliquidacion para ese mes
			query := "ContratoId:" + strconv.Itoa(contrato[0].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(i) + ",PreliquidacionId.Ano:" + strconv.Itoa(sucesor.FechaInicio.Year())
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &detalles)
					for j := 0; j < len(detalles); j++ {
						if detalles[j].ConceptoNominaId.Id == 87 && i == mesInicio {
							valorDia = detalles[j].ValorCalculado / detalles[j].DiasLiquidados
						}
						if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalles[j].Id), "DELETE", &aux, nil); err == nil {
							fmt.Println("Detalle eliminado con exito")
						} else {
							fmt.Println("Error al eliminar detalle: ", err)
						}
					}
					if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[0].Id), "DELETE", &aux, nil); err == nil {
						fmt.Println("contrato preliquidacion eliminado con exito")
					} else {
						fmt.Println("Error al eliminar contrato preliquidacion: ", err)
					}
				} else {
					fmt.Println("Error al obtener detalles")
				}
			} else {
				fmt.Println("Error al obtener contrato_preliquidacion")
			}
		}

		//Obtener lo liquidado hasta el momento
		for i := int(contrato[0].FechaInicio.Month()); i < mesInicio; i++ {
			//Obtener contrat_preliquidacion para ese mes
			query := "ContratoId:" + strconv.Itoa(contrato[0].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(i) + ",PreliquidacionId.Ano:" + strconv.Itoa(sucesor.FechaInicio.Year())
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &detalles)
					diasLiquidados = diasLiquidados + int(detalles[0].DiasLiquidados)
					for j := 0; j < len(detalles); j++ {
						if detalles[j].ConceptoNominaId.Id == 87 {
							valorNuevo = valorNuevo + detalles[j].ValorCalculado
						}
					}
					if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[0].Id), "DELETE", &aux, nil); err == nil {
						fmt.Println("contrato preliquidacion eliminado con exito")
					} else {
						fmt.Println("Error al eliminar contrato preliquidacion: ", err)
					}
				} else {
					fmt.Println("Error al obtener detalles")
				}
			} else {
				fmt.Println("Error al obtener contrato_preliquidacion")
			}
		}

		//Agregar el nuevo contrato
		contratoNuevo.Id = 0
		contratoNuevo.Documento = sucesor.Documento
		contratoNuevo.NombreCompleto = sucesor.NombreCompleto
		contratoNuevo.FechaInicio = sucesor.FechaInicio
		contratoNuevo = registratContrato(contratoNuevo)

		contrato[0].FechaFin = sucesor.FechaInicio.Add(24 * time.Hour * -1)
		//Dias Restantes de contrato
		diasLiquidados = diasLiquidados + contrato[0].FechaFin.Day()
		diasRestantes = diasCompletos - diasLiquidados
		//Actualizar fecha de finalización del contrato
		if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato[0].Id), "PUT", &aux, contrato[0]); err == nil {
			fmt.Println("Contrato Actualizado")
		} else {
			fmt.Println("Error al ctualizar el contrato: ", err)
		}
		if contrato[0].FechaInicio.Month() != sucesor.FechaInicio.Month() {
			contrato[0].FechaInicio = time.Date(sucesor.FechaInicio.Year(), sucesor.FechaInicio.Month(), 1, 12, 0, 0, 0, time.UTC)
		}

		if sucesor.FechaInicio.Day() == 31 {
			contrato[0].ValorContrato = valorDia * 30
			contratoNuevo.ValorContrato = valorViejo - valorNuevo
		} else {
			contrato[0].ValorContrato = valorDia * float64(contrato[0].FechaFin.Day())
			contratoNuevo.ValorContrato = valorViejo - valorNuevo - valorDia*float64(contrato[0].FechaFin.Day())
		}
		if contrato[0].TipoNominaId == 411 {
			liquidarCPS(contrato[0], 0)
			liquidarCPS(contratoNuevo, diasRestantes)
		} else if contrato[0].TipoNominaId == 409 {
			liquidarHCH(contrato[0])
		} else if contrato[0].TipoNominaId == 410 {

		}
	} else {
		fmt.Println("Erro al unmarsahl de los datos del sucesor: ", err)
	}
}

// Post ...
// @Title Aplicar Oro sí
// @Description Aplicar novedad contratul de otro sí
// @Param	OtroSí		body 	models.OtroSi 	true	"Datos del sucesor del contrato"
// @Success 201 {object} models.Contrato
// @Failure 403 body is empty
// @router /otrosi_contrato [post]
func (c *NovedadController) AplicarOtrosi() {
	var aux map[string]interface{}
	var otro_si models.OtroSi
	var contrato []models.Contrato
	var contrato_preliquidacion []models.ContratoPreliquidacion
	var detalles []models.DetallePreliquidacion
	var diasNuevos float64
	var valorNuevo float64
	var mesesContrato int
	var valorMensual float64
	var honorarios float64
	var dias int
	var diasRestantes int
	var fechaRespaldo time.Time

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &otro_si); err == nil {
		//Traer el contrato a cancelar
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=NumeroContrato:"+otro_si.NumeroContrato, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contrato)
			fechaRespaldo = contrato[0].FechaFin
		} else {
			fmt.Println("Error al obtener el contrato: ", err)
		}
	} else {
		fmt.Println("Error al unmarshal de los datos del otro sí", err)
	}

	diasNuevos, _ = CalcularDias(contrato[0].FechaFin, otro_si.FechaFin)

	//Calcular Valor mensual del contrato
	_, mesesContrato = calcularDiasContratoCPS(contrato[0].FechaInicio, contrato[0].FechaFin)
	valorMensual = contrato[0].ValorContrato / float64(mesesContrato)

	//Eliminar el detalle del último mes
	query := "ContratoId:" + strconv.Itoa(contrato[0].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(int(contrato[0].FechaFin.Month())) + ",PreliquidacionId.Ano:" + strconv.Itoa(contrato[0].FechaFin.Year())
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &detalles)
			for j := 0; j < len(detalles); j++ {
				if detalles[j].ConceptoNominaId.Id == 87 {
					honorarios = detalles[j].ValorCalculado
					dias = int(detalles[j].DiasLiquidados)
				}
				if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalles[j].Id), "DELETE", &aux, nil); err == nil {
					fmt.Println("Detalle eliminado con exito")
				} else {
					fmt.Println("Error al eliminar detalle: ", err)
				}
			}
			if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[0].Id), "DELETE", &aux, nil); err == nil {
				fmt.Println("contrato preliquidacion eliminado con exito")
			} else {
				fmt.Println("Error al eliminar contrato preliquidacion: ", err)
			}
		} else {
			fmt.Println("Error al obtener detalles")
		}
	} else {
		fmt.Println("Error al obtener contrato_preliquidacion")
	}
	//Acutalizar los valores del contrato
	contrato[0].FechaFin = otro_si.FechaFin
	contrato[0].ValorContrato = contrato[0].ValorContrato + valorMensual*(diasNuevos/30)
	if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato[0].Id), "PUT", &aux, contrato[0]); err == nil {
		fmt.Println("Contrato Actualizado")
	} else {
		fmt.Println("Error al actualizar contrato: ", err)
	}
	//Calcular valor agregado para la liquidacion
	valorNuevo = (valorMensual * (diasNuevos / 30)) + honorarios
	diasRestantes = int(diasNuevos) + dias

	//Ajustar Datos para el la liquidación a partir de ahí
	contrato[0].FechaInicio = time.Date(fechaRespaldo.Year(), fechaRespaldo.Month(), 1, 12, 0, 0, 0, time.UTC)
	contrato[0].FechaFin = otro_si.FechaFin
	contrato[0].ValorContrato = valorNuevo
	fmt.Println("Contrato: ", contrato[0])
	//liquidar el contrato
	liquidarCPS(contrato[0], diasRestantes)

}

// Post ...
// @Title Suspender Contrato
// @Description Maneja la novedad contractual de Suspensión de Contrato
// @Param	Suspension	 body 	true	"Duración de la suspensión"
// @Success 201 {object} models.Contrato
// @Failure 403 body is empty
// @router /suspender_contrato [post]
func (c *NovedadController) SuspenderContrato() {
	var suspension models.Suspension
	var aux map[string]interface{}
	var contrato []models.Contrato
	var contratoNuevo models.Contrato
	var contrato_preliquidacion []models.ContratoPreliquidacion
	var detalles []models.DetallePreliquidacion
	var valorDia float64
	var valorNuevo float64
	var valorViejo float64
	var mesInicio int
	var mesFin int
	var diasRestantes int
	var diasCompletos int
	var diasLiquidados int
	var diasSuspension float64

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &suspension); err == nil {
		mesInicio = int(suspension.FechaInicio.Month())

		//Traer el contrato a cancelar
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=NumeroContrato:"+suspension.NumeroContrato, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contrato)
			contratoNuevo = contrato[0]
			valorViejo = contrato[0].ValorContrato

			if contrato[0].TipoNominaId == 411 {
				diasCompletos, _ = calcularDiasContratoCPS(contrato[0].FechaInicio, contrato[0].FechaFin)
			} else if contrato[0].TipoNominaId == 409 {
				//liquidarHCH(contrato[0])
			} else if contrato[0].TipoNominaId == 410 {

			}
		} else {
			fmt.Println("Error al obtener el contrato: ", err)
		}

		if contrato[0].FechaInicio.Year() != contrato[0].FechaFin.Year() {
			mesFin = 12
		} else {
			mesFin = int(contrato[0].FechaFin.Month())
		}

		//Eliminar los detalles y los contratos_preliquidacion
		for i := mesInicio; i <= mesFin; i++ {
			//Obtener contrat_preliquidacion para ese mes
			query := "ContratoId:" + strconv.Itoa(contrato[0].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(i) + ",PreliquidacionId.Ano:" + strconv.Itoa(suspension.FechaInicio.Year())
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &detalles)
					for j := 0; j < len(detalles); j++ {
						if detalles[j].ConceptoNominaId.Id == 87 && i == mesInicio {
							valorDia = detalles[j].ValorCalculado / detalles[j].DiasLiquidados
						}
						if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalles[j].Id), "DELETE", &aux, nil); err == nil {
							fmt.Println("Detalle eliminado con exito")
						} else {
							fmt.Println("Error al eliminar detalle: ", err)
						}
					}
					if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[0].Id), "DELETE", &aux, nil); err == nil {
						fmt.Println("contrato preliquidacion eliminado con exito")
					} else {
						fmt.Println("Error al eliminar contrato preliquidacion: ", err)
					}
				} else {
					fmt.Println("Error al obtener detalles")
				}
			} else {
				fmt.Println("Error al obtener contrato_preliquidacion")
			}
		}

		//Obtener lo liquidado hasta el momento
		for i := int(contrato[0].FechaInicio.Month()); i < mesInicio; i++ {
			//Obtener contrat_preliquidacion para ese mes
			query := "ContratoId:" + strconv.Itoa(contrato[0].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(i) + ",PreliquidacionId.Ano:" + strconv.Itoa(suspension.FechaInicio.Year())
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &detalles)
					diasLiquidados = diasLiquidados + int(detalles[0].DiasLiquidados)
					for j := 0; j < len(detalles); j++ {
						if detalles[j].ConceptoNominaId.Id == 87 {
							valorNuevo = valorNuevo + detalles[j].ValorCalculado
						}
					}
					if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[0].Id), "DELETE", &aux, nil); err == nil {
						fmt.Println("contrato preliquidacion eliminado con exito")
					} else {
						fmt.Println("Error al eliminar contrato preliquidacion: ", err)
					}
				} else {
					fmt.Println("Error al obtener detalles")
				}
			} else {
				fmt.Println("Error al obtener contrato_preliquidacion")
			}
		}
		//Calcular los días de la suspension
		diasSuspension, _ = CalcularDias(suspension.FechaInicio, suspension.FechaFin)
		contratoNuevo.FechaInicio = suspension.FechaFin
		contratoNuevo.FechaFin = contratoNuevo.FechaFin.Add(24 * time.Hour * time.Duration(diasSuspension))

		contrato[0].FechaFin = suspension.FechaInicio.Add(24 * time.Hour * -1)
		//Dias Restantes de contrato
		diasLiquidados = diasLiquidados + contrato[0].FechaFin.Day()
		diasRestantes = diasCompletos - diasLiquidados
		//Actualizar fecha de finalización del contrato
		if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contratoNuevo.Id), "PUT", &aux, contratoNuevo); err == nil {
			fmt.Println("Contrato Actualizado")
		} else {
			fmt.Println("Error al ctualizar el contrato: ", err)
		}

		//Actualizar Fecha de inicio para liquidar el mes en el que se aplicó la suspensión
		if contrato[0].FechaInicio.Month() != suspension.FechaInicio.Month() {
			contrato[0].FechaInicio = time.Date(suspension.FechaInicio.Year(), suspension.FechaInicio.Month(), 1, 0, 0, 0, 0, time.UTC)
		}

		if suspension.FechaInicio.Day() == 31 {
			contrato[0].ValorContrato = valorDia * 30
			contratoNuevo.ValorContrato = valorViejo - valorNuevo
		} else {
			contrato[0].ValorContrato = valorDia * float64(contrato[0].FechaFin.Day())
			contratoNuevo.ValorContrato = valorViejo - valorNuevo - valorDia*float64(contrato[0].FechaFin.Day())
		}

		if contrato[0].TipoNominaId == 411 {
			liquidarCPS(contrato[0], 0)
			liquidarCPS(contratoNuevo, diasRestantes)
		} else if contrato[0].TipoNominaId == 409 {
			liquidarHCH(contrato[0])
		} else if contrato[0].TipoNominaId == 410 {

		}
	} else {
		fmt.Println("Erro al unmarsahl de los datos del sucesor: ", err)
	}
}
