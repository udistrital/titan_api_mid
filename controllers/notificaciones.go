package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/models"
)

type NotificacionesController struct {
	beego.Controller
}

func (c *NotificacionesController) URLMapping() {
	c.Mapping("EnviarNotificacion", c.EnviarNotificacion)
}

// Get ...
// @Title Enviar Mensaje
// @Description Enviar notificacion de aprobación de nómina
// @Param	dependencia		path 	string	true		"Nómina a la que pertence la preliquidación"
// @Param	mes		path 	string	true		"Mes de la preliquidación"
// @Param	ano		path 	string	true		"Año de la preliquidación"
// @Success 201 {object} models.Detalle
// @Failure 403 body is empty
// @router /enviar_notificacion/:dependencia/:mes/:ano [get]
func (c *NotificacionesController) EnviarNotificacion() {

	dependencia := c.Ctx.Input.Param(":dependencia")
	mes := c.Ctx.Input.Param(":mes")
	ano := c.Ctx.Input.Param(":ano")

	var mensaje models.Mensaje

	mensaje.ArnTopic = "arn:aws:sns:us-east-1:699001025740:test-Titan"
	mensaje.Asunto = "Aprobación de Nómina"

	mensaje.Atributos["dependencia"] = dependencia

	mensaje.DestinatarioId[0] = 0

	mensaje.Mensaje = "Buen día. el ordenador del gasto ha aprobado la nómina para " + dependencia + "del mes de " + mes + "del año " + ano

	mensaje.RemitenteId = "Titan"

	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": nil}

	c.ServeJSON()
}
