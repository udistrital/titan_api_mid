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
	c.Mapping("AplicarOtrosi", c.AplicarOtrosi)
	c.Mapping("SuspenderContrato", c.SuspenderContrato)
}

// Post ...
// @Title Agregar Novedad
// @Description Agregar Novedad a un contrato y liquidar nuevamente
// @Param	novedad		body 	models.Novedad 	true	"Cuerpo de la novedad a guardar"
// @Success 201 {object} models.Novedad
// @Failure 400 the request contains incorrect syntax
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
	var mesIterativo int
	var anoIterativo int
	var posible bool
	posible = true

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &novedad); err == nil {
		fmt.Println("Cuotas: ", novedad.Cuotas)
		if novedad.Cuotas >= 0 {
			//Traer el concepto de la novedad para validar si es devengo o descuento
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/concepto_nomina?limit=-1&query=Id:"+strconv.Itoa(novedad.ConceptoNominaId.Id), &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &auxConcepto)
				novedad.ConceptoNominaId = &auxConcepto[0]
				//Si es devengo se añade el valor de una vez
				if auxConcepto[0].NaturalezaConceptoNominaId == 423 {
					fmt.Println("Devengo: ")
					if int(fecha_actual.Month())+novedad.Cuotas-1 > 12 {
						mesFin = int(fecha_actual.Month()) + novedad.Cuotas - 13
						anoFin = fecha_actual.Year() + 1
					} else {
						mesFin = int(fecha_actual.Month()) + novedad.Cuotas - 1
						anoFin = fecha_actual.Year()
					}

					if mesFin == 2 {
						novedad.FechaFin = time.Date(anoFin, time.Month(mesFin), 28, 0, 0, 0, 0, time.Local)
					} else {
						novedad.FechaFin = time.Date(anoFin, time.Month(mesFin), 30, 0, 0, 0, 0, time.Local)
					}
					novedad.FechaInicio = fecha_actual

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
						c.Data["message"] = "Error, no se pudo guardar la novedad" + err.Error()
						c.Abort("400")
					}
					//si es desuento
				} else if auxConcepto[0].NaturalezaConceptoNominaId == 424 {
					fmt.Println("Descuento: ")

					if int(fecha_actual.Month())+novedad.Cuotas-1 > 12 {
						mesFin = int(fecha_actual.Month()) + novedad.Cuotas - 13
						anoFin = fecha_actual.Year() + 1
					} else {
						mesFin = int(fecha_actual.Month()) + novedad.Cuotas - 1
						anoFin = fecha_actual.Year()
					}
					novedad.FechaInicio = fecha_actual
					if mesFin == 2 {
						novedad.FechaFin = time.Date(anoFin, time.Month(mesFin), 28, 0, 0, 0, 0, time.Local)
					} else {
						novedad.FechaFin = time.Date(anoFin, time.Month(mesFin), 30, 0, 0, 0, 0, time.Local)
					}

					//Verificar que las cuotas no se pasen del tiempo restante del contrato

					//Obtener el contrato
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=Id:"+strconv.Itoa(novedad.ContratoId.Id), &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &contrato)
						if contrato[0].FechaFin.Year() == novedad.FechaInicio.Year() {
							if int(contrato[0].FechaFin.Month())-int(fecha_actual.Month())+1 < novedad.Cuotas {
								fmt.Println("Las cuotas superan los meses")
								c.Data["Message"] = "Las cuotas superan los meses"
								c.Abort("400")
								posible = false
							}
						} else {

							if int(contrato[0].FechaFin.Month())+13-int(fecha_actual.Month()) < novedad.Cuotas {
								fmt.Println("Las cuotas superan los meses")
								c.Data["Message"] = "Las cuotas superan los meses"
								c.Abort("400")
								posible = false
							}
						}
					} else {
						fmt.Println("Error al obtener el contrato: ", err)
						c.Data["message"] = "Error, no existe el contrato solicitado, por favor verifique: " + err.Error()
						c.Abort("400")
					}

					mesIterativo = int(novedad.FechaInicio.Month())
					anoIterativo = novedad.FechaInicio.Year()
					auxCuotas := novedad.Cuotas
					for {
						//Obtener el valor de los honorarios de ese mes
						fmt.Println("Mes: ", mesIterativo)
						fmt.Println("Año: ", anoIterativo)
						var query = "ContratoId:" + strconv.Itoa(novedad.ContratoId.Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo)
						if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
							LimpiezaRespuestaRefactor(aux, &contratoPreliquidacion)
							query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:87"
							if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
								LimpiezaRespuestaRefactor(aux, &auxDetalle)
								honorarios = auxDetalle[0].ValorCalculado
								fmt.Println("honorarios: ", honorarios)
							} else {
								fmt.Println("Error al obtener el valor de los honorarios: ", err)
								c.Data["message"] = "Error al obtener el valor de los honorarios: " + err.Error()
								c.Abort("400")
							}
						} else {
							fmt.Println("Error al obtener el contrato peliquidacion: ", err)
							c.Data["message"] = "El contrato no está vigente para mes solicitado " + err.Error()
							c.Abort("400")
						}
						//Obtener el Valor de los descuentos de ese mes
						query = "ContratoId:" + strconv.Itoa(novedad.ContratoId.Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo)
						if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
							LimpiezaRespuestaRefactor(aux, &contratoPreliquidacion)
							query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:573"
							if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
								LimpiezaRespuestaRefactor(aux, &auxDetalle)
								descuentos = auxDetalle[0].ValorCalculado
							} else {
								fmt.Println("Error al obtener el valor de los descuentos: ", err)
								c.Data["message"] = "No se pudieron obtener los descuentos" + err.Error()
								c.Abort("400")
							}
						} else {
							fmt.Println("Error al obtener el contrato peliquidacion: ", err)
							c.Data["message"] = "El contrato no está vigente para mes solicitado " + err.Error()
							c.Abort("400")
						}
						//Verificar que el valor de los descuentos no supera la mitad de los honorarios
						if auxConcepto[0].TipoConceptoNominaId == 419 {
							if descuentos+novedad.Valor > (honorarios / 2) {
								fmt.Println("Las cuotas superan el valor de los honorarios")
								c.Data["message"] = "Las cuotas superan el valor de los honorarios, por favor verifique el valor de la novedad"
								c.Abort("400")
								posible = false
							}
						} else if auxConcepto[0].TipoConceptoNominaId == 420 {
							if descuentos+((novedad.Valor/100)*honorarios) > (honorarios / 2) {
								fmt.Println("Las cuotas superan el valor de los honorarios")
								c.Data["message"] = "Las cuotas superan el valor de los honorarios, por favor verifique el valor de la novedad"
								c.Abort("400")
								posible = false
							}
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

					if posible {
						if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/novedad", "POST", &aux, novedad); err == nil {
							fmt.Println("Novedad Registrada ")
							LimpiezaRespuestaRefactor(aux, &novedad)
							fmt.Println(novedad.ConceptoNominaId.TipoConceptoNominaId)
							if novedad.ConceptoNominaId.TipoConceptoNominaId == 419 || novedad.ConceptoNominaId.TipoConceptoNominaId == 420 {
								//Agregar el Valor al detalle
								mensaje, err := AgregarValorNovedad(novedad)
								if err == nil {
									c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": novedad}
								} else {
									c.Data["mesaage"] = mensaje + err.Error()
									c.Abort("404")
								}
							}
						} else {
							fmt.Println("No se pudo guardar la novedad", err)
							c.Data["mesaage"] = "Error, no se pudo guardar la novedad: " + err.Error()
							c.Abort("400")
						}
					} else {
						fmt.Println("No se cumplen los requisitos para guardar la novedad")
						c.Data["mesaage"] = "No se cumplen los requisitos para guardar la novedad, por favor revise cuotas o valor"
						c.Abort("400")
					}

				} else {
					fmt.Println("Error al obtener el concepto: ", err)
					c.Data["mesaage"] = "Error, el concepto no existe: " + err.Error()
					c.Abort("400")
				}

			}
		} else {
			fmt.Println("Número de cuotas inválido porque es menor que 0 ")
			c.Data["mesaage"] = "Error, Por favor revise el número de cuotas, no puede ser menor o igual que 0"
			c.Abort("400")
		}

	} else {
		fmt.Println("Error al Unmarshal de novedad: ", err)
		c.Data["mesaage"] = "Error, el JSON enviado contiene un parámetro incorrecto: " + err.Error()
		c.Abort("400")
	}
	c.ServeJSON()
}

// Get ...
// @Title Eliminar Novedad
// @Description Eliminar Novedad Novedad a un contrato y liquidar nuevamente
// @Param	id		path 	true	"Id de la novedad que se va retirar"
// @Success 201 models.Novedad.Id
// @Failure 400 the request contains incorrect syntax
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
			mensaje, err := EliminarValorNovedad(novedad[0], fecha_actual)
			if err == nil {
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": novedad[0].Id}
			} else {
				c.Data["message"] = mensaje + err.Error()
				c.Abort("404")
			}
		} else {
			c.Data["mesaage"] = "Error, no se pudo eliminar la novedad: " + err.Error()
			c.Abort("404")
		}
	} else {
		fmt.Println("Error al unmarshal de la novedad ", err)
		c.Data["mesaage"] = "Error, no existe la novedad: " + err.Error()
		c.Abort("404")
	}
	c.ServeJSON()
}

// Get ...
// @Title Registrar Cancelacion
// @Description Maneja la novedad contractual de cancelación
// @Param	NumeroContrato		path 	true	"Numero del contrato que se va a cancelar"
// @Param	Vigencia		path 	true	"Vigencia del contrato que se va a cancelar"
// @Param	Documento		path 	true	"Documento del contratista"
// @Success 201 {object} models.Contrato
// @Failure 400
// @router /cancelar_contrato/:NumeroContrato/:Vigencia/:Documento [get]
func (c *NovedadController) CancelarContrato() {
	var numero = c.Ctx.Input.Param(":NumeroContrato")
	var vigencia = c.Ctx.Input.Param(":Vigencia")
	var documento = c.Ctx.Input.Param(":Documento")
	var fecha_actual = time.Now()
	var aux map[string]interface{}
	var contrato []models.Contrato
	var contrato_preliquidacion []models.ContratoPreliquidacion
	var valorDia float64
	var detalles []models.DetallePreliquidacion
	var mesInicio = int(fecha_actual.Month())
	var mesFin int

	//Traer el contrato a cancelar
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=NumeroContrato:"+numero+",Vigencia:"+vigencia+",Documento:"+documento, &aux); err == nil {
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
	//Actualizar fecha de finalización del contrato
	if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato[0].Id), "PUT", &aux, contrato[0]); err == nil {
		fmt.Println("Contrato Actualizado")
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Registration successful", "Data": contrato[0]}
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
		liquidarCPS(contrato[0])
	} else if contrato[0].TipoNominaId == 409 {
		liquidarHCH(contrato[0])
	}

	c.ServeJSON()
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
	var mesIterativo int
	var anoIterativo int

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &sucesor); err == nil {

		//Traer el contrato a cancelar
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=NumeroContrato:"+sucesor.NumeroContrato+",Vigencia:"+strconv.Itoa(sucesor.Vigencia)+",Documento:"+sucesor.DocumentoActual, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contrato)
			contratoNuevo = contrato[0]
			valorViejo = contrato[0].ValorContrato

			if sucesor.FechaInicio.Before(contrato[0].FechaFin) {
				mesIterativo = int(sucesor.FechaInicio.Month())
				anoIterativo = sucesor.FechaInicio.Year()
				//Eliminar los detalles y los contratos_preliquidacion
				for {
					fmt.Println("Mes:", mesIterativo)
					fmt.Println("Ano:", anoIterativo)
					//Obtener contrato_preliquidacion para ese mes
					query := "ContratoId:" + strconv.Itoa(contrato[0].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo)
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
						if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
							LimpiezaRespuestaRefactor(aux, &detalles)
							for j := 0; j < len(detalles); j++ {
								if detalles[j].ConceptoNominaId.Id == 87 && mesIterativo == int(sucesor.FechaInicio.Month()) {
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

					if mesIterativo == int(contrato[0].FechaFin.Month()) && anoIterativo == contrato[0].FechaFin.Year() {
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

				//Obtener lo liquidado hasta el momento
				fmt.Println("Obteniendo:")
				fmt.Println("Fecha Inicio:", sucesor.FechaInicio.Month(), " ", sucesor.FechaInicio.Year())
				query := "ContratoPreliquidacionId.ContratoId.NumeroContrato:" + contrato[0].NumeroContrato
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &detalles)
					for j := 0; j < len(detalles); j++ {
						if detalles[j].ConceptoNominaId.Id == 87 {
							valorNuevo = valorNuevo + detalles[j].ValorCalculado
						}
					}
				} else {
					fmt.Println("Error al obtener detalles")
				}

				//Agregar el nuevo contrato
				contratoNuevo.Id = 0
				contratoNuevo.Documento = sucesor.DocumentoNuevo
				contratoNuevo.NombreCompleto = sucesor.NombreCompleto
				contratoNuevo.FechaInicio = sucesor.FechaInicio
				contratoNuevo = registratContrato(contratoNuevo)
				fmt.Println("Contrato nuevo: ", contratoNuevo)

				contrato[0].FechaFin = sucesor.FechaInicio.Add(24 * time.Hour * -1)
				//Actualizar fecha de finalización del contrato
				if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato[0].Id), "PUT", &aux, contrato[0]); err == nil {
					fmt.Println("Contrato Actualizado")
				} else {
					fmt.Println("Error al actualizar el contrato: ", err)
				}

				if contrato[0].FechaInicio.Month() != sucesor.FechaInicio.Month() && contrato[0].FechaInicio.Year() != sucesor.FechaInicio.Year() {
					contrato[0].FechaInicio = time.Date(sucesor.FechaInicio.Year(), sucesor.FechaInicio.Month(), 1, 12, 0, 0, 0, time.UTC)
				} else {
					contrato[0].FechaInicio = time.Date(sucesor.FechaInicio.Year(), sucesor.FechaInicio.Month(), contrato[0].FechaInicio.Day(), 12, 0, 0, 0, time.UTC)
				}

				//El valor del contrato nuevo es lo que queda del contrato pasado (no se almacena únicamente para cálculos)
				if sucesor.FechaInicio.Day() == 1 {
					contratoNuevo.ValorContrato = valorViejo - valorNuevo
					fmt.Println("Liquidando nuevo:", contratoNuevo.NumeroContrato, " de ", contratoNuevo.NombreCompleto)
					liquidarCPS(contratoNuevo)
				} else {
					contrato[0].ValorContrato = valorDia * float64(contrato[0].FechaFin.Day())
					contratoNuevo.ValorContrato = valorViejo - valorNuevo - valorDia*float64(contrato[0].FechaFin.Day())
					fmt.Println("Liquidando actual:", contrato[0].NumeroContrato, " de ", contrato[0].NombreCompleto)
					liquidarCPS(contrato[0])
					fmt.Println("Liquidando nuevo:", contratoNuevo.NumeroContrato, " de ", contratoNuevo.NombreCompleto)
					liquidarCPS(contratoNuevo)
				}
			} else {
				fmt.Println("El contrato no puede ser cedido, debido que la fecha fin está antes de la fecha inicio")
				c.Data["mesaage"] = "El contrato no puede ser cedido, por favor verifique fechas"
				c.Abort("400")
			}
		} else {
			fmt.Println("Error al obtener el contrato: ", err)
			c.Data["mesaage"] = "Error, el contrato soilicitado no existe: " + err.Error()
			c.Abort("400")
		}
	} else {
		fmt.Println("Error al unmarsahl de los datos del sucesor: ", err)
		c.Data["mesaage"] = "Error service POST: The request contains an incorrect data type or an invalid parameter: " + err.Error()
		c.Abort("400")
	}
}

// Post ...
// @Title Aplicar Oro sí
// @Description Aplicar novedad contractual de otro sí
// @Param	OtroSí		body 	models.OtroSi 	true	"Documentos de contratista Nueva Fecha de finalización del contrato"
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
	var mesesContrato float64
	var valorMensual float64
	var honorarios float64
	var fechaRespaldo time.Time

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &otro_si); err == nil {
		//Traer el contrato
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=NumeroContrato:"+otro_si.NumeroContrato+",Vigencia:"+strconv.Itoa(otro_si.Vigencia)+",Documento:"+otro_si.Documento, &aux); err == nil {
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
	_, mesesContrato = CalcularDias(contrato[0].FechaInicio, contrato[0].FechaFin)
	mesesContrato = float64(int(mesesContrato + 1))
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
	contrato[0].ValorContrato = Roundf(contrato[0].ValorContrato + valorMensual*(diasNuevos/30))
	if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato[0].Id), "PUT", &aux, contrato[0]); err == nil {
		fmt.Println("Contrato Actualizado")
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Registration successful", "Data": contrato[0]}
	} else {
		fmt.Println("Error al actualizar contrato: ", err)
		c.Data["mesaage"] = "Error service POST: The request contains an incorrect data type or an invalid parameter: " + err.Error()
		c.Abort("400")
	}
	//Calcular valor agregado para la liquidacion
	valorNuevo = (valorMensual * (diasNuevos / 30)) + honorarios
	//diasRestantes = int(diasNuevos) + dias

	//Ajustar Datos para el la liquidación a partir de ahí
	contrato[0].FechaInicio = time.Date(fechaRespaldo.Year(), fechaRespaldo.Month(), 1, 12, 0, 0, 0, time.UTC)
	contrato[0].FechaFin = otro_si.FechaFin
	contrato[0].ValorContrato = valorNuevo
	//liquidar el contrato
	liquidarCPS(contrato[0])
	c.ServeJSON()
}

// Post ...
// @Title Suspender Contrato
// @Description Maneja la novedad contractual de Suspensión de Contrato
// @Param	Suspension	 body  models.Suspension	true	"Duración de la suspensión"
// @Success 201 {object} models.Contrato
// @Failure 400 the request contains incorrect syntax
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
	var mesIterativo int
	var anoIterativo int
	var diasSuspension float64

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &suspension); err == nil {
		//Traer el contrato a cancelar
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=NumeroContrato:"+suspension.NumeroContrato+",Vigencia:"+strconv.Itoa(suspension.Vigencia)+",Documento:"+suspension.Documento, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contrato)
			contratoNuevo = contrato[0]
			valorViejo = contrato[0].ValorContrato
		} else {
			fmt.Println("Error al obtener el contrato: ", err)
		}

		mesIterativo = int(suspension.FechaInicio.Month())
		anoIterativo = suspension.FechaInicio.Year()

		//Eliminar los detalles y los contratos_preliquidacion
		for {
			//Obtener contrato_preliquidacion para ese mes
			query := "ContratoId:" + strconv.Itoa(contrato[0].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo)
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &detalles)
					for j := 0; j < len(detalles); j++ {
						if detalles[j].ConceptoNominaId.Id == 87 && mesIterativo == int(suspension.FechaInicio.Month()) {
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

			if mesIterativo == int(contrato[0].FechaFin.Month()) && anoIterativo == contrato[0].FechaFin.Year() {
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

		mesIterativo = int(contrato[0].FechaInicio.Month())
		anoIterativo = contrato[0].FechaInicio.Year()
		fmt.Println("Obteniendo: ")
		//Obtener lo liquidado hasta el momento
		for {
			//Obtener contrato_preliquidacion para ese mes
			fmt.Println("Mes: ", mesIterativo)
			fmt.Println("Año: ", anoIterativo)

			if mesIterativo == int(suspension.FechaInicio.Month()) && anoIterativo == suspension.FechaInicio.Year() {
				break
			}

			query := "ContratoId:" + strconv.Itoa(contrato[0].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo)
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &detalles)
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

			if mesIterativo == 12 {
				mesIterativo = 1
				anoIterativo = anoIterativo + 1
			} else {
				mesIterativo = mesIterativo + 1
			}

		}
		//Calcular los días de la suspension
		diasSuspension, _ = CalcularDias(suspension.FechaInicio, suspension.FechaFin)
		contratoNuevo.FechaInicio = suspension.FechaFin.Add(24 * time.Hour)

		//Desface de dos días para febrero
		if int(suspension.FechaFin.Month()) == 2 {
			contratoNuevo.FechaFin = contratoNuevo.FechaFin.Add(24 * time.Hour * time.Duration(diasSuspension+2))
		} else {
			contratoNuevo.FechaFin = contratoNuevo.FechaFin.Add(24 * time.Hour * time.Duration(diasSuspension))
		}

		contrato[0].FechaFin = contratoNuevo.FechaFin
		//Actualizar fecha de finalización del contrato
		if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato[0].Id), "PUT", &aux, contrato[0]); err == nil {
			fmt.Println("Contrato Actualizado")
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Registration successful", "Data": contrato[0]}
		} else {
			fmt.Println("Error al ctualizar el contrato: ", err)
			c.Data["mesaage"] = "Error service POST: The request contains an incorrect data type or an invalid parameter: " + err.Error()
			c.Abort("400")
		}
		contrato[0].FechaFin = suspension.FechaInicio.Add(24 * time.Hour * -1)

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
			liquidarCPS(contrato[0])
			liquidarCPS(contratoNuevo)
		} else if contrato[0].TipoNominaId == 409 {
			liquidarHCH(contrato[0])
		}
	} else {
		fmt.Println("Error al unmarsahl de los datos del sucesor: ", err)
		c.Data["mesaage"] = "Error service POST: The request contains an incorrect data type or an invalid parameter"
		c.Abort("400")
	}
}
