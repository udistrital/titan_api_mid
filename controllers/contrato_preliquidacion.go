package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

type CumplidoController struct {
	beego.Controller
}

func (c *CumplidoController) URLMapping() {
	c.Mapping("ActualizarCumplido", c.ActualizarCumplido)
	c.Mapping("ActualizarPreliquidado", c.ActualizarPreliquidado)
}

// Get ...
// @Title Actualizar Cumplido
// @Description Actualizar cumplido del contrato
// @Param	ano		path 	string	true		"Año de la preliquidación"
// @Param	mes		path 	string	true		"Mes de la preliquidación"
// @Param	contrato		path 	string	true		"Contrato a buscar"
// @Param	vigencia		path 	string	true		"vigencia del contrato"
// @Success 201 {object} models.ContratoPreliquidacion
// @Failure 403 body is empty
// @router /cumplido/:ano/:mes/:contrato/:vigencia [get]
func (c *CumplidoController) ActualizarCumplido() {
	var aux map[string]interface{}
	var contrato []models.Contrato
	var contrato_preliquidacion []models.ContratoPreliquidacion
	var id int
	cumplidoCompleto := true
	ano := c.Ctx.Input.Param(":ano")
	mes := c.Ctx.Input.Param(":mes")
	numeroContrato := c.Ctx.Input.Param(":contrato")
	vigencia := c.Ctx.Input.Param((":vigencia"))

	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=NumeroContrato:"+numeroContrato+",Vigencia:"+vigencia, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contrato)
		//Obtener contrato Preliquidacion para ese mes
		query := "PreliquidacionId.Ano:" + ano + ",PreliquidacionId.Mes:" + mes + ",ContratoId.NumeroContrato:" + numeroContrato
		fmt.Println(beego.AppConfig.String("UrlTitanCrud") + "/contrato_preliquidacion?limit=-1&query=" + query)
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
			//actualiar cumplido
			contrato_preliquidacion[0].Cumplido = true
			if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[0].Id), "PUT", &aux, contrato_preliquidacion[0]); err == nil {
				fmt.Println("Cumplido actualizado")
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": contrato_preliquidacion}
				query := "ContratoId.Documento:" + contrato[0].Documento + ",PreliquidacionId.Mes:" + mes + ",PreliquidacionId.Ano:" + ano
				if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
					LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
					if contrato_preliquidacion[0].Id != 0 {
						for i := 0; i < len(contrato_preliquidacion); i++ {
							if contrato_preliquidacion[i].ContratoId.NumeroContrato != "GENERAL"+mes {
								if !contrato_preliquidacion[i].Cumplido {
									cumplidoCompleto = false
								}
							} else {
								id = i
							}
						}
						if cumplidoCompleto {
							//Actualizar el cumplido del contrato General
							fmt.Println("Actualizando cumplido general")
							contrato_preliquidacion[id].Cumplido = true
							if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[id].Id), "PUT", &aux, contrato_preliquidacion[id]); err == nil {
								fmt.Println("contrato general actualizado")
							} else {
								fmt.Println("Error al actualizar cumplido general: ", err)
								c.Data["mesaage"] = "Error al actualizar cumplido general: " + err.Error()
								c.Abort("404")
							}
						}
					} else {
						fmt.Println("Error al obetener cumplidos: ", err)
						c.Data["mesaage"] = "Error al obtener cumplidos: " + err.Error()
						c.Abort("404")
					}
				} else {
					fmt.Println("Error al actualizar cumplido: ", err)
					c.Data["mesaage"] = "Error al actualizar cumplido: " + err.Error()
					c.Abort("404")
				}
			} else {
				fmt.Println("Error al actualizar cumplido: ", err)
				c.Data["mesaage"] = "Error al actualizar cumplido: " + err.Error()
				c.Abort("404")
			}

		} else {
			fmt.Println("Error al obtener contrato preliquidación")
			c.Data["mesaage"] = "Error al obtener contrato preliquidacion " + err.Error()
			c.Abort("404")
		}
	} else {
		fmt.Println("Error al obtener contrato")
		c.Data["mesaage"] = "Error al obtener contrato " + err.Error()
		c.Abort("404")
	}
	c.ServeJSON()
}

// Get ...
// @Title Actualizar Preliquidado
// @Description Actualizar preliquidado del contrato
// @Param	ano		path 	string	true		"Año de la preliquidación"
// @Param	mes		path 	string	true		"Mes de la preliquidación"
// @Param	contrato		path 	string	true		"Contrato a buscar"
// @Param	vigencia		path 	string	true		"vigencia del contrato"
// @Success 201 {object} models.ContratoPreliquidacion
// @Failure 403 body is empty
// @router /preliquidado/:ano/:mes/:contrato/:vigencia [get]
func (c *CumplidoController) ActualizarPreliquidado() {
	var aux map[string]interface{}
	var contrato []models.Contrato
	var contrato_preliquidacion []models.ContratoPreliquidacion
	ano := c.Ctx.Input.Param(":ano")
	mes := c.Ctx.Input.Param(":mes")
	numeroContrato := c.Ctx.Input.Param(":contrato")
	vigencia := c.Ctx.Input.Param((":vigencia"))

	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato?limit=-1&query=NumeroContrato:"+numeroContrato+",Vigencia:"+vigencia, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &contrato)
		//Obtener contrato Preliquidacion para ese mes
		query := "PreliquidacionId.Ano:" + ano + ",PreliquidacionId.Mes:" + mes + ",ContratoId.NumeroContrato:" + numeroContrato
		fmt.Println(beego.AppConfig.String("UrlTitanCrud") + "/contrato_preliquidacion?limit=-1&query=" + query)
		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query="+query, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contrato_preliquidacion)
			//actualiar preliquidado
			contrato_preliquidacion[0].Preliquidado = true
			if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[0].Id), "PUT", &aux, contrato_preliquidacion[0]); err == nil {
				fmt.Println("Cumplido actualizado")
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": contrato_preliquidacion}
				c.ServeJSON()
			}
		} else {
			fmt.Println("Error al obtener contrato preliquidación")
			c.Data["mesaage"] = "Error al obtener contrato preliquidacion " + err.Error()
			c.Abort("404")
		}
	} else {
		fmt.Println("Error al obtener contrato")
		c.Data["mesaage"] = "Error al obtener contrato " + err.Error()
		c.Abort("404")
	}
	c.ServeJSON()
}
