package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
	var desagregado models.DesagregadoContratoHCS

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &vinculacion); err == nil {
		desagregado = Desagregar(vinculacion)

		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": desagregado}
	} else {
		fmt.Println("Error al obtener detalles del contrato")
		c.Data["mesaage"] = "Error al obtener datos del contrato seleccionado: " + err.Error()
		c.Abort("400")
	}
	c.ServeJSON()
}

func Desagregar(vinculacion models.DatosVinculacion) (desagregado models.DesagregadoContratoHCS) {

	var predicados []models.Predicado
	var lowCategoria string  //Categoría en minúscula
	var lowDedicacion string //Dedicacion en minúscula

	if vinculacion.Dedicacion == "HCP" {
		if vinculacion.NivelAcademico == "POSGRADO" {
			lowDedicacion = "hcpos"
		} else {
			lowDedicacion = "hcpre"
		}
		predicados = append(predicados, models.Predicado{Nombre: "aplica_prima(0)."})
	} else {
		lowDedicacion = strings.ToLower(vinculacion.Dedicacion)
		if vinculacion.Dedicacion == "MTO" {
			predicados = append(predicados, models.Predicado{Nombre: "aplica_prima(0)."})
		} else {
			predicados = append(predicados, models.Predicado{Nombre: "aplica_prima(1)."})
		}
		if vinculacion.Cancelacion {
			predicados = append(predicados, models.Predicado{Nombre: "cancelacion(1)."})
		} else {
			predicados = append(predicados, models.Predicado{Nombre: "cancelacion(0)."})
		}
	}

	lowCategoria = strings.ToLower(vinculacion.Categoria)
	predicados = append(predicados, models.Predicado{Nombre: "horas_semanales(" + strconv.Itoa(vinculacion.HorasSemanales) + ")."})
	predicados = append(predicados, models.Predicado{Nombre: "duracion_contrato(" + vinculacion.Documento + "," + strconv.Itoa(vinculacion.NumeroSemanas) + "," + strconv.Itoa(vinculacion.Vigencia) + ")."})
	predicados = append(predicados, models.Predicado{Nombre: "valor_punto(" + strconv.Itoa(vinculacion.Vigencia) + "," + strconv.Itoa(int(vinculacion.PuntoSalarial)) + ")."})
	reglasbase := cargarReglasBase("HCS") + FormatoReglas(predicados)
	desagregado = golog.DesagregarContrato(reglasbase, lowCategoria, vinculacion.Documento, lowDedicacion, strconv.Itoa(vinculacion.Vigencia))
	desagregado.NumeroContrato = vinculacion.NumeroContrato
	desagregado.Vigencia = vinculacion.Vigencia
	return desagregado
}
