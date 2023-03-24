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

type NovedadVEController struct {
	beego.Controller
}

func (c *NovedadVEController) URLMapping() {
	c.Mapping("VerficarDescuentos", c.VerificarDescuentos)
	c.Mapping("AgregarNovedad", c.AgregarNovedad)
	c.Mapping("EliminarNovedad", c.EliminarNovedad)
	c.Mapping("GenerarAdicion", c.GenerarAdicion)
	c.Mapping("GenerarAnulacion", c.AplicarAnulacion)
}

// Post ...
// @Title Verificar Novedad
// @Description Verificar los valores de una novedad para ver si superan los descuentos
// @Param	novedad		body 	models.Novedad 	true	"Cuerpo de la novedad a guardar"
// @Success 201 {object} models.MensajeNovedad
// @Failure 400 the request contains incorrect syntax
// @router /verificar_descuentos [post]
func (c *NovedadVEController) VerificarDescuentos() {

	var aux map[string]interface{}
	var novedad models.Novedad
	var auxConcepto []models.ConceptoNomina
	var contrato []models.Contrato
	var auxDetalle []models.DetallePreliquidacion
	var contratoPreliquidacion []models.ContratoPreliquidacion
	var res models.MensajeNovedad
	var honorarios float64
	var descuentos float64
	var idHonorarios int
	var fecha_actual time.Time
	var anoFin int
	var mesFin int
	var mesIterativo int
	var anoIterativo int

	const CUOTAS_SUPERADAS = 1
	const DESCUENTOS_SUPERADOS = 2
	const CONCEPTO_EXISTENTE = 3
	const SIN_PROBLEMA = 4

	res.Estado = SIN_PROBLEMA
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &novedad); err == nil {

		if novedad.FechaInicio == time.Date(0001, 01, 01, 00, 00, 00, 0000, time.UTC) {
			fmt.Println("Fecha vacía")
			fecha_actual = time.Now()
		} else {
			fmt.Println("Fecha provista")
			fecha_actual = novedad.FechaInicio
		}

		if novedad.Cuotas >= 0 {
			//Traer el concepto de la novedad para validar si es devengo o descuento
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/concepto_nomina?limit=-1&query=Id:"+strconv.Itoa(novedad.ConceptoNominaId.Id), &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &auxConcepto)
				novedad.ConceptoNominaId = &auxConcepto[0]
				//traer el contrato
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=Id:"+strconv.Itoa(novedad.ContratoId.Id), &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contrato)
					if contrato[0].TipoNominaId == 410 {
						idHonorarios = 152
					} else {
						idHonorarios = 87
					}
					novedad.ContratoId = &contrato[0]

					if auxConcepto[0].NaturalezaConceptoNominaId == 424 {
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
						/*if contrato[0].FechaFin.Year() == novedad.FechaInicio.Year() {
							if int(contrato[0].FechaFin.Month())-int(fecha_actual.Month())+1 < novedad.Cuotas {
								fmt.Println("Las cuotas superan los meses", int(contrato[0].FechaFin.Month()), int(fecha_actual.Month()), int(contrato[0].FechaFin.Month())-int(fecha_actual.Month())+1, novedad.Cuotas, 1)
								res.Mensaje = "Las cuotas superan los meses"
								res.Estado = CUOTAS_SUPERADAS
								c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "successful", "Data": res}
							}
						} else {
							if int(contrato[0].FechaFin.Month())+13-int(fecha_actual.Month()) < novedad.Cuotas {
								fmt.Println("Las cuotas superan los meses, ", int(contrato[0].FechaFin.Month())+13-int(fecha_actual.Month()), novedad.Cuotas, 2)
								res.Mensaje = "Las cuotas superan los meses"
								res.Estado = CUOTAS_SUPERADAS
								c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "successful", "Data": res}
							}
						}*/

						if res.Estado != 1 {
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
									query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id)
									if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
										LimpiezaRespuestaRefactor(aux, &auxDetalle)
										if auxDetalle[0].Id != 0 {
											descuentos = 0
											for i := 0; i < len(auxDetalle); i++ {
												if auxDetalle[i].ConceptoNominaId.Id == novedad.ConceptoNominaId.Id {
													res.Mensaje = "El concepto ya existe en el mes " + strconv.Itoa(mesIterativo) + " del año " + strconv.Itoa(anoIterativo) + " no es posible agregarlo"
													res.Estado = CONCEPTO_EXISTENTE
													c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "successful", "Data": res}
													break
												} else {
													if auxDetalle[i].ConceptoNominaId.Id == idHonorarios {
														honorarios = auxDetalle[i].ValorCalculado
													} else if auxDetalle[i].ConceptoNominaId.NaturalezaConceptoNominaId == 424 && auxDetalle[i].ConceptoNominaId.Id != 64 {
														descuentos = descuentos + auxDetalle[i].ValorCalculado
													}
												}
											}
											//Si es fijo
											if novedad.ConceptoNominaId.TipoConceptoNominaId == 419 {
												if (novedad.Valor + descuentos) > (honorarios / 2) {
													res.Mensaje = "Se superan el tope de descuentos del mes " + strconv.Itoa(mesIterativo) + " del año " + strconv.Itoa(anoIterativo)
													res.Estado = DESCUENTOS_SUPERADOS
													c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "successful", "Data": res}
													break
												}
												//Si es porcentual
											} else if novedad.ConceptoNominaId.TipoConceptoNominaId == 420 {
												if (honorarios*(novedad.Valor/100) + descuentos) > (honorarios / 2) {
													res.Mensaje = "Se superan el tope de descuentos del mes " + strconv.Itoa(mesIterativo) + " del año " + strconv.Itoa(anoIterativo)
													res.Estado = DESCUENTOS_SUPERADOS
													c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "successful", "Data": res}
													break
												}
											}
										} else {
											fmt.Println("Error al obtener el valor de los honorarios: ", err)
											c.Data["mesaage"] = "No se encontraron detalles para el contrato seleccionado"
											c.Abort("400")
										}
									} else {
										fmt.Println("Error al obtener el valor de los honorarios: ", err)
										c.Data["mesaage"] = "Error al obtener el valor de los honorarios: " + err.Error()
										c.Abort("400")
									}
								} else {
									fmt.Println("Error al obtener el contrato peliquidacion: ", err)
									c.Data["mesaage"] = "El contrato no está vigente para mes solicitado " + err.Error()
									c.Abort("400")
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
						}

					} else {
						fmt.Println("Error al obtener el concepto: ", err)
						c.Data["mesaage"] = "Error, el concepto no existe: " + err.Error()
						c.Abort("400")
					}
				} else {
					fmt.Println("Error al obtener contrato: ", err)
					c.Data["mesaage"] = "Error, el contrato no existe: " + err.Error()
					c.Abort("400")
				}
			} else {
				fmt.Println("Error al obtener el concepto de nómina: ", err)
				c.Data["mesaage"] = "Error, no se encontró el concepto de nómina " + err.Error()
				c.Abort("400")
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

	if res.Estado == SIN_PROBLEMA {
		res.Mensaje = "No hay Problema para agregar la novedad"
		res.Estado = SIN_PROBLEMA
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "successful", "Data": res}
	}
	c.ServeJSON()
}

// Post ...
// @Title Agregar Novedad
// @Description Agregar Novedad a un contrato y liquidar nuevamente
// @Param	novedad		body 	models.Novedad 	true	"Cuerpo de la novedad a guardar"
// @Success 201 {object} models.Novedad
// @Failure 400 the request contains incorrect syntax
// @router /agregar_novedad [post]
func (c *NovedadVEController) AgregarNovedad() {
	var aux map[string]interface{}
	var auxNovedad []models.Novedad
	var contrato []models.Contrato
	var auxConcepto []models.ConceptoNomina
	var novedad models.Novedad
	var fecha_actual time.Time
	var mesFin int
	var anoFin int

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &novedad); err == nil {
		if novedad.FechaInicio == time.Date(0001, 01, 01, 00, 00, 00, 0000, time.UTC) {
			fmt.Println("Fecha vacía")
			fecha_actual = time.Now()
		} else {
			fmt.Println("Fecha provista")
			fecha_actual = novedad.FechaInicio
		}

		if novedad.Cuotas >= 0 {
			//Ajustar Fecha inicio y Fin de la novedad
			if int(fecha_actual.Month())+novedad.Cuotas-1 > 12 {
				mesFin = int(fecha_actual.Month()) + novedad.Cuotas - 13
				anoFin = fecha_actual.Year() + 1
			} else {
				mesFin = int(fecha_actual.Month()) + novedad.Cuotas - 1
				anoFin = fecha_actual.Year()
			}
			novedad.FechaInicio = fecha_actual
			if mesFin == 2 {
				novedad.FechaFin = time.Date(anoFin, time.Month(mesFin), 28, 12, 0, 0, 0, time.Local)
			} else {
				novedad.FechaFin = time.Date(anoFin, time.Month(mesFin), 30, 12, 0, 0, 0, time.Local)
			}

			//Traer el concepto de la novedad
			if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/concepto_nomina?limit=-1&query=Id:"+strconv.Itoa(novedad.ConceptoNominaId.Id), &aux); err == nil {
				LimpiezaRespuestaRefactor(aux, &auxConcepto)
				novedad.ConceptoNominaId = &auxConcepto[0]
				//traer el contrato
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=Id:"+strconv.Itoa(novedad.ContratoId.Id), &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contrato)
					novedad.ContratoId = &contrato[0]
					if novedad.ConceptoNominaId.TipoConceptoNominaId == 419 || novedad.ConceptoNominaId.TipoConceptoNominaId == 420 {
						if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/novedad", "POST", &aux, novedad); err == nil {
							LimpiezaRespuestaRefactor(aux, &novedad)
							if novedad.Id != 0 {
								if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/novedad?limit=-1&query=Id:"+strconv.Itoa(novedad.Id), &aux); err == nil {
									LimpiezaRespuestaRefactor(aux, &auxNovedad)
									novedad = auxNovedad[0]
									//Agregar el Valor al detalle
									fmt.Println("Novedad a enviar: ", novedad)
									mensaje, ids_detalles, err := AgregarValorNovedad(novedad)
									if err == nil {
										c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": novedad}
									} else {
										//Se hace rollback de lo agregado previamente
										if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/novedad/"+strconv.Itoa(novedad.Id), "DELETE", &aux, nil); err == nil {
											for i := 0; i < len(ids_detalles); i++ {
												if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(ids_detalles[i].Id), "DELETE", &aux, nil); err == nil {
													fmt.Println("Concepto Eliminado")
												} else {
													c.Data["mesaage"] = "Error al eliminar detalle preliquidacion de la novedad" + err.Error()
													c.Abort("404")
												}
											}
											c.Data["mesaage"] = mensaje + err.Error()
											c.Abort("404")
										} else {
											c.Data["mesaage"] = "Error al eliminar la novedad" + err.Error()
											c.Abort("404")
										}
										c.Data["mesaage"] = mensaje + err.Error()
										c.Abort("404")
									}
								} else {
									fmt.Println("No se pudo obtener la novedad", err)
									c.Data["mesaage"] = "No se pudo obtener la novedad " + err.Error()
									c.Abort("400")
								}
							} else {
								fmt.Println("No se pudo guardar la novedad", err)
								c.Data["mesaage"] = "Error, no se pudo guardar la novedad: " + err.Error()
								c.Abort("400")
							}
						} else {
							fmt.Println("No se pudo guardar la novedad", err)
							c.Data["mesaage"] = "Error, no se pudo guardar la novedad: " + err.Error()
							c.Abort("400")
						}
					} else {
						fmt.Println("Se intentó cargar una novedade de Seguridad Social")
						c.Data["mesaage"] = "Novedades de Seguridad Social no implementadas "
						c.Abort("400")
					}
				} else {
					fmt.Println("Error al obtener contrato: ", err)
					c.Data["mesaage"] = "Error, el contrato no existe: " + err.Error()
					c.Abort("400")
				}
			} else {
				fmt.Println("Error al obtener el concepto de nómina: ", err)
				c.Data["mesaage"] = "Error al obtener el concepto de nómina: " + err.Error()
				c.Abort("400")
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
func (c *NovedadVEController) EliminarNovedad() {
	var id = c.Ctx.Input.Param(":id")
	var fecha_actual = time.Now()
	var aux map[string]interface{}
	var novedad []models.Novedad
	//Verificar cómo se van a enviar los datos del contrato al trabajar el front
	fmt.Println("novedad_Id: ", id)
	//Buscar la novedad
	fmt.Println(beego.AppConfig.String("UrlTitanCrud") + "/novedad?limit=-1&query=Id:" + id)
	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/novedad?limit=-1&query=Id:"+id, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &novedad)
		if novedad[0].Id != 0 {
			mensaje, err := EliminarValorNovedad(novedad[0], fecha_actual)
			if err == nil {
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": novedad[0].Id}
			} else {
				c.Data["mesaage"] = mensaje + err.Error()
				c.Abort("404")
			}
		} else {
			fmt.Println("no se pudo eliminar la novedad: ", novedad[0])
			c.Data["mesaage"] = "Error, No se encontró la novedad: " + err.Error()
			c.Abort("404")
		}
	} else {
		fmt.Println("Error al unmarshal de la novedad ", err)
		c.Data["mesaage"] = "Error al unmarshal de la novedad: " + err.Error()
		c.Abort("404")
	}
	c.ServeJSON()
}

// Post ...
// @Title Adicionar horas VE
// @Description Maneja la novedad contractual de Adición de horas para contratos de docentes de VE
// @Param	Adición	 body  models.Adicion	true	"Datos de la adición"
// @Success 201 {object} models.Contrato
// @Failure 400 the request contains incorrect syntax
// @router /generar_adicion [post]
func (c *NovedadVEController) GenerarAdicion() {

	var aux map[string]interface{}
	var adicion models.Adicion
	var contrato []models.Contrato
	var contratoNuevo models.Contrato
	var desagregado models.DesagregadoContratoHCS
	var datosVinculacion models.DatosVinculacion
	var mensaje string

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &adicion); err == nil {
		//Asignar datos para el desagregado
		datosVinculacion.NumeroContrato = adicion.NumeroContrato
		datosVinculacion.Vigencia = adicion.Vigencia
		datosVinculacion.Documento = adicion.Documento
		datosVinculacion.HorasSemanales = adicion.HorasSemanales
		datosVinculacion.NumeroSemanas = adicion.NumeroSemanas
		datosVinculacion.Dedicacion = adicion.Dedicacion
		datosVinculacion.NivelAcademico = adicion.NivelAcademico
		datosVinculacion.Categoria = adicion.Categoria

		//Obtener valores desagregados
		desagregado = Desagregar(datosVinculacion)
		//Obtener los datos del contrato para generar adición
		query := "NumeroContrato:" + adicion.NumeroContrato + ",Vigencia:" + strconv.Itoa(adicion.Vigencia) + ",Documento:" + adicion.Documento + ",Rp:" + strconv.Itoa(adicion.RpActual) + ",Cdp:" + strconv.Itoa(adicion.Cdp)
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contrato)
			if contrato[0].Id != 0 {
				//Asignar nuevos valores al contrato
				contratoNuevo = contrato[0]
				contratoNuevo.Id = 0
				contratoNuevo.Rp = adicion.RpNuevo
				contratoNuevo.ValorContrato = desagregado.SueldoBasico
				contratoNuevo.FechaInicio = adicion.FechaInicio
				contratoNuevo.Unico = false
				contratoNuevo.Completo = true
				contratoNuevo, err = registrarContrato(contratoNuevo)
				if err == nil {
					mensaje, err = liquidarHCS(contratoNuevo, false, 0, contratoNuevo.Vigencia, 0, 0, false)
					if err == nil {
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": contratoNuevo}
					} else {
						if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contratoNuevo.Id), "DELETE", &aux, nil); err == nil {
							fmt.Println("Contrato de adición eliminado")
						} else {
							fmt.Println("Error al eliminar contrato nuevo", err)
							c.Data["mesaage"] = "Error al eliminar adición " + err.Error()
							c.Abort("404")
						}
						fmt.Println("Error al agregar adición: ", err)
						c.Data["mesaage"] = mensaje + err.Error()
						c.Abort("404")
					}
				} else {
					fmt.Println("No se encontró el contrato:", err)
					c.Data["mesaage"] = "Error, no se encontró el contrato: " + err.Error()
					c.Abort("404")
				}
			} else {
				fmt.Println("No se encontró el contrato:", err)
				c.Data["mesaage"] = "Error, no se encontró el contrato: " + err.Error()
				c.Abort("404")
			}
		} else {
			fmt.Println("Error al consultar contratos: ", err)
			c.Data["mesaage"] = "Error al consultar contratos: " + err.Error()
			c.Abort("404")
		}
	} else {
		fmt.Println("Error al unmarshal de la novedad: ", err)
		c.Data["mesaage"] = "Error, no existe la novedad: " + err.Error()
		c.Abort("404")
	}
	c.ServeJSON()
}

// Post ...
// @Title Aplicar anulación
// @Description Maneja la novedad contractual de anulación de contratos de docentes de VE
// @Param	Adición	 body  models.Anulacion	true	"Datos de la anulacion"
// @Success 201 {object} models.Contrato
// @Failure 400 the request contains incorrect syntax
// @router /aplicar_anulacion [post]
func (c *NovedadVEController) AplicarAnulacion() {
	var anulacion models.Anulacion
	var aux map[string]interface{}
	var contrato []models.Contrato
	var contratoOriginal models.Contrato
	var contrato_preliquidacion []models.ContratoPreliquidacion
	var valorDia float64
	var semanasContrato int
	var semanasTotales int
	var detalles []models.DetallePreliquidacion
	var codAux int
	var DetallesAux []models.DetallePreliquidacion
	var sumaContratosTemp float64
	var mensaje string //Mensaje de error

	//Traer el contrato a cancelar
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &anulacion); err == nil {
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=Rp:"+strconv.Itoa(anulacion.Rp)+",Vigencia:"+strconv.Itoa(anulacion.Vigencia)+",Documento:"+anulacion.Documento+",Activo:true", &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contrato)
			if contrato[0].Id != 0 {

				//Ordenar los contratos para tomar el más reciente
				for i := 0; i < len(contrato); i++ {
					if contrato[0].Id < contrato[i].Id {
						auxContrato := contrato[0]
						contrato[0] = contrato[i]
						contrato[i] = auxContrato
					}
				}

				contratoOriginal = contrato[0]

				anoIterativo := anulacion.FechaAnulacion.Year()
				mesIterativo := int(anulacion.FechaAnulacion.Month())

				//Eliminar los detalles y los contratos_preliquidacion
				for {
					//Obtener contrato_preliquidacion para ese mes
					query := "ContratoId:" + strconv.Itoa(contrato[0].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo)
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
						if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
							LimpiezaRespuestaRefactor(aux, &detalles)
							for j := 0; j < len(detalles); j++ {
								if contrato[0].TipoNominaId == 410 {
									if detalles[j].ConceptoNominaId.Id == 152 {
										valorDia = detalles[j].ValorCalculado / detalles[j].DiasLiquidados
									}
								} else {
									if detalles[j].ConceptoNominaId.Id == 87 {
										valorDia = detalles[j].ValorCalculado / detalles[j].DiasLiquidados
									}
								}
								if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalles[j].Id), "DELETE", &aux, nil); err == nil {
									fmt.Println("Detalle eliminado con éxito")
								} else {
									fmt.Println("Error al Eliminar Detalles:", err)
									c.Data["mesaage"] = "Error al Eliminar Detalles: " + err.Error()
									c.Abort("400")
								}
							}
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[0].Id), "DELETE", &aux, nil); err == nil {
								fmt.Println("contrato preliquidacion eliminado con éxito")
							} else {
								fmt.Println("Error al eliminar contrato preliquidacion: ", err)
								c.Data["mesaage"] = "Error al Eliminar Detalles: " + err.Error()
								c.Abort("400")
							}
						} else {
							fmt.Println("Error al obtener detalles")
							c.Data["mesaage"] = "Error al Eliminar Detalles: " + err.Error()
							c.Abort("400")
						}
					} else {
						fmt.Println("Error al obtener contrato_preliquidacion")
						c.Data["mesaage"] = "Error al Eliminar Detalles: " + err.Error()
						c.Abort("400")
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
				fmt.Println("Valor día: ", valorDia)
				contrato[0].FechaFin = anulacion.FechaAnulacion
				//Actualizar fecha de finalización del contrato
				contratoOriginal.Activo = false
				contratoOriginal.Id = 0
				// CREA CONTRATO NUEVO CON LA INFORMACIÓN DEL CONTRATO ORIGINAL CON CAMPO ACTIVO FALSE
				if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato", "POST", &aux, contratoOriginal); err == nil {
					if contrato[0].FechaInicio.Month() != anulacion.FechaAnulacion.Month() || contrato[0].FechaInicio.Year() != anulacion.FechaAnulacion.Year() {
						semanasTotales = int(calcularSemanasContratoDVE(contrato[0].FechaInicio, anulacion.FechaAnulacion))
						semanasContrato = int(calcularSemanasContratoDVE(time.Date(anulacion.FechaAnulacion.Year(), anulacion.FechaAnulacion.Month(), 1, 12, 0, 0, 0, time.UTC), anulacion.FechaAnulacion))
						contrato[0].FechaInicio = time.Date(anulacion.FechaAnulacion.Year(), anulacion.FechaAnulacion.Month(), 1, 12, 0, 0, 0, time.UTC)
						contrato[0].ValorContrato = valorDia * float64(semanasContrato)
					} else {
						fmt.Println("Anulación el mismo mes de inicio")
						diaAux := contrato[0].FechaInicio.AddDate(0, 0, 1)
						semanasContrato = int(calcularSemanasContratoDVE(diaAux, contrato[0].FechaFin))
						semanasTotales = semanasContrato
						contrato[0].ValorContrato = valorDia * float64(semanasContrato)
					}

					if contrato[0].TipoNominaId == 409 {
						codAux = 87
					} else if contrato[0].TipoNominaId == 410 {
						codAux = 152
					}

					// Trae todos los detalles preliquidación no eliminados para calcular el nuevo valor del contrato
					query := "ContratoPreliquidacionId__ContratoId__Id:" + strconv.Itoa(contrato[0].Id) + ",ConceptoNominaId:" + strconv.Itoa(codAux)
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &DetallesAux)
						for i := 0; i < len(DetallesAux); i++ {
							fmt.Println(DetallesAux[i].ValorCalculado)
							sumaContratosTemp += DetallesAux[i].ValorCalculado
						}
					} else {
						fmt.Println("Error al obtener detalles preliquidacion ", err)
						c.Data["mesaage"] = "Error al obtener detalles preliquidacion " + err.Error()
						c.Abort("400")
					}

					contratoAux := contrato[0]
					contratoAux.FechaInicio = contratoOriginal.FechaInicio
					sumaContratosTemp += contrato[0].ValorContrato
					contratoAux.ValorContrato = sumaContratosTemp
					// Actualiza los datos del contrato: Fecha fin y valor
					if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato[0].Id), "PUT", &aux, contratoAux); err == nil {
						if anulacion.FechaAnulacion.Day() != 30 {
							contrato[0].ValorContrato = Roundf(contrato[0].ValorContrato)
							if contrato[0].TipoNominaId == 409 {
								mensaje, err = liquidarHCH(contrato[0], false, 0, contrato[0].Vigencia)
								anularEnGenerales(contratoOriginal, anulacion.FechaAnulacion, anulacion.Vigencia)
							} else if contrato[0].TipoNominaId == 410 {
								mensaje, err = liquidarHCS(contrato[0], false, 0, contrato[0].Vigencia, semanasTotales, valorDia, true)
								anularEnGenerales(contratoOriginal, anulacion.FechaAnulacion, anulacion.Vigencia)
							}

							if err == nil {
								c.Ctx.Output.SetStatus(201)
								c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Registration successful", "Data": contrato[0]}
							} else {
								fmt.Println("Error al cancelar contrato: ", err)
								c.Data["mesaage"] = mensaje + err.Error()
								c.Abort("400")
							}
						}
					} else {
						fmt.Println("Error al crear el contrato nuevo ", err)
						c.Data["mesaage"] = "Error al crear el contrato nuevo: " + err.Error()
						c.Abort("400")
					}

				} else {
					fmt.Println("Error al ctualizar el contrato: ", err)
					c.Data["mesaage"] = "Error al Eliminar Detalles: " + err.Error()
					c.Abort("400")
				}
			} else {
				fmt.Println("Error al obtener el contrato: ", err)
				c.Data["mesaage"] = "Error al obtener el contrato:" + err.Error()
				c.Abort("400")
			}
		} else {
			fmt.Println("Error al obtener el contrato: ", err)
			c.Data["mesaage"] = "Error al Eliminar Detalles: " + err.Error()
			c.Abort("400")
		}
	} else {
		fmt.Println("Error al unmarshal del body: ", err)
		c.Data["mesaage"] = "Error al unmarshal del body:" + err.Error()
		c.Abort("400")

	}
	c.ServeJSON()
}

/*
// Post ...
// @Title Aplicar Reducción
// @Description Maneja la novedad contractual de Adición de horas para contratos de docentes de VE
// @Param	Adición	 body  models.Adicion	true	"Datos de la adición"
// @Success 201 {object} models.Contrato
// @Failure 400 the request contains incorrect syntax
// @router /aplicar_reduccion [post]
func (c *NovedadVEController) aplicarReduccion() {
}*/
