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
// @router /:ano/:mes/:contrato/:vigencia [get]
func (c *CumplidoController) ActualizarCumplido() {
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
			//actualiar cumplido
			contrato_preliquidacion[0].Cumplido = true
			if err := request.SendJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion/"+strconv.Itoa(contrato_preliquidacion[0].Id), "PUT", &aux, contrato_preliquidacion[0]); err == nil {
				fmt.Println("Cumplido actualizado")
				c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": contrato_preliquidacion}
				c.ServeJSON()
			}
		} else {
			fmt.Println("Error al obtener contrato preliquidación")
		}
	} else {
		fmt.Println("Error al obtener contrato")
	}
}
