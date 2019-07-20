package controllers

import (
	"encoding/json"
	//"strconv"
	"fmt"
	"github.com/astaxie/beego"
	//"github.com/manucorporat/try"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/titan_api_mid/models"
)

// ServicesController operations for Services
type ServicesController struct {
	beego.Controller
}

// URLMapping ...
func (c *ServicesController) URLMapping() {
	c.Mapping("DesagregacionContratoHCS", c.DesagregacionContratoHCS)

}

// DesagregadoContratoHCS ...
// @Title DesagregadoContratoHCS
// @Description Dado un valor de contrato para docente de hora cátedra salarios, su fecha de inicio y fin y una vigencia, se retorna el valor por concepto que le será pagado en la totalidad de su vinculación con la universidad. Los conceptos a mostrar son: sueldo básico, vacaciones, prima de vacaciones, prima de servicios, intereses sobre cesantías y cesantías.
// @Param	body		body 	models.InformacionContratoDocente	true		"body for Services content"
// @Success 201 {int} models.DesagregadoContratoHCS
// @Failure 403 body is empty
// @router /desagregacion_contrato_hcs [post]
func (c *ServicesController) DesagregacionContratoHCS() {
			var v models.InformacionContratoDocente
				var h []models.ConceptoNominaPorPersona
			if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
					fmt.Println("info cont doc:", v)
					if err := request.GetJson("http://"+beego.AppConfig.String("Urlcrud")+":"+beego.AppConfig.String("Portcrud")+"/"+beego.AppConfig.String("Nscrud")+"/concepto_nomina_por_persona?limit=-1", &h); err == nil {
						fmt.Println("prueba")
					}
			}else{
				var e models.Alert
				e.Body = err.Error();
			}
}
