package controllers

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GestionOpsController operations for GestionOps
type GestionOpsController struct {
	beego.Controller
}

// URLMapping ...
func (c *GestionOpsController) URLMapping() {
	c.Mapping("GenerarOrdenPago", c.GenerarOrdenPago)

}

// Get ...
// @Enviar Mensaje
// @Description Enviar orden de pago a financiera
// @Param	id		path 	string	true		"Id de la preliquidacion"
// @Param	mes		path 	string	true		"Mes de la preliquidación"
// @Param	ano		path 	string	true		"Año de la preliquidación"
// @Success 201 {object} models.Preliquidacion
// @Failure 403 body is empty
// @router /generar_op/:id [get]
func (c *GestionOpsController) GenerarOrdenPago() {

	id := c.Ctx.Input.Param(":id")
	var aux map[string]interface{}
	var preliquicion []models.Preliquidacion
	var contratos []models.ContratoPreliquidacion
	//Obtener la preliquidación de ese mes

	if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/preliquicion?limit=-1&query=Id:"+id, &aux); err == nil {
		LimpiezaRespuestaRefactor(aux, &preliquicion)
		//Verificar si quedan personas pendientes por cumplido

		if err := request.GetJson(beego.AppConfig.String("UrlTitanCrud")+"/contrato_preliquidacion?limit=-1&query=PreliquidacionId.Id:"+strconv.Itoa(preliquicion[0].Id), &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contratos)
			//Se cierra directamente la preliquidación
			preliquicion[0].EstadoPreliquidacionId = 402

			//Verificar si todos tienen cumplidos o no, pues en caso contrario se cambia a pendientes
			for i := 0; i < len(contratos); i++ {
				if !contratos[i].Cumplido {
					preliquicion[0].EstadoPreliquidacionId = 405
					break
				}
			}

			//Actualizar el estado de la preliquidación
			if err := request.SendJson((beego.AppConfig.String("UrlTitanCrud") + "/preliquidacion/" + strconv.Itoa(preliquicion[0].Id)), "PUT", &aux, preliquicion[0]); err == nil {
				fmt.Println("Preliquidación actualizada con éxito")
				//Solicitar orden de pago

				if err := request.GetJsonWSO2(beego.AppConfig.String("UrlArgoColas")+"/liquidacion/"+strconv.Itoa(preliquicion[0].Ano)+"/"+strconv.Itoa(preliquicion[0].Mes), &aux); err == nil {
					fmt.Println("Orden de pago generada")
					c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": preliquicion}

				} else {
					fmt.Println("Error al generar orden de pago: ", err)
					c.Data["message"] = "Error al generar orden de pago " + err.Error()
					c.Abort("404")
				}
			} else {
				fmt.Println("No se pudo actualizar la preliquidacion: ", err)
				c.Data["message"] = "Error al actualizar la preliquidacion " + err.Error()
				c.Abort("404")
			}

		} else {
			fmt.Println("Error al obtener contratos: ", err)
			c.Data["message"] = "Error al obtener los cumplidos de los contratos " + err.Error()
			c.Abort("404")
		}

	} else {
		fmt.Println("Error al obtener la preliquidación: ", err)
		c.Data["message"] = "Error al obtener la preliquidación " + err.Error()
		c.Abort("404")
	}

	c.ServeJSON()

}
