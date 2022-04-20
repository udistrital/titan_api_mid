package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
)

type DesagregadoHCSController struct {
	beego.Controller
}

func (c *DesagregadoHCSController) URLMapping() {
	c.Mapping("ObtenerDesagregado", c.ObtenerDesagregado)
}

// Get ...
// @Title Obtener Desagregado HCS
// @Description Obtener valores desagregados de los contratos de VE para Salarios
// @Param	body		body 	models.DatosVinculacion		true		"Dettales de la vinculación del contrato"
// @Success 201 {object} models.DesagregadoContratoHCS
// @Failure 403 body is empty
// @router / [post]
func (c *DesagregadoHCSController) ObtenerDesagregado() {

	var vinculacion models.DatosVinculacion
	var predicados []models.Predicado
	var desagregado models.DesagregadoContratoHCS

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &vinculacion); err == nil {
		predicados = append(predicados, models.Predicado{Nombre: "horas_semanales(" + strconv.Itoa(vinculacion.HorasSemanales) + ")."})
		predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + strconv.Itoa(vinculacion.Documento) + "," + strconv.Itoa(vinculacion.NumeroSemanas) + "," + strconv.Itoa(vinculacion.Vigencia) + ")."})
		reglasbase := cargarReglasBase("HCS") + FormatoReglas(predicados)
		desagregado = golog.DesagregarContrato(reglasbase, vinculacion.Categoria, strconv.Itoa(vinculacion.Documento), vinculacion.Dedicacion, strconv.Itoa(vinculacion.Vigencia))
		desagregado.NumeroContrato = vinculacion.NumeroContrato
		desagregado.Vigencia = vinculacion.Vigencia
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": desagregado}
	} else {
		fmt.Println("Error al obtener detalles del contrato")
		c.Data["mesaage"] = "Error al obtener datos del contrato seleccionado: " + err.Error()
		c.Abort("400")
	}
	c.ServeJSON()
}
