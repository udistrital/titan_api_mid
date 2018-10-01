package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
)

//  Concepto_nomina_por_personaController operations for Concepto_nomina_por_persona
type Concepto_nomina_por_personaController struct {
	beego.Controller
}

// URLMapping ...
func (c *Concepto_nomina_por_personaController) URLMapping() {
	c.Mapping("TrRegistroIncapacidades", c.TrRegistroIncapacidades)
}

// Post ...
// @Title tr_registro_incapacidades
// @Description create tr_registro_incapacidades
// @Param	body		body 	[]map[string]interface{}	true		"body for Concepto_nomina_por_persona content"
// @Success 201 {int} models.Concepto_nomina_por_persona
// @Failure 403 body is empty
// @router /tr_registro_incapacidades [post]
func (c *Concepto_nomina_por_personaController) TrRegistroIncapacidades() {
	var (
		incapacidades map[string][]map[string]interface{} // parámetro
		apiResponse   interface{}                         // respuesta del api de titan y seguridad social respectivamente
	)
	try.This(func() {
		json.Unmarshal(c.Ctx.Input.RequestBody, &incapacidades)
		err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+
			"/"+beego.AppConfig.String("Nscrud")+"/concepto_nomina_por_persona/TrConceptosPorPersona", "POST", &apiResponse, &incapacidades)
		if err != nil {
			panic(apiResponse)
		}

		idNovedades := apiResponse.(map[string]interface{})["Body"].([]interface{})
		for i, id := range idNovedades {
			incapacidades["Conceptos"][i]["Id"] = int(id.(float64))
		}
		infoNovedades := incapacidades["Conceptos"]

		err = sendJson("http://"+beego.AppConfig.String("UrlSScrud")+":"+beego.AppConfig.String("PortSS")+
			"/"+beego.AppConfig.String("NSSS")+"/detalle_novedad_seguridad_social/tr_registrar_detalle", "POST", &apiResponse, &infoNovedades)
		if err != nil {
			for _, id := range idNovedades {
				aux := int(id.(float64))
				err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+
					"/"+beego.AppConfig.String("Nscrud")+"/concepto_nomina_por_persona/"+strconv.Itoa(aux), "DELETE", &apiResponse, nil)
				if err != nil {
					panic(err.Error())
				}
			}
			panic(err.Error())
		}

		c.Data["json"] = apiResponse

	}).Catch(func(e try.E) {
		beego.Error("error en TrRegistroIncapacidades: ", e)
		c.Data["json"] = e
	})

	c.ServeJSON()
}
