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

type NovedadCPSController struct {
	beego.Controller
}

func (c *NovedadCPSController) URLMapping() {
	c.Mapping("CancelarContrato", c.CancelarContrato)
	c.Mapping("CederContrato", c.CederContrato)
	c.Mapping("AplicarOtrosi", c.AplicarOtrosi)
	c.Mapping("SuspenderContrato", c.SuspenderContrato)
}

// Get ...
// @Title Registrar Cancelacion
// @Description Maneja la novedad contractual de cancelación
// @Param	novedad		body 	models.Cancelacion 	true	"Datos del contrato a cancelar"
// @Success 201 {object} models.Contrato
// @Failure 400
// @router /cancelar_contrato [post]
func (c *NovedadCPSController) CancelarContrato() {

	var cancelacion models.Cancelacion
	var aux map[string]interface{}
	var contrato []models.Contrato
	var contrato_preliquidacion []models.ContratoPreliquidacion
	var valorDia float64
	var detalles []models.DetallePreliquidacion
	var mensaje string //Mensaje de error

	//Traer el contrato a cancelar
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &cancelacion); err == nil {
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=NumeroContrato:"+cancelacion.NumeroContrato+",Vigencia:"+strconv.Itoa(cancelacion.Vigencia)+",Documento:"+cancelacion.Documento, &aux); err == nil {
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

				anoIterativo := cancelacion.FechaCancelacion.Year()
				mesIterativo := int(cancelacion.FechaCancelacion.Month())

				//Eliminar los detalles y los contratos_preliquidacion
				for {
					//Obtener contrato_preliquidacion para ese mes
					query := "ContratoId:" + strconv.Itoa(contrato[0].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo)
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
						if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
							LimpiezaRespuestaRefactor(aux, &detalles)
							for j := 0; j < len(detalles); j++ {
								if detalles[j].ConceptoNominaId.Id == 87 {
									valorDia = detalles[j].ValorCalculado / detalles[j].DiasLiquidados
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
				contrato[0].FechaFin = cancelacion.FechaCancelacion
				//Actualizar fecha de finalización del contrato
				if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato[0].Id), "PUT", &aux, contrato[0]); err == nil {
					fmt.Println("Contrato Actualizado")

					if contrato[0].FechaInicio.Month() != cancelacion.FechaCancelacion.Month() || contrato[0].FechaInicio.Year() != cancelacion.FechaCancelacion.Year() {
						contrato[0].FechaInicio = time.Date(cancelacion.FechaCancelacion.Year(), cancelacion.FechaCancelacion.Month(), 1, 12, 0, 0, 0, time.UTC)
						contrato[0].ValorContrato = valorDia * float64(contrato[0].FechaFin.Day())
					} else {
						contrato[0].ValorContrato = valorDia * float64(contrato[0].FechaFin.Day()-contrato[0].FechaInicio.Day()+1)
					}

					if cancelacion.FechaCancelacion.Day() != 30 {

						contrato[0].ValorContrato = Roundf(contrato[0].ValorContrato)
						mensaje, err = liquidarCPS(contrato[0])

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

// Post ...
// @Title Ceder Contrato
// @Description Maneja la novedad contractual de cesión de contrato
// @Param	Sucesor		body 	models.Sucesor 	true	"Datos del sucesor del contrato"
// @Success 201 {object} models.Contrato
// @Failure 403 body is empty
// @router /ceder_contrato [post]
func (c *NovedadCPSController) CederContrato() {
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
		//No se requiere el numero de documento para CPS
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=NumeroContrato:"+sucesor.NumeroContrato+",Vigencia:"+strconv.Itoa(sucesor.Vigencia)+",TipoNominaId:411", &aux); err == nil {
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
									if detalles[j].ConceptoNominaId.Id == 87 && detalles[j].DiasLiquidados == 30 {
										valorDia = detalles[j].ValorCalculado / detalles[j].DiasLiquidados
									}
									if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalles[j].Id), "DELETE", &aux, nil); err == nil {
										fmt.Println("Detalle eliminado con éxito")
									} else {
										fmt.Println("Error al eliminar detalle: ", err)
									}
								}
								if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[0].Id), "DELETE", &aux, nil); err == nil {
									fmt.Println("contrato preliquidacion eliminado con éxito")
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
					valorNuevo = 0
					query := "ContratoPreliquidacionId.ContratoId.Id:" + strconv.Itoa(contrato[0].Id) + ",ContratoPreliquidacionId.ContratoId.Vigencia:" + strconv.Itoa(contrato[0].Vigencia)
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &detalles)
						if sucesor.FechaInicio.Month() == contrato[0].FechaInicio.Month() && sucesor.FechaInicio.Year() == contrato[0].FechaInicio.Year() {
							valorNuevo = 0
							fmt.Println("Se cede el mismo mes, no hay nada pago")
						} else {
							for j := 0; j < len(detalles); j++ {
								if detalles[j].ConceptoNominaId.Id == 87 {
									valorNuevo = valorNuevo + detalles[j].ValorCalculado
								}
							}
						}

						//Agregar el nuevo contrato
						contratoNuevo.Id = 0
						contratoNuevo.Documento = sucesor.DocumentoNuevo
						contratoNuevo.NombreCompleto = sucesor.NombreCompleto
						contratoNuevo.Vigencia = contrato[0].Vigencia
						contratoNuevo.FechaInicio = sucesor.FechaInicio

						//El valor del contrato nuevo es lo que queda del contrato pasado (no se almacena, es únicamente para cálculos)
						if sucesor.FechaInicio.Day() == 1 {
							contrato[0].ValorContrato = valorNuevo
							contratoNuevo.ValorContrato = valorViejo - valorNuevo
							if int(sucesor.FechaInicio.Month())-1 == 2 {
								contrato[0].FechaFin = time.Date(sucesor.FechaInicio.Year(), sucesor.FechaInicio.Month()-1, 28, 12, 0, 0, 0, time.UTC)
							} else {
								contrato[0].FechaFin = time.Date(sucesor.FechaInicio.Year(), sucesor.FechaInicio.Month()-1, 30, 12, 0, 0, 0, time.UTC)
							}
						} else {
							contrato[0].FechaFin = sucesor.FechaInicio.Add(24 * time.Hour * -1)
							if contrato[0].FechaInicio.Month() != sucesor.FechaInicio.Month() && contrato[0].FechaInicio.Year() != sucesor.FechaInicio.Year() {
								contrato[0].ValorContrato = valorNuevo + valorDia*float64(contrato[0].FechaFin.Day())
							} else {
								contrato[0].ValorContrato = valorNuevo + valorDia*float64(contrato[0].FechaFin.Day()-contrato[0].FechaInicio.Day()+1)
							}
							contratoNuevo.ValorContrato = valorViejo - contrato[0].ValorContrato
						}

						contratoNuevo, _ = registrarContrato(contratoNuevo)

						//Actualizar fecha de finalización del contrato
						if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato[0].Id), "PUT", &aux, contrato[0]); err == nil {
							fmt.Println("Contrato Actualizado")

							if contrato[0].FechaInicio.Day() == 31 {
								contratoNuevo.FechaFin = contratoNuevo.FechaFin.Add(24 * time.Hour)
								contrato[0].FechaInicio = contrato[0].FechaInicio.Add(24 * time.Hour)
							}

							if contrato[0].FechaInicio.Month() != sucesor.FechaInicio.Month() && contrato[0].FechaInicio.Year() != sucesor.FechaInicio.Year() {
								contrato[0].FechaInicio = time.Date(sucesor.FechaInicio.Year(), sucesor.FechaInicio.Month(), 1, 12, 0, 0, 0, time.UTC)
							} else {

								contrato[0].FechaInicio = time.Date(sucesor.FechaInicio.Year(), sucesor.FechaInicio.Month(), 1, 12, 0, 0, 0, time.UTC)

							}

							//El valor del contrato nuevo es lo que queda del contrato pasado (no se almacena, es únicamente para cálculos)
							fmt.Println("valorNuevo:", valorNuevo)
							fmt.Println("valorViejo:", valorViejo)
							fmt.Println("valorDia:", valorDia)
							if sucesor.FechaInicio.Day() == 1 {
								contratoNuevo.ValorContrato = valorViejo - valorNuevo
								fmt.Println("Liquidando nuevo:", contratoNuevo.NumeroContrato, " de ", contratoNuevo.NombreCompleto)
								liquidarCPS(contratoNuevo)
							} else {
								contrato[0].ValorContrato = valorDia * float64(contrato[0].FechaFin.Day()-contrato[0].FechaInicio.Day()+1)
								contratoNuevo.ValorContrato = valorViejo - valorNuevo - contrato[0].ValorContrato
								fmt.Println("Liquidando actual:", contrato[0])
								liquidarCPS(contrato[0])
								fmt.Println("Liquidando nuevo:", contratoNuevo)
								liquidarCPS(contratoNuevo)
							}
						} else {
							fmt.Println("Error al actualizar el contrato: ", err)
							c.Data["mesaage"] = "Error al actualizar el contrato: " + err.Error()
							c.Abort("400")
						}
					} else {
						fmt.Println("Error al obtener detalles: ", err)
						c.Data["mesaage"] = "Error al obtener lo liquidado previamente: " + err.Error()
						c.Abort("400")
					}

				} else {
					fmt.Println("El contrato no puede ser cedido, debido que la fecha fin está antes de la fecha inicio")
					c.Data["mesaage"] = "El contrato no puede ser cedido, por favor verifique fechas"
					c.Abort("400")
				}
			} else {
				fmt.Println("Error al obtener el contrato: ", err)
				c.Data["mesaage"] = "Error, el contrato solicitado no existe: " + err.Error()
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
// @Param	OtroSí		body 	models.OtroSi 	true	"Documentos de contratista con Nueva Fecha de finalización del contrato"
// @Success 201 {object} models.Contrato
// @Failure 403 Fallo al obtener algún dato
// @router /otrosi_contrato [post]
func (c *NovedadCPSController) AplicarOtrosi() {
	var aux map[string]interface{}
	var otro_si models.OtroSi
	var contrato []models.Contrato
	var contratoNuevo models.Contrato
	var mensaje string

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &otro_si); err == nil {
		//traer el contrato
		//para CPS solo se requiere el numero de contrato y la vigencia
		query := "NumeroContrato:" + otro_si.NumeroContrato + ",Vigencia:" + strconv.Itoa(otro_si.Vigencia) + ",TipoNominaId:411"
		fmt.Println(beego.AppConfig.String("UrlTitanCrud") + "/contrato?limit=-1&query=" + query)
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query="+query, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contrato)
			if contrato[0].Id != 0 {
				//Organizar para tomar el más reciente
				for i := 0; i < len(contrato); i++ {
					if contrato[0].Id < contrato[i].Id {
						auxContrato := contrato[0]
						contrato[0] = contrato[i]
						contrato[i] = auxContrato
					}
				}
				contratoNuevo = contrato[0]
				//Seteamos nuevos valores
				contratoNuevo.Id = 0
				if contratoNuevo.FechaFin.Before(otro_si.FechaFin) {
					//si el contrato origninal termina el 30 se corre uno el mes y se pone el dia al 1ro
					if contrato[0].FechaFin.Day() == 30 {
						contratoNuevo.FechaInicio = time.Date(contrato[0].FechaFin.Year(), contrato[0].FechaFin.Month()+1, 01, 12, 0, 0, 0, time.Local)
					} else {
						contratoNuevo.FechaInicio = contrato[0].FechaFin.Add(24 * time.Hour)
					}
					contratoNuevo.FechaFin = otro_si.FechaFin
					contratoNuevo.Rp = otro_si.Rp
					contratoNuevo.Cdp = otro_si.Cdp
					contratoNuevo.ValorContrato = otro_si.Valor
					//Guardamos el nuevo contrato
					contratoNuevo, err = registrarContrato(contratoNuevo)
					if err == nil {
						mensaje, err = liquidarCPS(contratoNuevo)
						if err == nil {
							fmt.Println("Novedad Aplicada")
							c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Registration successful", "Data": contrato[0]}
						} else {
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contratoNuevo.Id), "DELETE", &aux, nil); err == nil {
								fmt.Println("Error al liquidar el nuevo contrato, contrato eliminado", err)
								c.Data["mesaage"] = mensaje + err.Error()
								c.Abort("403")
							} else {
								fmt.Println("Error al eliminar el contrato sin liquidar", err)
								c.Data["mesaage"] = "Error al eliminar el contrato sin liquidar: " + err.Error()
								c.Abort("403")
							}
						}
					} else {
						fmt.Println("Error al registrar el nuevo contrato", err)
						c.Data["mesaage"] = "Error al registrar el nuevo contrato: " + err.Error()
						c.Abort("403")
					}
				} else {
					fmt.Println("FECHAS INCORRECTAS", err)
					c.Data["mesaage"] = "FECHAS INCORRECTAS: " + err.Error()
					c.Abort("403")
				}
			} else {
				fmt.Println("No se encontró el contrato", err)
				c.Data["mesaage"] = "No se encontró el contrato."
				c.Abort("403")
			}
		} else {
			fmt.Println("Error al traer el contrato", err)
			c.Data["mesaage"] = "Error al traer el contrato" + err.Error()
			c.Abort("403")
		}
	} else {
		fmt.Println("Error al unmarshal de los datos del otro sí", err)
		c.Data["mesaage"] = "Los parámetros están más escritos" + err.Error()
		c.Abort("403")
	}
	c.ServeJSON()
}

// Post ...
// @Title Suspender Contrato
// @Description Maneja la novedad contractual de Suspensión de Contrato
// @Param	Suspension	 body  models.Suspension	true	"Duración de la suspensión"
// @Success 201 {object} models.Contrato
// @Failure 400 the request contains incorrect syntax
// @router /suspender_contrato [post]
func (c *NovedadCPSController) SuspenderContrato() {
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

	var mensaje string

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &suspension); err == nil {
		//Traer el contrato a cancelar
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=NumeroContrato:"+suspension.NumeroContrato+",Vigencia:"+strconv.Itoa(suspension.Vigencia)+",Documento:"+suspension.Documento, &aux); err == nil {
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
				fmt.Println("Selecciono el contrato con id: ", contrato[0].Id)

				//Igualamos el contrato Nuevo al contrato viejo
				contratoNuevo = contrato[0]
				//Valor total del contrato
				valorViejo = contrato[0].ValorContrato

				mesIterativo = int(suspension.FechaInicio.Month())
				anoIterativo = suspension.FechaInicio.Year()

				//Eliminar los detalles y los contratos_preliquidacion
				for {
					//Obtener contrato_preliquidacion para ese mes
					query := "ContratoId:" + strconv.Itoa(contrato[0].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo)
					// fmt.Println(beego.AppConfig.String("UrlTitanCrud") + "/contrato_preliquidacion?limit=-1&query=" + query)
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
						// fmt.Println(beego.AppConfig.String("UrlTitanCrud") + "/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:" + strconv.Itoa(contrato_preliquidacion[0].Id))
						if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
							LimpiezaRespuestaRefactor(aux, &detalles)
							for j := 0; j < len(detalles); j++ {
								if detalles[j].ConceptoNominaId.Id == 87 {
									valorDia = detalles[j].ValorCalculado / detalles[j].DiasLiquidados
								}
								if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalles[j].Id), "DELETE", &aux, nil); err == nil {
									fmt.Println("Detalle eliminado con éxito")
								} else {
									fmt.Println("Error al eliminar detalle: ", err)
									c.Data["mesaage"] = "Error al eliminar detalle: " + err.Error()
									c.Abort("400")
								}
							}
							fmt.Println("Valor día: ", valorDia)
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[0].Id), "DELETE", &aux, nil); err == nil {
								fmt.Println("contrato preliquidacion eliminado con éxito")
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
							} else {
								fmt.Println("Error al eliminar contrato preliquidacion: ", err)
								c.Data["mesaage"] = "Error al eliminar el contrato preliquidacion: " + err.Error()
								c.Abort("400")
							}
						} else {
							fmt.Println("Error al obtener detalles")
							c.Data["mesaage"] = "Error al obtener los detalles: " + err.Error()
							c.Abort("400")
						}
					} else {
						fmt.Println("Error al obtener contrato_preliquidacion")
						c.Data["mesaage"] = "Error al obtener el contrato preliquidacion " + err.Error()
						c.Abort("400")
					}
				}
				//Traer lo que se ha pagado hasta el momento
				query := "ContratoPreliquidacionId.ContratoId.Id:" + strconv.Itoa(contrato[0].Id) + ",ConceptoNominaId.Id:87"
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &detalles)
					if suspension.FechaInicio.Month() == contrato[0].FechaInicio.Month() && suspension.FechaInicio.Year() == contrato[0].FechaInicio.Year() {
						valorNuevo = 0
					} else {
						//Sumar lo pagado hasta el momento
						for j := 0; j < len(detalles); j++ {
							valorNuevo = valorNuevo + detalles[j].ValorCalculado
						}
					}
					fmt.Println("Liquidado hasta el momento: ", valorNuevo)

					//Calcular los días de la suspension
					diasSuspension, _ = CalcularDias(suspension.FechaInicio, suspension.FechaFin)
					diasSuspension = diasSuspension + 1 //Dia inclusive
					fmt.Println("Dias suspensión: ", diasSuspension)

					//Fecha de reanudación
					if suspension.FechaFin.Day() == 30 || suspension.FechaFin.Day() == 31 {
						aux := suspension.FechaFin.Add(24 * 30 * time.Hour)
						contratoNuevo.FechaInicio = time.Date(suspension.FechaFin.Year(), aux.Month(), 1, 12, 0, 0, 0, time.UTC)
					} else {
						contratoNuevo.FechaInicio = suspension.FechaFin.Add(24 * time.Hour)
					}

					fmt.Println("Fecha Inicio de nuevo contrato: ", contratoNuevo.FechaInicio)

					//Fecha Finalización nueva
					//Desface de dos días para febrero
					if int(suspension.FechaFin.Month()) == 2 {
						fmt.Println("Fecha anterior:", contratoNuevo.FechaFin)
						contratoNuevo.FechaFin = contratoNuevo.FechaFin.Add(24 * time.Hour * time.Duration(diasSuspension+2))
						fmt.Println("Fecha nuevo:", contratoNuevo.FechaFin)
					} else {
						fmt.Println("Fecha anterior:", contratoNuevo.FechaFin)
						contratoNuevo.FechaFin = contratoNuevo.FechaFin.Add(24 * time.Hour * time.Duration(diasSuspension))
						fmt.Println("Fecha nuevo:", contratoNuevo.FechaFin)
					}

					//Ajustar Fecha fin del contrato original
					contrato[0].FechaFin = suspension.FechaInicio.Add(24 * time.Hour * -1)

					if suspension.FechaInicio.Day() == 31 || suspension.FechaInicio.Day() == 1 {
						fmt.Println("Entro a 1")
						contrato[0].ValorContrato = valorNuevo
						contratoNuevo.ValorContrato = valorViejo - contrato[0].ValorContrato
					} else {
						fmt.Println("Entro a 2")
						fmt.Println(contrato[0].FechaFin.Day() - contrato[0].FechaInicio.Day() + 1)
						fmt.Println("Día de Fecha fin: ", contrato[0].FechaFin.Day())
						if suspension.FechaInicio.Month() == contrato[0].FechaInicio.Month() && suspension.FechaInicio.Year() == contrato[0].FechaInicio.Year() {
							contrato[0].ValorContrato = valorDia * float64(contrato[0].FechaFin.Day()-contrato[0].FechaInicio.Day()+1)
						} else {
							contrato[0].ValorContrato = valorNuevo + valorDia*float64(contrato[0].FechaFin.Day())
						}
						contratoNuevo.ValorContrato = valorViejo - contrato[0].ValorContrato
						contratoNuevo.ValorContrato = Roundf(contratoNuevo.ValorContrato)
						contrato[0].ValorContrato = Roundf(contrato[0].ValorContrato)
					}

					//Actualizar fecha de finalización y valor del contrato previo a la suspension
					if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato[0].Id), "PUT", &aux, contrato[0]); err == nil {
						fmt.Println("Contrato Actualizado")

						//Actualizar Fecha de inicio para liquidar el mes en el que se aplicó la suspensión
						if contrato[0].FechaInicio.Month() == suspension.FechaInicio.Month() && contrato[0].FechaInicio.Year() == suspension.FechaInicio.Year() {
							contrato[0].ValorContrato = valorDia * float64(contrato[0].FechaFin.Day()-contrato[0].FechaInicio.Day()+1)
						} else {
							contrato[0].FechaInicio = time.Date(suspension.FechaInicio.Year(), suspension.FechaInicio.Month(), 1, 12, 0, 0, 0, time.UTC)
							contrato[0].ValorContrato = valorDia * float64(contrato[0].FechaFin.Day())
						}

						if (suspension.FechaFin.Day() != 30 || suspension.FechaFin.Day() != 31) && suspension.FechaInicio.Day() != 1 {
							mensaje, err = liquidarCPS(contrato[0])
							if err == nil {
								fmt.Println("Contrato previo a sucesión liquidado")
							} else {
								fmt.Println("Error al liquidar el contrato anterior", err)
								c.Data["mesaage"] = mensaje + err.Error()
								c.Abort("400")
							}
						}

						contratoNuevo.Id = 0
						contratoNuevo, err = registrarContrato(contratoNuevo)
						if err == nil {
							mensaje, err = liquidarCPS(contratoNuevo)
							if err == nil {
								fmt.Println("Contrato nuevo liquidado")
								c.Ctx.Output.SetStatus(201)
								c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Registration successful", "Data": contrato[0]}
							} else {
								fmt.Println("Error al liquidar contrato luego de suspensión: ", err)
								c.Data["mesaage"] = mensaje + err.Error()
								c.Abort("400")
							}
						} else {
							fmt.Println("Error al guardar el contrato", err)
							c.Data["mesaage"] = "Error al guardar el contrato" + err.Error()
							c.Abort("400")
						}
					} else {
						fmt.Println("Error al ctualizar el contrato: ", err)
						c.Data["mesaage"] = "Error al actualizar contrato " + err.Error()
						c.Abort("400")
					}
				} else {
					fmt.Println("Error al obtener detalles")
					c.Data["mesaage"] = "Error al obtener detalles: " + err.Error()
					c.Abort("400")
				}
			} else {
				fmt.Println("Error al obtener el contrato: ", err)
				c.Data["mesaage"] = "El contrato no existe: " + err.Error()
				c.Abort("400")
			}
		} else {
			fmt.Println("Error al obtener el contrato: ", err)
			c.Data["mesaage"] = "Error al obtener el contrato: " + err.Error()
			c.Abort("400")
		}
	} else {
		fmt.Println("Error al unmarsahl de los datos del sucesor: ", err)
		c.Data["mesaage"] = "Error service POST: The request contains an incorrect data type or an invalid parameter"
		c.Abort("400")
	}
	c.ServeJSON()
}

// Post ...
// @Title Reiniciar Contrato
// @Description Maneja la novedad contractual de Reinicio de contrato
// @Param	Reinicio	 body  models.Reinicio	true	"Reinicio del contrato"
// @Success 201 {object} models.Contrato
// @Failure 400 the request contains incorrect syntax
// @router /reiniciar_contrato [post]
func (c *NovedadCPSController) ReiniciarContrato() {
	var reinicio models.Reinicio
	var aux map[string]interface{}
	var contrato []models.Contrato
	var contratoNuevo models.Contrato
	var contrato_preliquidacion []models.ContratoPreliquidacion
	var detalles []models.DetallePreliquidacion
	var valorDia float64
	var diasRestantesSuspension float64
	var mesIterativo int
	var anoIterativo int
	var mensaje string

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &reinicio); err == nil {
		//Traer el contrato a reiniciar
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=NumeroContrato:"+reinicio.NumeroContrato+",Vigencia:"+strconv.Itoa(reinicio.Vigencia)+",Documento:"+reinicio.Documento, &aux); err == nil {
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

				//Igualamos el contrato Nuevo al contrato viejo
				contratoNuevo = contrato[0]
				mesIterativo = int(contrato[0].FechaInicio.Month())
				anoIterativo = contrato[0].FechaInicio.Year()

				// Revisa si la fecha de reinicio es la misma de finalización de suspensión
				if contratoNuevo.FechaInicio.Day() == reinicio.FechaReinicio.Day() && contratoNuevo.FechaInicio.Month() == reinicio.FechaReinicio.Month() {
					// No es necesario hacer modificaciones de preliquidaciones en Titan, retorna el mismo contrato
					c.Ctx.Output.SetStatus(201)
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Registration successful", "Data": contratoNuevo}
				} else {
					//Eliminar los detalles y los contratos_preliquidacion
					for {
						//Obtener contrato_preliquidacion para ese mes
						query := "ContratoId:" + strconv.Itoa(contrato[0].Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(mesIterativo) + ",PreliquidacionId.Ano:" + strconv.Itoa(anoIterativo)
						if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
							LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
							if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query=ContratoPreliquidacionId:"+strconv.Itoa(contrato_preliquidacion[0].Id), &aux); err == nil {
								LimpiezaRespuestaRefactor(aux, &detalles)
								for j := 0; j < len(detalles); j++ {
									if detalles[j].ConceptoNominaId.Id == 87 {
										valorDia = detalles[j].ValorCalculado / detalles[j].DiasLiquidados
									}
									if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion/"+strconv.Itoa(detalles[j].Id), "DELETE", &aux, nil); err == nil {
										fmt.Println("Detalle eliminado con éxito")
									} else {
										fmt.Println("Error al eliminar detalle: ", err)
										c.Data["mesaage"] = "Error al eliminar detalle: " + err.Error()
										c.Abort("400")
									}
								}
								fmt.Println("Valor día: ", valorDia)
								if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[0].Id), "DELETE", &aux, nil); err == nil {
									fmt.Println("contrato preliquidacion eliminado con éxito")
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
								} else {
									fmt.Println("Error al eliminar contrato preliquidacion: ", err)
									c.Data["mesaage"] = "Error al eliminar el contrato preliquidacion: " + err.Error()
									c.Abort("400")
								}
							} else {
								fmt.Println("Error al obtener detalles")
								c.Data["mesaage"] = "Error al obtener los detalles: " + err.Error()
								c.Abort("400")
							}
						} else {
							fmt.Println("Error al obtener contrato_preliquidacion")
							c.Data["mesaage"] = "Error al obtener el contrato preliquidacion " + err.Error()
							c.Abort("400")
						}
					}
					//Fecha de reinicio
					diasRestantesSuspension, _ = CalcularDias(reinicio.FechaReinicio, contratoNuevo.FechaInicio)
					contratoNuevo.FechaInicio = reinicio.FechaReinicio
					contratoNuevo.FechaFin = contratoNuevo.FechaFin.Add(24 * time.Hour * (-time.Duration(diasRestantesSuspension)))

					if contratoNuevo.FechaFin.Day() == 31 {
						aux := contratoNuevo.FechaFin.Add(24 * 30 * time.Hour)
						contratoNuevo.FechaFin = time.Date(contratoNuevo.FechaFin.Year(), aux.Month(), 1, 12, 0, 0, 0, time.UTC)
					}

					mensaje, err = liquidarCPS(contratoNuevo)
					if err == nil {
						fmt.Println("Contrato despues de reinicio liquidado")
						c.Ctx.Output.SetStatus(201)
						c.Data["json"] = map[string]interface{}{"Success": true, "Status": "201", "Message": "Registration successful", "Data": contratoNuevo}
					} else {
						fmt.Println("Error al liquidar contrato luego del reinicio: ", err)
						c.Data["mesaage"] = mensaje + err.Error()
						c.Abort("400")
					}

					//Actualización de fechas despues del reinicio
					if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato/"+strconv.Itoa(contrato[0].Id), "PUT", &aux, contratoNuevo); err == nil {
						fmt.Println("Contrato Actualizado")

					} else {
						fmt.Println("Error al actualizar el contrato: ", err)
						c.Data["mesaage"] = "Error al actualizar el contrato: " + err.Error()
						c.Abort("400")
					}
				}
			} else {
				fmt.Println("Error al obtener el contrato: ", err)
				c.Data["mesaage"] = "El contrato no existe: " + err.Error()
				c.Abort("400")
			}
		} else {
			fmt.Println("Error al obtener el contrato: ", err)
			c.Data["mesaage"] = "Error al obtener el contrato: " + err.Error()
			c.Abort("400")
		}
	} else {
		fmt.Println("Error al unmarsahl de los datos del sucesor: ", err)
		c.Data["mesaage"] = "Error service POST: The request contains an incorrect data type or an invalid parameter"
		c.Abort("400")
	}
	c.ServeJSON()
}
