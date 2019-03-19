package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/manucorporat/try"
	"github.com/udistrital/utils_oas/request"
)

// Concepto_nomina_por_personaController operations for Concepto_nomina_por_persona
type Concepto_nomina_por_personaController struct {
	beego.Controller
}

// URLMapping ...
func (c *Concepto_nomina_por_personaController) URLMapping() {
	c.Mapping("TrRegistroIncapacidades", c.TrRegistroIncapacidades)
	c.Mapping("TrRegistroProrrogaIncapacidad", c.TrRegistroProrrogaIncapacidad)
}

// TrRegistroIncapacidades ...
// @Title tr_registro_incapacidades
// @Description create tr_registro_incapacidades
// @Param	body		body 	models.ConceptoNominaPorPersona	true		"body for Concepto_nomina_por_persona content"
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
		err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+
			"/"+beego.AppConfig.String("Nscrud")+"/concepto_nomina_por_persona/TrConceptosPorPersona", "POST", &apiResponse, &incapacidades)
		if err != nil {
			panic(err.Error())
		}

		idNovedades := apiResponse.(map[string]interface{})["Body"].([]interface{})
		for i, id := range idNovedades {
			incapacidades["Conceptos"][i]["Id"] = int(id.(float64))
		}
		infoNovedades := incapacidades["Conceptos"]

		err = request.SendJson("http://"+beego.AppConfig.String("UrlSScrud")+":"+beego.AppConfig.String("PortSS")+
			"/"+beego.AppConfig.String("NSSS")+"/detalle_novedad_seguridad_social/tr_registrar_detalle", "POST", &apiResponse, &infoNovedades)
		if err != nil {
			for _, id := range idNovedades {
				aux := int(id.(float64))
				err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+
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

// TrRegistroProrrogaIncapacidad ...
// @Title TrRegistroProrrogaIncapacidad
// @Description Recibe un objeto con la estructura de concepto_nomina_por_persona,
// lo envía a una transacción del crud que se encarga de cambiar el estado del registro y también crea uno nuevo,
// si la transacción del crud es correcta, se envía a un servicio del ss_crud_api que se encarga de registrar
// una nueva información en la tabla detalle_novedad_seguridad_social. En caso de que ésta última sea correcta
// finaliza la transacción, de lo contarió lo envía a una transacción del crud que se encarga de devolver el estado
// al concepto_nomina_por_persona anterior y luego elimina el nuevo registro realizado al comienzo.
// @Param	body		body 	models.ConceptoNominaPorPersona   	true		"body for Concepto_nomina_por_persona content"
// @Success 201 {int} models.Concepto_nomina_por_persona
// @Failure 403 body is empty
// @router /tr_registro_prorroga_incapacidad [post]
func (c *Concepto_nomina_por_personaController) TrRegistroProrrogaIncapacidad() {
	var (
		incapacidad,
		apiResponse map[string]interface{}
		detalleNovedad []map[string]interface{}
	)
	try.This(func() {
		json.Unmarshal(c.Ctx.Input.RequestBody, &incapacidad)
		aux := make(map[string][]map[string]interface{})
		aux["Conceptos"] = append(aux["Conceptos"], incapacidad)
		err := request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+
			"/"+beego.AppConfig.String("Nscrud")+"/concepto_nomina_por_persona/TrActualizarIncapacidadProrroga", "POST", &apiResponse, &aux)
		if err != nil {
			panic(apiResponse)
		}

		apiResponse["Body"].(map[string]interface{})["Descripcion"] = aux["Conceptos"][0]["Descripcion"]
		detalleNovedad = append(detalleNovedad, apiResponse["Body"].(map[string]interface{}))

		err = request.SendJson("http://"+beego.AppConfig.String("UrlSScrud")+":"+beego.AppConfig.String("PortSS")+
			"/"+beego.AppConfig.String("NSSS")+"/detalle_novedad_seguridad_social/tr_registrar_detalle", "POST", &apiResponse, &detalleNovedad)
		if err != nil {

			request.SendJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+
				"/"+beego.AppConfig.String("Nscrud")+"/concepto_nomina_por_persona/TrEliminarIncapacidadProrroga", "POST", &apiResponse, &detalleNovedad)
			if err != nil {
				panic(apiResponse)
			}

			panic(apiResponse)
		}

		c.Data["json"] = apiResponse
	}).Catch(func(e try.E) {
		beego.Error("Error en TrRegistroProrrogaIncapacidad(): ", e)
		c.Data["json"] = e
	})
	c.ServeJSON()
}
