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

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &novedad); err == nil {
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/concepto_nomina?limit=-1&query=Id:"+strconv.Itoa(novedad.ConceptoNominaId.Id), &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &auxConcepto)
		} else {
			fmt.Println("Error al obtener el concepto: ", err)
		}
		posible = true
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
			novedad.FechaFin = time.Date(anoFin, time.Month(mesFin), 30, 0, 0, 0, 0, time.Local)
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
		} else if auxConcepto[0].NaturalezaConceptoNominaId == 424 {
			if int(fecha_actual.Month())+novedad.Cuotas-1 > 12 {
				fmt.Println(int(fecha_actual.Month()))
				mesFin = int(fecha_actual.Month()) + novedad.Cuotas - 13
				anoFin = fecha_actual.Year() + 1
			} else {
				mesFin = int(fecha_actual.Month()) + novedad.Cuotas - 1
				anoFin = fecha_actual.Year()
			}
			novedad.FechaInicio = fecha_actual
			novedad.FechaFin = time.Date(anoFin, time.Month(mesFin), 30, 0, 0, 0, 0, time.Local)
			iteracion = int(novedad.FechaInicio.Month())
			for iteracion < int(novedad.FechaInicio.Month())+novedad.Cuotas {
				//Obtener el valor de los honorarios de ese mes
				var query = "ContratoId:" + strconv.Itoa(novedad.ContratoId.Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(iteracion) + ",PreliquidacionId.Ano:" + strconv.Itoa(fecha_actual.Year())
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contratoPreliquidacion)
					query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:10"
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &auxDetalle)
						honorarios = auxDetalle[0].ValorCalculado
						fmt.Println("Honorarios: ", honorarios)
					} else {
						fmt.Println("Error al obtener el valor de los honorarios: ", err)
					}
				} else {
					fmt.Println("Error al obtener el contrato peliquidacion: ", err)
				}
				//Obtener el Valor de los honorarios
				query = "ContratoId:" + strconv.Itoa(novedad.ContratoId.Id) + ",PreliquidacionId.Mes:" + strconv.Itoa(iteracion) + ",PreliquidacionId.Ano:" + strconv.Itoa(fecha_actual.Year())
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contratoPreliquidacion)
					query = "ContratoPreliquidacionId:" + strconv.Itoa(contratoPreliquidacion[0].Id) + ",ConceptoNominaId:2352"
					if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/detalle_preliquidacion?limit=-1&query="+query, &aux); err == nil {
						LimpiezaRespuestaRefactor(aux, &auxDetalle)
						descuentos = auxDetalle[0].ValorCalculado
						fmt.Println("Descuentos: ", descuentos)
					} else {
						fmt.Println("Error al obtener el valor de los descuentos: ", err)
					}
				} else {
					fmt.Println("Error al obtener el contrato peliquidacion: ", err)
				}
				//Verificar que el valor de los descuentos no supera la mitad de los honorarios
				if auxConcepto[0].TipoConceptoNominaId == 419 {
					fmt.Println(descuentos + novedad.Valor)
					if descuentos+novedad.Valor > (honorarios / 2) {
						posible = false
					}
				} else if auxConcepto[0].TipoConceptoNominaId == 420 {
					fmt.Println(descuentos + novedad.Valor)
					if descuentos+((novedad.Valor/100)*honorarios) > (honorarios / 2) {
						posible = false
					}
				}
				if iteracion == 12 {
					break
				} else {
					iteracion = iteracion + 1
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
