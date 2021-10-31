package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/golog"
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

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {

		// 1. Se cargan reglas para HCS
		reglasbase := cargarReglasBase("HCS") //funcion general para dar formato a reglas cargadas desde el ruler
		reglasbase = reglasbase + cargarReglasSS()

		totales := golog.CalcularTotalesContratoHCS(v.NumDocumento, v.MesesContrato, v.VigenciaContrato, v.ValorContrato, reglasbase)

		c.Data["json"] = totales

	} else {
		var e models.Alert
		e.Body = err.Error()
		c.Data["json"] = e
	}

	c.ServeJSON()
}
