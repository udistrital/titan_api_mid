package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/udistrital/titan_api_mid/golog"
	"github.com/udistrital/titan_api_mid/models"
	"github.com/udistrital/utils_oas/request"
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

	if vinculacion.ObjetoNovedad != nil {
		var aux map[string]interface{}
		var contratoOriginal []models.Contrato
		// Primero se obtienen los porcentajes que se usaron en la vinculacion original
		query := "numero_contrato:" + vinculacion.ObjetoNovedad.VinculacionOriginal + ",vigencia:" + strconv.Itoa(vinculacion.ObjetoNovedad.VigenciaVinculacionOriginal)
		urlContrato := beego.AppConfig.String("UrlTitanCrud") + "/contrato?query=" + query
		log.Printf("Desagregar -> consultando contrato original en Titan: %s", urlContrato)
		if err := request.GetJson(urlContrato, &aux); err == nil {
			LimpiezaRespuestaRefactor(aux, &contratoOriginal)
		} else {
			log.Printf("Desagregar -> error (%s): %v", query, err)
		}

		// Si no se obtiene contrato original se aplican los valores de la vigencia actual para el calculo del desagregado
		if contratoOriginal[0].Id == 0 {
			if _, porcentajesID := ObtenerReglasPrestaciones(false); porcentajesID != 0 {
				contratoOriginal = []models.Contrato{{PorcentajesDesagregadoId: porcentajesID}}
				log.Printf("Desagregar -> no se obtuvo contrato original, usando contrato generico con PorcentajesDesagregadoId=%d", porcentajesID)
			} else {
				log.Printf("Desagregar -> no fue posible determinar un PorcentajesDesagregadoId vigente para la vinculacion original (%s)", query)
			}
		}

		switch vinculacion.ObjetoNovedad.TipoResolucion {
		case "RCAN":
			// Se debe calcular el desagregado estableciendo los predicados correspondientes de prestaciones
			// segun la cantidad de semanas de la cancelacion
			semanasOriginales := vinculacion.NumeroSemanas + vinculacion.ObjetoNovedad.SemanasNuevas
			semanasRestantes := vinculacion.ObjetoNovedad.SemanasNuevas

			reglasOriginal, lowCategoria, lowDedicacion := ConstruirReglasDesagregado(vinculacion, semanasOriginales, contratoOriginal...)
			reglasRestante, _, _ := ConstruirReglasDesagregado(vinculacion, semanasRestantes, contratoOriginal...)

			desagregadoOriginal := golog.DesagregarContrato(reglasOriginal, lowCategoria, vinculacion.Documento, lowDedicacion, strconv.Itoa(vinculacion.Vigencia))
			desagregadoRestante := golog.DesagregarContrato(reglasRestante, lowCategoria, vinculacion.Documento, lowDedicacion, strconv.Itoa(vinculacion.Vigencia))

			// Calcular la diferencia
			desagregado.NumeroContrato = vinculacion.NumeroContrato
			desagregado.Vigencia = vinculacion.Vigencia
			desagregado.SueldoBasico = desagregadoOriginal.SueldoBasico - desagregadoRestante.SueldoBasico
			desagregado.PrimaServicios = desagregadoOriginal.PrimaServicios - desagregadoRestante.PrimaServicios
			desagregado.PrimaVacaciones = desagregadoOriginal.PrimaVacaciones - desagregadoRestante.PrimaVacaciones
			desagregado.InteresesCesantias = desagregadoOriginal.InteresesCesantias - desagregadoRestante.InteresesCesantias
			desagregado.Cesantias = desagregadoOriginal.Cesantias - desagregadoRestante.Cesantias
			desagregado.Vacaciones = desagregadoOriginal.Vacaciones - desagregadoRestante.Vacaciones
			desagregado.PrimaNavidad = desagregadoOriginal.PrimaNavidad - desagregadoRestante.PrimaNavidad
			desagregado.BonificacionServicios = desagregadoOriginal.BonificacionServicios - desagregadoRestante.BonificacionServicios

		case "RADD", "RRED":
			fmt.Println("Desagregado por RADD y RRED")
			// Para la adicion y reduccion solo se le pasa los datos del nuevo contrato
			// Si se aplicaron los porcentajes mayores en la original se deben mantener
			// Si se aplicaron los porcentajes menores en la original se deben mantener

			reglasBase, lowCategoria, lowDedicacion := ConstruirReglasDesagregado(vinculacion, vinculacion.NumeroSemanas, contratoOriginal...)
			desagregado = golog.DesagregarContrato(reglasBase, lowCategoria, vinculacion.Documento, lowDedicacion, strconv.Itoa(vinculacion.Vigencia))
			desagregado.NumeroContrato = vinculacion.NumeroContrato
			desagregado.Vigencia = vinculacion.Vigencia

		}

	} else {
		reglasBase, lowCategoria, lowDedicacion := ConstruirReglasDesagregado(vinculacion, vinculacion.NumeroSemanas)
		desagregado = golog.DesagregarContrato(reglasBase, lowCategoria, vinculacion.Documento, lowDedicacion, strconv.Itoa(vinculacion.Vigencia))
		desagregado.NumeroContrato = vinculacion.NumeroContrato
		desagregado.Vigencia = vinculacion.Vigencia
	}
	return desagregado
}
