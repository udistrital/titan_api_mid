package controllers

import (
	"encoding/json"

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
		// incapacidades                                    TrConceptosNomPersona // parámetro
		incapacidades                    map[string][]map[string]interface{} // parámetro
		titanCrudResponse, ssCrudReponse map[string]interface{}              // respuesta del api de titan y seguridad social respectivamente
	)
	try.This(func() {
		json.Unmarshal(c.Ctx.Input.RequestBody, &incapacidades)
		err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+
			"/"+beego.AppConfig.String("Nscrud")+"/concepto_nomina_por_persona/TrConceptosPorPersona", "POST", &titanCrudResponse, &incapacidades)
		if err != nil {
			beego.Error()
			panic(titanCrudResponse)
		}

		idNovedades := titanCrudResponse["Body"].([]interface{})
		beego.Info("idNovedades: ", idNovedades)
		for i, id := range idNovedades {
			incapacidades["Conceptos"][i]["Id"] = int(id.(float64))
		}
		infoNovedades := incapacidades["Conceptos"]
		beego.Info(infoNovedades)

		err = sendJson("http://"+beego.AppConfig.String("UrlSScrud")+":"+beego.AppConfig.String("PortSS")+
			"/"+beego.AppConfig.String("NSSS")+"/detalle_novedad_seguridad_social/tr_registrar_detalle", "POST", &ssCrudReponse, &infoNovedades)
		if err != nil {
			beego.Error("ss_crud_api error")
			for _, id := range idNovedades {
				err := sendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+
					"/"+beego.AppConfig.String("Nscrud")+"/concepto_nomina_por_persona", "DELETE", &titanCrudResponse, &id)
				if err != nil {
					panic(titanCrudResponse)
				}
			}
			panic(ssCrudReponse)
		}

		c.Data["json"] = ssCrudReponse

	}).Catch(func(e try.E) {
		beego.Error("error en TrRegistroIncapacidades: ", e)
		c.Data["json"] = incapacidades
	})

	c.ServeJSON()
}
